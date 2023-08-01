/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package launcher

import (
	"log"
	"strconv"
	"strings"

	"github.com/galasa.dev/cli/pkg/api"
	"github.com/galasa.dev/cli/pkg/embedded"
	"github.com/galasa.dev/cli/pkg/errors"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/props"
	"github.com/galasa.dev/cli/pkg/utils"
)

// TestLocation the user passes us this information in one string.
// We split it up into these useful chunks.
type TestLocation struct {
	OSGiBundleName         string
	QualifiedJavaClassName string
}

// We gather a summary of passes and failures from the
// test results we read from ras/<runId>/structure.json
type TestResultsSummary struct {
	MethodPasses int
	MethodFails  int
}

// JvmLauncher can act as a launcher, it's given test cases which need to
// be executed, and it launches them within a local JVM.
type JvmLauncher struct {
	// The fully-qualified path to JAVA_HOME where we can find the bin/java command.
	javaHome string

	// The parameters from the command-line.
	cmdParams RunsSubmitLocalCmdParameters

	// An abstraction of the environment, so we can look up things like JAVA_HOME
	env utils.Environment

	// An abstraction of the file system so we can mock it out easily for unit tests.
	fileSystem files.FileSystem

	// A location galasa can call home
	galasaHome utils.GalasaHome

	// A file system so we can get at embedded content if required.
	// (Like so we can unpack the boot.jar)
	embeddedFileSystem embedded.ReadOnlyFileSystem

	// The collection of tests which are running, or have completed.
	localTests []*LocalTest

	// This timer service can be interrupted when we don't want it to sleep.
	timeService utils.TimeService

	// A service which can create OS processes.
	processFactory ProcessFactory

	// A map of bootstrap properties
	bootstrapProps props.JavaProperties
}

// These parameters are gathered from the command-line and passed into the laucher.
type RunsSubmitLocalCmdParameters struct {

	// A list of OBRs, which we hope one of these contains the tests we want to run.
	Obrs []string

	// The remote maven repo, eg: maven central, where we can load the galasa uber-obr
	RemoteMaven string

	// The version of galasa we want to launch. This indicates which uber-obr will be
	// loaded.
	TargetGalasaVersion string

	// Should the JVM be launched in debug mode ?
	IsDebugEnabled bool

	// When launched in debug mode, which port should the JVM use to talk to the Java
	// debugger ? This port is either listened on, or attached to depending on the
	// DebugMode field.
	DebugPort uint32

	// A string indicating whether the test JVM should 'attach' to the debug port
	// to talk to the Java debugger (JDB), or whether it should 'listen' on a port
	// ready for the JDB to attach to.
	DebugMode string
}

const (
	DEBUG_PORT_DEFAULT uint32 = 2970
)

// -----------------------------------------------------------------------------
// Constructors
// -----------------------------------------------------------------------------

// NewJVMLauncher creates a JVM launcher. Primes it with references to services
// which can be used to launch JVM servers.
func NewJVMLauncher(
	bootstrapProps props.JavaProperties,
	env utils.Environment,
	fileSystem files.FileSystem,
	embeddedFileSystem embedded.ReadOnlyFileSystem,
	runsSubmitLocalCmdParams RunsSubmitLocalCmdParameters,
	timeService utils.TimeService,
	processFactory ProcessFactory,
	galasaHome utils.GalasaHome,
) (*JvmLauncher, error) {

	var (
		err      error        = nil
		launcher *JvmLauncher = nil
	)

	javaHome := env.GetEnv("JAVA_HOME")

	err = utils.ValidateJavaHome(fileSystem, javaHome)

	if err == nil {
		launcher = new(JvmLauncher)
		launcher.javaHome = javaHome
		launcher.cmdParams = runsSubmitLocalCmdParams
		launcher.env = env
		launcher.fileSystem = fileSystem
		launcher.embeddedFileSystem = embeddedFileSystem
		launcher.processFactory = processFactory
		launcher.galasaHome = galasaHome
		launcher.timeService = timeService
		launcher.bootstrapProps = bootstrapProps

		// Make sure the home folder has the boot jar unpacked and ready to invoke.
		err = utils.InitialiseGalasaHomeFolder(
			launcher.galasaHome,
			launcher.fileSystem,
			launcher.embeddedFileSystem,
		)
	}

	return launcher, err
}

