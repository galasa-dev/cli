/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"testing"

	// "github.com/stretchr/testify/mock"
	// "github.com/galasa.dev/cli/pkg/galasaapi"

	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCanWriteAndReadBackThrottleFile(t *testing.T) {

	mockFileSystem := utils.NewMockFileSystem()
	err := writeThrottleFile(mockFileSystem, "throttle", 101)
	if err != nil {
		assert.Fail(t, "Should not have failed to write a throttle file. "+err.Error())
	}

	isThrottleFileExists, err := mockFileSystem.Exists("throttle")
	if err != nil {
		assert.Fail(t, "Should not have failed to check for the existence of a throttle file. "+err.Error())
	}

	assert.True(t, isThrottleFileExists, "throttle file does not exist!")

	var readBackThrottle int
	readBackThrottle, err = readThrottleFile(mockFileSystem, "throttle")
	if err != nil {
		assert.Fail(t, "Should not have failed to read from a throttle file. "+err.Error())
	}
	assert.Equal(t, readBackThrottle, 101, "read back the wrong throttle value")
}

func TestReadBackThrottleFileFailsIfNoThrottleFileThere(t *testing.T) {

	var err error
	mockFileSystem := utils.NewMockFileSystem()

	_, err = readThrottleFile(mockFileSystem, "throttle")
	if err == nil {
		assert.Fail(t, "Should have failed to read from a throttle file. "+err.Error())
	}
	assert.Contains(t, err.Error(), "GAL1048", "Error returned should contain GAL1048 error indicating read throttle file failed."+err.Error())
}

func TestReadBackThrottleFileFailsIfFileContainsInvalidInt(t *testing.T) {

	var err error
	mockFileSystem := utils.NewMockFileSystem()

	mockFileSystem.WriteTextFile("throttle", "abc")

	_, err = readThrottleFile(mockFileSystem, "throttle")
	if err == nil {
		assert.Fail(t, "Should have failed to read from a throttle file. "+err.Error())
	}
	assert.Contains(t, err.Error(), "GAL1049E", "Error returned should contain GAL1049E error indicating read invalid throttle file content."+err.Error())
}

func TestUpdateThrottleFromFileIfDifferentChangesValueWhenDifferent(t *testing.T) {

	mockFileSystem := utils.NewMockFileSystem()

	mockFileSystem.WriteTextFile("throttle", "10")
	newValue, isLost := updateThrottleFromFileIfDifferent(mockFileSystem, "throttle", 20, false)

	assert.Equal(t, 10, newValue)
	assert.False(t, isLost)
}

func TestUpdateThrottleFromFileIfDifferentDoesntChangeIfFileMissing(t *testing.T) {

	mockFileSystem := utils.NewMockFileSystem()

	// mockFileSystem.WriteTextFile("throttle", "10") - file is missing now.
	newValue, isLost := updateThrottleFromFileIfDifferent(mockFileSystem, "throttle", 20, false)

	assert.Equal(t, 20, newValue)
	assert.True(t, isLost)
}

/*
func TestCanSubmitSmallPortfolio(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()

	portfolioText := ` Some portfolio text. tbd.
	`

	mockFileSystem.WriteTextFile("small.portfolio", portfolioText)

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
		isLocal:					   false,
	}

	mockTimeService := utils.NewMockTimeServiceAsMock()

	apiClient := &galasaapi.APIClient{}

	// RunsAPIApi.On("GetRunsGroup",nil, groupName).Return(nil)

	// When ...
	err := executeSubmitRemote(
		mockFileSystem,
		params,
		apiClient,
		mockTimeService)

	// Then...
	if err != nil {
		assert.Fail(t, "Not expecting error "+err.Error())
	}
	assert.Fail(t, "Exiting.")
}
*/
