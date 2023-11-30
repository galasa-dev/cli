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
//	properties get --namespace "framework" --prefix "pro" --sufix "ty"
//  And then display all properties filtered by either prefix, suffix or both, or empty if not found
//OR
//	properties get --namespace "framework" --name "hello"
//  And then display value of specified property or return empty if not found

// Variables set by cobra's command-line parsing.
type PropertiesGetCmdValues struct {
	propertiesPrefix       string
	propertiesSuffix       string
	propertiesInfix        string
	propertiesOutputFormat string
}

type PropertiesGetComamnd struct {
	values       *PropertiesGetCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewPropertiesGetCommand(factory Factory, propertiesCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {

	cmd := new(PropertiesGetComamnd)
	err := cmd.init(factory, propertiesCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesGetComamnd) Name() string {
	return COMMAND_NAME_PROPERTIES_GET
}

func (cmd *PropertiesGetComamnd) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *PropertiesGetComamnd) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *PropertiesGetComamnd) init(factory Factory, propertiesCommand GalasaCommand, rootCommand GalasaCommand) error {

	var err error = nil

	cmd.values = &PropertiesGetCmdValues{}
	cmd.cobraCommand = cmd.createCobraCommand(factory, cmd.values, propertiesCommand, rootCommand)

	return err
}

func (cmd *PropertiesGetComamnd) createCobraCommand(factory Factory, propertiesGetCmdValues *PropertiesGetCmdValues, propertiesCommand GalasaCommand, rootCommand GalasaCommand) *cobra.Command {

	propertiesGetCobraCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get the details of properties in a namespace.",
		Long:    "Get the details of all properties in a namespace, filtered with flags if present",
		Args:    cobra.NoArgs,
		Aliases: []string{"properties get"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePropertiesGet(factory, cmd, args, propertiesGetCmdValues,
				propertiesCommand.Values().(*PropertiesCmdValues), rootCommand.Values().(*RootCmdValues))
		},
	}

	formatters := properties.GetFormatterNamesString(properties.CreateFormatters())
	propertiesGetCobraCmd.PersistentFlags().StringVar(&propertiesGetCmdValues.propertiesPrefix, "prefix", "",
		"Prefix to match against the start of the property name within the namespace."+
			" Optional. Cannot be used in conjunction with the '--name' option.")
	propertiesGetCobraCmd.PersistentFlags().StringVar(&propertiesGetCmdValues.propertiesSuffix, "suffix", "",
		"Suffix to match against the end of the property name within the namespace."+
			" Optional. Cannot be used in conjunction with the '--name' option.")
	propertiesGetCobraCmd.PersistentFlags().StringVar(&propertiesGetCmdValues.propertiesInfix, "infix", "",
		"Infix(es) that could be part of the property name within the namespace."+
			" Multiple infixes can be supplied as a comma-separated list. "+
			" Optional. Cannot be used in conjunction with the '--name' option.")
	propertiesGetCobraCmd.PersistentFlags().StringVar(&propertiesGetCmdValues.propertiesOutputFormat, "format", "summary",
		"output format for the data returned. Supported formats are: "+formatters+".")

	// The namespace property is mandatory for get.
	addNamespaceFlag(propertiesGetCobraCmd, true, propertiesCommand.Values().(*PropertiesCmdValues))
	addPropertyNameFlag(propertiesGetCobraCmd, false, propertiesCommand.Values().(*PropertiesCmdValues))

	// Name field cannot be used in conjunction wiht the prefix, suffix or infix commands.
	propertiesGetCobraCmd.MarkFlagsMutuallyExclusive("name", "prefix")
	propertiesGetCobraCmd.MarkFlagsMutuallyExclusive("name", "suffix")
	propertiesGetCobraCmd.MarkFlagsMutuallyExclusive("name", "infix")

	propertiesCommand.CobraCommand().AddCommand(propertiesGetCobraCmd)

	return propertiesGetCobraCmd
}

func executePropertiesGet(
	factory Factory,
	cmd *cobra.Command,
	args []string,
	propertiesGetCmdValues *PropertiesGetCmdValues,
	propertiesCmdValues *PropertiesCmdValues,
	rootCmdValues *RootCmdValues,
) error {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err == nil {

		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Get ecosystem properties")

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
				err = properties.GetProperties(
					propertiesCmdValues.namespace,
					propertiesCmdValues.propertyName,
					propertiesGetCmdValues.propertiesPrefix,
					propertiesGetCmdValues.propertiesSuffix,
					propertiesGetCmdValues.propertiesInfix,
					apiClient,
					propertiesGetCmdValues.propertiesOutputFormat,
					console,
				)
			}
		}
	}
	return err
}
