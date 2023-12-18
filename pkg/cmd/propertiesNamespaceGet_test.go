/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"testing"

	// "github.com/galasa-dev/cli/pkg/utils"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPropertiesNamespaceGetCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesNamespaceGetCommand, err := commands.GetCommand(COMMAND_NAME_PROPERTIES_NAMESPACE_GET)
	assert.Nil(t, err)
	
	assert.Equal(t, COMMAND_NAME_PROPERTIES_NAMESPACE_GET, propertiesNamespaceGetCommand.Name())
	assert.NotNil(t, propertiesNamespaceGetCommand.Values())
	assert.IsType(t, &PropertiesNamespaceGetCmdValues{}, propertiesNamespaceGetCommand.Values())
	assert.NotNil(t, propertiesNamespaceGetCommand.CobraCommand())
}


func TestPropertiesNamespaceGetHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"properties", "namespaces", "get", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Displays the options for the 'properties namespaces get' command.")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}

func TestPropertiesNamespacesGetReturnsWithoutError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection, err := NewCommandCollection(factory)
	assert.Nil(t, err)
	
	var propertiesNamespacesGetCommand GalasaCommand
	propertiesNamespacesGetCommand, err = commandCollection.GetCommand("properties namespaces get")
	assert.Nil(t, err)
	assert.NotNil(t, propertiesNamespacesGetCommand)
	propertiesNamespacesGetCommand.CobraCommand().RunE = func(cobraCmd *cobra.Command, args []string) error { return nil }

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"properties", "namespaces", "get"}

	// When...
	err = commandCollection.Execute(args)

	// Then...
	// Check what the user saw is reasonable.
	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}

func TestPropertiesNamespacesGetFormatReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection, err := NewCommandCollection(factory)
	assert.Nil(t, err)

	var propertiesNamespacesGetCommand GalasaCommand
	propertiesNamespacesGetCommand, err = commandCollection.GetCommand("properties namespaces get")
	assert.Nil(t, err)
	assert.NotNil(t, propertiesNamespacesGetCommand)
	propertiesNamespacesGetCommand.CobraCommand().RunE = func(cobraCmd *cobra.Command, args []string) error { return nil }

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"properties", "namespaces", "get", "--format", "yaml"}

	// When...
	err = commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}
