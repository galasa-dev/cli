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
		Run:   executePropertiesGet,
	}

	// Variables set by cobra's command-line parsing.
	propertiesPrefix       string
	propertiesSuffix       string
	propertiesInfix        string
	propertiesOutputFormat string
)

func init() {
	formatters := properties.GetFormatterNamesString(properties.CreateFormatters())
	propertiesGetCmd.PersistentFlags().StringVar(&propertiesPrefix, "prefix", "", "returns properties from a specified namespace with the supplied prefix")
	propertiesGetCmd.PersistentFlags().StringVar(&propertiesSuffix, "suffix", "", "returns properties from a specified namespace with the supplied suffix")
	propertiesGetCmd.PersistentFlags().StringVar(&propertiesInfix, "infix", "", "returns properties from a specified namespace which contains the supplied infix or at least one of the supplied infixes if there are multiple")
	propertiesGetCmd.PersistentFlags().StringVar(&propertiesOutputFormat, "format", "summary", "output format for the data returned. Supported formats are: "+formatters+".")
	parentCommand := propertiesCmd
	parentCommand.AddCommand(propertiesGetCmd)
}

func executePropertiesGet(cmd *cobra.Command, args []string) {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, logFileName)
	if err != nil {
		panic(err)
	}
	isCapturingLogs = true

	log.Println("Galasa CLI - Get ecosystem properties")

	//Checks if --name has been provided with one or more of --prefix, --suffix, --infix as they are mutually exclusive
	if propertyName != "" && (propertiesPrefix != "" || propertiesSuffix != "" || propertiesInfix != "") {
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
		log.Printf("The API server is at '%s'\n", apiServerUrl)

		// Call to process the command in a unit-testable way.
		err = properties.GetProperties(namespace, propertyName, propertiesPrefix, propertiesSuffix, propertiesInfix, apiServerUrl, propertiesOutputFormat, console)
		if err != nil {
			panic(err)
		}
	}

}
