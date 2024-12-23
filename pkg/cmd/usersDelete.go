/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/users"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Objective: Allow user to do this:
//
//	users delete
type UsersDeleteCommand struct {
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewUsersDeleteCommand(
	factory spi.Factory,
	usersDeleteCommand spi.GalasaCommand,
	commsCmd spi.GalasaCommand,
) (spi.GalasaCommand, error) {

	cmd := new(UsersDeleteCommand)

	err := cmd.init(factory, usersDeleteCommand, commsCmd)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *UsersDeleteCommand) Name() string {
	return COMMAND_NAME_USERS_DELETE
}

func (cmd *UsersDeleteCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *UsersDeleteCommand) Values() interface{} {
	// There are no values.
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *UsersDeleteCommand) init(factory spi.Factory, usersCommand spi.GalasaCommand, commsCmd spi.GalasaCommand) error {
	var err error

	cmd.cobraCommand, err = cmd.createCobraCmd(factory, usersCommand, commsCmd)

	return err
}

func (cmd *UsersDeleteCommand) createCobraCmd(
	factory spi.Factory,
	usersCommand,
	commsCmd spi.GalasaCommand,
) (*cobra.Command, error) {

	var err error

	commsCmdValues := commsCmd.Values().(*CommsCmdValues)

	userCommandValues := usersCommand.Values().(*UsersCmdValues)
	usersDeleteCobraCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Deletes a user by login ID",
		Long:    "Deletes a single user by their login ID from the ecosystem",
		Aliases: []string{COMMAND_NAME_USERS_DELETE},
		RunE: func(cobraCommand *cobra.Command, args []string) error {
			executionFunc := func() error {
				return cmd.executeUsersDelete(factory, usersCommand.Values().(*UsersCmdValues), commsCmdValues)
			}
			return executeCommandWithRetries(factory, commsCmdValues, executionFunc)
		},
	}

	addLoginIdFlag(usersDeleteCobraCmd, true, userCommandValues)

	usersCommand.CobraCommand().AddCommand(usersDeleteCobraCmd)

	return usersDeleteCobraCmd, err
}

func (cmd *UsersDeleteCommand) executeUsersDelete(
	factory spi.Factory,
	userCmdValues *UsersCmdValues,
	commsCmdValues *CommsCmdValues,
) error {

	var err error
	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()
	byteReader := factory.GetByteReader()

	commsCmdValues.isCapturingLogs = true

	log.Println("Galasa CLI - Delete user from the ecosystem")

	// Get the ability to query environment variables.
	env := factory.GetEnvironment()

	var galasaHome spi.GalasaHome
	galasaHome, err = utils.NewGalasaHome(fileSystem, env, commsCmdValues.CmdParamGalasaHomePath)
	if err == nil {

		// Read the bootstrap users.
		var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
		var bootstrapData *api.BootstrapData
		bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, commsCmdValues.bootstrap, urlService)
		if err == nil {

			apiServerUrl := bootstrapData.ApiServerURL
			log.Printf("The API server is at '%s'\n", apiServerUrl)

			authenticator := factory.GetAuthenticator(
				apiServerUrl,
				galasaHome,
			)

			var apiClient *galasaapi.APIClient
			apiClient, err = authenticator.GetAuthenticatedAPIClient()

			if err == nil {
				// Call to process the command in a unit-testable way.
				err = users.DeleteUser(userCmdValues.name, apiClient, byteReader)
			}
		}
	}

	return err
}
