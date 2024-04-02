/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"encoding/json"
	"log"
	"path/filepath"
	"time"

	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/golang-jwt/jwt/v5"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

type BearerTokenJson struct {
	Jwt string `json:"jwt"`
}

const (
	TOKEN_EXPIRY_BUFFER_MINUTES = 10
)

// Writes a new bearer-token.json file containing a JWT in the following format:
//
//	{
//	  "jwt": "<bearer-token-here>"
//	}
func WriteBearerTokenJsonFile(fileSystem files.FileSystem, galasaHome utils.GalasaHome, jwt string) error {
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

// Gets the JWT from the bearer-token.json file if it exists, errors if the file does not exist or if the token is invalid
func GetBearerTokenFromTokenJsonFile(fileSystem files.FileSystem, galasaHome utils.GalasaHome, timeService utils.TimeService) (string, error) {
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
	} else {
		log.Printf("Validating bearer token retrieved from file '%s'", bearerTokenFilePath)
		if !IsBearerTokenValid(bearerToken, timeService) {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_BEARER_TOKEN)
		} else {
			log.Printf("Validated bearer token retrieved from file '%s' OK", bearerTokenFilePath)
		}
	}

	return bearerToken, err
}

// Checks whether a given bearer token is valid or not, returning true if it is valid and false otherwise
func IsBearerTokenValid(bearerTokenString string, timeService utils.TimeService) bool {
	var err error = nil
	var bearerToken *jwt.Token

	// Decode the bearer token without verifying its signature
	bearerToken, _, err = jwt.NewParser().ParseUnverified(bearerTokenString, jwt.MapClaims{})
	if err == nil {
		var tokenExpiry *jwt.NumericDate
		tokenExpiry, err = bearerToken.Claims.GetExpirationTime()
		if err == nil {
			// Add a buffer to the current time to make sure the bearer token does not expire within
			// this buffer (e.g. if the buffer is 10 mins, make sure the token doesn't expire within 10 mins)
			acceptableExpiryTime := timeService.Now().Add(time.Duration(TOKEN_EXPIRY_BUFFER_MINUTES) * time.Minute)
			if (tokenExpiry.Time).After(acceptableExpiryTime) {
				return true
			}
		}
	}
	return false
}
