/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/launcher"
	"github.com/galasa.dev/cli/pkg/props"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCanWriteAndReadBackThrottleFile(t *testing.T) {

	mockFileSystem := files.NewMockFileSystem()
	env := utils.NewMockEnv()
	mockLauncher := launcher.NewMockLauncher()
	mockTimeService := utils.NewMockTimeService()

	galasaHome, _ := utils.NewGalasaHome(mockFileSystem, env, "")

	console := utils.NewMockConsole()

	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)

	err := submitter.writeThrottleFile("throttle", 101)
	if err != nil {
		assert.Fail(t, "Should not have failed to write a throttle file. "+err.Error())
	}

	isThrottleFileExists, err := mockFileSystem.Exists("throttle")
	if err != nil {
		assert.Fail(t, "Should not have failed to check for the existence of a throttle file. "+err.Error())
	}

	assert.True(t, isThrottleFileExists, "throttle file does not exist!")

	var readBackThrottle int
	readBackThrottle, err = submitter.readThrottleFile("throttle")
	if err != nil {
		assert.Fail(t, "Should not have failed to read from a throttle file. "+err.Error())
	}
	assert.Equal(t, readBackThrottle, 101, "read back the wrong throttle value")
}

func TestReadBackThrottleFileFailsIfNoThrottleFileThere(t *testing.T) {

	var err error
	mockFileSystem := files.NewMockFileSystem()
	env := utils.NewMockEnv()
	mockLauncher := launcher.NewMockLauncher()
	mockTimeService := utils.NewMockTimeService()

	galasaHome, _ := utils.NewGalasaHome(mockFileSystem, env, "")

	console := utils.NewMockConsole()

	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)

	_, err = submitter.readThrottleFile("throttle")
	if err == nil {
		assert.Fail(t, "Should have failed to read from a throttle file. "+err.Error())
	}
	assert.Contains(t, err.Error(), "GAL1048", "Error returned should contain GAL1048 error indicating read throttle file failed."+err.Error())
}

func TestReadBackThrottleFileFailsIfFileContainsInvalidInt(t *testing.T) {

	var err error
	mockFileSystem := files.NewMockFileSystem()

	mockFileSystem.WriteTextFile("throttle", "abc")

	env := utils.NewMockEnv()
	mockLauncher := launcher.NewMockLauncher()
	mockTimeService := utils.NewMockTimeService()

	galasaHome, _ := utils.NewGalasaHome(mockFileSystem, env, "")
	console := utils.NewMockConsole()

	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)

	_, err = submitter.readThrottleFile("throttle")
	if err == nil {
		assert.Fail(t, "Should have failed to read from a throttle file. "+err.Error())
	}
	assert.Contains(t, err.Error(), "GAL1049E", "Error returned should contain GAL1049E error indicating read invalid throttle file content."+err.Error())
}

func TestUpdateThrottleFromFileIfDifferentChangesValueWhenDifferent(t *testing.T) {

	mockFileSystem := files.NewMockFileSystem()

	env := utils.NewMockEnv()
	mockLauncher := launcher.NewMockLauncher()
	mockTimeService := utils.NewMockTimeService()

	galasaHome, _ := utils.NewGalasaHome(mockFileSystem, env, "")
	console := utils.NewMockConsole()
	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)

	mockFileSystem.WriteTextFile("throttle", "10")
	newValue, isLost := submitter.updateThrottleFromFileIfDifferent("throttle", 20, false)

	assert.Equal(t, 10, newValue)
	assert.False(t, isLost)
}

func TestUpdateThrottleFromFileIfDifferentDoesntChangeIfFileMissing(t *testing.T) {

	mockFileSystem := files.NewMockFileSystem()

	env := utils.NewMockEnv()
	mockLauncher := launcher.NewMockLauncher()
	mockTimeService := utils.NewMockTimeService()

	galasaHome, _ := utils.NewGalasaHome(mockFileSystem, env, "")
	console := utils.NewMockConsole()
	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)

	// mockFileSystem.WriteTextFile("throttle", "10") - file is missing now.
	newValue, isLost := submitter.updateThrottleFromFileIfDifferent("throttle", 20, false)

	assert.Equal(t, 20, newValue)
	assert.True(t, isLost)
}

