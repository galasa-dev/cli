/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/runs"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Objective: Allow the user to do this:
//    run get --runname 12345
// And then show the results in a human-readable form.

// Variables set by cobra's command-line parsing.
type RunsGetCmdValues struct {
	runName            string
	age                string
	outputFormatString string
	requestor          string
	result             string
	isActiveRuns       bool
}

func createRunsGetCmd(factory Factory, parentCmd *cobra.Command, runsCmdValues *RunsCmdValues, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	var err error = nil

	runsGetCmdValues := &RunsGetCmdValues{}

	runsGetCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get the details of a test runname which ran or is running.",
		Long:    "Get the details of a test runname which ran or is running, displaying the results to the caller",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs get"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeRunsGet(factory, cmd, args, runsGetCmdValues, runsCmdValues, rootCmdValues)
		},
	}

	units := runs.GetTimeUnitsForErrorMessage()
	formatters := runs.GetFormatterNamesString(runs.CreateFormatters())
	runsGetCmd.PersistentFlags().StringVar(&runsGetCmdValues.runName, "name", "", "the name of the test run we want information about")
	runsGetCmd.PersistentFlags().StringVar(&runsGetCmdValues.age, "age", "", "the age of the test run(s) we want information about. Supported formats are: 'FROM' or 'FROM:TO', where FROM and TO are each ages,"+
		" made up of an integer and a time-unit qualifier. Supported time-units are "+units+". If missing, the TO part is defaulted to '0h'. Examples: '--age 1d',"+
		" '--age 6h:1h' (list test runs which happened from 6 hours ago to 1 hour ago)."+
		" The TO part must be a smaller time-span than the FROM part.")
	runsGetCmd.PersistentFlags().StringVar(&runsGetCmdValues.outputFormatString, "format", "summary", "output format for the data returned. Supported formats are: "+formatters+".")
	runsGetCmd.PersistentFlags().StringVar(&runsGetCmdValues.requestor, "requestor", "", "the requestor of the test run we want information about")
	runsGetCmd.PersistentFlags().StringVar(&runsGetCmdValues.result, "result", "", "A filter on the test runs we want information about. Optional. Default is to display test runs with any result. Case insensitive. Value can be a single value or a comma-separated list. For example \"--result Failed,Ignored,EnvFail\"")
	runsGetCmd.PersistentFlags().BoolVar(&runsGetCmdValues.isActiveRuns, "active", false, "parameter to retrieve runs that have not finished yet.")

	parentCmd.AddCommand(runsGetCmd)

	return runsGetCmd, err
}

func executeRunsGet(
	factory Factory,
	cmd *cobra.Command,
	args []string,
	runsGetCmdValues *RunsGetCmdValues,
	runsCmdValues *RunsCmdValues,
	rootCmdValues *RootCmdValues,
) error {

	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err == nil {
		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Get info about a run")

		// Get the ability to query environment variables.
		env := factory.GetEnvironment()

		var galasaHome utils.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Read the bootstrap properties.
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, runsCmdValues.bootstrap, urlService)
			if err == nil {

				var console = factory.GetStdOutConsole()

				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API server is at '%s'\n", apiServerUrl)

				timeService := utils.NewRealTimeService()

				// Call to process the command in a unit-testable way.
				err = runs.GetRuns(
					runsGetCmdValues.runName,
					runsGetCmdValues.age,
					runsGetCmdValues.requestor,
					runsGetCmdValues.result,
					runsGetCmdValues.isActiveRuns,
					runsGetCmdValues.outputFormatString,
					timeService,
					console,
					apiServerUrl,
				)
			}
		}
	}

	log.Printf("executeRunsGet returning %v", err)
	return err
}
