/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"context"
	"log"
	"strings"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/tokensformatter"
)

// GetTokens - performs all the logic to implement the `galasactl auth tokens get` command
func GetTokens(
	apiClient *galasaapi.APIClient,
	console spi.Console,
) error {

	authTokens, err := getAuthTokensFromRestApi(apiClient)

	if err == nil {
		summaryFormatter := tokensformatter.NewTokenSummaryFormatter()

		var outputText string
		outputText, err = summaryFormatter.FormatTokens(authTokens)

		if err == nil {
			console.WriteString(outputText)
		}
	}

	return err
}

func GetTokensByLoginId(
	apiClient *galasaapi.APIClient,
	console spi.Console,
	loginId string,
) error {

	authTokens, err := getAuthTokensByLoginIdFromRestApi(apiClient, loginId)

	if err == nil {
		summaryFormatter := tokensformatter.NewTokenSummaryFormatter()

		var outputText string
		outputText, err = summaryFormatter.FormatTokens(authTokens)

		if err == nil {
			console.WriteString(outputText)
		}
	}

	return err
}

func getAuthTokensByLoginIdFromRestApi(apiClient *galasaapi.APIClient, loginId string) ([]galasaapi.AuthToken, error) {
	var context context.Context = nil
	var authTokens []galasaapi.AuthToken

	validateLoginIdFlag(loginId)

	tokens, resp, err := apiClient.AuthenticationAPIApi.GetTokens(context).LoginId(loginId).Execute()

	if err != nil {
		log.Println("getAuthTokensFromRestApi - Failed to retrieve list of tokens from API server")
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_TOKEN_LIST_FROM_API_SERVER, err.Error())
	} else {
		defer resp.Body.Close()
		authTokens = tokens.GetTokens()
		log.Printf("getAuthTokensFromRestApi -  %v tokens collected", len(authTokens))
	}

	return authTokens, err
}

func getAuthTokensFromRestApi(apiClient *galasaapi.APIClient) ([]galasaapi.AuthToken, error) {
	var context context.Context = nil
	var authTokens []galasaapi.AuthToken

	tokens, resp, err := apiClient.AuthenticationAPIApi.GetTokens(context).Execute()

	if err != nil {
		log.Println("getAuthTokensFromRestApi - Failed to retrieve list of tokens from API server")
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_TOKEN_LIST_FROM_API_SERVER, err.Error())
	} else {
		defer resp.Body.Close()
		authTokens = tokens.GetTokens()
		log.Printf("getAuthTokensFromRestApi -  %v tokens collected", len(authTokens))
	}

	return authTokens, err
}

func validateLoginIdFlag(loginId string) (string, error) {

	var err error

	loginId = strings.TrimSpace(loginId)

	if loginId == "" {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_MISSING_USER_LOGIN_ID_FLAG)
	}

	return loginId, err
}
