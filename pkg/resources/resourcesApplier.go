/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package resources

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/files"
)

// ApplyResources - performs all the logic to implement the
// `galasactl resources apply/create/update` command,
// but in a unit-testable manner.
func ApplyResources(
	action string,
	filePath string,
	fileSystem files.FileSystem,
	apiServerUrl string,
) error {
	var err error
	var fileContent string
	var jsonBytes []byte

	err = validateFilePathExists(fileSystem, filePath)

	if err == nil {
		//read yaml file content
		fileContent, err = getYamlFileContent(fileSystem, filePath)

		if err == nil {
			jsonBytes, err = yamlToByteArray(fileContent, action)
		}

		if err == nil {
			err = sendResourcesRequestToServer(jsonBytes, apiServerUrl)
		}
	}
	return err
}

func sendResourcesRequestToServer(payloadJsonToSend []byte, apiServerUrl string) error {
	var err error
	var responseBody []byte
	resourcesApiServerUrl := apiServerUrl + "/resources/"

	var req *http.Request
	req, err = http.NewRequest("POST", resourcesApiServerUrl, bytes.NewBuffer(payloadJsonToSend))

	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Encoding", "gzip,deflate,br")

		log.Printf("sendResourcesRequestToServer url:%s - payload: '%s'", resourcesApiServerUrl, string(payloadJsonToSend))

		var resp *http.Response
		client := &http.Client{}

		// A non-2xx status code doesn't cause an error.
		// If there is an error, the response should be nil
		resp, err = client.Do(req)
		if err == nil {
			statusCode := resp.StatusCode

			defer resp.Body.Close()

			if statusCode != http.StatusOK && statusCode != http.StatusCreated {
				responseBody, err = io.ReadAll(resp.Body)

				log.Printf("sendResourcesRequestToServer Failed - HTTP response - status code:%v payload:%v", statusCode, string(responseBody))
				if err == nil {
					//Get an arraylist of errors for each resource (obtained from the yaml file given in the command)
					var apiErrors *galasaErrors.GalasaAPIErrorsArray
					apiErrors, err = galasaErrors.NewGalasaApiErrorsArray(responseBody)
					if err == nil {
						//Ensure that the conversion of the error doesn't raise another exception
						errMessages := apiErrors.GetErrorMessages()
						responseString := fmt.Sprint(strings.Join(errMessages, "\n"))

						if 400 <= statusCode && statusCode <= 499 {
							err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RESOURCES_RESP_CLIENT_ERROR, statusCode, responseString)
						} else if 500 <= statusCode && statusCode <= 599 {
							err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RESOURCES_RESP_SERVER_ERROR)
						} else {
							err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RESOURCES_RESP_UNEXPECTED_ERROR, statusCode, responseString)
						}
					} else {
						//error occurred when trying to retrieve the api error
						//unable to retrieve galasa apiError
						err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_GET_API_ERRORS_ARRAY, err)
					}
				} else {
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_READ_RESPONSE_BODY, err)
				}

			} else {
				log.Println("response Status:", resp.Status)
				log.Println("response Headers:", resp.Header)
			}
		}
	}

	return err
}
