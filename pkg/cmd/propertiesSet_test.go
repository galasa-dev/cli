/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPropertiesSetCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesSetCommand, err := commands.GetCommand(COMMAND_NAME_PROPERTIES_SET)
	assert.Nil(t, err)
	
	assert.Equal(t, COMMAND_NAME_PROPERTIES_SET, propertiesSetCommand.Name())
	assert.NotNil(t, propertiesSetCommand.Values())
	assert.IsType(t, &PropertiesSetCmdValues{}, propertiesSetCommand.Values())
	assert.NotNil(t, propertiesSetCommand.CobraCommand())
}


func TestPropertiesSetHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"properties", "set", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("Displays the options for the 'properties set' command.", "", "", factory, t)

	assert.Nil(t, err)
}

func TestPropertiesSetNoFlagsReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"properties", "set"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Check what the user saw is reasonable.
	checkOutput("", "Error: required flag(s) \"name\", \"namespace\", \"value\" not set", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "required flag(s) \"name\", \"namespace\", \"value\" not set")
}

func TestPropertiesSetNameNamespaceValueReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection, err := NewCommandCollection(factory)
	assert.Nil(t, err)

	var propertiesSetCommand GalasaCommand
	propertiesSetCommand, err = commandCollection.GetCommand("properties set")
	assert.Nil(t, err)
	propertiesSetCommand.CobraCommand().RunE = func(cobraCmd *cobra.Command, args []string) error { return nil }

	var args []string = []string{"properties", "set", "--namespace", "mince", "--name", "pies.are.so.tasty", "--value", "some kinda value"}

	// When...
	err = commandCollection.Execute(args)

	// Then...
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestPropertiesSetNamespaceOnlyReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"properties", "set", "--namespace", "sunshine"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Check what the user saw is reasonable.
	checkOutput("", "Error: required flag(s) \"name\", \"value\" not set", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "required flag(s) \"name\", \"value\" not set")
}

func TestPropertiesSetOnlyNameReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"properties", "set", "--name", "call.me.little.sunshine"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Check what the user saw is reasonable.
	checkOutput("", "Error: required flag(s) \"namespace\", \"value\" not set", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "required flag(s) \"namespace\", \"value\" not set")
}

func TestPropertiesOnlyValueReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"properties", "set", "--value", "ghost"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Check what the user saw is reasonable.
	checkOutput("", "Error: required flag(s) \"name\", \"namespace\" not set", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "required flag(s) \"name\", \"namespace\" not set")
}