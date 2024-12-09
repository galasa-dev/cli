/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/spi"

	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

type AuthLoginCmdValues struct {
	bootstrap string
}

type AuthLoginComamnd struct {
	values       *AuthLoginCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewAuthLoginCommand(factory spi.Factory, authCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) (spi.GalasaCommand, error) {
	cmd := new(AuthLoginComamnd)
	err := cmd.init(factory, authCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthLoginComamnd) Name() string {
	return COMMAND_NAME_AUTH_LOGIN
}

func (cmd *AuthLoginComamnd) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *AuthLoginComamnd) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthLoginComamnd) init(factory spi.Factory, authCommand spi.GalasaCommand, rootCmd spi.GalasaCommand) error {
	var err error

	cmd.values = &AuthLoginCmdValues{}

	cmd.cobraCommand, err = cmd.createCobraCommand(factory, authCommand, rootCmd)

	return err
}

func (cmd *AuthLoginComamnd) createCobraCommand(
	factory spi.Factory,
	authCommand spi.GalasaCommand,
	rootCmd spi.GalasaCommand,
) (*cobra.Command, error) {

	var err error
	authLoginCobraCmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to a Galasa ecosystem using an existing access token",
		Long: "Log in to a Galasa ecosystem using an existing access token stored in the 'galasactl.properties' file in your GALASA_HOME directory. " +
			"If you do not have an access token, request one through your ecosystem's web user interface " +
			"and follow the instructions on the web user interface to populate the 'galasactl.properties' file.",
		Args:    cobra.NoArgs,
		Aliases: []string{"auth login"},
		RunE: func(cobraCommand *cobra.Command, args []string) error {
			return cmd.executeAuthLogin(factory, rootCmd.Values().(*RootCmdValues))
		},
	}

	addBootstrapFlag(authLoginCobraCmd, &cmd.values.bootstrap)

	authCommand.CobraCommand().AddCommand(authLoginCobraCmd)

	return authLoginCobraCmd, err
}

func (cmd *AuthLoginComamnd) executeAuthLogin(
	factory spi.Factory,
	rootCmdValues *RootCmdValues,
) error {

	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err == nil {
		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Log in to an ecosystem")

		// Get the ability to query environment variables.
		env := factory.GetEnvironment()

		var galasaHome spi.GalasaHome
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

			authenticator := factory.GetAuthenticator(
				apiServerUrl,
				galasaHome,
			)
			err = authenticator.Login()
		}
	}
	return err
}