//-----------------------------------------------------------------------------
// Implementation of the Launcher interface
//-----------------------------------------------------------------------------

// SubmitTestRuns launch the test runs
//
// groupName - The run group ID. Used to group all the TestRuns together so we
// can query the results later.
//
// classNames - An array of strings in the form "<osgi-bundle-id>/<fully-qualified-java-classname>
// Note: There is no ".class" suffix needed for each entry. That is assumed.
//
// requestType - A metadata marker to indicate how the testRun was scheduled.
// requestor - Who wanted the testRun to launch.
// stream - The stream the test run is part of
// isTraceEnabled - True of the trace for the test run should be gathered.
// overrides - A map of overrides of key-value pairs.
func (launcher *JvmLauncher) SubmitTestRuns(
	groupName string,
	classNames []string,
	requestType string,
	requestor string,
	stream string,
	isTraceEnabled bool,
	overrides map[string]interface{},
) (*galasaapi.TestRuns, error) {

	log.Printf("JvmLauncher: SubmitTestRuns entered. group=%s classNames=%v "+
		"requestType=%s requestor=%s stream=%s isTraceEnabled=%v",
		groupName, classNames, requestType,
		requestor, stream, isTraceEnabled)

	testRuns := new(galasaapi.TestRuns)

	var obrs []utils.MavenCoordinates
	obrs, err := validateObrs(launcher.cmdParams.Obrs)
	if err == nil {

		var (
			overridesFilePath   string
			temporaryFolderPath string
		)
		temporaryFolderPath, overridesFilePath, err = prepareTempFiles(
			launcher.galasaHome, launcher.fileSystem, overrides)
		if err == nil {

			defer func() {
				deleteTempFiles(launcher.fileSystem, temporaryFolderPath)
			}()

			isComplete := false
			testRuns.Complete = &isComplete
			testRuns.Runs = make([]galasaapi.TestRun, 0)

			for _, classNameUserInput := range classNames {

				var testClassToLaunch *TestLocation
				testClassToLaunch, err = classNameUserInputToTestClassLocation(classNameUserInput)

				if err == nil {
					var (
						cmd  string
						args []string
					)
					cmd, args, err = getCommandSyntax(
						launcher.bootstrapProps,
						launcher.galasaHome,
						launcher.fileSystem, launcher.javaHome, obrs,
						*testClassToLaunch, launcher.cmdParams.RemoteMaven,
						launcher.cmdParams.TargetGalasaVersion, overridesFilePath,
						isTraceEnabled,
						launcher.cmdParams.IsDebugEnabled,
						launcher.cmdParams.DebugPort,
						launcher.cmdParams.DebugMode,
					)
					if err == nil {
						log.Printf("Launching command '%s' '%v'\n", cmd, args)
						localTest := NewLocalTest(launcher.timeService, launcher.fileSystem, launcher.processFactory)
						err = localTest.launch(cmd, args)

						if err == nil {
							// The JVM process started. Store away its' details
							launcher.localTests = append(launcher.localTests, localTest)

							localTest.testRun = new(galasaapi.TestRun)
							localTest.testRun.SetBundleName(testClassToLaunch.OSGiBundleName)
							localTest.testRun.SetStream(stream)
							localTest.testRun.SetGroup(groupName)
							localTest.testRun.SetRequestor(requestor)
							localTest.testRun.SetTrace(isTraceEnabled)
							localTest.testRun.SetType(requestType)
							localTest.testRun.SetName(localTest.runId)

							// The test run we started can be returned to the submitter.
							testRuns.Runs = append(testRuns.Runs, *localTest.testRun)
						}
					}
				}

				if err != nil {
					break
				}
			}
		}
	}

	return testRuns, err
}

func deleteTempFiles(fileSystem files.FileSystem, temporaryFolderPath string) {
	fileSystem.DeleteDir(temporaryFolderPath)
}

func prepareTempFiles(
	galasaHome utils.GalasaHome,
	fileSystem files.FileSystem,
	overrides map[string]interface{},
) (string, string, error) {

	var (
		temporaryFolderPath string = ""
		overridesFilePath   string = ""
		err                 error  = nil
	)

	// Create a temporary folder
	temporaryFolderPath, err = fileSystem.MkTempDir()
	if err == nil {
		overridesFilePath, err = createTemporaryOverridesFile(
			temporaryFolderPath, galasaHome, fileSystem, overrides)
	}

	// Clean up the temporary folder if we failed to create the props file.
	if err != nil {
		fileSystem.DeleteDir(temporaryFolderPath)
	}

	return temporaryFolderPath, overridesFilePath, err
}

