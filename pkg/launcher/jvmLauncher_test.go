/*
 * Copyright contributors to the Galasa project
 */
package launcher

import (
	// 	"os"
	"testing"

	"github.com/galasa.dev/cli/pkg/embedded"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	// 	"errors"
	// 	"github.com/galasa.dev/cli/pkg/embedded"
	// "github.com/galasa.dev/cli/pkg/utils"
	// "github.com/stretchr/testify/assert"
)

func TestCanCreateAJVMLauncher(t *testing.T) {

	env := utils.NewMockEnv()
	env.EnvVars["JAVA_HOME"] = "/java"

	fileSystem := utils.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fileSystem, "/java")

	jvmLaunchParams := getBasicJvmLaunchParams()

	launcher, err := NewJVMLauncher(env, fileSystem, embedded.GetEmbeddedFileSystem(), jvmLaunchParams)
	if err != nil {
		assert.Fail(t, "Constructor should not have failed but it did. error:%s", err.Error())
	}
	assert.NotNil(t, launcher, "Launcher reference was nil, shouldn't have been.")
}

func getBasicJvmLaunchParams() RunsSubmitLocalCmdParameters {
	return RunsSubmitLocalCmdParameters{
		Obrs:                nil,
		RemoteMaven:         "",
		TargetGalasaVersion: "",
	}
}

func TestCantCreateAJVMLauncherIfJVMHomeNotSet(t *testing.T) {

	env := utils.NewMockEnv()
	// env.EnvVars["JAVA_HOME"] = "/java"

	fileSystem := utils.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fileSystem, "/java")

	jvmLaunchParams := getBasicJvmLaunchParams()

	launcher, err := NewJVMLauncher(env, fileSystem, embedded.GetEmbeddedFileSystem(), jvmLaunchParams)
	if err == nil {
		assert.Fail(t, "Constructor should have failed but it did not.")
	}
	assert.Nil(t, launcher, "Launcher reference was not nil")
	assert.Contains(t, err.Error(), "GAL1050E")
}

// func TestJvmGetsLaunchedWithCorrectSyntax(t *testing.T) {

// 	// Given...
// 	mockFileSystem := utils.NewOSFileSystem()

// 	var testObrs []MavenCoordinates = []MavenCoordinates{
// 		{
// 			GroupId:    "dev.galasa.example.banking",
// 			ArtifactId: "dev.galasa.example.banking.obr",
// 			Version:    "0.0.1-SNAPSHOT",
// 		},
// 	}

// 	var testLocation TestLocation = TestLocation{
// 		OSGiBundleName: "dev.galasa.example.banking.payee",
// 		Class: JavaClassDef{
// 			PackageName: "dev.galasa.example.banking.payee",
// 			ClassName:   "TestPayee",
// 		},
// 	}

// 	javaHome := os.Getenv("JAVA_HOME")

// 	// remoteMaven := "https://repo.maven.apache.org/maven2"
// 	remoteMaven := "https://development.galasa.dev/main/maven-repo/obr/"

// 	// When
// 	err := executeTestInJVM(mockFileSystem, javaHome, testObrs, testLocation, remoteMaven)

// 	// Then
// 	if err != nil {
// 		assert.Fail(t, "Expecting no errors but there was one."+err.Error())
// 	}
// }

// func newLocalRunsSubmitCmdParameters() *utils.RunsSubmitCmdParameters {
// 	params := utils.RunsSubmitCmdParameters{
// 		PollIntervalSeconds:           1,
// 		NoExitCodeOnTestFailures:      true,
// 		ProgressReportIntervalMinutes: 1,
// 		Throttle:                      1,
// 		Trace:                         false,
// 		ReportYamlFilename:            "a.yaml",
// 		ReportJsonFilename:            "a.json",
// 		ReportJunitFilename:           "a.junit.xml",
// 		GroupName:                     "babe",
// 		PortfolioFileName:             "small.portfolio",
// 		IsLocal:                       true,
// 	}
// 	return &params
// }

// func TestLocalRunFailsIfJavaHomeNotSet(t *testing.T) {

// 	mockFileSystem := utils.NewMockFileSystem()
// 	params := *newLocalRunsSubmitCmdParameters()
// 	javaHome := "" // It's not set is the same as being blank.
// 	embeddedFileSystem := embedded.GetEmbeddedFileSystem()

// 	err := executeSubmitLocal(mockFileSystem, embeddedFileSystem, params, javaHome)

// 	if err == nil {
// 		assert.Fail(t, "Failed to detect that JAVA_HOME has not been set.")
// 	}
// }

// func TestLocalRunFailsIfJavaHomeFolderHasNoBinSubFolder(t *testing.T) {

// 	mockFileSystem := utils.NewMockFileSystem()
// 	params := *newLocalRunsSubmitCmdParameters()
// 	javaHome := "myJavaHome"
// 	embeddedFileSystem := embedded.GetEmbeddedFileSystem()

// 	err := executeSubmitLocal(mockFileSystem, embeddedFileSystem, params, javaHome)

// 	if err == nil {
// 		assert.Fail(t, "Failed to detect that JAVA_HOME has no sub-folder called 'bin'")
// 	}
// 	assert.Contains(t, err.Error(), "GAL1052E", "Returned error, but it was the wrong one !")
// }

