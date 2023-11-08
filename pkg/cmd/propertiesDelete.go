/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/properties"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

//Objective: Allow user to do this:
//	properties delete --namespace "framework" --name "hello"
//  And then display a successful message or error

func createPropertiesDeleteCmd(parentCmd *cobra.Command, propertiesCmdValues *PropertiesCmdValues, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	var err error = nil

	propertiesDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a property in a namespace.",
		Long:  "Delete a property and its value in a namespace",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			executePropertiesDelete(cmd, args, propertiesCmdValues, rootCmdValues)
		},
		Aliases: []string{"properties delete"},
	}

	parentCmd.AddCommand(propertiesDeleteCmd)

	addNameProperty(propertiesDeleteCmd, true, propertiesCmdValues)
	addNamespaceProperty(propertiesDeleteCmd, true, propertiesCmdValues)

	// There are no sub-commands to add to the tree.

	return propertiesDeleteCmd, err
}

func executePropertiesDelete(cmd *cobra.Command, args []string, propertiesCmdValues *PropertiesCmdValues, rootCmdValues *RootCmdValues) {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err != nil {
		panic(err)
	}
	rootCmdValues.isCapturingLogs = true

	log.Println("Galasa CLI - Delete ecosystem properties")

	// Get the ability to query environment variables.
	env := utils.NewEnvironment()

	galasaHome, err := utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
	if err != nil {
		panic(err)
	}

	// Read the bootstrap properties.
	var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
	var bootstrapData *api.BootstrapData
	bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, propertiesCmdValues.ecosystemBootstrap, urlService)
	if err != nil {
		panic(err)
	}

	apiServerUrl := bootstrapData.ApiServerURL
	log.Printf("The API server is at '%s'\n", apiServerUrl)

	// Call to process the command in a unit-testable way.
	err = properties.DeleteProperty(propertiesCmdValues.namespace, propertiesCmdValues.propertyName, apiServerUrl)
	if err != nil {
		panic(err)
	}

}
