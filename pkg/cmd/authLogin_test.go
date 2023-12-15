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

func TestAuthLoginCommandInCommandCollection(t *testing.T) {
	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	authCommand, err := commands.GetCommand(COMMAND_NAME_AUTH_LOGIN)
	assert.Nil(t, err)

	assert.NotNil(t, authCommand)
	assert.Equal(t, COMMAND_NAME_AUTH_LOGIN, authCommand.Name())
	assert.NotNil(t, authCommand.Values())
	assert.IsType(t, &AuthLoginCmdValues{}, authCommand.Values())
	assert.NotNil(t, authCommand.CobraCommand())
}

func TestAuthLoginHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()

	// Note: No --maven or --gradle flags here:
	var args []string = []string{"auth", "login", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Displays the options for the 'auth login' command.")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Empty(t, errText)

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}

func TestAuthLoginNoFlagsReturnsNoError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection, err := NewCommandCollection(factory)
	assert.Nil(t, err)

	authLoginCommand := commandCollection.GetCommand("auth login")
	authLoginCommand.CobraCommand().RunE = func(cobraCmd *cobra.Command, args []string) error { return nil }

	var args []string = []string{"auth", "login"}

	// When...
	err = commandCollection.Execute(args)

	// Then...
	// Check what the user saw is reasonable.
	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Equal(t, errText, "")

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.Nil(t, err)
}