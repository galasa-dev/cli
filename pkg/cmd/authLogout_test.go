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

	authLogoutommand := commands.GetCommand(COMMAND_NAME_AUTH_LOGOUT)
	assert.NotNil(t, authLogoutommand)
	assert.Equal(t, COMMAND_NAME_AUTH_LOGOUT, authLogoutommand.Name())
	assert.Nil(t, authLogoutommand.Values())
	assert.NotNil(t, authLogoutommand.CobraCommand())
}
