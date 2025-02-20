/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package api

import (
	"log"
	"time"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

type APICommsClientImpl struct {
	maxAttempts         int
	retryBackoffSeconds float64
	timeService         spi.TimeService

	bootstrapData *BootstrapData
	apiClient     *galasaapi.APIClient
	authenticator spi.Authenticator
}

type APICommsClient interface {
	RunCommandWithRateLimitRetries(commandExecutionFunc func() error) error
	RunAuthenticatedCommandWithRateLimitRetries(commandExecutionFunc func(apiClient *galasaapi.APIClient) error) error

	GetBootstrapData() *BootstrapData
	GetBearerToken() (string, error)
}

func NewAPICommsClient(
	bootstrap string,
	maxAttempts int,
	retryBackoffSeconds float64,
	factory spi.Factory,
	galasaHome spi.GalasaHome,
) (APICommsClient, error) {
	var err error
	commsClient := &APICommsClientImpl{
		maxAttempts:         maxAttempts,
		retryBackoffSeconds: retryBackoffSeconds,
		timeService:         factory.GetTimeService(),
	}

	env := factory.GetEnvironment()
	fileSystem := factory.GetFileSystem()

	var bootstrapData *BootstrapData
	err = commsClient.RunCommandWithRateLimitRetries(func() error {
		// Read the bootstrap properties, retrying if a rate limit has been exceeded
		var urlService *RealUrlResolutionService = new(RealUrlResolutionService)
		bootstrapData, err = LoadBootstrap(galasaHome, fileSystem, env, bootstrap, urlService)
		return err
	})

	if err == nil {
		apiServerUrl := bootstrapData.ApiServerURL
		log.Printf("The API server is at '%s'\n", apiServerUrl)

		authenticator := factory.GetAuthenticator(
			apiServerUrl,
			galasaHome,
		)

		commsClient.bootstrapData = bootstrapData
		commsClient.authenticator = authenticator
	}

	return commsClient, err
}

func (commsClient *APICommsClientImpl) GetBootstrapData() *BootstrapData {
	return commsClient.bootstrapData
}

func (commsClient *APICommsClientImpl) GetBearerToken() (string, error) {
	return commsClient.authenticator.GetBearerToken()
}

// RunAuthenticatedCommandWithRateLimitRetries tries to run a given execution function until we've tried enough, it worked,
// or it has failed too many times with rate limit or auth issues. If an unauthorized error is encountered,
// then this function will attempt to re-authenticate with the API server.
func (commsClient *APICommsClientImpl) RunAuthenticatedCommandWithRateLimitRetries(
	commandExecutionFunc func(apiClient *galasaapi.APIClient) error,
) error {
	var err error
	isDone := false
	attempt := 0
	maxAttempts := commsClient.maxAttempts
	retryBackoffSeconds := commsClient.retryBackoffSeconds
	timeService := commsClient.timeService

	for !isDone {

		var apiClientErr error
		if commsClient.apiClient == nil {
			commsClient.apiClient, apiClientErr = commsClient.authenticator.GetAuthenticatedAPIClient()
		}

		if apiClientErr == nil {
			err = commandExecutionFunc(commsClient.apiClient)
		} else {
			err = apiClientErr
		}

		isDone = true
		if err != nil {

			isRetryRequired := false

			// Try to convert the error received from the command into an API error
			galasaError, isGalasaError := err.(galasaErrors.GalasaCommsError)
			if isGalasaError {

				// If the command encountered an unauthorized error from the API server,
				// attempt to log in again to get a new JWT and use that in subsequent requests
				if galasaError.IsReauthRequired() {
					attempt++
					log.Printf("Reauthentication required. Login attempt %v/%v", attempt, maxAttempts)

					// Overwrite the API client being used to avoid having to re-authenticate again
					commsClient.apiClient, apiClientErr = commsClient.authenticator.GetAuthenticatedAPIClient()

					if apiClientErr == nil {
						isRetryRequired = true
					} else {
						err = apiClientErr
					}
				} else if galasaError.IsRateLimitedRetryRequired() {
					attempt++
					log.Printf("Rate limit exceeded on attempt %v/%v", attempt, maxAttempts)
					isRetryRequired = true
				}

				if isRetryRequired && attempt < maxAttempts {
					log.Printf("Retrying in %v second(s)", retryBackoffSeconds)
					timeService.Sleep(time.Duration(retryBackoffSeconds) * time.Second)
					isDone = false
				}
			}
		}
	}
	return err
}

// RunCommandWithRateLimitRetries keeps trying until we've tried enough, it worked,
// or it's failed too many times with rate limit issues.
func (commsClient *APICommsClientImpl) RunCommandWithRateLimitRetries(
	commandExecutionFunc func() error,
) error {
	var err error
	isDone := false
	attempt := 0
	maxAttempts := commsClient.maxAttempts
	retryBackoffSeconds := commsClient.retryBackoffSeconds
	timeService := commsClient.timeService

	for !isDone {

		err = commandExecutionFunc()

		isDone = true
		if err != nil {

			// Try to convert the error received from the command into an API error
			galasaError, isGalasaError := err.(galasaErrors.GalasaCommsError)
			if isGalasaError && galasaError.IsRateLimitedRetryRequired() {
				attempt++
				log.Printf("Rate limit exceeded on attempt %v/%v", attempt, maxAttempts)

				if attempt < maxAttempts {
					log.Printf("Retrying in %v second(s)", retryBackoffSeconds)
					timeService.Sleep(time.Duration(retryBackoffSeconds) * time.Second)
					isDone = false
				}
			}
		}
	}
	return err
}
