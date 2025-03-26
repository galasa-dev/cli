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

func TestCommandListContainsMonitorsEnableCommand(t *testing.T) {
    /// Given...
    factory := utils.NewMockFactory()
    commands, _ := NewCommandCollection(factory)

    // When...
    monitorsCommand, err := commands.GetCommand(COMMAND_NAME_MONITORS_ENABLE)
    assert.Nil(t, err)

    // Then...
    assert.NotNil(t, monitorsCommand)
    assert.Equal(t, COMMAND_NAME_MONITORS_ENABLE, monitorsCommand.Name())
    assert.NotNil(t, monitorsCommand.Values())
	assert.IsType(t, &MonitorsEnableCmdValues{}, monitorsCommand.Values())
}

func TestMonitorsEnableHelpFlagSetCorrectly(t *testing.T) {
    // Given...
    factory := utils.NewMockFactory()
	commandCollection, _ := setupTestCommandCollection(COMMAND_NAME_MONITORS_ENABLE, factory, t)

    var args []string = []string{"monitors", "enable", "--help"}

    // When...
    err := commandCollection.Execute(args)

    // Then...
    checkOutput("Enables a given monitor in the Galasa service", "", factory, t)

    assert.Nil(t, err)
}

func TestMonitorsEnableWithNameFlagReturnsOk(t *testing.T) {
	// Given...
	factory := utils.NewMockFactory()
	commandCollection, _ := setupTestCommandCollection(COMMAND_NAME_MONITORS_ENABLE, factory, t)

	var args []string = []string{"monitors", "enable", "--name", "myMonitor"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.Nil(t, err)

	// Check what the user saw is reasonable.
	checkOutput("", "", factory, t)
}

func TestMonitorsEnableNoFlagsReturnsErrorMessage(t *testing.T) {
	// Given...
	factory := utils.NewMockFactory()
	commandCollection, _ := setupTestCommandCollection(COMMAND_NAME_MONITORS_ENABLE, factory, t)

	var args []string = []string{"monitors", "enable"}

	// When...
	err := commandCollection.Execute(args)

	// Then...
	assert.NotNil(t, err)

	// Check what the user saw is reasonable.
	checkOutput("", `Error: required flag(s) "name" not set`, factory, t)
}

