/*
 * Copyright contributors to the Galasa project
 */
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

func InitialiseAPI(apiServerUrl string) *galasaapi.APIClient {
	// Calculate the bootstrap for this execution

	var apiClient *galasaapi.APIClient = nil

	cfg := galasaapi.NewConfiguration()
	cfg.Debug = false
	cfg.Servers = galasaapi.ServerConfigurations{{URL: apiServerUrl}}
	apiClient = galasaapi.NewAPIClient(cfg)

	return apiClient
}

func InitialiseAPIWithAuthHeader(apiServerUrl string, galasaHome utils.GalasaHome) (*galasaapi.APIClient, error) {
	var err error = nil
	var bearerToken string = ""

	apiClient := InitialiseAPI(apiServerUrl)
	cfg := apiClient.GetConfig()

	fileSystem := files.NewOSFileSystem()

	bearerToken, err = getBearerTokenFromTokenJsonFile(fileSystem, galasaHome)
	if err == nil {
		cfg.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}

	return apiClient, err
}

type BearerTokenJson struct {
	Jwt string `json:"jwt"`
}

// Gets the JWT from the bearer-token.json file if it exists, errors if the file does not exist
func getBearerTokenFromTokenJsonFile(fileSystem files.FileSystem, galasaHome utils.GalasaHome) (string, error) {
	var err error = nil
	var bearerToken string = ""

	bearerTokenFilePath := filepath.Join(galasaHome.GetNativeFolderPath(), "bearer-token.json")

	log.Printf("Retrieving JWT from bearer token file '%s'", bearerTokenFilePath)
	bearerTokenJsonContents, err := fileSystem.ReadTextFile(bearerTokenFilePath)
	if err == nil {
		var bearerTokenJson BearerTokenJson
		err = json.Unmarshal([]byte(bearerTokenJsonContents), &bearerTokenJson)
		if err == nil {
			bearerToken = bearerTokenJson.Jwt
			log.Printf("Retrieved JWT from bearer token file '%s' OK", bearerTokenFilePath)
		}
	}

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_BEARER_TOKEN_FROM_FILE, bearerTokenFilePath, err.Error())
	}
	return bearerToken, err
}