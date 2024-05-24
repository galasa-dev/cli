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
}

func NewMockAuthenticator() spi.Authenticator {

	authenticator := new(MockAuthenticator)

	return authenticator
}

func (authenticator *MockAuthenticator) GetBearerToken() (string, error) {
	bearerToken := ""
	var err error

	return bearerToken, err
}

// Gets a new authenticated API client, attempting to log in if a bearer token file does not exist
func (authenticator *MockAuthenticator) GetAuthenticatedAPIClient() (*galasaapi.APIClient, error) {
	var apiClient *galasaapi.APIClient = nil
	var err error
	return apiClient, err
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
