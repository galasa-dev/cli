/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"log"

	"github.com/galasa.dev/cli/pkg/api"
	"github.com/galasa.dev/cli/pkg/files"
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
	requestor          string
	result             string
)

func init() {
	units := runs.GetTimeUnitsForErrorMessage(runs.CreateTimeUnits())
	formatters := runs.GetFormatterNamesString(runs.CreateFormatters())
	runsGetCmd.PersistentFlags().StringVar(&runName, "name", "", "the name of the test run we want information about")
	runsGetCmd.PersistentFlags().StringVar(&age, "age", "", "the age of the test run(s) we want information about. Supported formats are: 'FROM' or 'FROM:TO', where FROM and TO are each ages,"+
		" made up of an integer and a time-unit qualifier. Supported time-units are "+units+". If missing, the TO part is defaulted to '0h'. Examples: '--age 1d' , '--age 6h:1h' ")
	runsGetCmd.PersistentFlags().StringVar(&outputFormatString, "format", "summary", "output format for the data returned. Supported formats are: "+formatters+".")
	runsGetCmd.PersistentFlags().StringVar(&requestor, "requestor", "", "the requestor of the test run we want information about")
	runsGetCmd.PersistentFlags().StringVar(&result, "result", "", "A filter on the test runs we want information about. Optional. Default is to display test runs with any result. Case insensitive. Value can be a single value or a comma-separated list. For example \"--result Failed,Ignored,EnvFail\"")
	parentCommand := runsCmd
	parentCommand.AddCommand(runsGetCmd)
}

func executeRunsGet(cmd *cobra.Command, args []string) {

	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

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
	err = runs.GetRuns(runName, age, requestor, result, outputFormatString, timeService, console, apiServerUrl, env)
	if err != nil {
		panic(err)
	}
}
