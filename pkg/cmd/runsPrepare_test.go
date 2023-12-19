/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunsPrepareCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	cmd, err := commands.GetCommand(COMMAND_NAME_RUNS_PREPARE)
	assert.Nil(t, err)
	
	assert.Equal(t, COMMAND_NAME_RUNS_PREPARE, cmd.Name())
	assert.NotNil(t, cmd.Values())
	assert.IsType(t, &RunsPrepareCmdValues{}, cmd.Values())
	assert.NotNil(t, cmd.CobraCommand())
}


func TestRunsPrepareHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	
	var args []string = []string{"runs", "prepare", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("Displays the options for the 'runs prepare' command.", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPrepareNoParametersReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	
	var args []string = []string{"runs", "prepare"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "Error: required flag(s) \"portfolio\" not set", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "required flag(s) \"portfolio\" not set")
}

func TestRunsPreparePortfolioFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "portfolio.file"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioAppendFlagReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--append"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}