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
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Objective: Allow user to do this:
//
//		auth tokens get
//	 And then display all tokens or returns empty
type AuthTokensGetCommand struct {
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewAuthTokensGetCommand(
	factory spi.Factory,
	authTokensCommand spi.GalasaCommand,
	rootCmd spi.GalasaCommand,
) (spi.GalasaCommand, error) {

	cmd := new(AuthTokensGetCommand)

	err := cmd.init(factory, authTokensCommand, rootCmd)
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
	// There are no values.
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthTokensGetCommand) init(factory spi.Factory, authTokensCommand spi.GalasaCommand, rootCmd spi.GalasaCommand) error {
	var err error

	cmd.cobraCommand, err = cmd.createCobraCmd(factory, authTokensCommand, rootCmd)

	return err
}

func (cmd *AuthTokensGetCommand) createCobraCmd(
	factory spi.Factory,
	authTokensCommand,
	rootCmd spi.GalasaCommand,
) (*cobra.Command, error) {

	var err error

	authTokensGetCommandValues := authTokensCommand.Values().(*AuthTokensCmdValues)
	authGetTokensCobraCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get a list of authentication tokens",
		Long:    "Get a list of tokens used for authentication with the Galasa API server",
		Aliases: []string{COMMAND_NAME_AUTH_TOKENS_GET},
		RunE: func(cobraCommand *cobra.Command, args []string) error {
			return cmd.executeAuthTokensGet(factory, authTokensCommand.Values().(*AuthTokensCmdValues), rootCmd.Values().(*RootCmdValues))
		},
	}

	addLoginIdFlagToAuthTokensGet(authGetTokensCobraCmd, authTokensGetCommandValues)
	authTokensCommand.CobraCommand().AddCommand(authGetTokensCobraCmd)

	return authGetTokensCobraCmd, err
}

func (cmd *AuthTokensGetCommand) executeAuthTokensGet(
	factory spi.Factory,
	authTokenCmdValues *AuthTokensCmdValues,
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

		var galasaHome spi.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Read the bootstrap properties.
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, authTokenCmdValues.bootstrap, urlService)
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
					err = auth.GetTokens(apiClient, console, authTokenCmdValues.loginId)
				}
			}
		}
	}

	return err
}

func addLoginIdFlagToAuthTokensGet(cmd *cobra.Command, authTokensGetCmdValues *AuthTokensCmdValues) {

	flagName := "user"
	var description string = "Optional. Retrieves a list of access tokens for the user with the given username."

	cmd.Flags().StringVar(&authTokensGetCmdValues.loginId, flagName, "", description)
}
