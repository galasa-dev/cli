/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"encoding/json"
	"log"
	"path/filepath"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/files"
)

type BearerTokenJson struct {
	Jwt string `json:"jwt"`
}

type BearerTokenFile interface {
	WriteJwt(jwt string) error
	ReadJwt() (string, error)
}

type BearerTokenFileImpl struct {
	fileSystem   files.FileSystem
	galasaHome   GalasaHome
	baseFileName string
	timeService  TimeService
}

func NewBearerTokenFile(fileSystem files.FileSystem, galasaHome GalasaHome, baseFileName string, timeService TimeService) BearerTokenFile {
	file := new(BearerTokenFileImpl)
	file.fileSystem = fileSystem
	file.galasaHome = galasaHome
	file.baseFileName = baseFileName
	file.timeService = timeService
	return file
}

// Writes a new bearer-token.json file containing a JWT in the following format:
//
//	{
//	  "jwt": "<bearer-token-here>"
//	}
func (file *BearerTokenFileImpl) WriteJwt(jwt string) error {
	bearerTokenFilePath := filepath.Join(file.galasaHome.GetNativeFolderPath(), "bearer-token.json")

	log.Printf("Writing bearer token to file '%s'", bearerTokenFilePath)

	json, err := buildBearerTokenFileContent(jwt)
	if err == nil {

		err = file.fileSystem.WriteTextFile(bearerTokenFilePath, json)

		if err == nil {
			log.Printf("Written bearer token to file '%s' OK", bearerTokenFilePath)
		} else {
			log.Printf("Failed to write bearer token file '%s'", bearerTokenFilePath)
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_WRITE_FILE, bearerTokenFilePath, err.Error())
		}
	}
	return err
}

// Pack the json string into a structure.
func buildBearerTokenFileContent(jwt string) (contentJson string, err error) {
	content := BearerTokenJson{
		Jwt: jwt,
	}
	var contentJsonBytes []byte
	contentJsonBytes, err = json.Marshal(content)
	if err == nil {
		contentJson = string(contentJsonBytes)
	}
	return contentJson, err
}

// Gets the JWT from the bearer-token.json file if it exists, errors if the file does not exist or if the token is invalid
func (file *BearerTokenFileImpl) ReadJwt() (string, error) {
	var err error = nil
	var bearerToken string = ""
	var bearerTokenJsonContents string = ""

	bearerTokenFilePath := filepath.Join(file.galasaHome.GetNativeFolderPath(), "bearer-token.json")

	log.Printf("Retrieving bearer token from file '%s'", bearerTokenFilePath)
	bearerTokenJsonContents, err = file.fileSystem.ReadTextFile(bearerTokenFilePath)
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
