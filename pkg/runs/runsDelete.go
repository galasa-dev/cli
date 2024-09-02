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

// RunsDelete - performs all the logic to implement the `galasactl runs delete` command,
// but in a unit-testable manner.
func RunsDelete(
	runName string,
	console spi.Console,
	apiServerUrl string,
	apiClient *galasaapi.APIClient,
	timeService spi.TimeService,
) error {
	var err error

	log.Printf("RunsDelete entered.")

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

		// 200-299 http status codes manifest in an error.
		if err != nil {
			if httpResponse == nil {
				// We never got a response, error sending it or something ?
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SERVER_DELETE_RUNS_FAILED, err.Error())
			} else {

				contentType := httpResponse.Header.Get("Content-Type")
				if contentType == "" {
					// There is no content in the response
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SERVER_DELETE_RUNS_FAILED, strconv.Itoa(httpResponse.StatusCode))
				} else {
					// There is content in the response.
					// Process the error response payload.

					// TODO: Fix this bit.
					// httpResponse.Body
					// GetApiErrorFromResponse()
				}
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

// func convertToGalasaError(
// 	response *http.Response,
// 	identifier string,
// 	errorMsgUnexpectedStatusCodeNoResponseBody *galasaErrors.MessageType,
// 	errorMsgUnableToReadResponseBody *galasaErrors.MessageType,
// 	msg3 *galasaErrors.MessageType,
// 	errorMsgResponsePayloadInWrongFormat *galasaErrors.MessageType,
// ) error {
// 	defer response.Body.Close()
// 	var err error
// 	var responseBodyBytes []byte
// 	statusCode := response.StatusCode

// 	if response.ContentLength == 0 {
// 		log.Printf("Failed - HTTP response - status code: '%v'\n", statusCode)
// 		err = galasaErrors.NewGalasaError(errorMsgUnexpectedStatusCodeNoResponseBody, identifier, statusCode)
// 	} else {

// 		responseBodyBytes, err = io.ReadAll(response.Body)
// 		if err != nil {
// 			err = galasaErrors.NewGalasaError(errorMsgUnableToReadResponseBody, identifier, statusCode, err.Error())
// 		} else {

// 			var errorFromServer *galasaErrors.GalasaAPIError
// 			errorFromServer, err = galasaErrors.GetApiErrorFromResponse(responseBodyBytes)

// 			if err != nil {
// 				//unable to parse response into api error. It should have been json.
// 				log.Printf("Failed - HTTP response - status code: '%v' payload in response is not json: '%v' \n", statusCode, string(responseBodyBytes))
// 				err = galasaErrors.NewGalasaError(errorMsgResponsePayloadInWrongFormat, identifier, statusCode)
// 			} else {
// 				// server returned galasa api error structure we understand.
// 				log.Printf("Failed - HTTP response - status code: '%v' server responded with error message: '%v' \n", statusCode,errorMsg)
// 				err = galasaErrors.NewGalasaError(errorMsg, identifier, errorFromServer.Message)
// 			}
// 		}
// 	}
// 	return err
// }
