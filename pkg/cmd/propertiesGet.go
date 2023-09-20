/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"log"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/properties"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

//Objective: Allow user to do this:
//	properties get --namespace "framework" --prefix "pro" --sufix "ty"
//  And then display all properties filtered by either prefix, suffix or both, or empty if not found
//OR
//	properties get --namespace "framework" --name "hello"
//  And then display value of specified property or return empty if not found

var (
	propertiesGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get the details of properties in a namespace.",
		Long:  "Get the details of all properties in a namespace, filtered with flags if present",
		Args:  cobra.NoArgs,
		Run:   executepropertiesGet,
	}

	// Variables set by cobra's command-line parsing.
	propertiesPrefix       string
	propertiesSuffix       string
	propertiesOutputFormat string
)

func init() {
	formatters := properties.GetFormatterNamesString(properties.CreateFormatters())
	propertiesGetCmd.PersistentFlags().StringVar(&propertiesPrefix, "prefix", "", "the name of properties from a specified namespace with the provided prefix")
	propertiesGetCmd.PersistentFlags().StringVar(&propertiesSuffix, "suffix", "", "the name of properties from a specified namespace with the provided suffix")
	propertiesGetCmd.PersistentFlags().StringVar(&propertiesOutputFormat, "format", "summary", "output format for the data returned. Supported formats are: "+formatters+".")
	parentCommand := propertiesCmd
	parentCommand.AddCommand(propertiesGetCmd)
}

func executepropertiesGet(cmd *cobra.Command, args []string) {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, logFileName)
	if err != nil {
		panic(err)
	}
	isCapturingLogs = true

	log.Println("Galasa CLI - Get ecosystem properties")

	if propertyName != "" && (propertiesPrefix != "" || propertiesSuffix != "") {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_PROPERTIES_FLAG_COMBINATION)
	} else {
		// Get the ability to query environment variables.
		env := utils.NewEnvironment()

		galasaHome, err := utils.NewGalasaHome(fileSystem, env, CmdParamGalasaHomePath)
		if err != nil {
			panic(err)
		}

		// Read the bootstrap properties.
		var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
		var bootstrapData *api.BootstrapData
		bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, bootstrap, urlService)
		if err != nil {
			panic(err)
		}

		var console = utils.NewRealConsole()

		apiServerUrl := bootstrapData.ApiServerURL
		log.Printf("The API sever is at '%s'\n", apiServerUrl)

		// Call to process the command in a unit-testable way.
		err = properties.GetProperties(namespace, propertyName, propertiesPrefix, propertiesSuffix, apiServerUrl, propertiesOutputFormat, console)
		if err != nil {
			panic(err)
		}
	}

}
