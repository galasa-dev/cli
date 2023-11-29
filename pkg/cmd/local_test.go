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

func TestCommandListContainsLocalCommand(t *testing.T) {
	/// Given...
	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	// When...
	localCommand := commands.GetCommand(COMMAND_NAME_LOCAL)

	// Then...
	assert.NotNil(t, localCommand)
	assert.Equal(t, COMMAND_NAME_LOCAL, localCommand.GetName())
	assert.Nil(t, localCommand.GetValues())
}
