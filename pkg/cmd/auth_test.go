/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	authCommand, err := commands.GetCommand(COMMAND_NAME_AUTH)
	assert.Nil(t, err)

	assert.NotNil(t, authCommand)
	assert.Equal(t, COMMAND_NAME_AUTH, authCommand.Name())
	assert.Nil(t, authCommand.Values())
	assert.NotNil(t, authCommand.CobraCommand())
}

func TestAuthHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	var args []string = []string{"auth", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Check what the user saw is reasonable.
	checkOutput("Displays the options for the 'auth' command.", "", factory, t)

	assert.Nil(t, err)
}

func TestAuthNoCommandsProducesUsageReport(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	var args []string = []string{"auth"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Check what the user saw is reasonable.
	checkOutput("Usage:\n  galasactl auth [command]", "", factory, t)

	assert.Nil(t, err)
}
