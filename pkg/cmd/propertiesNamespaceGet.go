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
//	properties namespaces get
//  And then display all namespaces in the cps or returns empty

var (
	propertiesNamespaceGetCmd = &cobra.Command{
		Use:   "namespaces get",
		Short: "Get a list of namespaces.",
		Long:  "Get a list of namespaces within the CPS",
		Args:  cobra.NoArgs,
		Run:   executePropertiesNamespaceGet,
	}

	namespaceOutputFormat string
)

func init() {
	formatters := properties.GetFormatterNamesString(properties.CreateFormatters())
	propertiesNamespaceGetCmd.PersistentFlags().StringVar(&namespaceOutputFormat, "format", "summary", "output format for the data returned. Supported formats are: "+formatters+".")
	parentCommand := propertiesCmd
	parentCommand.AddCommand(propertiesNamespaceGetCmd)
}

func executePropertiesNamespaceGet(cmd *cobra.Command, args []string) {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, logFileName)
	if err != nil {
		panic(err)
	}
	isCapturingLogs = true

	log.Println("Galasa CLI - Get ecosystem namespaces")

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
	err = properties.GetNamespaceProperties(apiServerUrl, console)
	if err != nil {
		panic(err)
	}

}
