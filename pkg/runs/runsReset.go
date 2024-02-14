/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	galasaapi "github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
)

var (
	RESET_STATUS = "queued"
	RESET_RESULT = ""
)

func ResetRun(
	runName string,
	timeService utils.TimeService,
	console utils.Console,
	apiServerUrl string,
	apiClient *galasaapi.APIClient,
) error {
	var err error
	var runId string

	log.Printf("ResetRun entered.")

	if runName == "" {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_MISSING_NAME_FLAG, runName)
	}

	if (err == nil) && (runName != "") {
		err = ValidateRunName(runName)
	}

	if err == nil {

		runId, err = getRunIdFromRunName(runName, timeService, apiClient)

		if err == nil {

			updateRunStatusRequest := createUpdateRunStatusRequest(RESET_STATUS, RESET_RESULT)

			err = resetRun(runName, runId, updateRunStatusRequest, apiClient)

			if err == nil {
				consoleErr := console.WriteString(fmt.Sprintf(galasaErrors.GALASA_INFO_RUNS_RESET_SUCCESS.Template, runName))

				// Console error is not as important to report as the original error if there was one.
				if consoleErr != nil && err == nil {
					err = consoleErr
				}
			}

		}

	}

	log.Printf("ResetRun exiting. err is %v", err)
	return err
}

func resetRun(runName string,
	runId string,
	runStatusUpdateRequest *galasaapi.UpdateRunStatusRequest,
	apiClient *galasaapi.APIClient,
) error {
	var err error = nil
	var resp *http.Response
	var context context.Context = nil
	var restApiVersion string
	var responseBody []byte

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {

		_, resp, err = apiClient.ResultArchiveStoreAPIApi.PutRasRunStatusById(context, runId).
			UpdateRunStatusRequest(*runStatusUpdateRequest).
			ClientApiVersion(restApiVersion).Execute()

		if (resp != nil) && (resp.StatusCode != http.StatusOK) {
			defer resp.Body.Close()

			responseBody, err = io.ReadAll(resp.Body)
			log.Printf("putRasRunStatusById Failed - HTTP Response - Status Code: '%v' Payload: '%v'", resp.StatusCode, string(responseBody))

			if err == nil {
				var errorFromServer *galasaErrors.GalasaAPIError
				errorFromServer, err = galasaErrors.GetApiErrorFromResponse(responseBody)

				if err == nil {
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RESET_RUN_FAILED, runName, errorFromServer.Message)
				} else {
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RESET_RUN_RESPONSE_PARSING)
				}

			} else {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_READ_RESPONSE_BODY, err.Error())
			}
		}

	}

	return err
}

func getRunIdFromRunName(runName string,
	timeService utils.TimeService,
	apiClient *galasaapi.APIClient,
) (string, error) {
	var err error
	var runs []galasaapi.Run
	var runId string

	requestorParameter := ""
	resultParameter := ""
	fromAgeHours := 0
	toAgeHours := 0
	shouldGetActive := true

	runs, err = GetRunsFromRestApi(runName, requestorParameter, resultParameter, fromAgeHours, toAgeHours, shouldGetActive, timeService, apiClient)

	if err == nil {

		if len(runs) > 1 {

			// More than 1 active run has been found with this runName, as runs might be stuck in active state like ending...
			
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_MULTIPLE_ACTIVE_RUNS_WITH_RUNNAME, runName)

		} else if len(runs) == 1 {

			runId = runs[0].GetRunId()

		} else {

			log.Printf("No active runs found matching run name: '%s'", runName)
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_NO_ACTIVE_RUNS_WITH_RUNNAME, runName)

		}
	}

	return runId, err
}

func createUpdateRunStatusRequest(status string, result string) *galasaapi.UpdateRunStatusRequest {
	var updateRunStatusRequest = galasaapi.NewUpdateRunStatusRequest()

	updateRunStatusRequest.SetStatus(status)
	updateRunStatusRequest.SetResult(result)

	return updateRunStatusRequest
}
