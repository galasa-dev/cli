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
	jwtCache JwtCache,
) utils.Authenticator {

	authenticator := new(authenticatorImpl)

	authenticator.apiServerUrl = apiServerUrl
	authenticator.timeService = timeService
	authenticator.galasaHome = galasaHome
	authenticator.fileSystem = fileSystem
	authenticator.env = env

	authenticator.cache = jwtCache

	return authenticator
}

func (authenticator *authenticatorImpl) GetBearerToken() (string, error) {

	var bearerToken string
	var err error
	var galasaTokenValue string

	_, galasaTokenValue, err = GetAuthProperties(authenticator.fileSystem, authenticator.galasaHome, authenticator.env)
	if err == nil {
		bearerToken, err = authenticator.cache.Get(authenticator.apiServerUrl, galasaTokenValue)
		if err != nil {
			// Attempt to log in
			log.Printf("Logging in to the Galasa Ecosystem at '%s'", authenticator.apiServerUrl)
			err = authenticator.Login()
			if err == nil {
				log.Printf("Logged in to the Galasa Ecosystem at '%s' OK", authenticator.apiServerUrl)
				bearerToken, err = authenticator.cache.Get(authenticator.apiServerUrl, galasaTokenValue)
			}
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
	var err error
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
