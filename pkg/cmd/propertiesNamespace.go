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
func NewPropertiesNamespaceCommand(factory Factory, propertiesCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(PropertiesNamespaceCommand)

	err := cmd.init(factory, propertiesCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesNamespaceCommand) GetName() string {
	return COMMAND_NAME_PROPERTIES_NAMESPACE
}

func (cmd *PropertiesNamespaceCommand) GetCobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *PropertiesNamespaceCommand) GetValues() interface{} {
	// There are no values.
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *PropertiesNamespaceCommand) init(factory Factory, propertiesCommand GalasaCommand, rootCommand GalasaCommand) error {
	var err error
	cmd.cobraCommand, err = cmd.createPropertiesNamespaceCobraCmd(factory, propertiesCommand.GetCobraCommand(), propertiesCommand.GetValues().(*PropertiesCmdValues), rootCommand.GetValues().(*RootCmdValues))
	return err
}

func (cmd *PropertiesNamespaceCommand) createPropertiesNamespaceCobraCmd(factory Factory, propertiesCmd *cobra.Command, propertiesCmdValues *PropertiesCmdValues, rootCmdValues *RootCmdValues) (*cobra.Command, error) {

	var err error
	propertiesNamespaceCmd := &cobra.Command{
		Use:   "namespaces",
		Short: "Queries namespaces in an ecosystem",
		Long:  "Allows interaction with the CPS to query namespaces in Galasa Ecosystem",
		Args:  cobra.NoArgs,
	}

	propertiesCmd.AddCommand(propertiesNamespaceCmd)

	return propertiesNamespaceCmd, err
}
