/*
 * Copyright contributors to the Galasa project
 */
package launcher

import (
	"bytes"
	"embed"
	"errors"
	"log"
	"os/exec"
	"strings"

	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

type TestLocation struct {
	OSGiBundleName         string
	QualifiedJavaClassName string
}

type TestResultsSummary struct {
	MethodPasses int
	MethodFails  int
}

type JvmLauncher struct {
	// The fully-qualified path to JAVA_HOME where we can find the bin/java command.
	javaHome string

	// The parameters from the command-line.
	cmdParams RunsSubmitLocalCmdParameters

	env                utils.Environment
	fileSystem         utils.FileSystem
	embeddedFileSystem embed.FS

	// The collection of tests which are running, or have completed.
	localTests []*LocalTest
}

type RunsSubmitLocalCmdParameters struct {
	Obrs                []string
	RemoteMaven         string
	TargetGalasaVersion string
}

type LocalTest struct {
	jvmProcess *exec.Cmd
	stdout     *JVMOutputProcessor
	stderr     *bytes.Buffer

	reportingChannel chan string

	// What runId is this test using ?
	// We don't initially know it. This info is extracted from the JVM trace.
	runId string

	// Where is the RAS folder storing results for this test ?
	// We don't initially know it. This info is extracted from the JVM trace.
	rasFolderPath string

	testRun *galasaapi.TestRun
}

// -----------------------------------------------------------------------------
// Constructors
// -----------------------------------------------------------------------------
func NewJVMLauncher(
	env utils.Environment,
	fileSystem utils.FileSystem,
	embeddedFileSystem embed.FS,
	runsSubmitLocalCmdParams RunsSubmitLocalCmdParameters,
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

		// Make sure the home folder has the boot jar unpacked and ready to invoke.
		err = utils.InitialiseGalasaHomeFolder(launcher.fileSystem, launcher.embeddedFileSystem)
	}

	return launcher, err
}

//-----------------------------------------------------------------------------
// Implementation of the Launcher interface
//-----------------------------------------------------------------------------

