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
//    runs reset --name U123
// And then galasactl resets the run by requeuing it.

type RunsResetCommand struct {
	values       *RunsResetCmdValues
	cobraCommand *cobra.Command
}

type RunsResetCmdValues struct {
	runName string
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewRunsResetCommand(factory spi.Factory, runsCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) (spi.GalasaCommand, error) {
	cmd := new(RunsResetCommand)
	err := cmd.init(factory, runsCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsResetCommand) Name() string {
	return COMMAND_NAME_RUNS_RESET
}

func (cmd *RunsResetCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsResetCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsResetCommand) init(factory spi.Factory, runsCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) error {
	var err error
	cmd.values = &RunsResetCmdValues{}
	cmd.cobraCommand, err = cmd.createRunsResetCobraCmd(
		factory,
		runsCommand,
		rootCommand.Values().(*RootCmdValues),
	)
	return err
}

func (cmd *RunsResetCommand) createRunsResetCobraCmd(factory spi.Factory,
	runsCommand spi.GalasaCommand,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {

	var err error
	runsCmdValues := runsCommand.Values().(*RunsCmdValues)

	runsResetCmd := &cobra.Command{
		Use:     "reset",
		Short:   "reset an active run in the ecosystem",
		Long:    "Reset an active test run in the ecosystem if it is stuck or looping.",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs reset"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.executeReset(factory, runsCmdValues, rootCmdValues)
		},
	}

	runsResetCmd.PersistentFlags().StringVar(&cmd.values.runName, "name", "", "the name of the test run to reset")

	runsResetCmd.MarkPersistentFlagRequired("name")

	runsCommand.CobraCommand().AddCommand(runsResetCmd)

	return runsResetCmd, err
}

func (cmd *RunsResetCommand) executeReset(
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

		log.Println("Galasa CLI - Reset an active run by requeuing it.")

		// Get the ability to query environment variables.
		env := factory.GetEnvironment()

		var galasaHome spi.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Read the bootstrap properties
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, runsCmdValues.bootstrap, urlService)
			if err == nil {

				var console = factory.GetStdOutConsole()
				timeService := factory.GetTimeService()

				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API Server is at '%s'\n", apiServerUrl)

				var apiClient *galasaapi.APIClient
				authenticator := factory.GetAuthenticator(
					apiServerUrl,
					galasaHome,
				)
				apiClient, err = authenticator.GetAuthenticatedAPIClient()

				if err == nil {
					// Call to process command in unit-testable way.
					err = runs.ResetRun(
						cmd.values.runName,
						timeService,
						console,
						apiServerUrl,
						apiClient,
					)
				}
			}
		}
	}

	log.Printf("executeRunsReset returning %v\n", err)
	return err
}
