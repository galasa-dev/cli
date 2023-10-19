/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"log"

	"github.com/galasa.dev/cli/pkg/api"
	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/properties"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

//Objective: Allow user to do this:
//	properties set --namespace "framework" --name "hello" --value "newValue"
//  And then display a successful message or error

var (
	propertiesSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set the details of properties in a namespace.",
		Long: "Set the details of a property in a namespace. " +
			"If the property does not exist, a new property is created, otherwise the value for that property will be updated.",
		Args: cobra.NoArgs,
		Run:  executePropertiesSet,
	}

	// Variables set by cobra's command-line parsing.
	propertyValue string
)

func init() {
	propertiesSetCmd.PersistentFlags().StringVar(&propertyValue, "value", "", "the value of the property you want to create")
	parentCommand := propertiesCmd
	propertiesSetCmd.MarkFlagRequired("value")
	propertiesSetCmd.MarkPersistentFlagRequired("name")
	parentCommand.AddCommand(propertiesSetCmd)
}

func executePropertiesSet(cmd *cobra.Command, args []string) {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, logFileName)
	if err != nil {
		panic(err)
	}
	isCapturingLogs = true

	log.Println("Galasa CLI - Set ecosystem properties")

	// Get the ability to query environment variables.
	env := utils.NewEnvironment()

	galasaHome, err := utils.NewGalasaHome(fileSystem, env, CmdParamGalasaHomePath)
	if err != nil {
		panic(err)
	}

	// Read the bootstrap properties.
	var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
	var bootstrapData *api.BootstrapData
	bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, ecosystemBootstrap, urlService)
	if err != nil {
		panic(err)
	}

	var console = utils.NewRealConsole()

	apiServerUrl := bootstrapData.ApiServerURL
	log.Printf("The API server is at '%s'\n", apiServerUrl)

	// Call to process the command in a unit-testable way.
	err = properties.SetProperty(namespace, propertyName, propertyValue, apiServerUrl, console)
	if err != nil {
		panic(err)
	}

}
