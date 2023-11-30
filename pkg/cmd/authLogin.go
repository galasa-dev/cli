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
func NewAuthLoginCommand(factory Factory, authCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(AuthLoginComamnd)
	err := cmd.init(factory, authCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthLoginComamnd) GetName() string {
	return COMMAND_NAME_AUTH_LOGIN
}

func (cmd *AuthLoginComamnd) GetCobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *AuthLoginComamnd) GetValues() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *AuthLoginComamnd) init(factory Factory, authCommand GalasaCommand, rootCommand GalasaCommand) error {
	var err error = nil

	authLoginCmdValues := &AuthLoginCmdValues{}

	authLoginCmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to a Galasa ecosystem using an existing access token",
		Long: "Log in to a Galasa ecosystem using an existing access token stored in the 'galasactl.properties' file in your GALASA_HOME directory. " +
			"If you do not have an access token, request one through your ecosystem's web user interface " +
			"and follow the instructions on the web user interface to populate the 'galasactl.properties' file.",
		Args:    cobra.NoArgs,
		Aliases: []string{"auth login"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAuthLogin(factory, cmd, args, authLoginCmdValues, rootCommand.GetValues().(*RootCmdValues))
		},
	}

	// TODO: This bootstrap flag is defined in several places. De-duplication is needed. Put it into a function and call it from where it's needed.
	authLoginCmd.PersistentFlags().StringVarP(&authLoginCmdValues.bootstrap, "bootstrap", "b", "",
		"Bootstrap URL. Should start with 'http://' or 'file://'. "+
			"If it starts with neither, it is assumed to be a fully-qualified path. "+
			"If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. "+
			"Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties")

	authCommand.GetCobraCommand().AddCommand(authLoginCmd)

	cmd.values = authLoginCmdValues
	cmd.cobraCommand = authLoginCmd

	return err
}

func executeAuthLogin(
	factory Factory,
	cmd *cobra.Command,
	args []string,
	authLoginCmdValues *AuthLoginCmdValues,
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

		var galasaHome utils.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err != nil {
			panic(err)
		}

		// Read the bootstrap properties.
		var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
		var bootstrapData *api.BootstrapData
		bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, authLoginCmdValues.bootstrap, urlService)
		if err == nil {
			apiServerUrl := bootstrapData.ApiServerURL
			log.Printf("The API server is at '%s'\n", apiServerUrl)

			// Call to process the command in a unit-testable way.
			err = auth.Login(
				apiServerUrl,
				fileSystem,
				galasaHome,
			)
		}
	}
	return err
}
