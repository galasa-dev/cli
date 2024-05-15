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

func TestAuthTokensGetCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	AuthTokensGetCommand, err := commands.GetCommand(COMMAND_NAME_AUTH_TOKENS_GET)
	assert.Nil(t, err)

	assert.Equal(t, COMMAND_NAME_AUTH_TOKENS_GET, AuthTokensGetCommand.Name())
	assert.NotNil(t, AuthTokensGetCommand.Values())
	assert.IsType(t, &AuthTokensGetCmdValues{}, AuthTokensGetCommand.Values())
	assert.NotNil(t, AuthTokensGetCommand.CobraCommand())
}

func TestAuthTokensGetHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	var args []string = []string{"auth", "tokens", "get", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("Displays the options for the 'auth tokens get' command.", "", factory, t)

	assert.Nil(t, err)
}

func TestAuthTokenssGetReturnsWithoutError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection, _ := setupTestCommandCollection(COMMAND_NAME_AUTH_TOKENS_GET, factory, t)

	var args []string = []string{"auth", "tokens", "get"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	// Check what the user saw is reasonable.
	checkOutput("", "", factory, t)

	assert.Nil(t, err)
}

// TO DO: implement different formats
// func TestAuthTokenssGetFormatReturnsOk(t *testing.T) {
// 	// Given...
// 	factory := NewMockFactory()
// 	commandCollection, cmd := setupTestCommandCollection(COMMAND_NAME_AUTH_TOKENS_GET, factory, t)

// 	var args []string = []string{"auth", "tokens", "get", "--format", "yaml"}

// 	// When...
// 	err := commandCollection.Execute(args)

// 	// Then...
// 	assert.Nil(t, err)

// 	// Check what the user saw is reasonable.
// 	checkOutput("", "", factory, t)

// 	assert.Contains(t, cmd.Values().(*AuthTokensGetCmdValues).tokensOutputFormat, "yaml")
// }
