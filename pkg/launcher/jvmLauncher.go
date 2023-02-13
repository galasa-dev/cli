/*
 * Copyright contributors to the Galasa project
 */
package launcher

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
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
}

type launchedJvm struct {
	jvmProcess *exec.Cmd
}

type RunsSubmitLocalCmdParameters struct {
	Obrs                []string
	RemoteMaven         string
	TargetGalasaVersion string
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

	utils.InitialiseGalasaHomeFolder(launcher.fileSystem, launcher.embeddedFileSystem)

	var err error = nil

	// We expect a parameter to be of the form:
	// mvn:dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOT/obr
	var obrs []utils.MavenCoordinates
	for _, obr := range launcher.cmdParams.Obrs {
		parts := strings.Split(obr, "/")
		if len(parts) <= 0 {
			err = errors.New("badly formed OBR parameter. " + obr)
		} else if len(parts) > 4 {
			err = errors.New("badly formed OBR parameter. " + obr)
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

	for _, classNameUserInput := range classNames {

		var testClassToLaunch *TestLocation
		testClassToLaunch, err = classNameUserInputToTestClassLocation(classNameUserInput)

		if err == nil {
			err = executeTestInJVM(launcher.fileSystem, launcher.javaHome,
				obrs, testClassToLaunch, launcher.cmdParams.RemoteMaven, launcher.cmdParams.TargetGalasaVersion)
		}

		if err != nil {
			break
		}
	}

	return nil, err
}

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

	var isComplete = true
	var testRuns = galasaapi.TestRuns{
		Complete: &isComplete,
		Runs:     []galasaapi.TestRun{},
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

//-----------------------------------------------------------------------------
// Local functions
//-----------------------------------------------------------------------------

func executeTestInJVM(
	fileSystem utils.FileSystem,
	javaHome string,
	testObrs []utils.MavenCoordinates,
	testLocation *TestLocation,
	remoteMaven string,
	targetGalasaVersion string,
) error {
	cmd, args, err := getCommandSyntax(fileSystem, javaHome, testObrs, *testLocation, remoteMaven, targetGalasaVersion)
	if err == nil {
		log.Printf("Launching command '%s' '%v'\n", cmd, args)
		jvmProcess := exec.Command(cmd, args...)

		outStream := bytes.NewBuffer([]byte{})
		jvmProcess.Stdout = outStream

		errStream := bytes.NewBuffer([]byte{})
		jvmProcess.Stderr = errStream

		err = jvmProcess.Start()
		if err != nil {
			log.Printf("Failed to start the JVM. %s\n", err.Error())
			log.Printf("Failing command is %s %v\n", cmd, args)
		} else {

			// Wait for it to complete.
			err = jvmProcess.Wait()
			if err != nil {
				log.Printf("Failed to wait for the JVM. %s\n", err.Error())
				log.Printf("Failing command is %s %v\n", cmd, args)
			} else {

				log.Printf("jvm standard output is %s\n", outStream.String())
				log.Printf("jvm error output is %s\n", errStream.String())

				results := summariseJVMOutput(outStream.Bytes())
				log.Printf("Results: Method passes: %d, fails: %d", results.MethodPasses, results.MethodFails)
			}
		}
	}
	return err
}

func summariseJVMOutput(jvmOutput []byte) TestResultsSummary {
	jvmOutputStr := string(jvmOutput[:])

	results := TestResultsSummary{MethodPasses: 0, MethodFails: 0}

	results.MethodPasses = strings.Count(jvmOutputStr, "*** Passed - Test method ")
	results.MethodFails = strings.Count(jvmOutputStr, "*** Failed - Test method ")

	return results
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

// executeSubmitLocal attempts to launch a JVM locally in which the Galasa framework is booted
// which will run the test(s) we wish to use. The tests execute locally, but may use values
// from the Galasa ecosystem, if the bootstrap points to it.
//
// JAVA_HOME must be set.
func executeSubmitLocal(
	fileSystem utils.FileSystem,
	embeddedFileSystem embed.FS,
	javaHome string) error {

	var err error = nil

	err = utils.ValidateJavaHome(fileSystem, javaHome)

	if err == nil {
		err = utils.InitialiseGalasaHomeFolder(fileSystem, embeddedFileSystem)
	}

	if err == nil {
		err = executeTestsInJvm(fileSystem, javaHome)
	}

	return err
}

func executeTestsInJvm(fileSystem utils.FileSystem, javaHome string) error {
	var err error = nil

	// The Output method runs the command, waits for it to finish and collects its standard output. If there were no errors, dateOut will hold bytes with the date info.

	// dateOut, err := dateCmd.Output()
	// if err != nil {
	//     panic(err)
	// }
	// fmt.Println("> date")
	// fmt.Println(string(dateOut))
	// Output and other methods of Command will return *exec.Error if there was a problem executing the command (e.g. wrong path), and *exec.ExitError if the command ran but exited with a non-zero return code.

	_, err = exec.Command("date", "-x").Output()
	if err != nil {
		switch e := err.(type) {
		case *exec.Error:
			fmt.Println("failed executing:", err)
		case *exec.ExitError:
			fmt.Println("command exit rc =", e.ExitCode())
		default:
			panic(err)
		}
	}

	return err
}
