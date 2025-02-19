/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

type MockAuthenticator struct {
	apiClient *galasaapi.APIClient
}

func NewMockAuthenticator() spi.Authenticator {
	return NewMockAuthenticatorWithAPIClient(nil)
}

func NewMockAuthenticatorWithAPIClient(apiClient *galasaapi.APIClient) spi.Authenticator {

	authenticator := new(MockAuthenticator)
	authenticator.apiClient = apiClient
	return authenticator
}

func (authenticator *MockAuthenticator) GetBearerToken() (string, error) {
	bearerToken := ""
	var err error

	return bearerToken, err
}

// Gets a new authenticated API client, attempting to log in if a bearer token file does not exist
func (authenticator *MockAuthenticator) GetAuthenticatedAPIClient() (*galasaapi.APIClient, error) {
	var err error
	return authenticator.apiClient, err
}

// Login - performs all the logic to implement the `galasactl auth login` command
func (authenticator *MockAuthenticator) Login() error {
	var err error
	return err
}

// Login - performs all the logic to implement the `galasactl auth login` command
func (authenticator *MockAuthenticator) LogoutOfEverywhere() error {
	var err error
	return err
}
