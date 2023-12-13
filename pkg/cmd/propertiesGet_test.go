/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"testing"

	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestPropertiesGetCommandInCommandCollectionHasName(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesGetCommand := commands.GetCommand(COMMAND_NAME_PROPERTIES_GET)
	assert.Equal(t, COMMAND_NAME_PROPERTIES_GET, propertiesGetCommand.Name())
	assert.NotNil(t, propertiesGetCommand.Values())
	assert.IsType(t, &PropertiesGetCmdValues{}, propertiesGetCommand.Values())
	assert.NotNil(t, propertiesGetCommand.CobraCommand())
}

func TestPropertiesGetHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"properties", "get", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Displays the options for the 'properties get' command.")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}

func TestPropertiesGetNoArgsReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"properties", "get"}
	// When...
	err := Execute(factory, args)

	// Then...

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Contains(t, errText, "Error: required flag(s) \"namespace\" not set")

	assert.NotNil(t, err)
}