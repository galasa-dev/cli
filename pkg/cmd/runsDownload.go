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
//    runs download --name U123 [--force]
// And then galasactl downloads the artifacts for the given run.

// Variables set by cobra's command-line parsing.
type RunsDownloadCmdValues struct {
	runNameDownload         string
	runForceDownload        bool
	runDownloadTargetFolder string
}

func createRunsDownloadCmd(factory Factory, parentCmd *cobra.Command, runsCmdValues *RunsCmdValues, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	var err error = nil

	runsDownloadCmdValues := &RunsDownloadCmdValues{}

	runsDownloadCmd := &cobra.Command{
		Use:     "download",
		Short:   "Download the artifacts of a test run which ran.",
		Long:    "Download the artifacts of a test run which ran and store them in a directory within the current working directory",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs download"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeRunsDownload(factory, cmd, args, runsDownloadCmdValues, runsCmdValues, rootCmdValues)
		},
	}

	runsDownloadCmd.PersistentFlags().StringVar(&runsDownloadCmdValues.runNameDownload, "name", "", "the name of the test run we want information about")
	runsDownloadCmd.PersistentFlags().BoolVar(&runsDownloadCmdValues.runForceDownload, "force", false, "force artifacts to be overwritten if they already exist")
	runsDownloadCmd.MarkPersistentFlagRequired("name")
	runsDownloadCmd.PersistentFlags().StringVar(&runsDownloadCmdValues.runDownloadTargetFolder, "destination", ".",
		"The folder we want to download test run artifacts into. Sub-folders will be created within this location",
	)

	parentCmd.AddCommand(runsDownloadCmd)

	// There are no children commands of this command to add to the command tree.

	return runsDownloadCmd, err
}

func executeRunsDownload(
	factory Factory,
	cmd *cobra.Command,
	args []string,
	runsDownloadCmdValues *RunsDownloadCmdValues,
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

				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API server is at '%s'\n", apiServerUrl)

				timeService := utils.NewRealTimeService()
				// Call to process the command in a unit-testable way.
				err = runs.DownloadArtifacts(
					runsDownloadCmdValues.runNameDownload,
					runsDownloadCmdValues.runForceDownload,
					fileSystem,
					timeService,
					console,
					apiServerUrl,
					runsDownloadCmdValues.runDownloadTargetFolder,
				)
			}
		}
	}
	return err
}
