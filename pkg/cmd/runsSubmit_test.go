/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"testing"

	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestRunsSubmitCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	cmd := commands.GetCommand(COMMAND_NAME_RUNS_SUBMIT)
	assert.Equal(t, COMMAND_NAME_RUNS_SUBMIT, cmd.Name())
	assert.NotNil(t, cmd.Values())
	assert.IsType(t, &utils.RunsSubmitCmdValues{}, cmd.Values())
	assert.NotNil(t, cmd.CobraCommand())
}
