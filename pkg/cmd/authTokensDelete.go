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
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Objective: Allow user to do this:
// auth tokens delete --tokenid xxx
//
//	And delete the token of that id
type AuthTokensDeleteCommand struct {
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewAuthTokensDeleteCommand(
	factory spi.Factory,
	authTokensCommand spi.GalasaCommand,
	rootCmd spi.GalasaCommand,
) (spi.GalasaCommand, error) {

	cmd := new(AuthTokensDeleteCommand)

	err := cmd.init(factory, authTokensCommand, rootCmd)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthTokensDeleteCommand) Name() string {
	return COMMAND_NAME_AUTH_TOKENS_DELETE
}

func (cmd *AuthTokensDeleteCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *AuthTokensDeleteCommand) Values() interface{} {
	// There are no values.
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthTokensDeleteCommand) init(factory spi.Factory, authTokensCommand spi.GalasaCommand, rootCmd spi.GalasaCommand) error {
	var err error

	cmd.cobraCommand, err = cmd.createCobraCmd(factory, authTokensCommand, rootCmd)

	return err
}

func (cmd *AuthTokensDeleteCommand) createCobraCmd(
	factory spi.Factory,
	authTokensCommand,
	rootCmd spi.GalasaCommand,
) (*cobra.Command, error) {

	var err error

	authDeleteTokensCobraCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Deletes a personal access token",
		Long:    "Deletes a token used for authentication with the Galasa API server through the provided token id",
		Aliases: []string{COMMAND_NAME_AUTH_TOKENS_DELETE},
		RunE: func(cobraCommand *cobra.Command, args []string) error {
			return cmd.executeAuthTokensDelete(factory, authTokensCommand.Values().(*AuthTokensCmdValues), rootCmd.Values().(*RootCmdValues))
		},
	}

	authTokensCommand.CobraCommand().AddCommand(authDeleteTokensCobraCmd)

	return authDeleteTokensCobraCmd, err
}

func (cmd *AuthTokensDeleteCommand) executeAuthTokensDelete(
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

		log.Println("Galasa CLI - Delete a token from the ecosystem")

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
					// Call to process the command in a unit-testable way.
					//err = auth.DeleteTokens(apiClient, console)
					log.Print(console, apiClient)
				}
			}
		}
	}

	return err
}