/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"github.com/spf13/cobra"
)

//Objective: Allow user to do this:
//	auth tokens ...

type AuthTokensCmdValues struct {
	bootstrap string
}

type AuthTokensCommand struct {
	values       *AuthTokensCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewAuthTokensCommand(authCommand GalasaCommand, rootCmd GalasaCommand) (GalasaCommand, error) {
	cmd := new(AuthTokensCommand)

	err := cmd.init(authCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthTokensCommand) Name() string {
	return COMMAND_NAME_AUTH_TOKENS
}

func (cmd *AuthTokensCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *AuthTokensCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthTokensCommand) init(authCommand GalasaCommand) error {
	var err error
	cmd.values = &AuthTokensCmdValues{}
	cmd.cobraCommand, err = cmd.createAuthTokensCobraCmd(authCommand)
	return err
}

func (cmd *AuthTokensCommand) createAuthTokensCobraCmd(
	authCommand GalasaCommand,
) (*cobra.Command, error) {

	var err error
	authTokensCmd := &cobra.Command{
		Use:   "tokens",
		Short: "Queries tokens in an ecosystem",
		Long:  "Allows interaction to query tokens in Galasa Ecosystem",
		Args:  cobra.NoArgs,
	}

	addBootstrapFlag(authTokensCmd, &cmd.values.bootstrap)
	authCommand.CobraCommand().AddCommand(authTokensCmd)

	return authTokensCmd, err
}