// createTemporaryOverridesFile Gathers up all the overrides properties and puts
// them into a temporary file in a temporary folder.
//
// Returns:
// - the full path to the new overrides file
// - error if there was one.
func createTemporaryOverridesFile(
	temporaryFolderPath string,
	galasaHome utils.GalasaHome,
	fileSystem files.FileSystem,
	overrides map[string]interface{},
) (string, error) {
	overrides = addStandardOverrideProperties(galasaHome, fileSystem, overrides)

	// Write the properties to a file
	overridesFilePath := temporaryFolderPath + "overrides.properties"
	err := props.WritePropertiesFile(fileSystem, overridesFilePath, overrides)
	return overridesFilePath, err
}

func addStandardOverrideProperties(
	galasaHome utils.GalasaHome,
	fileSystem files.FileSystem,
	overrides map[string]interface{},
) map[string]interface{} {

	overrideRasStoreProperty(galasaHome, overrides)
	overrideLocalRunIdPrefixProperty(overrides)
	override3270TerminalOutputFormat(overrides)

	return overrides
}

func override3270TerminalOutputFormat(overrides map[string]interface{}) {
	// Force the launched runs to use the "L" prefix in their runids.
	const OVERRIDE_PROPERTY_3270_TERMINAL_OUTPUT_FORMAT = "zos3270.terminal.output"

	// Only set this property if it's not already set by the user, or in the users' override file.
	_, isPropAlreadySet := overrides[OVERRIDE_PROPERTY_3270_TERMINAL_OUTPUT_FORMAT]
	if !isPropAlreadySet {
		overrides[OVERRIDE_PROPERTY_3270_TERMINAL_OUTPUT_FORMAT] = "json,png"
	}
}

func overrideLocalRunIdPrefixProperty(overrides map[string]interface{}) {
	// Force the launched runs to use the "L" prefix in their runids.
	const OVERRIDE_PROPERTY_LOCAL_RUNID_PREFIX = "framework.request.type.LOCAL.prefix"

	// Only set this property if it's not already set by the user, or in the users' override file.
	_, isPropAlreadySet := overrides[OVERRIDE_PROPERTY_LOCAL_RUNID_PREFIX]
	if !isPropAlreadySet {
		overrides[OVERRIDE_PROPERTY_LOCAL_RUNID_PREFIX] = "L"
	}
}

func overrideRasStoreProperty(galasaHome utils.GalasaHome, overrides map[string]interface{}) {
	// Set the ras location to be local disk always.
	const OVERRIDE_PROPERTY_FRAMEWORK_RESULT_STORE = "framework.resultarchive.store"

	// Only set this property if it's not already set by the user, or in the users' override file.
	{
		_, isRasPropAlreadySet := overrides[OVERRIDE_PROPERTY_FRAMEWORK_RESULT_STORE]
		if !isRasPropAlreadySet {
			rasPathUri := "file:///" + galasaHome.GetUrlFolderPath() + "/ras"
			overrides[OVERRIDE_PROPERTY_FRAMEWORK_RESULT_STORE] = rasPathUri
		}
	}
}

func (launcher *JvmLauncher) GetRunsByGroup(groupName string) (*galasaapi.TestRuns, error) {
	log.Printf("JvmLauncher: GetRunsByGroup(groupName=%s) entered. ", groupName)

	var isAllComplete = true
	var testRuns = galasaapi.TestRuns{
		Complete: &isAllComplete,
		Runs:     []galasaapi.TestRun{},
	}

	for _, localTest := range launcher.localTests {

		testName := localTest.testRun.GetName()

		if localTest.isCompleted() {
			log.Printf("GetRunsByGroup: localTest %s is complete.\n", testName)
		} else {
			log.Printf("GetRunsByGroup: localTest %s is not yet complete.\n", testName)
			isAllComplete = false
		}

		var testRun *galasaapi.TestRun
		if localTest.testRun != nil {
			testRun = localTest.testRun
		} else {
			testRun = createSimulatedTestRun(testName)
		}

		testRuns.Runs = append(testRuns.Runs, *testRun)
	}

	log.Printf("JvmLauncher: GetRunsByGroup(groupName=%s) exiting. isComplete:%v testRuns returned:%v\n", groupName, *testRuns.Complete, len(testRuns.Runs))
	for _, testRun := range testRuns.Runs {
		log.Printf("JvmLauncher: GetRunsByGroup test name:%s status:%s result:%s\n", *testRun.Name, *testRun.Status, *testRun.Result)
	}
	return &testRuns, nil
}