func TestOverridesReadFromOverridesFile(t *testing.T) {

	fileProps := make(map[string]interface{})
	fileProps["c"] = "d"

	mockFileSystem := files.NewMockFileSystem()
	props.WritePropertiesFile(mockFileSystem, "/tmp/temp.properties", fileProps)

	commandParameters := utils.RunsSubmitCmdParameters{
		Overrides:        []string{"a=b"},
		OverrideFilePath: "/tmp/temp.properties",
	}

	env := utils.NewMockEnv()
	mockLauncher := launcher.NewMockLauncher()
	mockTimeService := utils.NewMockTimeService()

	galasaHome, _ := utils.NewGalasaHome(mockFileSystem, env, "")
	console := utils.NewMockConsole()
	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)

	overrides, err := submitter.buildOverrideMap(commandParameters)

	assert.Nil(t, err)
	assert.NotNil(t, overrides)
	assert.Contains(t, overrides, "a", "command-line override wasn't used.")
	assert.Equal(t, overrides["a"], "b", "command-line override not passed correctly.")
	assert.Contains(t, overrides, "c", "file-based override wasn't used")
	assert.Equal(t, overrides["c"], "d", "file-based override value wasn't passed correctly.")
}

func TestOverridesFileSpecifiedButDoesNotExist(t *testing.T) {

	fileProps := make(map[string]interface{})
	fileProps["c"] = "d"

	mockFileSystem := files.NewMockFileSystem()
	props.WritePropertiesFile(mockFileSystem, "/tmp/temp.properties", fileProps)

	commandParameters := utils.RunsSubmitCmdParameters{
		Overrides:        []string{"a=b"},
		OverrideFilePath: "/tmp/temp.wrong.file.properties",
	}

	env := utils.NewMockEnv()
	mockLauncher := launcher.NewMockLauncher()
	mockTimeService := utils.NewMockTimeService()

	galasaHome, _ := utils.NewGalasaHome(mockFileSystem, env, "")
	console := utils.NewMockConsole()
	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)
	overrides, err := submitter.buildOverrideMap(commandParameters)

	assert.NotNil(t, err)
	assert.Nil(t, overrides)

	assert.Contains(t, err.Error(), "GAL1059")
}

func TestOverrideFileCorrectedWhenDefaultedAndOverridesFileNotExists(t *testing.T) {

	mockFileSystem := files.NewMockFileSystem()
	env := utils.NewMockEnv()
	galasaHome, err := utils.NewGalasaHome(mockFileSystem, env, "")
	if err != nil {
		assert.Fail(t, "Should not have failed! message = %s", err.Error())
	}

	commandParameters := utils.RunsSubmitCmdParameters{
		Overrides:        []string{"a=b"},
		OverrideFilePath: "",
	}

	mockLauncher := launcher.NewMockLauncher()
	mockTimeService := utils.NewMockTimeService()
	console := utils.NewMockConsole()
	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)
	err = submitter.correctOverrideFilePathParameter(&commandParameters)

	if err != nil {
		assert.Fail(t, "Should not have failed! message = %s", err.Error())
	}
	// We expect the default behaviour with missing command-line parameter, and missing overrides file in ~/.galasa to
	// result in an ignored overrideFilePath.
	assert.Equal(t, commandParameters.OverrideFilePath, "-")
}

func TestOverrideFileCorrectedWhenDefaultedAndNoOverridesFileDoesExist(t *testing.T) {

	mockFileSystem := files.NewMockFileSystem()
	env := utils.NewMockEnv()
	galasaHome, err := utils.NewGalasaHome(mockFileSystem, env, "")
	if err != nil {
		assert.Fail(t, "Should not have failed! message = %s", err.Error())
	}

	// A dummy overrides file in .galasa
	home, _ := mockFileSystem.GetUserHomeDirPath()
	separator := mockFileSystem.GetFilePathSeparator()
	path := home + separator + ".galasa" + separator + "overrides.properties"
	fileProps := make(map[string]interface{})
	fileProps["c"] = "d"
	props.WritePropertiesFile(mockFileSystem, path, fileProps)

	commandParameters := utils.RunsSubmitCmdParameters{
		Overrides:        []string{"a=b"},
		OverrideFilePath: "",
	}

	mockLauncher := launcher.NewMockLauncher()
	mockTimeService := utils.NewMockTimeService()
	console := utils.NewMockConsole()
	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)
	err = submitter.correctOverrideFilePathParameter(&commandParameters)

	if err != nil {
		assert.Fail(t, "Should not have failed! message = %s", err.Error())
	}
	assert.Equal(t, commandParameters.OverrideFilePath, path, "Wrong path of overrides file set. Expected %s", path)
}

