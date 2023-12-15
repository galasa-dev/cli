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

	authCommand, err := commands.GetCommand(COMMAND_NAME_AUTH)
	assert.Nil(t, err)
	
	assert.NotNil(t, authCommand)
	assert.Equal(t, COMMAND_NAME_AUTH, authCommand.Name())
	assert.Nil(t, authCommand.Values())
	assert.NotNil(t, authCommand.CobraCommand())
}
