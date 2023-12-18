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

func TestRunsDownloadCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	runsDownloadCommand, err := commands.GetCommand(COMMAND_NAME_RUNS_DOWNLOAD)
	assert.Nil(t, err)
	
	assert.Equal(t, COMMAND_NAME_RUNS_DOWNLOAD, runsDownloadCommand.Name())
	assert.NotNil(t, runsDownloadCommand.Values())
	assert.IsType(t, &RunsDownloadCmdValues{}, runsDownloadCommand.Values())
	assert.NotNil(t, runsDownloadCommand.CobraCommand())
}


func TestRunsDownloadHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"runs", "download", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("Displays the options for the 'runs download' command.", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsDownloadNoFlagsReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"runs", "download"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Check what the user saw is reasonable.
	checkOutput("", "Error: required flag(s) \"name\" not set", "", factory, t)

	assert.NotNil(t, err)
}