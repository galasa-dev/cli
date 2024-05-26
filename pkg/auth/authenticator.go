/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

type authenticatorImpl struct {
	apiServerUrl string
	fileSystem   spi.FileSystem
	galasaHome   spi.GalasaHome
	timeService  spi.TimeService
	env          spi.Environment
	cache        JwtCache
}

func NewAuthenticator(
	apiServerUrl string,
	fileSystem spi.FileSystem,
	galasaHome spi.GalasaHome,
	timeService spi.TimeService,
	env spi.Environment,
	jwtCache JwtCache,
) spi.Authenticator {

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

	log.Printf("GetBearerToken entered.\n")

	_, galasaTokenValue, err = getAuthProperties(authenticator.fileSystem, authenticator.galasaHome, authenticator.env)
	if err == nil {
		bearerToken, err = authenticator.cache.Get(authenticator.apiServerUrl, galasaTokenValue)
		if err == nil {
			if bearerToken == "" {
				// Attempt to log in
				log.Printf("Logging in to the Galasa Ecosystem at '%s'", authenticator.apiServerUrl)
				err = authenticator.Login()
				if err == nil {
					log.Printf("Logged in to the Galasa Ecosystem at '%s' OK", authenticator.apiServerUrl)
					bearerToken, err = authenticator.cache.Get(authenticator.apiServerUrl, galasaTokenValue)
				}
			}
		}
	}

	log.Printf("GetBearerToken exiting. length of bearerToken: %v err:%v\n", len(bearerToken), err)

	return bearerToken, err
}

// Gets a new authenticated API client, attempting to log in if a bearer token file does not exist
func (authenticator *authenticatorImpl) GetAuthenticatedAPIClient() (*galasaapi.APIClient, error) {

	log.Printf("GetAuthenticatedAPIClient entered.\n")

	bearerToken, err := authenticator.GetBearerToken()
	var apiClient *galasaapi.APIClient
	if err == nil {
		apiClient = api.InitialiseAuthenticatedAPI(authenticator.apiServerUrl, bearerToken)
	}

	log.Printf("GetAuthenticatedAPIClient exiting. err: %v\n", err)

	return apiClient, err
}

// Login - performs all the logic to implement the `galasactl auth login` command
func (authenticator *authenticatorImpl) Login() error {
	var err error
	var authProperties galasaapi.AuthProperties
	var galasaTokenValue string

	log.Printf("Login entered.\n")

	authProperties, galasaTokenValue, err = getAuthProperties(authenticator.fileSystem, authenticator.galasaHome, authenticator.env)
	if err == nil {
		var jwt string
		jwt, err = GetJwtFromRestApi(authenticator.apiServerUrl, authProperties)
		if err == nil {
			err = authenticator.cache.Put(authenticator.apiServerUrl, galasaTokenValue, jwt)
		}
	}

	log.Printf("Login exiting. %v\n", err)

	return err
}

// Logout - performs all the logout to implement the `galasactl auth login` command
func (authenticator *authenticatorImpl) LogoutOfEverywhere() error {
	var err error

	log.Printf("LogoutOfEverywhere entered.")

	authenticator.cache.ClearAll()
	return err
}
