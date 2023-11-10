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

	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/utils"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

type BearerTokenJson struct {
    Jwt string `json:"jwt"`
}

// Writes a new bearer-token.json file containing a JWT with the following format:
// {
//   "jwt": "<bearer-token-here>"
// }
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

// Gets the JWT from the bearer-token.json file if it exists, errors if the file does not exist
func GetBearerTokenFromTokenJsonFile(fileSystem files.FileSystem, galasaHome utils.GalasaHome) (string, error) {
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
