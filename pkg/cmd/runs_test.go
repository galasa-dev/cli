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

func TestCommandListContainsRunsCommand(t *testing.T) {
	/// Given...
	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	// When...
	runsCommand, err := commands.GetCommand(COMMAND_NAME_RUNS)
	assert.Nil(t, err)

	// Then...
	assert.NotNil(t, runsCommand)
	assert.Equal(t, COMMAND_NAME_RUNS, runsCommand.Name())
	assert.NotNil(t, runsCommand.Values())
	assert.IsType(t, &RunsCmdValues{}, runsCommand.Values())
}


func TestRunsHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"runs", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Displays the options for the 'runs' command.")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}