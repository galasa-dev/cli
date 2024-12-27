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
	"github.com/galasa-dev/cli/pkg/properties"
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

//Objective: Allow user to do this:
//	properties namespaces get
//  And then display all namespaces in the cps or returns empty

type PropertiesNamespaceGetCmdValues struct {
	namespaceOutputFormat string
}

type PropertiesNamespaceGetCommand struct {
	values       *PropertiesNamespaceGetCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewPropertiesNamespaceGetCommand(
	factory spi.Factory,
	propertiesNamespaceCommand spi.GalasaCommand,
	propertiesCommand spi.GalasaCommand,
	commsFlagSet GalasaFlagSet,
) (spi.GalasaCommand, error) {

	cmd := new(PropertiesNamespaceGetCommand)

	err := cmd.init(factory, propertiesNamespaceCommand, commsFlagSet)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesNamespaceGetCommand) Name() string {
	return COMMAND_NAME_PROPERTIES_NAMESPACE_GET
}

func (cmd *PropertiesNamespaceGetCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *PropertiesNamespaceGetCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesNamespaceGetCommand) init(factory spi.Factory, propertiesNamespaceCommand spi.GalasaCommand, commsFlagSet GalasaFlagSet) error {
	var err error
	cmd.values = &PropertiesNamespaceGetCmdValues{}
	cmd.cobraCommand, err = cmd.createCobraCommand(factory, propertiesNamespaceCommand, commsFlagSet)
	return err
}

func (cmd *PropertiesNamespaceGetCommand) createCobraCommand(
	factory spi.Factory,
	propertiesNamespaceCommand spi.GalasaCommand,
	commsFlagSet GalasaFlagSet,
) (*cobra.Command, error) {

	var err error
	commsFlagSetValues := commsFlagSet.Values().(*CommsFlagSetValues)

	propertieNamespaceGetCobraCommand := &cobra.Command{
		Use:   "get",
		Short: "Get a list of namespaces.",
		Long:  "Get a list of namespaces within the CPS",
		Args:  cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			executionFunc := func() error {
				return cmd.executePropertiesNamespaceGet(factory, commsFlagSetValues)
			}
			return executeCommandWithRetries(factory, commsFlagSetValues, executionFunc)
		},
		Aliases: []string{"namespaces get"},
	}

	namespaceHasYamlFormat := false
	formatters := properties.GetFormatterNamesString(properties.CreateFormatters(namespaceHasYamlFormat))
	propertieNamespaceGetCobraCommand.PersistentFlags().StringVar(&cmd.values.namespaceOutputFormat, "format", "summary", "output format for the data returned. Supported formats are: "+formatters+".")

	propertiesNamespaceCommand.CobraCommand().AddCommand(propertieNamespaceGetCobraCommand)

	return propertieNamespaceGetCobraCommand, err
}

func (cmd *PropertiesNamespaceGetCommand) executePropertiesNamespaceGet(
	factory spi.Factory,
	commsFlagSetValues *CommsFlagSetValues,
) error {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	commsFlagSetValues.isCapturingLogs = true

	log.Println("Galasa CLI - Get ecosystem namespaces")

	// Get the ability to query environment variables.
	env := factory.GetEnvironment()

	var galasaHome spi.GalasaHome
	galasaHome, err = utils.NewGalasaHome(fileSystem, env, commsFlagSetValues.CmdParamGalasaHomePath)
	if err == nil {

		// Read the bootstrap properties.
		var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
		var bootstrapData *api.BootstrapData
		bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, commsFlagSetValues.bootstrap, urlService)
		if err == nil {

			var console = factory.GetStdOutConsole()

			apiServerUrl := bootstrapData.ApiServerURL
			log.Printf("The API server is at '%s'\n", apiServerUrl)

			var apiClient *galasaapi.APIClient
			authenticator := factory.GetAuthenticator(
				apiServerUrl,
				galasaHome,
			)
			apiClient, err = authenticator.GetAuthenticatedAPIClient()

			if err == nil {
				// Call to process the command in a unit-testable way.
				err = properties.GetPropertiesNamespaces(apiClient, cmd.values.namespaceOutputFormat, console)
			}
		}
	}
	return err
}
