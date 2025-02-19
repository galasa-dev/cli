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
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Objective: Allow the user to do this:
//    runs download --name U123 [--force]
// And then galasactl downloads the artifacts for the given run.

type RunsDownloadCommand struct {
	values       *RunsDownloadCmdValues
	cobraCommand *cobra.Command
}

// Variables set by cobra's command-line parsing.
type RunsDownloadCmdValues struct {
	runNameDownload         string
	runForceDownload        bool
	runDownloadTargetFolder string
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewRunsDownloadCommand(factory spi.Factory, runsCommand spi.GalasaCommand, commsFlagSet GalasaFlagSet) (spi.GalasaCommand, error) {
	cmd := new(RunsDownloadCommand)
	err := cmd.init(factory, runsCommand, commsFlagSet)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsDownloadCommand) Name() string {
	return COMMAND_NAME_RUNS_DOWNLOAD
}

func (cmd *RunsDownloadCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsDownloadCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *RunsDownloadCommand) init(factory spi.Factory, runsCommand spi.GalasaCommand, commsFlagSet GalasaFlagSet) error {
	var err error
	cmd.values = &RunsDownloadCmdValues{}
	cmd.cobraCommand, err = cmd.createRunsDownloadCobraCmd(factory,
		runsCommand,
		commsFlagSet.Values().(*CommsFlagSetValues),
	)
	return err
}

func (cmd *RunsDownloadCommand) createRunsDownloadCobraCmd(
	factory spi.Factory,
	runsCommand spi.GalasaCommand,
	commsFlagSetValues *CommsFlagSetValues,
) (*cobra.Command, error) {

	var err error

	runsDownloadCobraCmd := &cobra.Command{
		Use:     "download",
		Short:   "Download the artifacts of a test run which ran.",
		Long:    "Download the artifacts of a test run which ran and store them in a directory within the current working directory",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs download"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.executeRunsDownload(factory, commsFlagSetValues)
		},
	}

	runsDownloadCobraCmd.PersistentFlags().StringVar(&cmd.values.runNameDownload, "name", "", "the name of the test run we want information about")
	runsDownloadCobraCmd.PersistentFlags().BoolVar(&cmd.values.runForceDownload, "force", false, "force artifacts to be overwritten if they already exist")
	runsDownloadCobraCmd.MarkPersistentFlagRequired("name")
	runsDownloadCobraCmd.PersistentFlags().StringVar(&cmd.values.runDownloadTargetFolder, "destination", ".",
		"The folder we want to download test run artifacts into. Sub-folders will be created within this location",
	)

	runsCommand.CobraCommand().AddCommand(runsDownloadCobraCmd)

	return runsDownloadCobraCmd, err
}

func (cmd *RunsDownloadCommand) executeRunsDownload(
	factory spi.Factory,
	commsFlagSetValues *CommsFlagSetValues,
) error {

	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, commsFlagSetValues.logFileName)
	if err == nil {
		commsFlagSetValues.isCapturingLogs = true
	
		log.Println("Galasa CLI - Download artifacts for a run")
	
		// Get the ability to query environment variables.
		env := factory.GetEnvironment()
	
		var galasaHome spi.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, commsFlagSetValues.CmdParamGalasaHomePath)
		if err == nil {

			timeService := factory.GetTimeService()
			commsRetrier := api.NewCommsRetrier(commsFlagSetValues.maxRetries, commsFlagSetValues.retryBackoffSeconds, timeService)

			// Read the bootstrap properties, retrying if a rate limit has been exceeded
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			loadBootstrapWithRetriesFunc := func() error {
				bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, commsFlagSetValues.bootstrap, urlService)
				return err
			}

			err = commsRetrier.ExecuteCommandWithRateLimitRetries(loadBootstrapWithRetriesFunc)
			if err == nil {
	
				var console = factory.GetStdOutConsole()
	
				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API server is at '%s'\n", apiServerUrl)
	
				authenticator := factory.GetAuthenticator(
					apiServerUrl,
					galasaHome,
				)
	
				commsRetrier, err = api.NewCommsRetrierWithAPIClient(
					commsFlagSetValues.maxRetries,
					commsFlagSetValues.retryBackoffSeconds,
					timeService,
					authenticator,
				)
	
				if err == nil {
					// Call to process the command in a unit-testable way.
					err = runs.DownloadArtifacts(
						cmd.values.runNameDownload,
						cmd.values.runForceDownload,
						fileSystem,
						timeService,
						console,
						commsRetrier,
						cmd.values.runDownloadTargetFolder,
					)
				}
			}
		}
	}

	return err
}
