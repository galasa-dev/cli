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

func TestAuthTokensCommandInCommandCollection(t *testing.T) {
	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	authCommand, err := commands.GetCommand(COMMAND_NAME_AUTH_TOKENS)
	assert.Nil(t, err)

	assert.NotNil(t, authCommand)
	assert.Equal(t, COMMAND_NAME_AUTH_TOKENS, authCommand.Name())
	assert.NotNil(t, authCommand.Values())
	assert.IsType(t, &AuthTokensCmdValues{}, authCommand.Values())
	assert.NotNil(t, authCommand.CobraCommand())
}

func TestAuthTokensHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	var args []string = []string{"auth", "tokens", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Check what the user saw is reasonable.
	checkOutput("Displays the options for the 'auth tokens' command.", "", factory, t)

	assert.Nil(t, err)
}

func TestAuthTokensNoFlagsReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection, _ := setupTestCommandCollection(COMMAND_NAME_AUTH_TOKENS, factory, t)

	var args []string = []string{"auth", "tokens"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw is reasonable.
	checkOutput("", "", factory, t)
}