func TestOverridesWithDashFileDontReadFromAnyFile(t *testing.T) {

	mockFileSystem := files.NewMockFileSystem()
	env := utils.NewMockEnv()
	galasaHome, err := utils.NewGalasaHome(mockFileSystem, env, "")

	commandParameters := utils.RunsSubmitCmdParameters{
		Overrides:        []string{"a=b"},
		OverrideFilePath: "-",
	}

	mockLauncher := launcher.NewMockLauncher()
	mockTimeService := utils.NewMockTimeService()
	console := utils.NewMockConsole()
	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)
	overrides, err := submitter.buildOverrideMap(commandParameters)

	if err != nil {
		assert.Fail(t, "Should not have failed! message = %s", err.Error())
	}
	assert.NotNil(t, overrides)
	assert.Contains(t, overrides, "a", "command-line override wasn't used.")
	assert.Equal(t, overrides["a"], "b", "command-line override not passed correctly.")
}

func TestValidateAndCorrectParametersSetsDefaultOverrideFile(t *testing.T) {

	mockFileSystem := files.NewMockFileSystem()
	env := utils.NewMockEnv()
	env.SetUserName("myuserid")

	galasaHome, err := utils.NewGalasaHome(mockFileSystem, env, "")
	if err != nil {
		assert.Fail(t, "Should not have failed! message = %s", err.Error())
	}

	commandParameters := &utils.RunsSubmitCmdParameters{
		Overrides:        []string{"a=b"},
		OverrideFilePath: "",
	}

	regexSelectValue := false
	submitSelectionFlags := &TestSelectionFlags{
		bundles:     new([]string),
		packages:    new([]string),
		tests:       new([]string),
		tags:        new([]string),
		classes:     new([]string),
		stream:      "myStream",
		regexSelect: &regexSelectValue,
	}

	mockLauncher := launcher.NewMockLauncher()

	mockTimeService := utils.NewMockTimeService()
	console := utils.NewMockConsole()
	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)
	err = submitter.validateAndCorrectParams(commandParameters, submitSelectionFlags)

	assert.Nil(t, err)
	assert.NotEmpty(t, commandParameters.OverrideFilePath)
}

func TestLocalLaunchCanUseAPortfolioOk(t *testing.T) {

	mockFileSystem := files.NewMockFileSystem()

	obrName := "myobr"
	bundleName := "myBundle"
	className := "myClass"

	portfolioFilePath := "myportfolio.yaml"
	_ = createTestPortfolioFile(t, mockFileSystem, portfolioFilePath, bundleName, className, "", obrName)

	env := utils.NewMockEnv()
	env.SetUserName("myuserid")

	galasaHome, err := utils.NewGalasaHome(mockFileSystem, env, "")
	if err != nil {
		assert.Fail(t, "Should not have failed! message = %s", err.Error())
	}

	commandParameters := &utils.RunsSubmitCmdParameters{}
	commandParameters.PortfolioFileName = portfolioFilePath

	regexSelectValue := false
	submitSelectionFlags := &TestSelectionFlags{
		bundles:     new([]string),
		packages:    new([]string),
		tests:       new([]string),
		tags:        new([]string),
		classes:     new([]string),
		stream:      "",
		regexSelect: &regexSelectValue,
	}

	mockLauncher := launcher.NewMockLauncher()

	mockTimeService := utils.NewMockTimeService()
	console := utils.NewMockConsole()
	submitter := NewSubmitter(
		galasaHome,
		mockFileSystem,
		mockLauncher,
		mockTimeService,
		env,
		console,
	)
	// Do the launching of the tests.
	err = submitter.ExecuteSubmitRuns(
		*commandParameters,
		submitSelectionFlags,
	)

	assert.Nil(t, err)

	launchesRecorded := mockLauncher.GetRecordedLaunchRecords()

	assert.Equal(t, 1, len(launchesRecorded))
	if len(launchesRecorded) > 0 {
		assert.Equal(t, obrName, launchesRecorded[0].ObrFromPortfolio)
		assert.Equal(t, bundleName+"/"+className, launchesRecorded[0].ClassName)
	}

}
