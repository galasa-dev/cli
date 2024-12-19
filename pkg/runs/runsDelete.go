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

	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

// ---------------------------------------------------

// RunsDelete - performs all the logic to implement the `galasactl runs delete` command,
// but in a unit-testable manner.
func RunsDelete(
	runName string,
	console spi.Console,
	apiServerUrl string,
	apiClient *galasaapi.APIClient,
	timeService spi.TimeService,
	byteReader spi.ByteReader,
) error {
	var err error

	log.Printf("RunsDelete entered.")

	if runName != "" {
		// Validate the runName as best we can without contacting the ecosystem.
		err = ValidateFlagValue(runName)
	}

	if err == nil {

		requestorParameter := ""
		resultParameter := ""
		fromAgeHours := 0
		toAgeHours := 0
		group := ""
		shouldGetActive := false
		var runs []galasaapi.Run
		runs, err = GetRunsFromRestApi(runName, requestorParameter, resultParameter, fromAgeHours, toAgeHours, shouldGetActive, timeService, apiClient, group)

		if err == nil {

			if len(runs) == 0 {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SERVER_DELETE_RUN_NOT_FOUND, runName)
			} else {
				err = deleteRuns(runs, apiClient, byteReader)
			}
		}

		if err != nil {
			console.WriteString(err.Error())
		}
	}
	log.Printf("RunsDelete exiting. err is %v\n", err)
	return err
}

func deleteRuns(
	runs []galasaapi.Run,
	apiClient *galasaapi.APIClient,
	byteReader spi.ByteReader,
) error {
	var err error

	var restApiVersion string
	var context context.Context = nil

	var httpResponse *http.Response

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()
	if err == nil {
		for _, run := range runs {
			runId := run.GetRunId()
			runName := *run.GetTestStructure().RunName

			apicall := apiClient.ResultArchiveStoreAPIApi.DeleteRasRunById(context, runId).ClientApiVersion(restApiVersion)
			httpResponse, err = apicall.Execute()

			if httpResponse != nil {
				defer httpResponse.Body.Close()
			}

			// 200-299 http status codes manifest in an error.
			if err != nil {
				if httpResponse == nil {
					// We never got a response, error sending it or something ?
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SERVER_DELETE_RUNS_FAILED, err.Error())
				} else {
					err = galasaErrors.HttpResponseToGalasaError(
						httpResponse,
						runName,
						byteReader,
						galasaErrors.GALASA_ERROR_DELETE_RUNS_NO_RESPONSE_CONTENT,
						galasaErrors.GALASA_ERROR_DELETE_RUNS_RESPONSE_PAYLOAD_UNREADABLE,
						galasaErrors.GALASA_ERROR_DELETE_RUNS_UNPARSEABLE_CONTENT,
						galasaErrors.GALASA_ERROR_DELETE_RUNS_SERVER_REPORTED_ERROR,
						galasaErrors.GALASA_ERROR_DELETE_RUNS_EXPLANATION_NOT_JSON,
					)
				}
			}

			if err != nil {
				break
			} else {
				log.Printf("Run with runId '%s' and runName '%s', was deleted OK.\n", runId, run.TestStructure.GetRunName())
			}
		}
	}

	return err
}
