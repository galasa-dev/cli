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
//    runs delete --name 12345
// And then show the results in a human-readable form.

// Variables set by cobra's command-line parsing.
type RunsDeleteCmdValues struct {
	runName string
}

type RunsDeleteCommand struct {
	values       *RunsDeleteCmdValues
	cobraCommand *cobra.Command
}

func NewRunsDeleteCommand(factory spi.Factory, runsCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) (spi.GalasaCommand, error) {
	cmd := new(RunsDeleteCommand)
	err := cmd.init(factory, runsCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsDeleteCommand) Name() string {
	return COMMAND_NAME_RUNS_DELETE
}

func (cmd *RunsDeleteCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsDeleteCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *RunsDeleteCommand) init(factory spi.Factory, runsCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) error {
	var err error
	cmd.values = &RunsDeleteCmdValues{}
	cmd.cobraCommand, err = cmd.createCobraCommand(factory, runsCommand, rootCommand.Values().(*RootCmdValues))
	return err
}

func (cmd *RunsDeleteCommand) createCobraCommand(
	factory spi.Factory,
	runsCommand spi.GalasaCommand,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {

	var err error
	runsCmdValues := runsCommand.Values().(*RunsCmdValues)

	runsDeleteCobraCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a named test run.",
		Long:    "Delete a named test run.",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs delete"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.executeRunsDelete(factory, runsCmdValues, rootCmdValues)
		},
	}

	runsDeleteCobraCmd.Flags().StringVar(&cmd.values.runName, "name", "", "the name of the test run we want to delete.")

	runsDeleteCobraCmd.MarkFlagRequired("name")

	runsCommand.CobraCommand().AddCommand(runsDeleteCobraCmd)

	return runsDeleteCobraCmd, err
}

func (cmd *RunsDeleteCommand) executeRunsDelete(
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

		log.Println("Galasa CLI - Delete runs about to execute")

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

				byteReader := factory.GetByteReader()

				if err == nil {
					// Call to process the command in a unit-testable way.
					err = runs.RunsDelete(
						cmd.values.runName,
						console,
						apiServerUrl,
						apiClient,
						timeService,
						byteReader,
					)
				}
			}
		}
	}

	log.Printf("executeRunsDelete returning %v", err)
	return err
}
