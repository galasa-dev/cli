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
//    run get --runname 12345
// And then show the results in a human-readable form.

var (
	runsGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get the details of a test runname which ran or is running.",
		Long:  "Get the details of a test runname which ran or is running, displaying the results to the caller",
		Args:  cobra.NoArgs,
		Run:   executeRunsGet,
	}

	// Variables set by cobra's command-line parsing.
	runname string
	output  string
)

func init() {
	runsGetCmd.PersistentFlags().StringVar(&runname, "runname", "", "the runname of the test run we want information about")
	runsGetCmd.PersistentFlags().StringVar(&output, "output", "summary", "output format for the data returned (default : summary)")
	runsGetCmd.MarkPersistentFlagRequired("runname")

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

	var testRunIDs []string
	var testRunDetail *galasaapi.Run

	testRunIDs, _, err = apiClient.ResultArchiveStoreAPIApi.GetRasRunIDsByName(nil, runname).Execute()

	if err == nil {
		for indx, val := range testRunIDs {
			testRunDetail, _, err = apiClient.ResultArchiveStoreAPIApi.GetRasRunById(nil, val).Execute()

			if err == nil {
				results := testRunDetail.GetTestStructure()
				status := results.GetStatus()
				result := results.GetResult()
				log.Printf("runname:'%s' status:'%s' result:'%s'\n", runname, status, result)
			} else {
				log.Printf("Failed to get the details of the runname '%s'. List Item '%d' \n Server Id: '%s' \n Reason: %s", runname, indx, val, err.Error())
			}
		}
	} else {
		log.Printf("Failed to get the details of the runname '%s'. Reason: %s", runname, err.Error())
	}
}
