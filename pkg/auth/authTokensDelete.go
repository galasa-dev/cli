/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"context"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

const (
	TOKEN_ID_PATTERN = "^[a-zA-Z0-9\\_\\-]+$"
)

// DeleteToken - performs all the logic to implement the `galasactl auth tokens delete --tokenid xxx` command
func DeleteToken(
	tokenId string,
	apiClient *galasaapi.APIClient,
	console spi.Console,
) error {
	var err error

	err = validateTokenId(tokenId)
	if err == nil {
		log.Print("DeleteToken - valid token id provided")
		err = deleteTokenFromRestApi(tokenId, apiClient)
	}

	return err
}

func validateTokenId(tokenId string) error {

	validTokenIdFormat, err := regexp.Compile(TOKEN_ID_PATTERN)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_COMPILE_TOKEN_ID_REGEX, err.Error())
	} else {
		// Check if the token ID format is valid
		if !validTokenIdFormat.MatchString(tokenId) {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_TOKEN_ID_FORMAT, tokenId)
		}
	}

	return err
}

func deleteTokenFromRestApi(tokenId string, apiClient *galasaapi.APIClient) error {
	var err error
	var context context.Context
	var resp *http.Response
	var responseBody []byte
	
	var restApiVersion string
	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {
		_, resp, err = apiClient.AuthenticationAPIApi.DeleteToken(context, tokenId).ClientApiVersion(restApiVersion).Execute()
	
		if err != nil {
			// Try to get the error returned from the API server and return that message
			if (resp != nil) && (resp.StatusCode != http.StatusOK) {
				defer resp.Body.Close()
					responseBody, err = io.ReadAll(resp.Body)
					log.Printf("deleteTokenFromRestApi - HTTP response - Status Code: '%v' Payload: '%v' ", resp.StatusCode, string(responseBody))
			
					if err == nil {
						var errorFromServer *galasaErrors.GalasaAPIError
						errorFromServer, err = galasaErrors.GetApiErrorFromResponse(responseBody)
			
						if err == nil {
							// Return a Galasa API error, because the status code is not 200 (OK)
							err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DELETE_TOKEN_FAILED, tokenId, errorFromServer.Message)
						}
					}
			} else {
				// No response was received from the API server, so something else may have gone wrong
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DELETE_TOKEN_FAILED, tokenId, err.Error())
			}
		}
	}
	return err
}
