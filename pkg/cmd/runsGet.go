/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/runs"
	"github.com/galasa-dev/cli/pkg/spi"
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

type RunsGetCommand struct {
	values       *RunsGetCmdValues
	cobraCommand *cobra.Command
}

func NewRunsGetCommand(factory spi.Factory, runsCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) (spi.GalasaCommand, error) {
	cmd := new(RunsGetCommand)
	err := cmd.init(factory, runsCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsGetCommand) Name() string {
	return COMMAND_NAME_RUNS_GET
}

func (cmd *RunsGetCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsGetCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *RunsGetCommand) init(factory spi.Factory, runsCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) error {
	var err error
	cmd.values = &RunsGetCmdValues{}
	cmd.cobraCommand, err = cmd.createCobraCommand(factory, runsCommand, rootCommand.Values().(*RootCmdValues))
	return err
}

func (cmd *RunsGetCommand) createCobraCommand(
	factory spi.Factory,
	runsCommand spi.GalasaCommand,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {

	var err error
	runsCmdValues := runsCommand.Values().(*RunsCmdValues)

	runsGetCobraCmd := &cobra.Command{
		Use:     "get",
		Short:   "Get the details of a test runname which ran or is running.",
		Long:    "Get the details of a test runname which ran or is running, displaying the results to the caller.",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs get"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.executeRunsGet(factory, runsCmdValues, rootCmdValues)
		},
	}

	units := runs.GetTimeUnitsForErrorMessage()
	formatters := runs.GetFormatterNamesString(runs.CreateFormatters())
	runsGetCobraCmd.PersistentFlags().StringVar(&cmd.values.runName, "name", "", "the name of the test run we want information about."+
		" Cannot be used in conjunction with --requestor, --result or --active flags")
	runsGetCobraCmd.PersistentFlags().StringVar(&cmd.values.age, "age", "", "the age of the test run(s) we want information about. Supported formats are: 'FROM' or 'FROM:TO', where FROM and TO are each ages,"+
		" made up of an integer and a time-unit qualifier. Supported time-units are "+units+". If missing, the TO part is defaulted to '0h'. Examples: '--age 1d',"+
		" '--age 6h:1h' (list test runs which happened from 6 hours ago to 1 hour ago)."+
		" The TO part must be a smaller time-span than the FROM part.")
	runsGetCobraCmd.PersistentFlags().StringVar(&cmd.values.outputFormatString, "format", "summary", "output format for the data returned. Supported formats are: "+formatters+".")
	runsGetCobraCmd.PersistentFlags().StringVar(&cmd.values.requestor, "requestor", "", "the requestor of the test run we want information about."+
		" Cannot be used in conjunction with --name flag.")
	runsGetCobraCmd.PersistentFlags().StringVar(&cmd.values.result, "result", "", "A filter on the test runs we want information about. Optional. Default is to display test runs with any result. Case insensitive. Value can be a single value or a comma-separated list. For example \"--result Failed,Ignored,EnvFail\"."+
		" Cannot be used in conjunction with --name or --active flag.")
	runsGetCobraCmd.PersistentFlags().BoolVar(&cmd.values.isActiveRuns, "active", false, "parameter to retrieve runs that have not finished yet."+
		" Cannot be used in conjunction with --name or --result flag.")

	runsGetCobraCmd.MarkFlagsMutuallyExclusive("name", "requestor")
	runsGetCobraCmd.MarkFlagsMutuallyExclusive("name", "result")
	runsGetCobraCmd.MarkFlagsMutuallyExclusive("name", "active")
	runsGetCobraCmd.MarkFlagsMutuallyExclusive("result", "active")

	runsCommand.CobraCommand().AddCommand(runsGetCobraCmd)

	return runsGetCobraCmd, err
}

func (cmd *RunsGetCommand) executeRunsGet(
	factory spi.Factory,
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

		var galasaHome spi.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Read the bootstrap properties.
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, runsCmdValues.bootstrap, urlService)
			if err == nil {

				var console = factory.GetStdOutConsole()
				timeService := factory.GetTimeService()

				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API server is at '%s'\n", apiServerUrl)

				authenticator := factory.GetAuthenticator(
					apiServerUrl,
					galasaHome,
				)

				var apiClient *galasaapi.APIClient
				apiClient, err = authenticator.GetAuthenticatedAPIClient()

				if err == nil {
					// Call to process the command in a unit-testable way.
					err = runs.GetRuns(
						cmd.values.runName,
						cmd.values.age,
						cmd.values.requestor,
						cmd.values.result,
						cmd.values.isActiveRuns,
						cmd.values.outputFormatString,
						timeService,
						console,
						apiServerUrl,
						apiClient,
					)
				}
			}
		}
	}

	log.Printf("executeRunsGet returning %v", err)
	return err
}
