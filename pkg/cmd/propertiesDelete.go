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
	"github.com/galasa-dev/cli/pkg/properties"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

//Objective: Allow user to do this:
//	properties delete --namespace "framework" --name "hello"
//  And then display a successful message or error

func createPropertiesDeleteCmd(factory Factory, parentCmd *cobra.Command, propertiesCmdValues *PropertiesCmdValues, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	var err error = nil

	propertiesDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a property in a namespace.",
		Long:  "Delete a property and its value in a namespace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePropertiesDelete(factory, cmd, args, propertiesCmdValues, rootCmdValues)
		},
		Aliases: []string{"properties delete"},
	}

	parentCmd.AddCommand(propertiesDeleteCmd)

	addNameProperty(propertiesDeleteCmd, true, propertiesCmdValues)
	addNamespaceProperty(propertiesDeleteCmd, true, propertiesCmdValues)

	// There are no sub-commands to add to the tree.

	return propertiesDeleteCmd, err
}

func executePropertiesDelete(factory Factory, cmd *cobra.Command, args []string, propertiesCmdValues *PropertiesCmdValues, rootCmdValues *RootCmdValues) error {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err == nil {

		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Delete ecosystem properties")

		// Get the ability to query environment variables.
		env := factory.GetEnvironment()

		var galasaHome utils.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Read the bootstrap properties.
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, propertiesCmdValues.ecosystemBootstrap, urlService)
			if err == nil {

				timeService := factory.GetTimeService()

				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API server is at '%s'\n", apiServerUrl)

				apiClient := auth.GetAuthenticatedAPIClient(apiServerUrl, fileSystem, galasaHome, timeService)

				// Call to process the command in a unit-testable way.
				err = properties.DeleteProperty(propertiesCmdValues.namespace, propertiesCmdValues.propertyName, apiClient)
			}
		}
	}
	return err
}
