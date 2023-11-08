/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/files"
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

func createPropertiesGetCmd(parentCmd *cobra.Command, propertiesCmdValues *PropertiesCmdValues, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	var err error = nil

	propertiesGetCmdValues := &PropertiesGetCmdValues{}

	propertiesGetCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get the details of properties in a namespace.",
		Long:    "Get the details of all properties in a namespace, filtered with flags if present",
		Args:    cobra.NoArgs,
		Aliases: []string{"properties get"},
		Run: func(cmd *cobra.Command, args []string) {
			executePropertiesGet(cmd, args, propertiesGetCmdValues, propertiesCmdValues, rootCmdValues)
		},
	}

	formatters := properties.GetFormatterNamesString(properties.CreateFormatters())
	propertiesGetCmd.PersistentFlags().StringVar(&propertiesGetCmdValues.propertiesPrefix, "prefix", "",
		"Prefix to match against the start of the property name within the namespace."+
			" Optional. Cannot be used in conjunction with the '--name' option.")
	propertiesGetCmd.PersistentFlags().StringVar(&propertiesGetCmdValues.propertiesSuffix, "suffix", "",
		"Suffix to match against the end of the property name within the namespace."+
			" Optional. Cannot be used in conjunction with the '--name' option.")
	propertiesGetCmd.PersistentFlags().StringVar(&propertiesGetCmdValues.propertiesInfix, "infix", "",
		"Infix(es) that could be part of the property name within the namespace."+
			" Multiple infixes can be supplied as a comma-separated list. "+
			" Optional. Cannot be used in conjunction with the '--name' option.")
	propertiesGetCmd.PersistentFlags().StringVar(&propertiesGetCmdValues.propertiesOutputFormat, "format", "summary",
		"output format for the data returned. Supported formats are: "+formatters+".")

	// The namespace property is mandatory for get.
	addNamespaceProperty(propertiesGetCmd, true, propertiesCmdValues)
	addNameProperty(propertiesGetCmd, false, propertiesCmdValues)

	// Name field cannot be used in conjunction wiht the prefix, suffix or infix commands.
	propertiesGetCmd.MarkFlagsMutuallyExclusive("name", "prefix")
	propertiesGetCmd.MarkFlagsMutuallyExclusive("name", "suffix")
	propertiesGetCmd.MarkFlagsMutuallyExclusive("name", "infix")

	parentCmd.AddCommand(propertiesGetCmd)

	// There are no sub-command children to add to the command tree.

	return propertiesGetCmd, err
}

func executePropertiesGet(
	cmd *cobra.Command,
	args []string,
	propertiesGetCmdValues *PropertiesGetCmdValues,
	propertiesCmdValues *PropertiesCmdValues,
	rootCmdValues *RootCmdValues,
) {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err != nil {
		panic(err)
	}
	rootCmdValues.isCapturingLogs = true

	log.Println("Galasa CLI - Get ecosystem properties")

	// Get the ability to query environment variables.
	env := utils.NewEnvironment()

	galasaHome, err := utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
	if err != nil {
		panic(err)
	}

	// Read the bootstrap properties.
	var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
	var bootstrapData *api.BootstrapData
	bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, propertiesCmdValues.ecosystemBootstrap, urlService)
	if err != nil {
		panic(err)
	}

	var console = utils.NewRealConsole()

	apiServerUrl := bootstrapData.ApiServerURL
	log.Printf("The API server is at '%s'\n", apiServerUrl)

	// Call to process the command in a unit-testable way.
	err = properties.GetProperties(
		propertiesCmdValues.namespace,
		propertiesCmdValues.propertyName,
		propertiesGetCmdValues.propertiesPrefix,
		propertiesGetCmdValues.propertiesSuffix,
		propertiesGetCmdValues.propertiesInfix,
		apiServerUrl,
		propertiesGetCmdValues.propertiesOutputFormat,
		console)
	if err != nil {
		panic(err)
	}
}
