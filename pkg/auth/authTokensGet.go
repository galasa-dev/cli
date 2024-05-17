/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"context"
	"log"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/tokensformatter"
	"github.com/galasa-dev/cli/pkg/utils"
)

// GetTokens - performs all the logic to implement the `galasactl auth tokens get` command
func GetTokens(
	apiClient *galasaapi.APIClient,
	console utils.Console,
) error {
	var err error = nil
	var context context.Context = nil
	var authTokens *galasaapi.AuthTokens

	authTokens, _, err = apiClient.AuthenticationAPIApi.GetTokens(context).Execute()
	if err != nil {
		log.Println("GetTokens - Failed to retrieve list of tokens from API server")
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_TOKEN_LIST_FROM_API_SERVER, err.Error())
	} else {
		log.Printf("GetTokens -  %v tokens collected", len(authTokens.GetTokens()))

		summaryFormatter := tokensformatter.NewTokenSummaryFormatter()

		var outputText string
		outputText, err = summaryFormatter.FormatTokens(authTokens.GetTokens())

		if err == nil {
			console.WriteString(outputText)
		}

	}

	return err
}
