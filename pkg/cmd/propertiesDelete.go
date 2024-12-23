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

type PropertiesDeleteCommand struct {
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewPropertiesDeleteCommand(factory spi.Factory, propertiesCommand spi.GalasaCommand, commsCommand spi.GalasaCommand) (spi.GalasaCommand, error) {
	cmd := new(PropertiesDeleteCommand)

	err := cmd.init(factory, propertiesCommand, commsCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesDeleteCommand) Name() string {
	return COMMAND_NAME_PROPERTIES_DELETE
}

func (cmd *PropertiesDeleteCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *PropertiesDeleteCommand) Values() interface{} {
	// There are no values.
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesDeleteCommand) init(factory spi.Factory, propertiesCommand spi.GalasaCommand, commsCmd spi.GalasaCommand) error {
	var err error
	cmd.cobraCommand, err = cmd.createPropertiesDeleteCobraCmd(factory, propertiesCommand, commsCmd)
	return err
}

//Objective: Allow user to do this:
//	properties delete --namespace "framework" --name "hello"
//  And then display a successful message or error

func (cmd *PropertiesDeleteCommand) createPropertiesDeleteCobraCmd(
	factory spi.Factory,
	propertiesCommand spi.GalasaCommand,
	commsCmd spi.GalasaCommand) (*cobra.Command, error) {

	var err error
	propertiesCmdValues := propertiesCommand.Values().(*PropertiesCmdValues)
	commsCmdValues := commsCmd.Values().(*CommsCmdValues)

	propertiesDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a property in a namespace.",
		Long:  "Delete a property and its value in a namespace",
		Args:  cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			executionFunc := func() error {
				return cmd.executePropertiesDelete(factory, propertiesCmdValues, commsCmdValues)
			}
			return executeCommandWithRetries(factory, commsCmdValues, executionFunc)
		},
		Aliases: []string{"properties delete"},
	}

	propertiesCommand.CobraCommand().AddCommand(propertiesDeleteCmd)

	addPropertyNameFlag(propertiesDeleteCmd, true, propertiesCmdValues)
	addNamespaceFlag(propertiesDeleteCmd, true, propertiesCmdValues)

	// There are no sub-commands to add to the tree.

	return propertiesDeleteCmd, err
}

func (cmd *PropertiesDeleteCommand) executePropertiesDelete(factory spi.Factory, propertiesCmdValues *PropertiesCmdValues, commsCmdValues *CommsCmdValues) error {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	commsCmdValues.isCapturingLogs = true

	log.Println("Galasa CLI - Delete ecosystem properties")

	// Get the ability to query environment variables.
	env := factory.GetEnvironment()

	var galasaHome spi.GalasaHome
	galasaHome, err = utils.NewGalasaHome(fileSystem, env, commsCmdValues.CmdParamGalasaHomePath)
	if err == nil {

		// Read the bootstrap properties.
		var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
		var bootstrapData *api.BootstrapData
		bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, commsCmdValues.bootstrap, urlService)
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
				err = properties.DeleteProperty(propertiesCmdValues.namespace, propertiesCmdValues.propertyName, apiClient)
			}
		}
	}
	return err
}
