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
	factory Factory,
	propertiesNamespaceCommand GalasaCommand,
	propertiesCommand GalasaCommand,
	rootCommand GalasaCommand,
) (GalasaCommand, error) {

	cmd := new(PropertiesNamespaceGetCommand)

	err := cmd.init(factory, propertiesNamespaceCommand, propertiesCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesNamespaceGetCommand) GetName() string {
	return COMMAND_NAME_PROPERTIES_NAMESPACE_GET
}

func (cmd *PropertiesNamespaceGetCommand) GetCobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *PropertiesNamespaceGetCommand) GetValues() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesNamespaceGetCommand) init(factory Factory, propertiesNamespaceCommand GalasaCommand, propertiesCommand GalasaCommand, rootCommand GalasaCommand) error {
	var err error
	cmd.values = &PropertiesNamespaceGetCmdValues{}
	cmd.cobraCommand, err = cmd.createPropertiesNamespaceGetCobraCmd(
		factory, cmd.values, propertiesNamespaceCommand.GetCobraCommand(), propertiesCommand.GetValues().(*PropertiesCmdValues), rootCommand.GetValues().(*RootCmdValues))
	return err
}

func (cmd *PropertiesNamespaceGetCommand) createPropertiesNamespaceGetCobraCmd(
	factory Factory,
	propertiesNamespaceGetCmdValues *PropertiesNamespaceGetCmdValues,
	propertiesNamespaceCmd *cobra.Command,
	propertiesCmdValues *PropertiesCmdValues,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {

	var err error = nil

	propertieNamespaceGetCobraCommand := &cobra.Command{
		Use:   "get",
		Short: "Get a list of namespaces.",
		Long:  "Get a list of namespaces within the CPS",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePropertiesNamespaceGet(factory, cmd, args, propertiesNamespaceGetCmdValues, propertiesCmdValues, rootCmdValues)
		},
		Aliases: []string{"namespaces get"},
	}

	formatters := properties.GetFormatterNamesString(properties.CreateFormatters())
	propertieNamespaceGetCobraCommand.PersistentFlags().StringVar(&propertiesNamespaceGetCmdValues.namespaceOutputFormat, "format", "summary", "output format for the data returned. Supported formats are: "+formatters+".")
	parentCommand := propertiesNamespaceCmd
	parentCommand.AddCommand(propertieNamespaceGetCobraCommand)

	return propertieNamespaceGetCobraCommand, err
}

func executePropertiesNamespaceGet(
	factory Factory,
	cmd *cobra.Command,
	args []string,
	propertiesNamespaceGetCmdValues *PropertiesNamespaceGetCmdValues,
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

		var galasaHome utils.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Read the bootstrap properties.
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, propertiesCmdValues.ecosystemBootstrap, urlService)
			if err == nil {

				var console = factory.GetStdOutConsole()
				timeService := factory.GetTimeService()

				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API server is at '%s'\n", apiServerUrl)

				apiClient := auth.GetAuthenticatedAPIClient(apiServerUrl, fileSystem, galasaHome, timeService)

				// Call to process the command in a unit-testable way.
				err = properties.GetNamespaceProperties(apiClient, console)
			}
		}
	}
	return err
}
