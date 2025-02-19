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
	"github.com/galasa-dev/cli/pkg/roles"
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

type RolesGetCmdValues struct {
	outputFormat string
}

type RolesGetCommand struct {
	values       *RolesGetCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewRolesGetCommand(
	factory spi.Factory,
	rolesGetCommand spi.GalasaCommand,
	commsFlagSet GalasaFlagSet,
) (spi.GalasaCommand, error) {

	cmd := new(RolesGetCommand)

	err := cmd.init(factory, rolesGetCommand, commsFlagSet)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *RolesGetCommand) Name() string {
	return COMMAND_NAME_ROLES_GET
}

func (cmd *RolesGetCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RolesGetCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *RolesGetCommand) init(factory spi.Factory, RolesCommand spi.GalasaCommand, commsFlagSet GalasaFlagSet) error {
	var err error

	cmd.values = &RolesGetCmdValues{}
	cmd.cobraCommand, err = cmd.createCobraCmd(factory, RolesCommand, commsFlagSet.Values().(*CommsFlagSetValues))

	return err
}

func (cmd *RolesGetCommand) createCobraCmd(
	factory spi.Factory,
	RolesCommand spi.GalasaCommand,
	commsFlagSetValues *CommsFlagSetValues,
) (*cobra.Command, error) {

	var err error

	RolesCommandValues := RolesCommand.Values().(*RolesCmdValues)
	RolesGetCobraCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get Roles used in a Galasa service",
		Long:    "Get a list of Roles from a Galasa service",
		Aliases: []string{COMMAND_NAME_ROLES_GET},
		RunE: func(cobraCommand *cobra.Command, args []string) error {
			return cmd.executeRolesGet(factory, RolesCommand.Values().(*RolesCmdValues), commsFlagSetValues)
		},
	}

	addRolesNameFlag(RolesGetCobraCmd, false, RolesCommandValues)

	formatters := roles.GetFormatterNamesAsString()
	RolesGetCobraCmd.Flags().StringVar(&cmd.values.outputFormat, "format", "summary", "the output format of the returned Roles. Supported formats are: "+formatters+".")

	RolesCommand.CobraCommand().AddCommand(RolesGetCobraCmd)

	return RolesGetCobraCmd, err
}

func (cmd *RolesGetCommand) executeRolesGet(
	factory spi.Factory,
	RolesCmdValues *RolesCmdValues,
	commsFlagSetValues *CommsFlagSetValues,
) error {

	var err error
	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, commsFlagSetValues.logFileName)
	if err == nil {
		commsFlagSetValues.isCapturingLogs = true
	
		log.Println("Galasa CLI - Get Roles from the ecosystem")
	
		env := factory.GetEnvironment()
	
		var galasaHome spi.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, commsFlagSetValues.CmdParamGalasaHomePath)
		if err == nil {

			timeService := factory.GetTimeService()
			commsRetrier := api.NewCommsRetrier(commsFlagSetValues.maxRetries, commsFlagSetValues.retryBackoffSeconds, timeService)

			// Read the bootstrap properties, retrying if a rate limit has been exceeded
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			loadBootstrapWithRetriesFunc := func() error {
				bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, commsFlagSetValues.bootstrap, urlService)
				return err
			}

			err = commsRetrier.ExecuteCommandWithRateLimitRetries(loadBootstrapWithRetriesFunc)
			if err == nil {
	
				var console = factory.GetStdOutConsole()
	
				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API server is at '%s'\n", apiServerUrl)
	
				authenticator := factory.GetAuthenticator(
					apiServerUrl,
					galasaHome,
				)

				commsRetrier, err = api.NewCommsRetrierWithAPIClient(
					commsFlagSetValues.maxRetries,
					commsFlagSetValues.retryBackoffSeconds,
					timeService,
					authenticator,
				)

				byteReader := factory.GetByteReader()
	
				if err == nil {
					getRolesFunc := func(apiClient *galasaapi.APIClient) error {
						return roles.GetRoles(RolesCmdValues.name, cmd.values.outputFormat, console, apiClient, byteReader)
					}
					err = commsRetrier.ExecuteCommandWithRetries(getRolesFunc)
				}
			}
		}
	}

	return err
}
