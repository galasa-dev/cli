/*
 * Copyright contributors to the Galasa project
 */
package launcher

import (
	"embed"
	"errors"
	"log"
	"strings"

	"github.com/galasa.dev/cli/pkg/galasaapi"
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
	fileSystem utils.FileSystem

	// A file system so we can get at embedded content if required.
	// (Like so we can unpack the boot.jar)
	embeddedFileSystem embed.FS

	// The collection of tests which are running, or have completed.
	localTests []*LocalTest

	// This timer service can be interrupted when we don't want it to sleep.
	timeService utils.TimeService

	// A service which can create OS processes.
	processFactory ProcessFactory
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
}

// -----------------------------------------------------------------------------
// Constructors
// -----------------------------------------------------------------------------

// NewJVMLauncher creates a JVM launcher. Primes it with references to services
// which can be used to launch JVM servers.
func NewJVMLauncher(
	env utils.Environment,
	fileSystem utils.FileSystem,
	embeddedFileSystem embed.FS,
	runsSubmitLocalCmdParams RunsSubmitLocalCmdParameters,
	timeService utils.TimeService,
	processFactory ProcessFactory,
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

		launcher.timeService = timeService

		// Make sure the home folder has the boot jar unpacked and ready to invoke.
		err = utils.InitialiseGalasaHomeFolder(launcher.fileSystem, launcher.embeddedFileSystem)
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

	obrs, err := validateObrs(launcher.cmdParams.Obrs)

	testRuns := new(galasaapi.TestRuns)

	isComplete := false
	testRuns.Complete = &isComplete
	testRuns.Runs = make([]galasaapi.TestRun, 0)

	if err == nil {
		for _, classNameUserInput := range classNames {

			var testClassToLaunch *TestLocation
			testClassToLaunch, err = classNameUserInputToTestClassLocation(classNameUserInput)

			if err == nil {
				var (
					cmd  string
					args []string
				)
				cmd, args, err = getCommandSyntax(
					launcher.fileSystem, launcher.javaHome, obrs,
					*testClassToLaunch, launcher.cmdParams.RemoteMaven,
					launcher.cmdParams.TargetGalasaVersion)
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

	return testRuns, err
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

			// Update the test status by reading the json file if we can.
			localTest.updateTestStatusFromRasFile()

			if localTest.isCompleted() {
				log.Printf("GetRunsByGroup: localTest %s is not yet complete.\n", testName)
			} else {
				log.Printf("GetRunsByGroup: localTest read status and it is finished.")
				isAllComplete = false
			}
		}
		testRuns.Runs = append(testRuns.Runs, *localTest.testRun)
	}

	log.Printf("JvmLauncher: GetRunsByGroup(groupName=%s) exiting. ", groupName)
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
		if len(parts) <= 0 {
			err = errors.New("badly formed OBR parameter. Expected it to be of the form mvn:<GROUP_ID>/<ARTIFACT_ID>/<VERSION>/obr " + obr)
		} else if len(parts) > 4 {
			err = errors.New("badly formed OBR parameter. Expected it to be of the form mvn:<GROUP_ID>/<ARTIFACT_ID>/<VERSION>/obr " + obr)
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
//	java -jar /Users/mcobbett/builds/galasa/code/external/galasa-dev/cli/pkg/embedded/templates/galasahome/lib/galasa-boot-0.24.0.jar \
//	    --localmaven file:/Users/mcobbett/.m2/repository/ \
//	    --remotemaven https://development.galasa.dev/main/maven-repo/obr/ \
//	    --bootstrap file:/Users/mcobbett/.galasa/bootstrap.properties \
//	    --overrides file:/Users/mcobbett/.galasa/overrides.properties \
//	    --obr mvn:dev.galasa/dev.galasa.uber.obr/0.25.0/obr \
//	    --obr mvn:dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOT/obr \
//	    --test dev.galasa.example.banking.payee/dev.galasa.example.banking.payee.TestPayee
func getCommandSyntax(
	fileSystem utils.FileSystem,
	javaHome string,
	testObrs []utils.MavenCoordinates,
	testLocation TestLocation,
	remoteMaven string,
	galasaVersionToRun string,
) (string, []string, error) {

	var cmd string = ""
	var args []string = make([]string, 0)

	bootJarPath, err := utils.GetGalasaBootJarPath(fileSystem)
	if err == nil {

		cmd = javaHome + utils.FILE_SYSTEM_PATH_SEPARATOR + "bin" +
			utils.FILE_SYSTEM_PATH_SEPARATOR + "java"

		args = append(args, "-jar")
		args = append(args, bootJarPath)

		args = append(args, "-Dfile.encoding=UTF-8")

		var userHome string
		userHome, err = fileSystem.GetUserHomeDir()

		// --localmaven file:${M2_PATH}/repository/
		args = append(args, "--localmaven")
		localMavenPath := "file:" + userHome + utils.FILE_SYSTEM_PATH_SEPARATOR +
			".m2" + utils.FILE_SYSTEM_PATH_SEPARATOR + "repository"
		args = append(args, localMavenPath)

		// --remotemaven $REMOTE_MAVEN
		args = append(args, "--remotemaven")
		args = append(args, remoteMaven)

		// --bootstrap file:${HOME}/.galasa/bootstrap.properties
		args = append(args, "--bootstrap")
		bootstrapPath := "file:" + userHome + utils.FILE_SYSTEM_PATH_SEPARATOR +
			".galasa" + utils.FILE_SYSTEM_PATH_SEPARATOR + "bootstrap.properties"
		args = append(args, bootstrapPath)

		// --overrides file:${HOME}/.galasa/overrides.properties
		args = append(args, "--overrides")
		overridesPath := "file:" + userHome + utils.FILE_SYSTEM_PATH_SEPARATOR +
			".galasa" + utils.FILE_SYSTEM_PATH_SEPARATOR + "overrides.properties"
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

	}

	return cmd, args, err
}

// User input is expected of the form osgiBundleName/qualifiedJavaClassName
// So split the two pieces apart to help validate them.
func classNameUserInputToTestClassLocation(classNameUserInput string) (*TestLocation, error) {

	var (
		err               error         = nil
		testClassToLaunch *TestLocation = nil
	)

	parts := strings.Split(classNameUserInput, "/")
	if len(parts) <= 0 {
		err = errors.New("error! - Bad class format. Should be osgiBundleName/qualifiedJavaClassName - slash is missing. not handled yet")
	} else if len(parts) > 2 {
		err = errors.New("error - too many segments. Not handled yet")
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