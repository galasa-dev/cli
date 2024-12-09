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
//	properties set --namespace "framework" --name "hello" --value "newValue"
//  And then display a successful message or error

type PropertiesSetCmdValues struct {
	// Variables set by cobra's command-line parsing.
	propertyValue string
}

type PropertiesSetCommand struct {
	values       *PropertiesSetCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewPropertiesSetCommand(factory spi.Factory, propertiesCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) (spi.GalasaCommand, error) {

	cmd := new(PropertiesSetCommand)
	err := cmd.init(factory, propertiesCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesSetCommand) Name() string {
	return COMMAND_NAME_PROPERTIES_SET
}

func (cmd *PropertiesSetCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *PropertiesSetCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *PropertiesSetCommand) init(factory spi.Factory, propertiesCommand spi.GalasaCommand, rootCmd spi.GalasaCommand) error {

	var err error
	cmd.values = &PropertiesSetCmdValues{}

	cmd.cobraCommand, err = cmd.createCobraCommand(factory, propertiesCommand, rootCmd.Values().(*RootCmdValues))

	return err
}

func (cmd *PropertiesSetCommand) createCobraCommand(
	factory spi.Factory,
	propertiesCommand spi.GalasaCommand,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {

	var err error
	propertiesCmdValues := propertiesCommand.Values().(*PropertiesCmdValues)

	propertiesSetCobraCmd := &cobra.Command{
		Use:   "set",
		Short: "Set the details of properties in a namespace.",
		Long: "Set the details of a property in a namespace. " +
			"If the property does not exist, a new property is created, otherwise the value for that property will be updated.",
		Args:    cobra.NoArgs,
		Aliases: []string{"properties set"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.executePropertiesSet(factory, propertiesCmdValues, rootCmdValues)
		},
	}

	propertiesSetCobraCmd.PersistentFlags().StringVarP(&cmd.values.propertyValue, "value", "v", "", "A mandatory flag indicating the value of the property you want to create. "+
		"Empty values and values with spaces must be put in quotation marks.")
	propertiesSetCobraCmd.MarkPersistentFlagRequired("value")

	propertiesCommand.CobraCommand().AddCommand(propertiesSetCobraCmd)

	// The name & namespace properties are mandatory for set.
	addNamespaceFlag(propertiesSetCobraCmd, true, propertiesCmdValues)
	addPropertyNameFlag(propertiesSetCobraCmd, true, propertiesCmdValues)

	return propertiesSetCobraCmd, err
}

func (cmd *PropertiesSetCommand) executePropertiesSet(
	factory spi.Factory,
	propertiesCmdValues *PropertiesCmdValues,
	rootCmdValues *RootCmdValues,
) error {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err == nil {

		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Set ecosystem properties")

		// Get the ability to query environment variables.
		env := factory.GetEnvironment()

		var galasaHome spi.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Read the bootstrap properties.
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, propertiesCmdValues.ecosystemBootstrap, urlService)
			if err == nil {

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
					err = properties.SetProperty(
						propertiesCmdValues.namespace,
						propertiesCmdValues.propertyName,
						cmd.values.propertyValue,
						apiClient)
				}
			}
		}
	}
	return err
}
