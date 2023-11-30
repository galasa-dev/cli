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

func TestRunsGetCommandInCommandCollectionIsAsExpected(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	runsGetCommand := commands.GetCommand(COMMAND_NAME_RUNS_GET)
	assert.Equal(t, COMMAND_NAME_RUNS_GET, runsGetCommand.Name())
	assert.NotNil(t, runsGetCommand.Values())
	assert.IsType(t, &RunsGetCmdValues{}, runsGetCommand.Values())
	assert.NotNil(t, runsGetCommand.CobraCommand())
}
