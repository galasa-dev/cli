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
	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
)

// Login - performs all the logic to implement the `galasactl auth login` command
func Login(apiServerUrl string, fileSystem files.FileSystem, galasaHome utils.GalasaHome, env utils.Environment) error {

	var err error = nil
	var authProperties galasaapi.AuthProperties
	authProperties, err = GetAuthProperties(fileSystem, galasaHome, env)
	if err == nil {
		var jwt string
		jwt, err = GetJwtFromRestApi(apiServerUrl, authProperties)
		if err == nil {
			err = WriteBearerTokenJsonFile(fileSystem, galasaHome, jwt)
		}
	}
	return err
}

// Gets a JSON Web Token (JWT) from the API server's /auth endpoint
func GetJwtFromRestApi(apiServerUrl string, authProperties galasaapi.AuthProperties) (string, error) {
	var err error = nil
	var context context.Context = nil
	var jwtJsonStr string
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
			var tokenResponseJson []byte
			tokenResponseJson, err = tokenResponse.MarshalJSON()
			jwtJsonStr = string(tokenResponseJson)
			log.Println("Bearer token received from API server OK")
		}
	}

	return jwtJsonStr, err
}

// Gets a new authenticated API client, attempting to log in if a bearer token file does not exist
func GetAuthenticatedAPIClient(
	apiServerUrl string,
	fileSystem files.FileSystem,
	galasaHome utils.GalasaHome,
	timeService utils.TimeService,
	env utils.Environment,
) *galasaapi.APIClient {
	bearerToken, err := GetBearerTokenFromTokenJsonFile(fileSystem, galasaHome, timeService)
	if err != nil {
		// Attempt to log in
		log.Printf("Logging in to the Galasa Ecosystem at '%s'", apiServerUrl)
		err = Login(apiServerUrl, fileSystem, galasaHome, env)
		if err == nil {
			log.Printf("Logged in to the Galasa Ecosystem at '%s' OK", apiServerUrl)
			bearerToken, err = GetBearerTokenFromTokenJsonFile(fileSystem, galasaHome, timeService)
		}
	}

	var apiClient *galasaapi.APIClient
	if err == nil {
		apiClient = api.InitialiseAuthenticatedAPI(apiServerUrl, bearerToken)
	} else {
		// Temporary code to allow the CLI to continue running commands while authentication is being implemented.
		// Remove this once authentication has been delivered and users can authenticate against an ecosystem
		apiClient = api.InitialiseAPI(apiServerUrl)
	}
	return apiClient
}
