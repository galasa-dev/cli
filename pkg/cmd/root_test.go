/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"errors"
	"testing"

	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCommandsCollectionHasARootCommand(t *testing.T) {
	factory := NewMockFactory()
	commands, err := NewCommandCollection(factory)
	assert.Nil(t, err)
	rootCommand := commands.GetCommand(COMMAND_NAME_ROOT)
	assert.NotNil(t, rootCommand)
}

func TestRootCommandInCommandCollectionHasAName(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	// When...
	commands, err := NewCommandCollection(factory)
	// Then...
	assert.Nil(t, err)
	rootCommand := commands.GetCommand(COMMAND_NAME_ROOT)

	assert.Equal(t, rootCommand.Name(), COMMAND_NAME_ROOT)
}

func TestRootCommandInCommandCollectionHasACobraCommand(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	// When...
	commands, err := NewCommandCollection(factory)
	// Then...
	assert.Nil(t, err)
	rootCommand := commands.GetRootCommand()

	assert.NotNil(t, rootCommand.CobraCommand())
}

func TestRootCommandInCommandCollectionHasAValuesStructure(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	// When...
	commands, err := NewCommandCollection(factory)
	// Then...
	assert.Nil(t, err)
	rootCommand := commands.GetRootCommand()

	values := rootCommand.Values()
	assert.NotNil(t, values)
	assert.IsType(t, &RootCmdValues{}, values)
}

func TestVersionFromCommandLine(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	var args []string = make([]string, 0)
	args = append(args, "--version")

	// When...
	err := Execute(factory, args)

	// Then...
	assert.Nil(t, err)

	// Lets check that the version came out.
	console := factory.GetStdOutConsole().(*utils.MockConsole)
	text := console.ReadText()
	assert.Contains(t, text, "galasactl version")
	versionString, _ := embedded.GetGalasaCtlVersion()
	assert.Contains(t, text, versionString)

	// We expect the exit code for this to be 0, so the final word should be nil.
	mockFinalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	assert.Nil(t, mockFinalWordHandler.ReportedObject)
}

func TestNoParamsFromCommandLine(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	var args []string = make([]string, 0)
	args = append(args, "")

	// When...
	Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	console := factory.GetStdOutConsole().(*utils.MockConsole)
	text := console.ReadText()
	assert.Contains(t, text, "A tool for controlling Galasa resources")

	// We expect an exit code of 1 for this command.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)
}

func TestCanGetNormalExitCodeAndErrorTextFromAnError(t *testing.T) {
	var err = errors.New("my text")
	errorText, exitCode, isStackTraceWanted := extractErrorDetails(err)
	assert.Equal(t, errorText, "my text", "Failed to extract the exit text from an error!")
	assert.Equal(t, 1, exitCode, "Wrong default exit code")
	assert.True(t, isStackTraceWanted, "We want stack trace from non-galasa errors")
}

func TestCanGetNormalExitCodeAndErrorTextFromAGalasaError(t *testing.T) {
	err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_MISSING_ACTION_FLAGS)
	errorText, exitCode, isStackTraceWanted := extractErrorDetails(err)
	assert.Contains(t, errorText, "GAL1009E", "Failed to extract the exit text from a galasa error!")
	assert.Equal(t, 1, exitCode, "Wrong default exit code")
	assert.False(t, isStackTraceWanted, "We don't want stack trace from galasa errors")
}

func TestCanGetTestsFailedExitCodeAndErrorTextFromATestFailedGalasaErrorPointer(t *testing.T) {
	err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_TESTS_FAILED, 14)
	errorText, exitCode, isStackTraceWanted := extractErrorDetails(err)
	assert.Contains(t, errorText, "GAL1017E", "Failed to extract the exit text from a galasa error!")
	assert.Equal(t, 2, exitCode, "Wrong default exit code")
	assert.False(t, isStackTraceWanted, "We don't want stack trace from galasa errors")
}


func TestRootHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Displays the options for the 'galasactl' command.")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}