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

func TestAuthogoutCommandInCommandCollectionHasName(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	authLogoutommand, err := commands.GetCommand(COMMAND_NAME_AUTH_LOGOUT)
	assert.Nil(t, err)
	
	assert.NotNil(t, authLogoutommand)
	assert.Equal(t, COMMAND_NAME_AUTH_LOGOUT, authLogoutommand.Name())
	assert.Nil(t, authLogoutommand.Values())
	assert.NotNil(t, authLogoutommand.CobraCommand())
}


func TestAuthLogoutHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	
	var args []string = []string{"auth", "logout", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("Displays the options for the 'auth logout' command.", "", "", factory, t)

	assert.Nil(t, err)
}
func TestAuthLogoutNoFlagsExecutesOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	var args []string = []string{"auth", "logout"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}