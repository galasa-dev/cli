/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
)

type Authenticator interface {
	// Gets a bearer token from the persistent cache if there is one, else logs into the server to get one.
	GetBearerToken() (string, error)

	// Gets a new authenticated API client, attempting to log in if a bearer token file does not exist
	GetAuthenticatedAPIClient() (*galasaapi.APIClient, error)

	// Logs into the server, saving the JWT token obtained in a persistent cache for later
	Login() error
}

type authenticatorImpl struct {
	apiServerUrl string
	fileSystem   files.FileSystem
	galasaHome   utils.GalasaHome
	timeService  utils.TimeService
	env          utils.Environment
	cache        JwtCache
}

func NewAuthenticator(
	apiServerUrl string,
	fileSystem files.FileSystem,
	galasaHome utils.GalasaHome,
	timeService utils.TimeService,
	env utils.Environment,
) Authenticator {

	authenticator := new(authenticatorImpl)

	authenticator.apiServerUrl = apiServerUrl
	authenticator.timeService = timeService
	authenticator.galasaHome = galasaHome
	authenticator.fileSystem = fileSystem
	authenticator.env = env

	authenticator.cache = NewJwtCache(fileSystem, galasaHome, timeService)

	return authenticator
}

func (authenticator *authenticatorImpl) GetBearerToken() (string, error) {

	bearerToken, err := utils.GetBearerTokenFromTokenJsonFile(authenticator.fileSystem, authenticator.galasaHome, authenticator.timeService)
	if err != nil {
		// Attempt to log in
		log.Printf("Logging in to the Galasa Ecosystem at '%s'", authenticator.apiServerUrl)
		err = authenticator.Login()
		if err == nil {
			log.Printf("Logged in to the Galasa Ecosystem at '%s' OK", authenticator.apiServerUrl)
			bearerToken, err = utils.GetBearerTokenFromTokenJsonFile(authenticator.fileSystem, authenticator.galasaHome, authenticator.timeService)
		}
	}
	return bearerToken, err
}

// Gets a new authenticated API client, attempting to log in if a bearer token file does not exist
func (authenticator *authenticatorImpl) GetAuthenticatedAPIClient() (*galasaapi.APIClient, error) {
	bearerToken, err := authenticator.GetBearerToken()
	var apiClient *galasaapi.APIClient
	if err == nil {
		apiClient = api.InitialiseAuthenticatedAPI(authenticator.apiServerUrl, bearerToken)
	}
	return apiClient, err
}

// Login - performs all the logic to implement the `galasactl auth login` command
func (authenticator *authenticatorImpl) Login() error {
	var err error = nil
	var authProperties galasaapi.AuthProperties
	var galasaTokenValue string
	authProperties, galasaTokenValue, err = GetAuthProperties(authenticator.fileSystem, authenticator.galasaHome, authenticator.env)
	if err == nil {
		var jwt string
		jwt, err = GetJwtFromRestApi(authenticator.apiServerUrl, authProperties)
		if err == nil {
			err = authenticator.cache.Put(authenticator.apiServerUrl, galasaTokenValue, jwt)
		}
	}
	return err
}
