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

func TestCommandListContainsRunsCommand(t *testing.T) {
	/// Given...
	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	// When...
	runsCommand := commands.GetCommand(COMMAND_NAME_RUNS)

	// Then...
	assert.NotNil(t, runsCommand)
	assert.Equal(t, COMMAND_NAME_RUNS, runsCommand.GetName())
	assert.NotNil(t, runsCommand.GetValues())
	assert.IsType(t, &RunsCmdValues{}, runsCommand.GetValues())
}
