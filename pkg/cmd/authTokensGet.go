/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/auth"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

//Objective: Allow user to do this:
//	auth tokens get
//  And then display all namespaces in the cps or returns empty

type AuthTokensGetCmdValues struct {
	tokensOutputFormat string
}
type AuthTokensGetCommand struct {
	values       *AuthTokensGetCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewAuthTokensGetCommand(
	factory Factory,
	authTokensCommand GalasaCommand,
	authTokens GalasaCommand,
	rootCmd GalasaCommand,
) (GalasaCommand, error) {

	cmd := new(AuthTokensGetCommand)

	err := cmd.init(factory, authTokensCommand, authTokens, rootCmd)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthTokensGetCommand) Name() string {
	return COMMAND_NAME_AUTH_TOKENS_GET
}

func (cmd *AuthTokensGetCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *AuthTokensGetCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthTokensGetCommand) init(factory Factory, authTokensCommand GalasaCommand, authLoginCommand GalasaCommand, rootCmd GalasaCommand) error {
	var err error

	cmd.values = &AuthTokensGetCmdValues{}
	cmd.cobraCommand, err = cmd.createCobraCmd(factory, authTokensCommand, authLoginCommand, rootCmd)

	return err
}

func (cmd *AuthTokensGetCommand) createCobraCmd(
	factory Factory,
	authTokensCommand,
	authLoginCmd GalasaCommand,
	rootCmd GalasaCommand,
) (*cobra.Command, error) {

	var err error

	authGetTokensCobraCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get a list of authentication tokens",
		Long:    "Get a list of tokens used for authenticating with the Galasa API server",
		Aliases: []string{COMMAND_NAME_AUTH_TOKENS_GET},
		RunE: func(cobraCommand *cobra.Command, args []string) error {
			return cmd.executeAuthTokensGet(factory, authLoginCmd.Values().(*AuthLoginCmdValues), rootCmd.Values().(*RootCmdValues))
		},
	}

	formatters := auth.GetFormatterNamesString(auth.CreateFormatters())
	authGetTokensCobraCmd.PersistentFlags().StringVar(&cmd.values.tokensOutputFormat, "format", "summary",
		"output format for the data returned. Supported formats are: "+formatters+".")

	authTokensCommand.CobraCommand().AddCommand(authGetTokensCobraCmd)

	return authGetTokensCobraCmd, err
}

func (cmd *AuthTokensGetCommand) executeAuthTokensGet(
	factory Factory,
	authLoginCmdValues *AuthLoginCmdValues,
	rootCmdValues *RootCmdValues,
) error {

	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err == nil {
		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Get tokens from the ecosystem")

		// Get the ability to query environment variables.
		env := factory.GetEnvironment()

		var galasaHome utils.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Read the bootstrap properties.
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, authLoginCmdValues.bootstrap, urlService)
			if err == nil {

				var console = factory.GetStdOutConsole()
				timeService := factory.GetTimeService()

				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API server is at '%s'\n", apiServerUrl)

				var apiClient *galasaapi.APIClient
				apiClient, err = auth.GetAuthenticatedAPIClient(apiServerUrl, fileSystem, galasaHome, timeService, env)

				if err == nil {
					// Call to process the command in a unit-testable way.
					err = auth.GetTokens(
						apiClient,
						cmd.values.tokensOutputFormat,
						console,
					)
				}
			}
		}
	}

	return err
}
