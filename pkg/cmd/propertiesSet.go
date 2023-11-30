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
func NewPropertiesSetCommand(factory Factory, propertiesCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {

	cmd := new(PropertiesSetCommand)
	err := cmd.init(factory, propertiesCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesSetCommand) GetName() string {
	return COMMAND_NAME_PROPERTIES_SET
}

func (cmd *PropertiesSetCommand) GetCobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *PropertiesSetCommand) GetValues() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesSetCommand) init(factory Factory, propertiesCommand GalasaCommand, rootCommand GalasaCommand) error {

	var err error = nil
	propertiesSetCmdValues := &PropertiesSetCmdValues{}

	propertiesSetCobraCmd := &cobra.Command{
		Use:   "set",
		Short: "Set the details of properties in a namespace.",
		Long: "Set the details of a property in a namespace. " +
			"If the property does not exist, a new property is created, otherwise the value for that property will be updated.",
		Args:    cobra.NoArgs,
		Aliases: []string{"properties set"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePropertiesSet(factory, cmd, args, propertiesSetCmdValues, propertiesCommand.GetValues().(*PropertiesCmdValues), rootCommand.GetValues().(*RootCmdValues))
		},
	}

	propertiesSetCobraCmd.PersistentFlags().StringVar(&propertiesSetCmdValues.propertyValue, "value", "", "the value of the property you want to create")

	propertiesSetCobraCmd.MarkPersistentFlagRequired("value")

	propertiesCommand.GetCobraCommand().AddCommand(propertiesSetCobraCmd)

	// The name & namespace properties are mandatory for set.
	addNamespaceProperty(propertiesSetCobraCmd, true, propertiesCommand.GetValues().(*PropertiesCmdValues))
	addNameProperty(propertiesSetCobraCmd, true, propertiesCommand.GetValues().(*PropertiesCmdValues))

	cmd.values = propertiesSetCmdValues
	cmd.cobraCommand = propertiesSetCobraCmd

	return err
}

func executePropertiesSet(
	factory Factory,
	cmd *cobra.Command,
	args []string,
	propertiesSetCmdValues *PropertiesSetCmdValues,
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

		galasaHome, err := utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
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
				err = properties.SetProperty(
					propertiesCmdValues.namespace,
					propertiesCmdValues.propertyName,
					propertiesSetCmdValues.propertyValue, apiClient)
			}
		}
	}
	return err
}
