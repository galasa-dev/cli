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

func TestProjectNoCommandsProducesUsageReport(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"project"}

	// When...
	Execute(factory, args)

	// Then...
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Usage:")
	assert.Contains(t, outText, "galasactl project [command]")

	// We expect an exit code of 0 for this command.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)
}

func TestCommandListContainsProjectCommand(t *testing.T) {
	/// Given...
	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	// When...
	projectCommand := commands.GetCommand(COMMAND_NAME_PROJECT)

	// Then...
	assert.NotNil(t, projectCommand)
	assert.Equal(t, COMMAND_NAME_PROJECT, projectCommand.Name())
	assert.Nil(t, projectCommand.Values())
}
