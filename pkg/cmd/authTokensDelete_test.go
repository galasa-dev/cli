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

func TestAuthTokensDeleteCommandInCommandCollection(t *testing.T) {

	factory := utils.NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	AuthTokensDeleteCommand, err := commands.GetCommand(COMMAND_NAME_AUTH_TOKENS_DELETE)
	assert.Nil(t, err)

	assert.Equal(t, COMMAND_NAME_AUTH_TOKENS_DELETE, AuthTokensDeleteCommand.Name())
	assert.Nil(t, AuthTokensDeleteCommand.Values())
	assert.NotNil(t, AuthTokensDeleteCommand.CobraCommand())
}

func TestAuthTokensDeleteHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := utils.NewMockFactory()

	var args []string = []string{"auth", "tokens", "get", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("Displays the options for the 'auth tokens get' command.", "", factory, t)

	assert.Nil(t, err)
}

func TestAuthTokensDeleteReturnsWithoutError(t *testing.T) {
	// Given...
	factory := utils.NewMockFactory()
	commandCollection, _ := setupTestCommandCollection(COMMAND_NAME_AUTH_TOKENS_DELETE, factory, t)

	var args []string = []string{"auth", "tokens", "delete"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	// Check what the user saw is reasonable.
	checkOutput("", "", factory, t)

	assert.Nil(t, err)
}
