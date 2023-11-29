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

func TestCommandListContainsProjectCommand(t *testing.T) {
	/// Given...
	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	// When...
	projectCommand := commands.GetCommand(COMMAND_NAME_PROJECT)

	// Then...
	assert.NotNil(t, projectCommand)
	assert.Equal(t, COMMAND_NAME_PROJECT, projectCommand.GetName())
	assert.Nil(t, projectCommand.GetValues())
}
