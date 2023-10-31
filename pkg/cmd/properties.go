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
		Long:  "Allows interaction with the CPS to create, query and maintain properties in Galasa Ecosystem",
	}
	ecosystemBootstrap string
	namespace          string
	propertyName       string
)

func init() {
	cmd := propertiesCmd
	parentCmd := RootCmd

	cmd.PersistentFlags().StringVarP(&ecosystemBootstrap, "bootstrap", "b", "",
		"Bootstrap URL. Should start with 'http://' or 'file://'. "+
			"If it starts with neither, it is assumed to be a fully-qualified path. "+
			"If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. "+
			"Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties")

	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "s", "",
		"Namespace. A mandatory flag that describes the container for a collection of properties. "+
			"It has no default value.")
	cmd.MarkPersistentFlagRequired("namespace")

	cmd.PersistentFlags().StringVarP(&propertyName, "name", "n", "",
		"Name of a property in the namespace. "+
			"It has no default value.")

	parentCmd.AddCommand(propertiesCmd)
}
