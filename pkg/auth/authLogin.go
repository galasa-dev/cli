/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/props"
	"github.com/galasa.dev/cli/pkg/utils"
)

// Login - performs all the logic to implement the `galasactl auth login` command
func Login(apiServerUrl string, fileSystem files.FileSystem, galasaHome utils.GalasaHome) error {

	var err error = nil
	var authProperties galasaapi.AuthProperties
	authProperties, err = getAuthProperties(fileSystem, galasaHome)
	if err == nil {
		var jwt string
		jwt, err = GetJwtFromRestApi(apiServerUrl, authProperties)
		if err == nil {
			err = writeBearerTokenJsonFile(fileSystem, galasaHome, jwt)
		}
	}

	return err
}

// Gets authentication properties from the user's galasactl.properties file
func getAuthProperties(fileSystem files.FileSystem, galasaHome utils.GalasaHome) (galasaapi.AuthProperties, error) {
	var err error = nil
	authProperties := galasaapi.NewAuthProperties()

	galasactlPropertiesFilePath := filepath.Join(galasaHome.GetNativeFolderPath(), "galasactl.properties")
	galasactlProperties, err := props.ReadPropertiesFile(fileSystem, galasactlPropertiesFilePath)
	if err == nil {
		clientIdProperty := "GALASA_CLIENT_ID"
		secretProperty := "GALASA_SECRET"
		accessTokenProperty := "GALASA_ACCESS_TOKEN"

		requiredAuthProperties := []string{clientIdProperty, secretProperty, accessTokenProperty}
		err = validateRequiredGalasactlProperties(requiredAuthProperties, galasactlProperties)

		if err == nil {
			authProperties.SetClientId(galasactlProperties[clientIdProperty])
			authProperties.SetSecret(galasactlProperties[secretProperty])
			authProperties.SetRefreshToken(galasactlProperties[accessTokenProperty])
		}

	} else {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_READ_FILE, galasactlPropertiesFilePath, err.Error())
	}

	return *authProperties, err
}

// Ensures the provided galasactl properties contain values for the required properties, returning an error if a property is missing
func validateRequiredGalasactlProperties(requiredProperties []string, galasactlProperties props.JavaProperties) error {
	var err error = nil
	for _, property := range requiredProperties {
		if galasactlProperties[property] == "" {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_MISSING_GALASACTL_PROPERTY, property)
			break
		}
	}
	return err
}

// Gets a JSON Web Token (JWT) from the API server's /auth endpoint
func GetJwtFromRestApi(apiServerUrl string, authProperties galasaapi.AuthProperties) (string, error) {
	var err error = nil
	var context context.Context = nil
	var jwtJsonStr string

	restClient := api.InitialiseAPI(apiServerUrl)

	tokenResponse, httpResponse, err := restClient.AuthenticationAPIApi.PostAuthenticate(context).
		AuthProperties(authProperties).
		Execute()
	defer httpResponse.Body.Close()

	if err != nil {
		log.Println("Failed to retrieve bearer token from API server")
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_BEARER_TOKEN_FROM_API_SERVER, err.Error())
	} else {
		var tokenResponseJson []byte
		tokenResponseJson, err = tokenResponse.MarshalJSON()
		jwtJsonStr = string(tokenResponseJson)
		log.Println("Bearer token received from API server OK")

	}
	return jwtJsonStr, err
}

// Writes a new bearer-token.json file containing a JWT with the following format:
// {
//   "jwt": "<bearer-token-here>"
// }
func writeBearerTokenJsonFile(fileSystem files.FileSystem, galasaHome utils.GalasaHome, jwt string) error {
	bearerTokenFilePath := filepath.Join(galasaHome.GetNativeFolderPath(), "bearer-token.json")

	log.Printf("Writing bearer token to file '%s'", bearerTokenFilePath)
	err := fileSystem.WriteTextFile(bearerTokenFilePath, jwt)

	if err == nil {
		log.Printf("Written bearer token to file '%s' OK", bearerTokenFilePath)
	} else {
		log.Printf("Failed to write bearer token file '%s'", bearerTokenFilePath)
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_WRITE_FILE, bearerTokenFilePath, err.Error())
	}
	return err
}

type BearerTokenJson struct {
	Jwt string `json:"jwt"`
}

// Gets the JWT from the bearer-token.json file if it exists, errors if the file does not exist
func getBearerTokenFromTokenJsonFile(fileSystem files.FileSystem, galasaHome utils.GalasaHome) (string, error) {
	var err error = nil
	var bearerToken string = ""
	var bearerTokenJsonContents string = ""

	bearerTokenFilePath := filepath.Join(galasaHome.GetNativeFolderPath(), "bearer-token.json")

	log.Printf("Retrieving bearer token from file '%s'", bearerTokenFilePath)
	bearerTokenJsonContents, err = fileSystem.ReadTextFile(bearerTokenFilePath)
	if err == nil {
		var bearerTokenJson BearerTokenJson
		err = json.Unmarshal([]byte(bearerTokenJsonContents), &bearerTokenJson)
		if err == nil {
			bearerToken = bearerTokenJson.Jwt
			log.Printf("Retrieved bearer token from file '%s' OK", bearerTokenFilePath)
		}
	}

	if err != nil {
		log.Printf("Could not retrieve bearer token from file '%s'", bearerTokenFilePath)
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_BEARER_TOKEN_FROM_FILE, bearerTokenFilePath, err.Error())
	}
	return bearerToken, err
}

// Gets a new authenticated API client, attempting to log in if a bearer token file does not exist
func GetAuthenticatedAPIClient(apiServerUrl string, fileSystem files.FileSystem, galasaHome utils.GalasaHome) (*galasaapi.APIClient, error) {
	bearerToken, err := getBearerTokenFromTokenJsonFile(fileSystem, galasaHome)
	if err != nil {
		// Attempt to log in
		log.Printf("Logging in to the Galasa Ecosystem at '%s'", apiServerUrl)
		err = Login(apiServerUrl, fileSystem, galasaHome)
		if err == nil {
			log.Printf("Logged in to the Galasa Ecosystem at '%s' OK", apiServerUrl)
			bearerToken, err = getBearerTokenFromTokenJsonFile(fileSystem, galasaHome)
		}

	}

	var apiClient *galasaapi.APIClient
	if err == nil {
		apiClient = api.InitialiseAuthenticatedAPI(apiServerUrl, bearerToken)
	} else {
		// Continue without a bearer token
		apiClient = api.InitialiseAPI(apiServerUrl)
		err = nil
	}
	return apiClient, err
}
