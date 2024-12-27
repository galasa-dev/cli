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
//	users get
type UsersGetCommand struct {
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewUsersGetCommand(
	factory spi.Factory,
	usersGetCommand spi.GalasaCommand,
	commsFlagSet GalasaFlagSet,
) (spi.GalasaCommand, error) {

	cmd := new(UsersGetCommand)

	err := cmd.init(factory, usersGetCommand, commsFlagSet)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *UsersGetCommand) Name() string {
	return COMMAND_NAME_USERS_GET
}

func (cmd *UsersGetCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *UsersGetCommand) Values() interface{} {
	// There are no values.
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *UsersGetCommand) init(factory spi.Factory, usersCommand spi.GalasaCommand, commsFlagSet GalasaFlagSet) error {
	var err error

	cmd.cobraCommand, err = cmd.createCobraCmd(factory, usersCommand, commsFlagSet)

	return err
}

func (cmd *UsersGetCommand) createCobraCmd(
	factory spi.Factory,
	usersCommand spi.GalasaCommand,
	commsFlagSet GalasaFlagSet,
) (*cobra.Command, error) {

	var err error

	commsFlagSetValues := commsFlagSet.Values().(*CommsFlagSetValues)

	userCommandValues := usersCommand.Values().(*UsersCmdValues)
	usersGetCobraCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get a list of users",
		Long:    "Get a list of users stored in the Galasa API server",
		Aliases: []string{COMMAND_NAME_USERS_GET},
		RunE: func(cobraCommand *cobra.Command, args []string) error {
			executionFunc := func() error {
				return cmd.executeUsersGet(factory, usersCommand.Values().(*UsersCmdValues), commsFlagSetValues)
			}
			return executeCommandWithRetries(factory, commsFlagSetValues, executionFunc)
		},
	}

	addLoginIdFlag(usersGetCobraCmd, false, userCommandValues)

	usersCommand.CobraCommand().AddCommand(usersGetCobraCmd)

	return usersGetCobraCmd, err
}

func (cmd *UsersGetCommand) executeUsersGet(
	factory spi.Factory,
	userCmdValues *UsersCmdValues,
	commsFlagSetValues *CommsFlagSetValues,
) error {

	var err error
	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	commsFlagSetValues.isCapturingLogs = true

	log.Println("Galasa CLI - Get users from the ecosystem")

	// Get the ability to query environment variables.
	env := factory.GetEnvironment()

	var galasaHome spi.GalasaHome
	galasaHome, err = utils.NewGalasaHome(fileSystem, env, commsFlagSetValues.CmdParamGalasaHomePath)
	if err == nil {

		// Read the bootstrap users.
		var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
		var bootstrapData *api.BootstrapData
		bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, commsFlagSetValues.bootstrap, urlService)
		if err == nil {

			var console = factory.GetStdOutConsole()

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
				err = users.GetUsers(userCmdValues.name, apiClient, console)
			}
		}
	}

	return err
}
