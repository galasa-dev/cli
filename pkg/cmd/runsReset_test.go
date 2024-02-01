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

func TestRunsResetCommand(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	runsResetCommand, err := commands.GetCommand(COMMAND_NAME_RUNS_RESET)
	assert.Nil(t, err)

	assert.Equal(t, COMMAND_NAME_RUNS_RESET, runsResetCommand.Name())
	assert.NotNil(t, runsResetCommand.Values())
	assert.IsType(t, &RunsResetCmdValues{}, runsResetCommand.Values())
	assert.NotNil(t, runsResetCommand.CobraCommand())
}
