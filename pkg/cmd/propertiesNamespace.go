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

var (
	propertiesNamespaceCmd = &cobra.Command{
		Use:   "namespaces",
		Short: "Queries namespaces in an ecosystem",
		Long:  "Allows interaction with the CPS to query namespaces in Galasa Ecosystem",
		Args:  cobra.NoArgs,
	}
)

func init() {

	parentCommand := propertiesCmd
	parentCommand.AddCommand(propertiesNamespaceCmd)
}
