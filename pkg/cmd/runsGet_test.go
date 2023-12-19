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

func TestRunsGetCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	runsGetCommand, err := commands.GetCommand(COMMAND_NAME_RUNS_GET)
	assert.Nil(t, err)
	
	assert.Equal(t, COMMAND_NAME_RUNS_GET, runsGetCommand.Name())
	assert.NotNil(t, runsGetCommand.Values())
	assert.IsType(t, &RunsGetCmdValues{}, runsGetCommand.Values())
	assert.NotNil(t, runsGetCommand.CobraCommand())
}


func TestRunsGetHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	
	var args []string = []string{"runs", "get", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("Displays the options for the 'runs get' command.", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsGetNoFlagsReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_GET, factory, t)

	var args []string = []string{"runs", "get"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw was reasonable
	checkOutput("", "", "", factory, t)
}

func TestRunsGetActiveFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_GET, factory, t)

	var args []string = []string{"runs", "get", "--active"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw was reasonable
	checkOutput("", "", "", factory, t)
}

func TestRunsGetRequestorFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_GET, factory, t)

	var args []string = []string{"runs", "get", "--requestor", "galasateam"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw was reasonable
	checkOutput("", "", "", factory, t)
}

func TestRunsGetResultFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_GET, factory, t)

	var args []string = []string{"runs", "get", "--result", "passed"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw was reasonable
	checkOutput("", "", "", factory, t)
}

func TestRunsGetNameFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_GET, factory, t)

	var args []string = []string{"runs", "get", "--name", "gerald"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw was reasonable
	checkOutput("", "", "", factory, t)
}

func TestRunsGetageFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_GET, factory, t)

	var args []string = []string{"runs", "get", "--age", "10h"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw was reasonable
	checkOutput("", "", "", factory, t)
}

func TestRunsGetFormatFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_GET, factory, t)

	var args []string = []string{"runs", "get", "--format", "yaml"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw was reasonable
	checkOutput("", "", "", factory, t)
}

func TestRunsGetMultipleNameFlagsReturnsOK(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_GET, factory, t)

	var args []string = []string{"runs", "get", "--name", "C2020", "--name", "C4091"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw was reasonable
	checkOutput("", "", "", factory, t)
}

func TestRunsGetMultipleResultFlagsReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_GET, factory, t)

	var args []string = []string{"runs", "get", "--result", "passed", "--result", "failed"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw was reasonable
	checkOutput("", "", "", factory, t)
}

func TestRunsGetMultipleRequestorFlagsReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_GET, factory, t)

	var args []string = []string{"runs", "get", "--requestor", "root", "--requestor", "galasa"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw was reasonable
	checkOutput("", "", "", factory, t)
}
// flags are:
//   --format
//   --active
//   --age
//   --name
//   --requestor
//   --result
