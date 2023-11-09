/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/properties"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

//Objective: Allow user to do this:
//	properties namespaces get
//  And then display all namespaces in the cps or returns empty

type PropertiesNamespaceCmdValues struct {
	namespaceOutputFormat string
}

func createPropertiesNamespaceGetCmd(
	factory Factory,
	propertiesNamespaceCmd *cobra.Command,
	propertiesCmdValues *PropertiesCmdValues,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {

	var err error = nil

	// Allocate a memory block into which the parsed values of the command-line parameters are stored.
	propertiesNamespaceCmdValues := &PropertiesNamespaceCmdValues{}

	propertiesNamespaceGetCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a list of namespaces.",
		Long:  "Get a list of namespaces within the CPS",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePropertiesNamespaceGet(factory, cmd, args, propertiesNamespaceCmdValues, propertiesCmdValues, rootCmdValues)
		},
		Aliases: []string{"namespaces get"},
	}

	formatters := properties.GetFormatterNamesString(properties.CreateFormatters())
	propertiesNamespaceGetCmd.PersistentFlags().StringVar(&propertiesNamespaceCmdValues.namespaceOutputFormat, "format", "summary", "output format for the data returned. Supported formats are: "+formatters+".")
	parentCommand := propertiesNamespaceCmd
	parentCommand.AddCommand(propertiesNamespaceGetCmd)

	return propertiesNamespaceGetCmd, err
}

func executePropertiesNamespaceGet(
	factory Factory,
	cmd *cobra.Command,
	args []string,
	propertiesNamespaceCmdValues *PropertiesNamespaceCmdValues,
	propertiesCmdValues *PropertiesCmdValues,
	rootCmdValues *RootCmdValues,
) error {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err == nil {

		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Get ecosystem namespaces")

		// Get the ability to query environment variables.
		env := factory.GetEnvironment()

		galasaHome, err := utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Read the bootstrap properties.
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, propertiesCmdValues.ecosystemBootstrap, urlService)
			if err == nil {

				var console = factory.GetConsole()

				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API server is at '%s'\n", apiServerUrl)

				// Call to process the command in a unit-testable way.
				err = properties.GetNamespaceProperties(apiServerUrl, console)
			}
		}
	}
	return err
}
