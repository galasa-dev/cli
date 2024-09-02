/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"context"
	"log"
	"net/http"
	"strconv"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

// ---------------------------------------------------

// GetRuns - performs all the logic to implement the `galasactl runs get` command,
// but in a unit-testable manner.
func RunsDelete(
	runName string,
	console spi.Console,
	apiServerUrl string,
	apiClient *galasaapi.APIClient,
	timeService spi.TimeService,
) error {
	var err error

	log.Printf("GetRuns entered.")

	if runName != "" {
		// Validate the runName as best we can without contacting the ecosystem.
		err = ValidateRunName(runName)
	}

	if err == nil {

		requestorParameter := ""
		resultParameter := ""
		fromAgeHours := 0
		toAgeHours := 0
		shouldGetActive := false
		var runs []galasaapi.Run
		runs, err = GetRunsFromRestApi(runName, requestorParameter, resultParameter, fromAgeHours, toAgeHours, shouldGetActive, timeService, apiClient)

		if err == nil {
			err = deleteRuns(runs, apiClient)
		}
	}
	log.Printf("RunsDelete exiting. err is %v\n", err)
	return err
}

func deleteRuns(
	runs []galasaapi.Run,
	apiClient *galasaapi.APIClient,
) error {
	var err error

	var restApiVersion string
	var context context.Context = nil

	var httpResponse *http.Response

	for _, run := range runs {
		runId := *run.RunId

		apicall := apiClient.ResultArchiveStoreAPIApi.DeleteRasRunById(context, runId).ClientApiVersion(restApiVersion)

		httpResponse, err = apicall.Execute()

		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DELETE_RUN_FAILED, err.Error())
		} else {
			if httpResponse.StatusCode != http.StatusNoContent {
				httpError := "\nhttp response status code: " + strconv.Itoa(httpResponse.StatusCode)
				errString := httpError
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SERVER_DELETE_RUNS_FAILED, errString)
			}
		}

		if err != nil {
			break
		} else {
			log.Printf("Run runId:%s runName: %s was deleted OK.\n", runId, run.TestStructure.GetRunName())
		}
	}

	return err
}
