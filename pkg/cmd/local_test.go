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

func TestLocalNoCommandsProducesUsageReport(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"local"}

	// When...
	Execute(factory, args)

	// Then...
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Usage:")
	assert.Contains(t, outText, "galasactl local [command]")

	// We expect an exit code of 0 for this command.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)
}

func TestCommandListContainsLocalCommand(t *testing.T) {
	/// Given...
	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	// When...
	localCommand := commands.GetCommand(COMMAND_NAME_LOCAL)

	// Then...
	assert.NotNil(t, localCommand)
	assert.Equal(t, COMMAND_NAME_LOCAL, localCommand.Name())
	assert.Nil(t, localCommand.Values())
}
