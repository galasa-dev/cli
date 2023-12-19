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

func TestRunsSubmitCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	cmd, err := commands.GetCommand(COMMAND_NAME_RUNS_SUBMIT)
	assert.Nil(t, err)
	
	assert.Equal(t, COMMAND_NAME_RUNS_SUBMIT, cmd.Name())
	assert.NotNil(t, cmd.Values())
	assert.IsType(t, &utils.RunsSubmitCmdValues{}, cmd.Values())
	assert.NotNil(t, cmd.CobraCommand())
}


func TestRunsSubmitHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	
	var args []string = []string{"runs", "submit", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("Displays the options for the 'runs submit' command.", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsSubmitWithoutFlagsErrors(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"runs", "submit", "local"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Should throw an error asking for flags to be set
	checkOutput("", "required flag(s) \"class\", \"obr\" not set", "", factory, t)

	assert.NotNil(t, err, "err should have been set!")
	assert.Contains(t, err.Error(), "required flag(s) \"class\", \"obr\" not set")
}

func TestRunsSubmitExecutesWithPortfolio(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"runs", "submit", "local", "--class", "osgi.bundle/class.path"}

	// When...
	err := Execute(factory, args)

	// Then...
	// check what the user sees is acceptible
	checkOutput("", "required flag(s) \"obr\" not set", "", factory, t)
	
	// Should throw an error asking for flags to be set
	assert.NotNil(t, err, "err should have been set!")
	assert.Contains(t, err.Error(), "required flag(s) \"obr\" not set")
}