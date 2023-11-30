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

func TestAuthLoginCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	authCommand := commands.GetCommand(COMMAND_NAME_AUTH_LOGIN)
	assert.NotNil(t, authCommand)
}

func TestAuthLoginCommandInCommandCollectionHasName(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	authCommand := commands.GetCommand(COMMAND_NAME_AUTH_LOGIN)
	assert.Equal(t, COMMAND_NAME_AUTH_LOGIN, authCommand.GetName())
	assert.NotNil(t, authCommand.GetValues())
	assert.IsType(t, &AuthLoginCmdValues{}, authCommand.GetValues())
	assert.NotNil(t, authCommand.GetCobraCommand())
}