// GetRunsById gets the Run information for the run with a specific run identifier
func (launcher *JvmLauncher) GetRunsById(runId string) (*galasaapi.Run, error) {
	log.Printf("JvmLauncher: GetRunsById entered. runId=%s", runId)
	return nil, nil
}

// GetStreams gets a list of streams available on this launcher
func (launcher *JvmLauncher) GetStreams() ([]string, error) {
	log.Printf("JvmLauncher: GetStreams entered.")
	return nil, nil
}

// GetTestCatalog gets the test catalog for a given stream.
func (launcher *JvmLauncher) GetTestCatalog(stream string) (TestCatalog, error) {
	log.Printf("JvmLauncher: GetTestCatalog entered. stream=%s", stream)
	return nil, nil
}

// -----------------------------------------------------------------------------
// Local functions
// -----------------------------------------------------------------------------

// We expect a parameter to be of the form:
// mvn:dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOT/obr
// Validate that the --obr parameter(s) passed by the user conform to this convention by splitting the
// input into pieces.
func validateObrs(obrInputs []string) ([]utils.MavenCoordinates, error) {

	var err error = nil
	obrs := make([]utils.MavenCoordinates, 0)

	for _, obr := range obrInputs {
		parts := strings.Split(obr, "/")
		if len(parts) < 4 {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OBR_NOT_ENOUGH_PARTS, obr)
		} else if len(parts) > 4 {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OBR_TOO_MANY_PARTS, obr)
		} else if !strings.HasPrefix(parts[0], "mvn:") {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OBR_NO_MVN_PREFIX, obr)
		} else if parts[3] != "obr" {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OBR_NO_OBR_SUFFIX, obr)
		} else {
			groupId := strings.ReplaceAll(parts[0], "mvn:", "")
			coordinates := utils.MavenCoordinates{
				GroupId:    groupId,
				ArtifactId: parts[1],
				Version:    parts[2],
				Classifier: parts[3],
			}

			obrs = append(obrs, coordinates)
		}
	}
	return obrs, err
}

