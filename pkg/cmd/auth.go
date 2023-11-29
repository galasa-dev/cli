/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

type AuthCommand struct {
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors
// ------------------------------------------------------------------------------------------------
func NewAuthCommand(factory Factory, rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(AuthCommand)

	cmd.init(factory, rootCommand)
	return cmd, nil
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthCommand) GetName() string {
	return COMMAND_NAME_AUTH
}

func (cmd *AuthCommand) GetCobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *AuthCommand) GetValues() interface{} {
	// There are no values.
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthCommand) init(factory Factory, rootCommand GalasaCommand) error {
	var err error
	cmd.cobraCommand, err = createAuthCmd(factory, rootCommand.GetCobraCommand(), rootCommand.GetValues().(*RootCmdValues))
	return err
}

func createAuthCmd(factory Factory, parentCmd *cobra.Command, rootCmdValues *RootCmdValues) (*cobra.Command, error) {

	var err error = nil

	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manages authentication with a Galasa ecosystem",
		Long: "Manages authentication with a Galasa ecosystem using access tokens, " +
			"enabling secure interactions with the ecosystem.",
	}

	parentCmd.AddCommand(authCmd)

	err = createAuthCmdChildren(factory, authCmd, rootCmdValues)

	return authCmd, err
}

func createAuthCmdChildren(factory Factory, authCmd *cobra.Command, rootCmdValues *RootCmdValues) error {
	_, err := createAuthLoginCmd(factory, authCmd, rootCmdValues)
	if err == nil {
		_, err = createAuthLogoutCmd(factory, authCmd, rootCmdValues)
	}
	return err
}
