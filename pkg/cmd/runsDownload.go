/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
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
//    runs download --name U123 [--force]
// And then galasactl downloads the artifacts for the given run.

var (
	runsDownloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download the artifacts of a test run which ran.",
		Long:  "Download the artifacts of a test run which ran and store them in a directory within the current working directory",
		Args:  cobra.NoArgs,
		Run:   executeRunsDownload,
	}

	// Variables set by cobra's command-line parsing.
	runNameDownload         string
	runForceDownload        bool
	runDownloadTargetFolder string
)

func init() {
	runsDownloadCmd.PersistentFlags().StringVar(&runNameDownload, "name", "", "the name of the test run we want information about")
	runsDownloadCmd.PersistentFlags().BoolVar(&runForceDownload, "force", false, "force artifacts to be overwritten if they already exist")
	runsDownloadCmd.MarkPersistentFlagRequired("name")
	runsDownloadCmd.PersistentFlags().StringVar(&runDownloadTargetFolder, "destination", ".",
		"The folder we want to download test run artifacts into. Sub-folders will be created within this location",
	)

	parentCommand := runsCmd
	parentCommand.AddCommand(runsDownloadCmd)
}

func executeRunsDownload(cmd *cobra.Command, args []string) {

	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, logFileName)
	if err != nil {
		panic(err)
	}
	isCapturingLogs = true

	log.Println("Galasa CLI - Download artifacts for a run")

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
	err = runs.DownloadArtifacts(
		runNameDownload,
		runForceDownload,
		fileSystem,
		timeService,
		console,
		apiServerUrl,
		runDownloadTargetFolder,
	)
	if err != nil {
		panic(err)
	}
}
