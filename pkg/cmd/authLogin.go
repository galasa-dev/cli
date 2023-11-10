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

func createAuthLoginCmd(factory Factory, parentCmd *cobra.Command, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	var err error = nil

	authLoginCmdValues := &AuthLoginCmdValues{}

	authLoginCmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate against a Galasa ecosystem",
		Long:  "Log in to a Galasa ecosystem using an existing access token",
		Args:  cobra.NoArgs,
		Aliases: []string{"auth login"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAuthLogin(factory, cmd, args, authLoginCmdValues, rootCmdValues)
		},
	}

	authLoginCmd.PersistentFlags().StringVarP(&authLoginCmdValues.bootstrap, "bootstrap", "b", "",
		"Bootstrap URL. Should start with 'http://' or 'file://'. "+
			"If it starts with neither, it is assumed to be a fully-qualified path. "+
			"If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. "+
			"Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties")

	parentCmd.AddCommand(authLoginCmd)

	// There are no sub-command children to add to the command tree.

	return authLoginCmd, err
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
