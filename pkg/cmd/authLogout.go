/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"

	"github.com/galasa-dev/cli/pkg/auth"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

type AuthLogoutCommand struct {
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewAuthLogoutCommand(factory Factory, authCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(AuthLogoutCommand)

	err := cmd.init(factory, authCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthLogoutCommand) Name() string {
	return COMMAND_NAME_AUTH_LOGOUT
}

func (cmd *AuthLogoutCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *AuthLogoutCommand) Values() interface{} {
	// There are no values.
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthLogoutCommand) init(factory Factory, authCommand GalasaCommand, rootCommand GalasaCommand) error {
	var err error
	cmd.cobraCommand, err = cmd.createAuthLogoutCobraCmd(factory, authCommand.CobraCommand(), rootCommand.Values().(*RootCmdValues))
	return err
}

func (cmd *AuthLogoutCommand) createAuthLogoutCobraCmd(factory Factory, parentCmd *cobra.Command, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	var err error

	authLogoutCmd := &cobra.Command{
		Use:     "logout",
		Short:   "Log out from a Galasa ecosystem",
		Long:    "Log out from a Galasa ecosystem that you have previously logged in to",
		Aliases: []string{"auth logout"},
		RunE: func(cobraCommand *cobra.Command, args []string) error {
			return cmd.executeAuthLogout(factory, rootCmdValues)
		},
	}

	parentCmd.AddCommand(authLogoutCmd)

	return authLogoutCmd, err
}

func (cmd *AuthLogoutCommand) executeAuthLogout(
	factory Factory,
	rootCmdValues *RootCmdValues,
) error {

	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err == nil {
		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Log out of an ecosystem")

		// Get the ability to query environment variables.
		env := factory.GetEnvironment()

		var galasaHome utils.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {
			err = auth.Logout(fileSystem, galasaHome)
		}
	}
	return err
}
