/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	propertiesCmd = &cobra.Command{
		Use:   "properties",
		Short: "Manages properties in an ecosystem",
		Long:  "Allows interaction with the CPS to Initiate, query and maintain properties in Galasa Ecosystem",
	}
	ecosystemBootstrap string
	namespace          string
)

func init() {
	cmd := propertiesCmd
	parentCmd := RootCmd

	cmd.PersistentFlags().StringVarP(&ecosystemBootstrap, "bootstrap", "b", "",
		"Bootstrap URL. Should start with 'http://' or 'file://'. "+
			"If it starts with neither, it is assumed to be a fully-qualified path. "+
			"If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. "+
			"Examples: http://galasa-cicsk8s.hursley.ibm.com/bootstrap , file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties")

	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "",
		"Namespace. A container for a collection of properties. "+
			"It has no default value.")
	cmd.MarkFlagRequired("namespace")

	parentCmd.AddCommand(runsCmd)
}
