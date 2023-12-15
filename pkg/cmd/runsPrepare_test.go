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

func TestRunsPrepareCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	cmd, err := commands.GetCommand(COMMAND_NAME_RUNS_PREPARE)
	assert.Equal(t, COMMAND_NAME_RUNS_PREPARE, cmd.Name())
	assert.NotNil(t, cmd.Values())
	assert.IsType(t, &RunsPrepareCmdValues{}, cmd.Values())
	assert.NotNil(t, cmd.CobraCommand())
	assert.Nil(t, err)
}
