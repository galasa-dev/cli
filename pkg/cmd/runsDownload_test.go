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

func TestRunsDownloadCommandInCommandCollectionIsAsExpected(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	runsDownloadCommand := commands.GetCommand(COMMAND_NAME_RUNS_DOWNLOAD)
	assert.Equal(t, COMMAND_NAME_RUNS_DOWNLOAD, runsDownloadCommand.GetName())
	assert.NotNil(t, runsDownloadCommand.GetValues())
	assert.IsType(t, &RunsDownloadCmdValues{}, runsDownloadCommand.GetValues())
	assert.NotNil(t, runsDownloadCommand.GetCobraCommand())
}
