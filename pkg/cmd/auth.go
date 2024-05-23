/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

type AuthCommand struct {
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors
// ------------------------------------------------------------------------------------------------
func NewAuthCommand(rootCmd utils.GalasaCommand) (utils.GalasaCommand, error) {
	cmd := new(AuthCommand)

	cmd.init(rootCmd)
	return cmd, nil
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthCommand) Name() string {
	return COMMAND_NAME_AUTH
}

func (cmd *AuthCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *AuthCommand) Values() interface{} {
	// There are no values.
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthCommand) init(rootCmd utils.GalasaCommand) error {
	var err error
	cmd.cobraCommand, err = cmd.createCobraCommand(rootCmd)
	return err
}

func (cmd *AuthCommand) createCobraCommand(rootCmd utils.GalasaCommand) (*cobra.Command, error) {

	var err error

	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manages authentication with a Galasa ecosystem",
		Long: "Manages authentication with a Galasa ecosystem using access tokens, " +
			"enabling secure interactions with the ecosystem.",
	}

	rootCmd.CobraCommand().AddCommand(authCmd)

	return authCmd, err
}
