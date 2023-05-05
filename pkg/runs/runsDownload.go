/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"context"
	"log"
	"net/http"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

// DownloadArtifacts - performs all the logic to implement the `galasactl runs download` command,
// but in a unit-testable manner.
func DownloadArtifacts(
	runName string,
	timeService utils.TimeService,
	console utils.Console,
	apiServerUrl string,
) error {

	var err error
	var runJson []galasaapi.Run
	var artifacts []galasaapi.Artifact

	if err == nil {

		runJson, err = GetRunsFromRestApi(runName, timeService, apiServerUrl)
		if err == nil {
			for _, run := range runJson {
				artifacts, err = GetArtifactIDsFromRestApi(run.GetRunId(), apiServerUrl)
				if err == nil {

				}
			}
		}
	}

	return err
}

func GetArtifactIDsFromRestApi(
	runID string,
	apiServerUrl string,
) ([]galasaapi.Artifact, error) {

	var err error = nil
	var results []galasaapi.Run = make([]galasaapi.Run, 0)

	var context context.Context = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	toTime := timeService.Now()
	var pageNumberWanted int32 = 1
	gotAllResults := false

	for (!gotAllResults) && (err == nil) {

		var runData *galasaapi.RunResults
		var httpResponse *http.Response
		log.Printf("Requesting page '%d' ", pageNumberWanted)
		runData, httpResponse, err = restClient.ResultArchiveStoreAPIApi.
			GetRasSearchRuns(context).
			To(toTime).
			Runname(runName).
			Page(pageNumberWanted).
			Sort("to:desc").
			Execute()

		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
		} else {
			if httpResponse.StatusCode != HTTP_STATUS_CODE_OK {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
			} else {

				// Copy the results from this page into our bigger list of results.
				runsOnThisPage := runData.GetRuns()
				// Add all the runs into our set of results.
				// Note: The ... syntax means 'all of the array', so they all get appended at once.
				results = append(results, runsOnThisPage...)

				// Have we processed the last page ?
				if pageNumberWanted == runData.GetNumPages() {
					gotAllResults = true
				}
			}
		}
	}

	return results, err
}

// Retrieves test runs from the ecosystem API that match a given runName.
// Multiple test runs can be returned as the runName is not unique.
func GetFileFromRestApi(
	runName string,
	outputFormat OutputFormat,
	timeService utils.TimeService,
	apiServerUrl string,
) ([]galasaapi.Run, error) {

	var err error = nil
	var results []galasaapi.Run = make([]galasaapi.Run, 0)

	var context context.Context = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	toTime := timeService.Now()
	var pageNumberWanted int32 = 1
	gotAllResults := false

	for (!gotAllResults) && (err == nil) {

		var runData *galasaapi.RunResults
		var httpResponse *http.Response
		log.Printf("Requesting page '%d' ", pageNumberWanted)
		runData, httpResponse, err = restClient.ResultArchiveStoreAPIApi.
			GetRasSearchRuns(context).
			To(toTime).
			Runname(runName).
			Page(pageNumberWanted).
			Sort("to:desc").
			Execute()

		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
		} else {
			if httpResponse.StatusCode != HTTP_STATUS_CODE_OK {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
			} else {

				// Copy the results from this page into our bigger list of results.
				runsOnThisPage := runData.GetRuns()
				// Add all the runs into our set of results.
				// Note: The ... syntax means 'all of the array', so they all get appended at once.
				results = append(results, runsOnThisPage...)

				// Have we processed the last page ?
				if pageNumberWanted == runData.GetNumPages() {
					gotAllResults = true
				}
			}
		}
	}

	return results, err
}
