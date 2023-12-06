/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"github.com/spf13/cobra"
)

type PropertiesCmdValues struct {
	ecosystemBootstrap string
	namespace          string
	propertyName       string
}

type PropertiesCommand struct {
	values       *PropertiesCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewPropertiesCommand(rootCmd GalasaCommand) (GalasaCommand, error) {

	cmd := new(PropertiesCommand)
	err := cmd.init(rootCmd)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesCommand) Name() string {
	return COMMAND_NAME_PROPERTIES
}

func (cmd *PropertiesCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *PropertiesCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *PropertiesCommand) init(rootCmd GalasaCommand) error {

	var err error = nil

	cmd.values = &PropertiesCmdValues{}
	cmd.cobraCommand = cmd.createCobraCommand(rootCmd)

	return err
}

func (cmd *PropertiesCommand) createCobraCommand(
	rootCommand GalasaCommand,
	) *cobra.Command {
	propertiesCobraCmd := &cobra.Command{
		Use:   "properties",
		Short: "Manages properties in an ecosystem",
		Long:  "Allows interaction with the CPS to create, query and maintain properties in Galasa Ecosystem",
	}

	addBootstrapFlag(propertiesCobraCmd, &cmd.values.ecosystemBootstrap)

	rootCommand.CobraCommand().AddCommand(propertiesCobraCmd)

	return propertiesCobraCmd
}

func addNamespaceFlag(cmd *cobra.Command, isMandatory bool, propertiesCmdValues *PropertiesCmdValues) {

	flagName := "namespace"
	var description string
	if isMandatory {
		description = "A mandatory flag that describes the container for a collection of properties."
	} else {
		description = "An optional flag that describes the container for a collection of properties."
	}

	cmd.PersistentFlags().StringVarP(&propertiesCmdValues.namespace, flagName, "s", "", description)

	if isMandatory {
		cmd.MarkPersistentFlagRequired(flagName)
	}

}

// Some sub-commands need a name field to be mandatory, some don't.
func addPropertyNameFlag(cmd *cobra.Command, isMandatory bool, propertiesCmdValues *PropertiesCmdValues) {
	flagName := "name"
	var description string
	if isMandatory {
		description = "A mandatory field indicating the name of a property in the namespace."
	} else {
		description = "An optional field indicating the name of a property in the namespace."
	}

	cmd.PersistentFlags().StringVarP(&propertiesCmdValues.propertyName, flagName, "n", "", description)

	if isMandatory {
		cmd.MarkPersistentFlagRequired(flagName)
	}
}