// SubmitTestRuns launch the test runs
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
				cmd, args, err := getCommandSyntax(launcher.fileSystem, launcher.javaHome, obrs, *testClassToLaunch, launcher.cmdParams.RemoteMaven, launcher.cmdParams.TargetGalasaVersion)
				if err == nil {
					log.Printf("Launching command '%s' '%v'\n", cmd, args)
					localTest := NewLocalTest()
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

// User input is expected of the form osgiBundleName/qualifiedJavaClassName
// So split the two pieces apart to help validate them.
func classNameUserInputToTestClassLocation(classNameUserInput string) (*TestLocation, error) {

	var err error = nil
	var testClassToLaunch *TestLocation = nil

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

func (launcher *JvmLauncher) GetRunsByGroup(groupName string) (*galasaapi.TestRuns, error) {
	log.Printf("JvmLauncher: GetRunsByGroup(groupName=%s) entered. ", groupName)

	var isAllComplete = true
	var testRuns = galasaapi.TestRuns{
		Complete: &isAllComplete,
		Runs:     []galasaapi.TestRun{},
	}

	for _, localTest := range launcher.localTests {

		// Update the test status by reading the json file if we can.
		localTest.updateTestStatusFromRasFile(launcher.fileSystem)

		testName := localTest.testRun.GetName()
		if localTest.isCompleted() {
			log.Printf("GetRunsByGroup: localTest %s is complete.\n", testName)
		} else {
			log.Printf("GetRunsByGroup: localTest %s is not yet complete.\n", testName)
			isAllComplete = false
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
// A structure which tells us all we know about a JVM process we launched.

func NewLocalTest() *LocalTest {

	localTest := new(LocalTest)

	localTest.jvmProcess = nil
	localTest.stdout = NewJVMOutputProcessor()
	localTest.stderr = bytes.NewBuffer([]byte{})
	localTest.runId = ""
	localTest.testRun = nil

	localTest.reportingChannel = make(chan string, 100)

	return localTest
}

// Launch a test within a JVM.
// Hang around waiting for the JVM to trace the runID and ras location.
func (localTest *LocalTest) launch(cmd string, args []string) error {
	localTest.jvmProcess = exec.Command(cmd, args...)
	localTest.jvmProcess.Stdout = localTest.stdout
	localTest.jvmProcess.Stderr = localTest.stderr

	err := localTest.jvmProcess.Start()
	if err != nil {
		log.Printf("Failed to start the JVM. %s\n", err.Error())
		log.Printf("Failing command is %s %v\n", cmd, args)
	} else {
		localTest.runId, err = waitForRunIdAllocation(localTest.stdout)
		if err == nil {

			localTest.rasFolderPath, err = waitForRasFolderPath(localTest.stdout)
			if err == nil {
				log.Printf("JVM test started OK. Spawning a go routine to wait for it to complete.\n")
				go localTest.waitForCompletion()
			}
		}
	}
	return err
}

// Block this thread until we can gather where the RAS folder is for this test.
// It is resolved within the JVM, and traced, where we pick it up from.
func waitForRasFolderPath(outputProcessor *JVMOutputProcessor) (string, error) {
	var err error = nil

	// BLOCK THREAD !
	// Wait for the runId to be detected in the JVM output.
	<-outputProcessor.publishResultChannel

	rasFolderPath := outputProcessor.detectedRasFolderPath

	if rasFolderPath == "" {
		// TODO: Better error message please.
		err = errors.New("rasFolderPath could not be detected")
	}

	return rasFolderPath, err
}

// Block this thread until we can gather what the RunId for this test is
// It is allocated within the JVM, and traced, where we pick it up from.
func waitForRunIdAllocation(outputProcessor *JVMOutputProcessor) (string, error) {
	var err error = nil

	// BLOCK THREAD !
	// Wait for the runId to be detected in the JVM output.
	<-outputProcessor.publishResultChannel

	runId := outputProcessor.detectedRunId

	if runId == "" {
		// TODO: Better error message please.
		err = errors.New("runid could not be detected")
	}

	return runId, err
}

// This method is called by the launching thread as a go routine.
// The go routine waits for the JVM to complete, then emits
// a 'DONE' message which can be recieved by the monitoring thread.
// This call always blocks waiting for the launched JVM to complete and exit.
func (localTest *LocalTest) waitForCompletion() error {

	log.Printf("waiting for the JVM to complete within a go routine.\n")

	err := localTest.jvmProcess.Wait()
	if err != nil {
		log.Printf("Failed to wait for the JVM test to complete. %s\n", err.Error())
	} else {
		log.Printf("JVM has completed. Detected by waiting go routine.\n")
	}

	// Tell any polling thread that the JVM is complete now.
	localTest.testRun.SetStatus("finished")
	localTest.reportingChannel <- "DONE"
	close(localTest.reportingChannel)

	return err
}

// If we can find it, read the status report for the test from the
// ras folder.
func (localTest *LocalTest) updateTestStatusFromRasFile(fileSystem utils.FileSystem) error {

	var err error = nil

	if localTest.runId == "" || localTest.rasFolderPath == "" {
		log.Printf("Don't have enough information to find the structure.json in the RAS folder.\n")
	} else {

		jsonFilePath := localTest.rasFolderPath + "/" + localTest.runId + "/structure.json"
		log.Printf("Reading latest test status from %s\n", jsonFilePath)

		var testRun *galasaapi.TestRun
		testRun, err = readTestRunFromJsonFile(fileSystem, jsonFilePath)

		if err == nil {
			localTest.testRun = testRun
		}
	}
	return err
}

// This method is called by a thread monitoring the state of the JVM.
// It can receive messages from the JVM launcher go routine.
// This call never blocks waiting for anything.
func (localTest *LocalTest) isCompleted() bool {

	log.Printf("Checking to see if local test is completed...")
	isComplete := false
	select {
	case msg := <-localTest.reportingChannel:
		log.Printf("Message received from JVM launch thread: %s\n", msg)
		if msg == "DONE" || msg == "" {
			isComplete = true
		}
	default:
		log.Printf("No message received from JVM launch thread. Would block. JVM is not finished.")
		isComplete = false
	}
	return isComplete
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
