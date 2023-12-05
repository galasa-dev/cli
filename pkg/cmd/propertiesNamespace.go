/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"github.com/spf13/cobra"
)

//Objective: Allow user to do this:
//	properties namespaces get
//  And then display all namespaces in the cps or returns empty

type PropertiesNamespaceCommand struct {
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewPropertiesNamespaceCommand(propertiesCommand GalasaCommand, rootCmd GalasaCommand) (GalasaCommand, error) {
	cmd := new(PropertiesNamespaceCommand)

	err := cmd.init(propertiesCommand, rootCmd)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesNamespaceCommand) Name() string {
	return COMMAND_NAME_PROPERTIES_NAMESPACE
}

func (cmd *PropertiesNamespaceCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *PropertiesNamespaceCommand) Values() interface{} {
	// There are no values.
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesNamespaceCommand) init(propertiesCommand GalasaCommand, rootCmd GalasaCommand) error {
	var err error
	cmd.cobraCommand, err = cmd.createPropertiesNamespaceCobraCmd(propertiesCommand, rootCmd)
	return err
}

func (cmd *PropertiesNamespaceCommand) createPropertiesNamespaceCobraCmd(
	propertiesCommand GalasaCommand,
	rootCmd GalasaCommand,
	) (*cobra.Command, error) {

	var err error
	propertiesNamespaceCmd := &cobra.Command{
		Use:   "namespaces",
		Short: "Queries namespaces in an ecosystem",
		Long:  "Allows interaction with the CPS to query namespaces in Galasa Ecosystem",
		Args:  cobra.NoArgs,
	}

	propertiesCommand.CobraCommand().AddCommand(propertiesNamespaceCmd)

	return propertiesNamespaceCmd, err
}
