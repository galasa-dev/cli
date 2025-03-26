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

func TestCommandListContainsMonitorsDisableCommand(t *testing.T) {
    /// Given...
    factory := utils.NewMockFactory()
    commands, _ := NewCommandCollection(factory)

    // When...
    monitorsCommand, err := commands.GetCommand(COMMAND_NAME_MONITORS_DISABLE)
    assert.Nil(t, err)

    // Then...
    assert.NotNil(t, monitorsCommand)
    assert.Equal(t, COMMAND_NAME_MONITORS_DISABLE, monitorsCommand.Name())
    assert.NotNil(t, monitorsCommand.Values())
	assert.IsType(t, &MonitorsDisableCmdValues{}, monitorsCommand.Values())
}

func TestMonitorsDisableHelpFlagSetCorrectly(t *testing.T) {
    // Given...
    factory := utils.NewMockFactory()
	commandCollection, _ := setupTestCommandCollection(COMMAND_NAME_MONITORS_DISABLE, factory, t)

    var args []string = []string{"monitors", "disable", "--help"}

    // When...
    err := commandCollection.Execute(args)

    // Then...
    checkOutput("Disables a monitor with the given name in the Galasa service", "", factory, t)

    assert.Nil(t, err)
}

func TestMonitorsDisableWithNameFlagReturnsOk(t *testing.T) {
	// Given...
	factory := utils.NewMockFactory()
	commandCollection, _ := setupTestCommandCollection(COMMAND_NAME_MONITORS_DISABLE, factory, t)

	var args []string = []string{"monitors", "disable", "--name", "myMonitor"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw is reasonable.
	checkOutput("", "", factory, t)
}

func TestMonitorsDisableNoFlagsReturnsErrorMessage(t *testing.T) {
	// Given...
	factory := utils.NewMockFactory()
	commandCollection, _ := setupTestCommandCollection(COMMAND_NAME_MONITORS_DISABLE, factory, t)

	var args []string = []string{"monitors", "disable"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.NotNil(t, err)

	// Check what the user saw is reasonable.
	checkOutput("", `Error: required flag(s) "name" not set`, factory, t)
}

