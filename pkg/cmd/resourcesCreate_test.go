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

func TestResourcesCreateCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	resourcesCreateCommand, err := commands.GetCommand(COMMAND_NAME_RESOURCES_CREATE)
	assert.Nil(t, err)
	
	assert.NotNil(t, resourcesCreateCommand)
	assert.Equal(t, COMMAND_NAME_RESOURCES_CREATE, resourcesCreateCommand.Name())
	assert.NotNil(t, resourcesCreateCommand.Values())
	assert.IsType(t, &ResourcesCreateCmdValues{}, resourcesCreateCommand.Values())
	assert.NotNil(t, resourcesCreateCommand.CobraCommand())
}

func TestResourcesCreateHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"resources", "create", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Displays the options for the 'resources create' command.")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}

func TestResourcesCreateNoFlagsReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"resources", "create"}

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

func TestResourcesCreateNameNamespaceValueReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection, err := NewCommandCollection(factory)
	assert.Nil(t, err)

	var resourcesCreateCommand GalasaCommand
	resourcesCreateCommand, err = commandCollection.GetCommand("resources create")
	assert.Nil(t, err)
	resourcesCreateCommand.CobraCommand().RunE = func(cobraCmd *cobra.Command, args []string) error { return nil }

	var args []string = []string{"resources", "create", "--file", "mince.yaml"}

	// When...
	err = commandCollection.Execute(args)

	// Then...
	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}
