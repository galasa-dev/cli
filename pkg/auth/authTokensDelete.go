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

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

// GetTokens - performs all the logic to implement the `galasactl auth tokens delete --tokenid xxx` command
func DeleteToken(
	tokenId string,
	apiClient *galasaapi.APIClient,
	console spi.Console,
) error {

	var err error

	err = deleteTokenFromRestApi(tokenId, apiClient)

	return err
}

func deleteTokenFromRestApi(tokenId string, apiClient *galasaapi.APIClient) error {
	var err error
	var context context.Context
	var resp *http.Response
	var responseBody []byte

	apiCall := apiClient.AuthenticationAPIApi.DeleteToken(context, tokenId)
	_, resp, err = apiCall.Execute()

	if (resp != nil) && (resp.StatusCode != http.StatusOK) {
		defer resp.Body.Close()

		responseBody, err = io.ReadAll(resp.Body)
		log.Printf("deleteTokenFromRestApi Failed - HTTP response - status code: '%v' payload: '%v' ", resp.StatusCode, string(responseBody))

		if err == nil {
			var errorFromServer *galasaErrors.GalasaAPIError
			errorFromServer, err = galasaErrors.GetApiErrorFromResponse(responseBody)

			if err == nil {
				//return galasa api error, because status code is not 200 (OK)
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DELETE_TOKEN_FAILED, tokenId, errorFromServer.Message)
			} else {
				//unable to parse response into api error
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DELETE_TOKEN_RESPONSE_PARSING, err.Error())
			}
		} else {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_READ_RESPONSE_BODY, err.Error())
		}

	}

	return err
}
