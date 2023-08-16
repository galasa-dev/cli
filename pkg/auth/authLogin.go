/*
 * Copyright contributors to the Galasa project
 */
package auth

import (
	"context"
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
		clientIdProperty := "auth.client.id"
		secretProperty := "auth.secret"
		accessTokenProperty := "auth.access.token"

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
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_BEARER_TOKEN_FROM_API_SERVER, err.Error())
		log.Printf("Failed to retrieve JWT from API server. %s", err.Error())
	} else {
		var tokenResponseJson []byte
		tokenResponseJson, err = tokenResponse.MarshalJSON()
		jwtJsonStr = string(tokenResponseJson)
		log.Println("JWT received from API server OK")

	}
	return jwtJsonStr, err
}

// Writes a new bearer-token.json file containing a JWT with the following format:
// {
//   "jwt": "<bearer-token-here>"
// }
func writeBearerTokenJsonFile(fileSystem files.FileSystem, galasaHome utils.GalasaHome, jwt string) error {
	bearerTokenFilePath := filepath.Join(galasaHome.GetNativeFolderPath(), "bearer-token.json")

	log.Printf("Writing JWT to bearer token file '%s'", bearerTokenFilePath)
	err := fileSystem.WriteTextFile(bearerTokenFilePath, jwt)

	if err == nil {
		log.Printf("Written JWT to bearer token file '%s' OK", bearerTokenFilePath)
	} else {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_WRITE_FILE, bearerTokenFilePath, err.Error())
		log.Printf("Failed to write bearer token file '%s'. %s", bearerTokenFilePath, err.Error())
	}
	return err
}
