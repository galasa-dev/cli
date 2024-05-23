/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
)

// Gets a JSON Web Token (JWT) from the API server's /auth endpoint
func GetJwtFromRestApi(apiServerUrl string, authProperties galasaapi.AuthProperties) (string, error) {
	var err error
	var context context.Context = nil
	var jwt string
	var restApiVersion string

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {
		apiClient := api.InitialiseAPI(apiServerUrl)

		var tokenResponse *galasaapi.TokenResponse
		var httpResponse *http.Response
		tokenResponse, httpResponse, err = apiClient.AuthenticationAPIApi.PostAuthenticate(context).
			AuthProperties(authProperties).
			ClientApiVersion(restApiVersion).
			Execute()
		if err != nil {
			log.Println("Failed to retrieve bearer token from API server")
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_BEARER_TOKEN_FROM_API_SERVER, err.Error())
		} else {
			defer httpResponse.Body.Close()
			log.Println("Bearer token received from API server OK")
			jwt = tokenResponse.GetJwt()
		}
	}

	return jwt, err
}
