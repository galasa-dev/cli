/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/galasa-dev/cli/pkg/embedded"
	"github.com/galasa-dev/cli/pkg/galasaapi"
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

	log.Printf("ResetRun entered.")

	if runName == "" {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_MISSING_NAME_FLAG, runName)
	}

	if err == nil {

		runId := "GETRUNIDROMRUNNAME"

		updateRunStatusRequest := createUpdateRunStatusRequest(RESET_STATUS, RESET_RESULT)

		err = resetRun(runId, updateRunStatusRequest, apiClient)

	}

	log.Printf("ResetRun exiting. err is %v", err)
	return err
}

func resetRun(runId string,
	runStatusUpdateRequest *galasaapi.UpdateRunStatusRequest,
	apiClient *galasaapi.APIClient,
) error {
	var err error = nil
	var resp *http.Response
	var context context.Context = nil
	var restApiVersion string

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {

		_, resp, err = apiClient.ResultArchiveStoreAPIApi.PutRasRunStatusById(context, runId).
			UpdateRunStatusRequest(*runStatusUpdateRequest).
			ClientApiVersion(restApiVersion).Execute()

		if (resp != nil) && (resp.StatusCode != http.StatusOK) {
			defer resp.Body.Close()

			responseBody, err := io.ReadAll(resp.Body)
			log.Printf("putRasRunStatusById Failed - HTTP Response - Status Code: '%v' Payload: '%v'", resp.StatusCode, string(responseBody))

			if err == nil {
				// err = galasaErrors.NewGalasaError(galasaErrors.)
			}

			// read response and log errors if needed
		}

	}

	return err
}

func createUpdateRunStatusRequest(status string, result string) *galasaapi.UpdateRunStatusRequest {
	var updateRunStatusRequest = galasaapi.NewUpdateRunStatusRequest()

	updateRunStatusRequest.SetStatus(status)
	updateRunStatusRequest.SetResults(result)

	return updateRunStatusRequest
}
