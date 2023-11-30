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

func TestAuthCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	authCommand := commands.GetCommand(COMMAND_NAME_AUTH)
	assert.NotNil(t, authCommand)
}

func TestAuthCommandInCommandCollectionHasName(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	authCommand := commands.GetCommand(COMMAND_NAME_AUTH)
	assert.Equal(t, COMMAND_NAME_AUTH, authCommand.GetName())
	assert.Nil(t, authCommand.GetValues())
	assert.NotNil(t, authCommand.GetCobraCommand())
}
