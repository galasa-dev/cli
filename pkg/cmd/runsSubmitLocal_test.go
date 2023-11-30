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

func TestRunsSubmitLocalCommandInCommandCollectionIsAsExpected(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	cmd := commands.GetCommand(COMMAND_NAME_RUNS_SUBMIT_LOCAL)
	assert.Equal(t, COMMAND_NAME_RUNS_SUBMIT_LOCAL, cmd.GetName())
	assert.NotNil(t, cmd.GetValues())
	assert.IsType(t, &RunsSubmitLocalCmdValues{}, cmd.GetValues())
	assert.NotNil(t, cmd.GetCobraCommand())
}
