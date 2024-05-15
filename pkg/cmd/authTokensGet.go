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
	bootstrap          string
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
func (cmd *AuthTokensGetCommand) init(factory Factory, authTokensCommand GalasaCommand, authCommand GalasaCommand, rootCmd GalasaCommand) error {
	var err error

	cmd.values = &AuthTokensGetCmdValues{}
	cmd.cobraCommand, err = cmd.createCobraCmd(factory, authTokensCommand, authCommand, rootCmd)

	return err
}

func (cmd *AuthTokensGetCommand) createCobraCmd(
	factory Factory,
	authTokensCommand,
	authCommand GalasaCommand,
	rootCmd GalasaCommand,
) (*cobra.Command, error) {

	var err error

	authGetTokensCobraCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get tokens from a Galasa ecosystem",
		Long:    "Get tokens from a Galasa ecosystem which you are logged in to",
		Aliases: []string{COMMAND_NAME_AUTH_TOKENS_GET},
		RunE: func(cobraCommand *cobra.Command, args []string) error {
			return cmd.executeAuthTokensGet(factory, rootCmd.Values().(*RootCmdValues))
		},
	}

	addBootstrapFlag(authGetTokensCobraCmd, &cmd.values.bootstrap)
	// TO DO: implement format flag

	authTokensCommand.CobraCommand().AddCommand(authGetTokensCobraCmd)

	return authGetTokensCobraCmd, err
}

func (cmd *AuthTokensGetCommand) executeAuthTokensGet(
	factory Factory,
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
		if err != nil {
			panic(err)
		}

		// Read the bootstrap properties.
		var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
		var bootstrapData *api.BootstrapData
		bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, cmd.values.bootstrap, urlService)
		if err == nil {
			apiServerUrl := bootstrapData.ApiServerURL
			log.Printf("The API server is at '%s'\n", apiServerUrl)

			// Call to process the command in a unit-testable way.
			err = auth.GetTokens(
				apiServerUrl,
				fileSystem,
				galasaHome,
				env,
			)
		}
	}
	return err
}