// getCommandSyntax From the parameters we aim to build a command-line incantation which would launch the test in a JVM...
// For example:
// java -jar ${BOOT_JAR_PATH} \
// --localmaven file:${M2_PATH}/repository/ \
// --remotemaven $REMOTE_MAVEN \
// --bootstrap file:${HOME}/.galasa/bootstrap.properties \
// --overrides file:${HOME}/.galasa/overrides.properties \
// --obr mvn:dev.galasa/dev.galasa.uber.obr/${OBR_VERSION}/obr \
// --obr mvn:${TEST_OBR_GROUP_ID}/${TEST_OBR_ARTIFACT_ID}/${TEST_OBR_VERSION}/obr \
// --test ${TEST_BUNDLE}/${TEST_JAVA_CLASS} | tee jvm-log.txt | grep "[*][*][*]" | grep -v "[*][*][*][*]" | sed -e "s/[--]*//g"
//
// For example:
//
//	java -jar /Users/mcobbett/builds/galasa/code/external/galasa-dev/cli/pkg/embedded/templates/galasahome/lib/galasa-boot-0.27.0.jar \
//	    --localmaven file:/Users/mcobbett/.m2/repository/ \
//	    --remotemaven https://development.galasa.dev/main/maven-repo/obr/ \
//	    --bootstrap file:/Users/mcobbett/.galasa/bootstrap.properties \
//	    --overrides file:/Users/mcobbett/.galasa/overrides.properties \
//	    --obr mvn:dev.galasa/dev.galasa.uber.obr/0.26.0/obr \
//	    --obr mvn:dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOT/obr \
//	    --test dev.galasa.example.banking.payee/dev.galasa.example.banking.payee.TestPayee
func getCommandSyntax(
	bootstrapProperties props.JavaProperties,
	galasaHome utils.GalasaHome,
	fileSystem files.FileSystem,
	javaHome string,
	testObrs []utils.MavenCoordinates,
	testLocation TestLocation,
	remoteMaven string,
	galasaVersionToRun string,
	overridesFilePath string,
	isTraceEnabled bool,
	isDebugEnabled bool,
	debugPort uint32,
	debugMode string,
) (string, []string, error) {

	var cmd string = ""
	var args []string = make([]string, 0)
	var err error
	var bootJarPath string

	// Gather any variable values we need to.
	debugMode, err = calculateDebugMode(debugMode, bootstrapProperties)
	if err == nil {
		debugPort, err = calculateDebugPort(debugPort, bootstrapProperties)
		if err == nil {
			bootJarPath, err = utils.GetGalasaBootJarPath(fileSystem, galasaHome)
		}
	}

	if err == nil {

		separator := fileSystem.GetFilePathSeparator()

		// Note: Even in windows, when the java executable is called 'java.exe'
		// You don't need to add the '.exe' extension it seems.
		cmd = javaHome + separator + "bin" +
			separator + "java"

		args = appendArgsDebugOptions(args, isDebugEnabled, debugMode, debugPort)

		args = appendArgsBootstrapJvmLaunchOptions(args, bootstrapProperties)

		args = append(args, "-jar")
		args = append(args, bootJarPath)

		args = append(args, "-Dfile.encoding=UTF-8")

		var userHome string
		userHome, err = fileSystem.GetUserHomeDirPath()

		nativeGalasaHomeFolderPath := galasaHome.GetNativeFolderPath()
		args = append(args, `-DGALASA_HOME="`+nativeGalasaHomeFolderPath+`"`)

		// --localmaven file://${M2_PATH}/repository/
		// Note: URLs always have forward-slashes
		args = append(args, "--localmaven")
		localMavenPath := "file:///" + strings.ReplaceAll(userHome, "\\", "/") + "/.m2/repository"
		args = append(args, localMavenPath)

		// --remotemaven $REMOTE_MAVEN
		args = append(args, "--remotemaven")
		args = append(args, remoteMaven)

		// --bootstrap file:${HOME}/.galasa/bootstrap.properties
		args = append(args, "--bootstrap")
		bootstrapPath := "file:///" + galasaHome.GetUrlFolderPath() + "/bootstrap.properties"
		args = append(args, bootstrapPath)

		// --overrides file:${HOME}/.galasa/overrides.properties
		args = append(args, "--overrides")
		// Note: We turn the file path provided into a URL so slashes always
		// go the same way.
		overridesPath := "file:///" + strings.ReplaceAll(overridesFilePath, "\\", "/")
		args = append(args, overridesPath)

		for _, obrCoordinate := range testObrs {
			// We are aiming for this:
			// mvn:${TEST_OBR_GROUP_ID}/${TEST_OBR_ARTIFACT_ID}/${TEST_OBR_VERSION}/obr
			args = append(args, "--obr")
			obrMvnPath := "mvn:" + obrCoordinate.GroupId + "/" +
				obrCoordinate.ArtifactId + "/" + obrCoordinate.Version + "/obr"
			args = append(args, obrMvnPath)
		}

		// --obr mvn:dev.galasa/dev.galasa.uber.obr/${OBR_VERSION}/obr
		args = append(args, "--obr")
		galasaUberObrPath := "mvn:dev.galasa/dev.galasa.uber.obr/" + galasaVersionToRun + "/obr"
		args = append(args, galasaUberObrPath)

		// --test ${TEST_BUNDLE}/${TEST_JAVA_CLASS}
		args = append(args, "--test")
		args = append(args, testLocation.OSGiBundleName+"/"+testLocation.QualifiedJavaClassName)

		if isTraceEnabled {
			args = append(args, "--trace")
		}

	}

	return cmd, args, err
}

func appendArgsDebugOptions(args []string, isDebugEnabled bool, debugMode string, debugPort uint32) []string {

	if isDebugEnabled {
		var buff strings.Builder

		buff.WriteString("-agentlib:jdwp=transport=dt_socket,address=*:")
		buff.WriteString(strconv.FormatUint(uint64(debugPort), 10))
		buff.WriteString(",server=")
		if debugMode == "listen" {
			buff.WriteString("y")
		} else {
			buff.WriteString("n")
		}
		buff.WriteString(",suspend=y")

		args = append(args, buff.String())
	}

	return args
}

