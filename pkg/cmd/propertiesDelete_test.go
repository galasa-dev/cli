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

func TestPropertiesDeleteCommandInCommandCollectionHasName(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesDeleteCommand, err := commands.GetCommand(COMMAND_NAME_PROPERTIES_DELETE)
	assert.Nil(t, err)
	
	assert.Equal(t, COMMAND_NAME_PROPERTIES_DELETE, propertiesDeleteCommand.Name())
	assert.Nil(t, propertiesDeleteCommand.Values())
	assert.NotNil(t, propertiesDeleteCommand.CobraCommand())
}

func TestPropertiesDeleteHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"properties", "delete", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Displays the options for the 'properties delete' command.")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}

func TestPropertiesDeleteNoArgsReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"properties", "delete"}

	// When...
	err := Execute(factory, args)

	// Then...

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Contains(t, errText, "Error: required flag(s) \"name\", \"namespace\" not set")

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.NotNil(t, err)
}

func TestPropertiesDeleteWithoutName(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"properties", "delete", "--namespace", "jitters"}

	// When...
	err := Execute(factory, args)

	// Then...

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Contains(t, errText, "Error: required flag(s) \"name\" not set")

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.NotNil(t, err)
}

func TestPropertiesDeleteWithoutNamespace(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"properties", "delete", "--name", "jeepers"}

	// When...
	err := Execute(factory, args)

	// Then...
	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Contains(t, errText, "Error: required flag(s) \"namespace\" not set")

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.NotNil(t, err)
}

func TestPropertiesDeleteWithNameAndNamespace(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection, err := NewCommandCollection(factory)
	assert.Nil(t, err)

	propertiesDeleteCommand := commandCollection.GetCommand("properties delete")
	propertiesDeleteCommand.CobraCommand().RunE = func(cobraCmd *cobra.Command, args []string) error { return nil }
	
	var args []string = []string{"properties", "delete", "--namespace", "gyro", "--name", "space.ball"}

	// When...
	err = commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Equal(t, outText, "")
}