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
//    runs cancel --name U123
// And then galasactl cancels the run by abandoning it.

type RunsCancelCommand struct {
	values       *RunsCancelCmdValues
	cobraCommand *cobra.Command
}

type RunsCancelCmdValues struct {
	runName string
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewRunsCancelCommand(factory spi.Factory, runsCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) (spi.GalasaCommand, error) {
	cmd := new(RunsCancelCommand)
	err := cmd.init(factory, runsCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsCancelCommand) Name() string {
	return COMMAND_NAME_RUNS_CANCEL
}

func (cmd *RunsCancelCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsCancelCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsCancelCommand) init(factory spi.Factory, runsCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) error {
	var err error
	cmd.values = &RunsCancelCmdValues{}
	cmd.cobraCommand, err = cmd.createRunsCancelCobraCmd(
		factory,
		runsCommand,
		rootCommand.Values().(*RootCmdValues),
	)
	return err
}

func (cmd *RunsCancelCommand) createRunsCancelCobraCmd(factory spi.Factory,
	runsCommand spi.GalasaCommand,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {

	var err error
	runsCmdValues := runsCommand.Values().(*RunsCmdValues)

	runsCancelCmd := &cobra.Command{
		Use:     "cancel",
		Short:   "cancel an active run in the ecosystem",
		Long:    "Cancel an active test run in the ecosystem if it is stuck or looping.",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs cancel"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.executeCancel(factory, runsCmdValues, rootCmdValues)
		},
	}

	runsCancelCmd.PersistentFlags().StringVar(&cmd.values.runName, "name", "", "the name of the test run to cancel")

	runsCancelCmd.MarkPersistentFlagRequired("name")

	runsCommand.CobraCommand().AddCommand(runsCancelCmd)

	return runsCancelCmd, err
}

func (cmd *RunsCancelCommand) executeCancel(
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

		log.Println("Galasa CLI - Cancel an active run by abandoning it.")

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

				authenticator := factory.GetAuthenticator(
					apiServerUrl,
					galasaHome,
				)

				var apiClient *galasaapi.APIClient
				apiClient, err = authenticator.GetAuthenticatedAPIClient()

				if err == nil {
					// Call to process command in unit-testable way.
					err = runs.CancelRun(
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

	log.Printf("executeRunsCancel returning %v\n", err)
	return err
}
