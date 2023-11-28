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

func createAuthLogoutCmd(factory Factory, parentCmd *cobra.Command, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	var err error

	authLogoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out from a Galasa ecosystem",
		Long:  "Log out from a Galasa ecosystem that you have previously logged in to",
		Aliases: []string{"auth logout"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAuthLogout(factory, cmd, args, rootCmdValues)
		},
	}

	parentCmd.AddCommand(authLogoutCmd)

	// There are no sub-command children to add to the command tree.

	return authLogoutCmd, err
}

func executeAuthLogout(
	factory Factory,
	cmd *cobra.Command,
	args []string,
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