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

// Objective: Allow user to update fields on an existing user record.

type UsersSetCmdValues struct {
	// The role field on the servers' user record is mutable.
	role string
}

type UsersSetCommand struct {
	values       *UsersSetCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewUsersSetCommand(
	factory spi.Factory,
	usersSetCommand spi.GalasaCommand,
	commsFlagSet GalasaFlagSet,
) (spi.GalasaCommand, error) {

	cmd := new(UsersSetCommand)

	err := cmd.init(factory, usersSetCommand, commsFlagSet)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *UsersSetCommand) Name() string {
	return COMMAND_NAME_USERS_SET
}

func (cmd *UsersSetCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *UsersSetCommand) Values() interface{} {
	// There are no values.
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *UsersSetCommand) init(factory spi.Factory, usersCommand spi.GalasaCommand, commsFlagSet GalasaFlagSet) error {
	var err error

	cmd.values = &UsersSetCmdValues{}
	cmd.cobraCommand, err = cmd.createCobraCmd(factory, usersCommand, commsFlagSet)

	return err
}

func (cmd *UsersSetCommand) createCobraCmd(
	factory spi.Factory,
	usersCommand spi.GalasaCommand,
	commsFlagSet GalasaFlagSet,
) (*cobra.Command, error) {

	var err error

	commsFlagSetValues := commsFlagSet.Values().(*CommsFlagSetValues)

	userCommandValues := usersCommand.Values().(*UsersCmdValues)
	usersSetCobraCmd := &cobra.Command{
		Use:     "set",
		Short:   "Set various mutable fields in a selected user record",
		Long:    "Set various mutable fields in a selected user record",
		Aliases: []string{COMMAND_NAME_USERS_GET},
		RunE: func(cobraCommand *cobra.Command, args []string) error {
			executionFunc := func() error {
				return cmd.executeUsersSet(factory, usersCommand.Values().(*UsersCmdValues), commsFlagSetValues)
			}
			return executeCommandWithRetries(factory, commsFlagSetValues, executionFunc)
		},
	}

	addLoginIdFlag(usersSetCobraCmd, MANDATORY_FLAG, userCommandValues)
	addRoleFlag(usersSetCobraCmd, cmd.values)

	usersCommand.CobraCommand().AddCommand(usersSetCobraCmd)

	return usersSetCobraCmd, err
}

func addRoleFlag(cmd *cobra.Command, userSetCmdValues *UsersSetCmdValues) {
	flagName := "role"
	description := "An optional field indicating the new role of the specified user."
	cmd.Flags().StringVar(&userSetCmdValues.role, flagName, "", description)
}

func (cmd *UsersSetCommand) executeUsersSet(
	factory spi.Factory,
	userCmdValues *UsersCmdValues,
	commsFlagSetValues *CommsFlagSetValues,
) error {

	var err error
	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	commsFlagSetValues.isCapturingLogs = true

	log.Println("Galasa CLI - Sets properties on an existibg user in the ecosystem")

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
				byteReader := factory.GetByteReader()
				err = users.SetUsers(userCmdValues.name, cmd.values.role, apiClient, console, byteReader)
			}
		}
	}

	return err
}
