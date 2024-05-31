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

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

const (
	TOKEN_ID_PATTERN = "^[a-zA-Z0-9]+$"
)

// GetTokens - performs all the logic to implement the `galasactl auth tokens delete --tokenid xxx` command
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

// token id are currently strictly alphanumerical
func validateTokenId(tokenId string) error {

	validTokenIdFormat, err := regexp.Compile(TOKEN_ID_PATTERN)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_COMPILE_TOKEN_ID_REGEX, err.Error())
	} else {
		//check if the token id format matches
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

	_, resp, err = apiClient.AuthenticationAPIApi.DeleteToken(context, tokenId).Execute()

	if (resp != nil) && (resp.StatusCode != http.StatusOK) {
		defer resp.Body.Close()

		responseBody, err = io.ReadAll(resp.Body)
		log.Printf("deleteTokenFromRestApi - HTTP response - status code: '%v' payload: '%v' ", resp.StatusCode, string(responseBody))

		if err == nil {
			//no error returned if a 404 (Not Found) payload is returned, as the ultimate goal of the token not being present is accomplished
			if resp.StatusCode == http.StatusNotFound {
				log.Printf("deleteTokenFromRestApi - token id '%s' was not found, and therefore cannot be deleted.", tokenId)
			} else {
				var errorFromServer *galasaErrors.GalasaAPIError
				errorFromServer, err = galasaErrors.GetApiErrorFromResponse(responseBody)

				if err == nil {
					//return galasa api error, because status code is not 200 (OK)
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DELETE_TOKEN_FAILED, tokenId, errorFromServer.Message)
				} else {
					//unable to parse response into api error
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DELETE_TOKEN_RESPONSE_PARSING, err.Error())
				}
			}
		} else {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_READ_RESPONSE_BODY, err.Error())
		}

	}

	return err
}
