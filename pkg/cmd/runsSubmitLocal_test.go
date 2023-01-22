/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func newLocalRunsSubmitCmdParameters() *RunsSubmitCmdParameters {
	params := RunsSubmitCmdParameters{
		pollIntervalSeconds:           1,
		noExitCodeOnTestFailures:      true,
		progressReportIntervalMinutes: 1,
		throttle:                      1,
		trace:                         false,
		reportYamlFilename:            "a.yaml",
		reportJsonFilename:            "a.json",
		reportJunitFilename:           "a.junit.xml",
		groupName:                     "babe",
		portfolioFileName:             "small.portfolio",
		isLocal:                       true,
	}
	return &params
}

func TestLocalRunFailsIfJavaHomeNotSet(t *testing.T) {

	mockFileSystem := utils.NewMockFileSystem()
	params := *newLocalRunsSubmitCmdParameters()
	javaHome := "" // It's not set is the same as being blank.

	err := executeSubmitLocal(mockFileSystem, params, javaHome)

	if err == nil {
		assert.Fail(t, "Failed to detect that JAVA_HOME has not been set.")
	}
}

func TestLocalRunFailsIfJavaHomeFolderHasNoBinSubFolder(t *testing.T) {

	mockFileSystem := utils.NewMockFileSystem()
	params := *newLocalRunsSubmitCmdParameters()
	javaHome := "myJavaHome"

	err := executeSubmitLocal(mockFileSystem, params, javaHome)

	if err == nil {
		assert.Fail(t, "Failed to detect that JAVA_HOME has no sub-folder called 'bin'")
	}
	assert.Contains(t, err.Error(), "GAL1052E", "Returned error, but it was the wrong one !")
}

func TestLocalRunFailsIfJavaHomeFolderHasNoBinSubFolderWithTrailingSlashShouldHaveBeenRemoved(t *testing.T) {

	mockFileSystem := utils.NewMockFileSystem()
	params := *newLocalRunsSubmitCmdParameters()
	javaHome := "myJavaHome" + string(os.PathSeparator) // It's got a trailing slash which needs stripping.

	err := executeSubmitLocal(mockFileSystem, params, javaHome)

	if err == nil {
		assert.Fail(t, "Failed to detect that JAVA_HOME has no sub-folder called 'bin'")
	}
	assert.Contains(t, err.Error(), "GAL1052E", "Returned error, but it was the wrong one !")
	assert.NotContains(t, err.Error(), "myJavaHome//bin", "Returned error, message was wrong. Slash should have been removed.")
}

func TestLocalRunFailsIfJavaHomeFolderCheckFails(t *testing.T) {

	mockFileSystem := utils.NewOverridableMockFileSystem()
	mockFileSystem.VirtualFunction_DirExists = func(path string) (bool, error) {
		return false, errors.New("Simulated DirExists failure")
	}
	params := *newLocalRunsSubmitCmdParameters()
	javaHome := "myJavaHome"

	err := executeSubmitLocal(mockFileSystem, params, javaHome)

	assert.NotNil(t, err, "Failed to detect that when checking JAVA_HOME for a sub-folder called 'bin' a failure happened.")
	assert.Contains(t, err.Error(), "GAL1051E", "Returned error, but it was the wrong one !")
}

func TestLocalRunJVMJavaProgramMissingFails(t *testing.T) {

	mockFileSystem := utils.NewOverridableMockFileSystem()
	javaHome := "myJavaHome"
	binFolder := javaHome + FILE_SYSTEM_PATH_SEPARATOR + "bin"
	mockFileSystem.MkdirAll(binFolder)
	// We don't create the java command wihthin $JAVA_HOME/bin to see if it gets detected prior to any attempt to launch a JVM.
	params := *newLocalRunsSubmitCmdParameters()

	err := executeSubmitLocal(mockFileSystem, params, javaHome)

	assert.NotNil(t, err, "Failed to detect that JAVA_HOMEA/bin/java is missing.")
	assert.Contains(t, err.Error(), "GAL1054E", "Returned error, but it was the wrong one !")
}

func TestLocalRunFailsIfJavaProgramPresenceCheckFails(t *testing.T) {

	mockFileSystem := utils.NewOverridableMockFileSystem()

	// Make sure JAVA_HOME/bin exists.
	javaHome := "myJavaHome"
	binFolder := javaHome + FILE_SYSTEM_PATH_SEPARATOR + "bin"
	mockFileSystem.MkdirAll(binFolder)

	mockFileSystem.VirtualFunction_Exists = func(path string) (bool, error) {
		return false, errors.New("Simulated Exists failure")
	}
	params := *newLocalRunsSubmitCmdParameters()

	err := executeSubmitLocal(mockFileSystem, params, javaHome)

	assert.NotNil(t, err, "Failed to detect that when checking for presence of JAVA_HOME/bin/java.")
	assert.Contains(t, err.Error(), "GAL1053E", "Returned error, but it was the wrong one !")
}
