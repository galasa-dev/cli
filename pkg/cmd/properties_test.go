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

func TestPropertiesCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesCommand, err := commands.GetCommand(COMMAND_NAME_PROPERTIES)
	assert.Nil(t, err)
	
	assert.NotNil(t, propertiesCommand)
	assert.Equal(t, COMMAND_NAME_PROPERTIES, propertiesCommand.Name())
	assert.NotNil(t, propertiesCommand.Values())
	assert.IsType(t, &PropertiesCmdValues{}, propertiesCommand.Values())
	assert.NotNil(t, propertiesCommand.CobraCommand())
}


func TestPropertiesHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"properties", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Displays the options for the 'properties' command.")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}

func TestPropertiesNoCommandsProducesUsageReport(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"properties"}

	// When...
	Execute(factory, args)

	// Then...
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Usage:")
	assert.Contains(t, outText, "galasactl properties [command]")

	// We expect an exit code of 0 for this command.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)
}