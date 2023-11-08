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

func createPropertiesNamespaceCmd(factory Factory, propertiesCmd *cobra.Command, propertiesCmdValues *PropertiesCmdValues, rootCmdValues *RootCmdValues) (*cobra.Command, error) {

	propertiesNamespaceCmd := &cobra.Command{
		Use:   "namespaces",
		Short: "Queries namespaces in an ecosystem",
		Long:  "Allows interaction with the CPS to query namespaces in Galasa Ecosystem",
		Args:  cobra.NoArgs,
	}

	propertiesCmd.AddCommand(propertiesNamespaceCmd)

	err := createChildCommands(factory, propertiesNamespaceCmd, propertiesCmdValues, rootCmdValues)

	return propertiesNamespaceCmd, err
}

func createChildCommands(factory Factory, propertiesNamespaceCmd *cobra.Command, propertiesCmdValues *PropertiesCmdValues, rootCmdValues *RootCmdValues) error {
	var err error

	_, err = createPropertiesNamespaceGetCmd(factory, propertiesNamespaceCmd, propertiesCmdValues, rootCmdValues)

	return err
}
