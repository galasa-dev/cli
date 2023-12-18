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

func TestResourcesApplyCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	resourcesApplyCommand, err := commands.GetCommand(COMMAND_NAME_RESOURCES_APPLY)
	assert.Nil(t, err)
	
	assert.NotNil(t, resourcesApplyCommand)
	assert.Equal(t, COMMAND_NAME_RESOURCES_APPLY, resourcesApplyCommand.Name())
	assert.NotNil(t, resourcesApplyCommand.Values())
	assert.IsType(t, &ResourcesApplyCmdValues{}, resourcesApplyCommand.Values())
	assert.NotNil(t, resourcesApplyCommand.CobraCommand())
}

func TestResourcesApplyHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"resources", "apply", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Displays the options for the 'resources apply' command.")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}

func TestResourcesApplyNoFlagsReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"resources", "apply"}

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

func TestResourcesApplyNameNamespaceValueReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection, err := NewCommandCollection(factory)
	assert.Nil(t, err)

	var resourcesApplyCommand GalasaCommand
	resourcesApplyCommand, err = commandCollection.GetCommand("resources apply")
	assert.Nil(t, err)
	resourcesApplyCommand.CobraCommand().RunE = func(cobraCmd *cobra.Command, args []string) error { return nil }

	var args []string = []string{"resources", "apply", "--file", "mince.yaml"}

	// When...
	err = commandCollection.Execute(args)

	// Then...
	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}