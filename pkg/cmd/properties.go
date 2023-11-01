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
	ecosystemBootstrap string
	namespace          string
	propertyName       string
)

func init() {

}

func createPropertiesCmd(parentCmd *cobra.Command) (*cobra.Command, error) {
	var err error = nil

	propertiesCmd := &cobra.Command{
		Use:   "properties",
		Short: "Manages properties in an ecosystem",
		Long:  "Allows interaction with the CPS to create, query and maintain properties in Galasa Ecosystem",
	}

	propertiesCmd.PersistentFlags().StringVarP(&ecosystemBootstrap, "bootstrap", "b", "",
		"Bootstrap URL. Should start with 'http://' or 'file://'. "+
			"If it starts with neither, it is assumed to be a fully-qualified path. "+
			"If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. "+
			"Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties")

	addNamespaceProperty(propertiesCmd)

	parentCmd.AddCommand(propertiesCmd)

	err = createPropertiesCmdChildren(propertiesCmd)

	return propertiesCmd, err
}

func createPropertiesCmdChildren(propertiesCmd *cobra.Command) error {
	_, err := createPropertiesGetCmd(propertiesCmd)
	if err == nil {
		_, err = createPropertiesSetCmd(propertiesCmd)
	}
	if err == nil {
		_, err = createPropertiesDeleteCmd(propertiesCmd)
	}
	return err
}

func addNamespaceProperty(cmd *cobra.Command) {

	flagName := "namespace"
	cmd.PersistentFlags().StringVarP(&namespace, flagName, "s", "",
		"Namespace. A mandatory flag that describes the container for a collection of properties.")

	cmd.MarkPersistentFlagRequired(flagName)
}

// Some sub-commands need a name field to be mandatory, some don't.
func addNameProperty(cmd *cobra.Command, isMandatory bool) {
	flagName := "name"
	var description string
	if isMandatory {
		description = "A mandatory field indicating the name of a property in the namespace."
	} else {
		description = "An optional field indicating the name of a property in the namespace."
	}

	cmd.PersistentFlags().StringVarP(&propertyName, flagName, "n", "", description)

	if isMandatory {
		cmd.MarkPersistentFlagRequired(flagName)
	}
}