// func TestLocalRunFailsIfJavaHomeFolderHasNoBinSubFolderWithTrailingSlashShouldHaveBeenRemoved(t *testing.T) {

// 	mockFileSystem := utils.NewMockFileSystem()
// 	params := *newLocalRunsSubmitCmdParameters()
// 	javaHome := "myJavaHome" + string(os.PathSeparator) // It's got a trailing slash which needs stripping.
// 	embeddedFileSystem := embedded.GetEmbeddedFileSystem()

// 	err := executeSubmitLocal(mockFileSystem, embeddedFileSystem, params, javaHome)

// 	if err == nil {
// 		assert.Fail(t, "Failed to detect that JAVA_HOME has no sub-folder called 'bin'")
// 	}
// 	assert.Contains(t, err.Error(), "GAL1052E", "Returned error, but it was the wrong one !")
// 	assert.NotContains(t, err.Error(), "myJavaHome//bin", "Returned error, message was wrong. Slash should have been removed.")
// }

// func TestLocalRunFailsIfJavaHomeFolderCheckFails(t *testing.T) {

// 	mockFileSystem := utils.NewOverridableMockFileSystem()
// 	mockFileSystem.VirtualFunction_DirExists = func(path string) (bool, error) {
// 		return false, errors.New("Simulated DirExists failure")
// 	}
// 	params := *newLocalRunsSubmitCmdParameters()
// 	javaHome := "myJavaHome"
// 	embeddedFileSystem := embedded.GetEmbeddedFileSystem()

// 	err := executeSubmitLocal(mockFileSystem, embeddedFileSystem, params, javaHome)

// 	assert.NotNil(t, err, "Failed to detect that when checking JAVA_HOME for a sub-folder called 'bin' a failure happened.")
// 	assert.Contains(t, err.Error(), "GAL1051E", "Returned error, but it was the wrong one !")
// }

// func TestLocalRunJVMJavaProgramMissingFails(t *testing.T) {

// 	mockFileSystem := utils.NewOverridableMockFileSystem()
// 	javaHome := "myJavaHome"
// 	binFolder := javaHome + utils.FILE_SYSTEM_PATH_SEPARATOR + "bin"
// 	mockFileSystem.MkdirAll(binFolder)
// 	// We don't create the java command wihthin $JAVA_HOME/bin to see if it gets detected prior to any attempt to launch a JVM.
// 	params := *newLocalRunsSubmitCmdParameters()
// 	embeddedFileSystem := embedded.GetEmbeddedFileSystem()

// 	err := executeSubmitLocal(mockFileSystem, embeddedFileSystem, params, javaHome)

// 	assert.NotNil(t, err, "Failed to detect that JAVA_HOMEA/bin/java is missing.")
// 	assert.Contains(t, err.Error(), "GAL1054E", "Returned error, but it was the wrong one !")
// }

// func TestLocalRunFailsIfJavaProgramPresenceCheckFails(t *testing.T) {

// 	mockFileSystem := utils.NewOverridableMockFileSystem()

// 	// Make sure JAVA_HOME/bin exists.
// 	javaHome := "myJavaHome"
// 	binFolder := javaHome + utils.FILE_SYSTEM_PATH_SEPARATOR + "bin"
// 	mockFileSystem.MkdirAll(binFolder)
// 	embeddedFileSystem := embedded.GetEmbeddedFileSystem()

// 	mockFileSystem.VirtualFunction_Exists = func(path string) (bool, error) {
// 		return false, errors.New("Simulated Exists failure")
// 	}
// 	params := *newLocalRunsSubmitCmdParameters()

// 	err := executeSubmitLocal(mockFileSystem, embeddedFileSystem, params, javaHome)

// 	assert.NotNil(t, err, "Failed to detect that when checking for presence of JAVA_HOME/bin/java.")
// 	assert.Contains(t, err.Error(), "GAL1053E", "Returned error, but it was the wrong one !")
// }

// func TestLocalRunGoldenPath(t *testing.T) {

// 	mockFileSystem := utils.NewOverridableMockFileSystem()

// 	// Make sure JAVA_HOME/bin exists.
// 	javaHome := "myJavaHome"
// 	binFolder := javaHome + utils.FILE_SYSTEM_PATH_SEPARATOR + "bin"
// 	mockFileSystem.MkdirAll(binFolder)
// 	// Create a mock java program so it gets checked.
// 	javaProgramPath := binFolder + utils.FILE_SYSTEM_PATH_SEPARATOR + "java"
// 	mockFileSystem.VirtualFunction_WriteTextFile(javaProgramPath, "dummy in memory file")

// 	embeddedFileSystem := embedded.GetEmbeddedFileSystem()

// 	params := *newLocalRunsSubmitCmdParameters()

// 	err := executeSubmitLocal(mockFileSystem, embeddedFileSystem, params, javaHome)

// 	if err != nil {
// 		assert.Fail(t, "Failed to detect that when checking for presence of JAVA_HOME/bin/java. "+err.Error())
// 	}

// 	// Sanity check that it creates the home folder if that doesn't already exist.
// 	home, _ := mockFileSystem.GetUserHomeDir()
// 	isExists, _ := mockFileSystem.Exists(home + utils.FILE_SYSTEM_PATH_SEPARATOR +
// 		".galasa" + utils.FILE_SYSTEM_PATH_SEPARATOR + "cps.properties")
// 	assert.True(t, isExists, "cps.properties was not created.")
// }
