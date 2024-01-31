/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"fmt"
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/auth"
	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/images"
	"github.com/galasa-dev/cli/pkg/runs"
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
func NewRunsDownloadCommand(factory Factory, runsCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(RunsDownloadCommand)
	err := cmd.init(factory, runsCommand, rootCommand)
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

func (cmd *RunsDownloadCommand) init(factory Factory, runsCommand GalasaCommand, rootCommand GalasaCommand) error {
	var err error
	cmd.values = &RunsDownloadCmdValues{}
	cmd.cobraCommand, err = cmd.createRunsDownloadCobraCmd(factory,
		runsCommand,
		rootCommand.Values().(*RootCmdValues),
	)
	return err
}

func (cmd *RunsDownloadCommand) createRunsDownloadCobraCmd(
	factory Factory,
	runsCommand GalasaCommand,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {

	var err error = nil
	runsCmdValues := runsCommand.Values().(*RunsCmdValues)

	runsDownloadCobraCmd := &cobra.Command{
		Use:     "download",
		Short:   "Download the artifacts of a test run which ran.",
		Long:    "Download the artifacts of a test run which ran and store them in a directory within the current working directory",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs download"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.executeRunsDownload(factory, runsCmdValues, rootCmdValues)
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
	factory Factory,
	runsCmdValues *RunsCmdValues,
	rootCmdValues *RootCmdValues,
) error {

	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err == nil {

		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Download artifacts for a run")

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
				timeService := factory.GetTimeService()

				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API server is at '%s'\n", apiServerUrl)

				apiClient := auth.GetAuthenticatedAPIClient(apiServerUrl, fileSystem, galasaHome, timeService, env)

				// Call to process the command in a unit-testable way.
				err = runs.DownloadArtifacts(
					cmd.values.runNameDownload,
					cmd.values.runForceDownload,
					fileSystem,
					timeService,
					console,
					apiClient,
					cmd.values.runDownloadTargetFolder,
				)

				if err == nil {
					// No error, so try to expand the files.
					embeddedFileSystem := embedded.GetReadOnlyFileSystem()
					renderer := images.NewImageRenderer(embeddedFileSystem)
					expander := images.NewImageExpander(fileSystem, renderer)

					err = expander.ExpandImages(cmd.values.runDownloadTargetFolder)
					if err == nil {

						// Write out a status string to the console about how many files were rendered.
						count := expander.GetExpandedImageFileCount()
						message := fmt.Sprintf(galasaErrors.GALASA_INFO_RENDERED_IMAGE_COUNT.Template, count)
						console.WriteString(message)
					}
				}
			}
		}
	}
	return err
}
