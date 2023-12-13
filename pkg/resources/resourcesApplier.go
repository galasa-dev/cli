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
	var reqMethod = "POST"

	err = validateFilePathExists(fileSystem, filePath)

	if err == nil {
		fileContent, err = getYamlFileContent(fileSystem, filePath)

		if err == nil {
			//convert resources in yaml file into a json payload
			jsonBytes, err = yamlToByteArray(fileContent, action)
		}

		if err == nil {
			if action == "delete"{
				reqMethod = "DELETE"
			}
			err = sendJsonToApi(reqMethod, jsonBytes, apiServerUrl)
		}
	}
	return err
}

func sendJsonToApi(reqMethod string, payload []byte, apiServerUrl string) error {

	var err error
	var responseBody []byte
	if err == nil {

		var req *http.Request
		req, err = http.NewRequest(reqMethod, apiServerUrl+"/resources/", bytes.NewBuffer(payload))

		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Accept-Encoding", "gzip,deflate,br")

			var resp *http.Response
			client := &http.Client{}

			// A non-2xx status code doesn't cause an error.
			// If there is an error, the response should be nil
			resp, err = client.Do(req)
			if err == nil {
				statusCode := resp.StatusCode

				defer resp.Body.Close()

				if statusCode != http.StatusOK && statusCode != http.StatusCreated {
					log.Printf("Response status code is not okay (200 or 201), but rather %v. Error: %v", statusCode, err)

					responseBody, err = io.ReadAll(resp.Body)
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
								err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RESOURCES_RESP_SERVER_ERROR, statusCode, responseString)
							} else {
								err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RESOURCES_RESP_UNEXPECTED_ERROR, statusCode, responseString)
							}
						}
					}

				} else {
					log.Println("response Status:", resp.Status)
					log.Println("response Headers:", resp.Header)
				}
			}
		}
	}
	return err
}
