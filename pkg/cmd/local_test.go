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

func TestCommandListContainsLocalCommand(t *testing.T) {
	/// Given...
	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	// When...
	localCommand, err := commands.GetCommand(COMMAND_NAME_LOCAL)
	assert.Nil(t, err)

	// Then...
	assert.NotNil(t, localCommand)
	assert.Equal(t, COMMAND_NAME_LOCAL, localCommand.Name())
	assert.Nil(t, localCommand.Values())
}

func TestLocalHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"local", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Displays the options for the 'local' command.")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}

func TestLocalNoCommandsProducesUsageReport(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"local"}

	// When...
	Execute(factory, args)

	// Then...
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Usage:")
	assert.Contains(t, outText, "galasactl local [command]")

	// We expect an exit code of 0 for this command.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)
}