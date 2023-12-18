/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"testing"

	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestResourcesUpdateCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	resourcesUpdateCommand, err := commands.GetCommand(COMMAND_NAME_RESOURCES_UPDATE)
	assert.Nil(t, err)
	
	assert.NotNil(t, resourcesUpdateCommand)
	assert.Equal(t, COMMAND_NAME_RESOURCES_UPDATE, resourcesUpdateCommand.Name())
	assert.NotNil(t, resourcesUpdateCommand.Values())
	assert.IsType(t, &ResourcesUpdateCmdValues{}, resourcesUpdateCommand.Values())
	assert.NotNil(t, resourcesUpdateCommand.CobraCommand())
}

func TestResourcesUpdateHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"resources", "update", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Displays the options for the 'resources update' command.")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}

func TestResourcesUpdateNoFlagsReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"resources", "update"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Check what the user saw is reasonable.

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Contains(t, errText, "Error: required flag(s) \"file\" not set")

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.NotNil(t, err)
}

func TestResourcesUpdateNameNamespaceValueReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection, err := NewCommandCollection(factory)
	assert.Nil(t, err)

	var resourcesUpdateCommand GalasaCommand
	resourcesUpdateCommand, err = commandCollection.GetCommand("resources update")
	assert.Nil(t, err)
	resourcesUpdateCommand.CobraCommand().RunE = func(cobraCmd *cobra.Command, args []string) error { return nil }

	var args []string = []string{"resources", "update", "--file", "mince.yaml"}

	// When...
	err = commandCollection.Execute(args)

	// Then...
	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}