func appendArgsBootstrapJvmLaunchOptions(args []string, bootstrapProperties props.JavaProperties) []string {
	// Append all the java launch properties explicitly spelt-out in the boostrap file.
	// The framework.jvm.local.launch.options bootstrap file property can add parameters to the commmand-line.
	// For example -Xmx80m and similar parameters.
	// Use a space-separated list of options and the JVM gets launched with those in front.
	jvmLaunchOptions, isOptionsPresent := bootstrapProperties[api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_OPTIONS]
	if isOptionsPresent {
		// strip off the leading and trailing whitespace.
		jvmLaunchOptions = strings.Trim(jvmLaunchOptions, " \t\n\r")

		// Split based on commas
		launchOptionParts := strings.Split(jvmLaunchOptions, api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_OPTIONS_SEPARATOR)

		// Add each piece to the list of args returned.
		args = append(args, launchOptionParts...)
	}

	return args
}

func calculateDebugPort(debugPort uint32, bootstrapProperties props.JavaProperties) (uint32, error) {
	var err error

	if debugPort == 0 {
		// Debug port was not set on the command-line.

		// Look in the bootstrap properties for a value.
		bootstrapPropsValue, isPresent := bootstrapProperties[api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_PORT]
		if isPresent {
			// Not specified on command line. Use value in bootstrap property instead.
			var debugPortU64 uint64
			debugPortU64, err = strconv.ParseUint(bootstrapPropsValue, 10, 32)
			if err != nil {
				err = errors.NewGalasaError(
					errors.GALASA_ERROR_BOOTSTRAP_BAD_DEBUG_PORT_VALUE,
					bootstrapPropsValue,
					api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_PORT,
					strconv.FormatUint(uint64(DEBUG_PORT_DEFAULT), 10),
				)
			} else {
				// Bootstrap property value is good.
				debugPort = uint32(debugPortU64)
			}
		} else {
			// Not specified on command-linem, nothing in bootstrap property.
			debugPort = DEBUG_PORT_DEFAULT
		}
	}
	return debugPort, err
}

func calculateDebugMode(debugMode string, bootstrapProperties props.JavaProperties) (string, error) {
	var err error

	if debugMode == "" {
		// The value hasn't been set on the command-line.

		// Look in the bootstrap properties for a value.
		bootstrapPropsValue, isPresent := bootstrapProperties[api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_MODE]
		if isPresent {
			debugMode = bootstrapPropsValue
			err = checkDebugModeValueIsValid(debugMode, errors.GALASA_ERROR_BOOTSTRAP_BAD_DEBUG_MODE_VALUE)
		} else {
			// Default to 'listen'
			debugMode = "listen"
		}
	}

	if err == nil {
		err = checkDebugModeValueIsValid(debugMode, errors.GALASA_ERROR_ARG_BAD_DEBUG_MODE_VALUE)
	}

	return debugMode, err
}

func checkDebugModeValueIsValid(debugMode string, errorMessageIfInvalid *errors.MessageType) error {
	var err error = nil

	lowerCaseDebugMode := strings.ToLower(debugMode)

	switch lowerCaseDebugMode {
	case "listen":
	case "attach":
		break
	default:
		err = errors.NewGalasaError(errorMessageIfInvalid, debugMode, api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_MODE)
	}

	return err
}

// User input is expected of the form osgiBundleName/qualifiedJavaClassName
// So split the two pieces apart to help validate them.
func classNameUserInputToTestClassLocation(classNameUserInput string) (*TestLocation, error) {

	var (
		err               error         = nil
		testClassToLaunch *TestLocation = nil
	)

	parts := strings.Split(classNameUserInput, "/")
	if len(parts) < 2 {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_CLASS_INPUT_NO_SLASH, classNameUserInput)
	} else if len(parts) > 2 {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_CLASS_TOO_MANY_SLASHES, classNameUserInput)
	} else if strings.HasSuffix(parts[1], ".class") {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_CLASS_SUFFIX_FOUND, classNameUserInput)
	} else {
		osgiBundleName := parts[0]
		qualifiedJavaClassName := parts[1]

		testClassToLaunch = &TestLocation{
			OSGiBundleName:         osgiBundleName,
			QualifiedJavaClassName: qualifiedJavaClassName,
		}
	}

	return testClassToLaunch, err
}
