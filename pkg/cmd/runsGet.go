/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"log"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/runs"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Objective: Allow the user to do this:
//    run get --runname 12345
// And then show the results in a human-readable form.

var (
	runsGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get the details of a test runname which ran or is running.",
		Long:  "Get the details of a test runname which ran or is running, displaying the results to the caller",
		Args:  cobra.NoArgs,
		Run:   executeRunsGet,
	}

	// Variables set by cobra's command-line parsing.
	runName            string
	age                string
	outputFormatString string
)

func init() {
	runsGetCmd.PersistentFlags().StringVar(&runName, "name", "", "the name of the test run we want information about")
	runsGetCmd.PersistentFlags().StringVar(&age, "age", "", "the age of the test run(s) we want information about. Supported formats are: --age FROM, --age FROM:TO. Supported units are days, weeks and hours. Examples: --age 1d, --age 6h:1h")
	runsGetCmd.PersistentFlags().StringVar(&outputFormatString, "format", "summary", "output format for the data returned. Supported formats are: summary")

	parentCommand := runsCmd
	parentCommand.AddCommand(runsGetCmd)
}

func executeRunsGet(cmd *cobra.Command, args []string) {

	var err error

	if (runName == "") && (age == "") {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_NO_RUNNAME_OR_AGE_SPECIFIED)
		panic(err)
	}
	// Operations on the file system will all be relative to the current folder.
	fileSystem := utils.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, logFileName)
	if err != nil {
		panic(err)
	}
	isCapturingLogs = true

	log.Println("Galasa CLI - Get info about a run")

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

	timeService := utils.NewRealTimeService()

	// Call to process the command in a unit-testable way.
	err = runs.GetRuns(runName, age, outputFormatString, timeService, console, apiServerUrl)
	if err != nil {
		panic(err)
	}
}
