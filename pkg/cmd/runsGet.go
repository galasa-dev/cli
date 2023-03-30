/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/galasa.dev/cli/pkg/api"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

// Objective: Allow the user to do this:
//    run get --runID 12345
// And then show the results in a human-readable form.

var (
	runsGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get the details of a test runid which ran or is running.",
		Long:  "Get the details of a test runid which ran or is running, displaying the results to the caller",
		Args:  cobra.NoArgs,
		Run:   executeRunsGet,
	}

	// Variables set by cobra's command-line parsing.
	runId string
)

func init() {
	runsGetCmd.PersistentFlags().StringVar(&runId, "runid", "", "the runid of the test run we want information about")
	runsGetCmd.MarkPersistentFlagRequired("runid")

	runsCmd.AddCommand(runsGetCmd)
}

func executeRunsGet(cmd *cobra.Command, args []string) {
	var err error

	utils.CaptureLog(logFileName)
	isCapturingLogs = true

	log.Println("Galasa CLI - Get info about a run")

	// Operations on the file system will all be relative to the current folder.
	fileSystem := utils.NewOSFileSystem()

	// Get the ability to query environment variables.
	env := utils.NewEnvironment()

	// Read the bootstrap properties.
	var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
	var bootstrapData *api.BootstrapData
	bootstrapData, err = api.LoadBootstrap(fileSystem, env, bootstrap, urlService)
	if err != nil {
		panic(err)
	}

	apiServerUrl := bootstrapData.ApiServerURL
	log.Printf("The API sever is at '%s'\n", apiServerUrl)

	// An HTTP client which can communicate with the api server in an ecosystem.
	apiClient := api.InitialiseAPI(apiServerUrl)

	var testRunNames *galasaapi.ResultNames.testRunNames
	var testRunDetail *galasaapi.Run

	testRunNames, _, err = apiClient.ResultArchiveStoreAPIApi.GetRasRunIDsByName(nil, runId).Execute()

	if err == nil {
		testRunDetail, _, err = apiClient.ResultArchiveStoreAPIApi.GetRasRunById(nil, testRunNames).Execute()
	}else{

	}
	if err == nil {

		results := testRunDetail.GetTestStructure()
		status := results.GetStatus()
		result := results.GetResult()

		log.Printf("Runid:'%s' status:'%s' result:'%s'\n", runId, status, result)
	} else {
		log.Printf("Failed to get the details of the runid '%s'. Reason: %s", runId, err.Error())
	}
}
