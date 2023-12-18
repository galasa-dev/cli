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
	checkOutput("Displays the options for the 'resources update' command", "", "", factory, t)

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
	checkOutput("", "Error: required flag(s) \"file\" not set", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "required flag(s) \"file\" not set")
}

func TestResourcesUpdateNameNamespaceValueReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RESOURCES_UPDATE, factory, t)

	var args []string = []string{"resources", "update", "--file", "mince.yaml"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}
