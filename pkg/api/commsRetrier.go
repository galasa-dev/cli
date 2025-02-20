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

type CommsRetrierImpl struct {
	maxAttempts int
	retryBackoffSeconds float64
	timeService spi.TimeService

	apiClient *galasaapi.APIClient
	authenticator spi.Authenticator
}

type CommsRetrier interface {
	ExecuteCommandWithRateLimitRetries(commandExecutionFunc func() error) error
	ExecuteCommandWithRetries(commandExecutionFunc func(apiClient *galasaapi.APIClient) error) error
}

func NewCommsRetrier(maxAttempts int, retryBackoffSeconds float64, timeService spi.TimeService) CommsRetrier {
	return &CommsRetrierImpl{
		maxAttempts: maxAttempts,
		retryBackoffSeconds: retryBackoffSeconds,
		timeService: timeService,
	}
}

func NewCommsRetrierWithAPIClient(
	maxAttempts int,
	retryBackoffSeconds float64,
	timeService spi.TimeService,
	authenticator spi.Authenticator,
) (CommsRetrier, error) {
	var err error
	var apiClient *galasaapi.APIClient
	var commsRetrier CommsRetrier

	apiClient, err = authenticator.GetAuthenticatedAPIClient()
	if err == nil {
		commsRetrier = &CommsRetrierImpl{
			maxAttempts: maxAttempts,
			retryBackoffSeconds: retryBackoffSeconds,
			timeService: timeService,
			authenticator: authenticator,
			apiClient: apiClient,
		}
	}
	return commsRetrier, err
}

// ExecuteCommandWithRetries tries to run a given execution function until we've tried enough, it worked, 
// or it has failed too many times with rate limit or auth issues. If an unauthorized error is encountered,
// then this function will attempt to re-authenticate with the API server.
func (retrier *CommsRetrierImpl) ExecuteCommandWithRetries(
	commandExecutionFunc func(apiClient *galasaapi.APIClient) error,
) error {
	var err error
    isDone := false
	attempt := 0
	maxAttempts := retrier.maxAttempts
	retryBackoffSeconds := retrier.retryBackoffSeconds
	timeService := retrier.timeService

    for !isDone {

        err = commandExecutionFunc(retrier.apiClient)

        isDone = true
        if err != nil {

			isRetryRequired := false

            // Try to convert the error received from the command into an API error
            galasaError, isGalasaError := err.(galasaErrors.GalasaCommsError)
            if isGalasaError {

				// If the command encountered an unauthorized error from the API server,
				// attempt to log in again to get a new JWT and use that in subsequent requests
				if galasaError.IsReauthRequired() {
					log.Printf("Reauthentication required. Login attempt %v/%v", (attempt + 1), maxAttempts)
					var newApiClient *galasaapi.APIClient
					var apiClientErr error
					newApiClient, apiClientErr = retrier.authenticator.GetAuthenticatedAPIClient()

					if apiClientErr == nil {
						// Overwrite the API client being used to avoid having to re-authenticate again
						retrier.apiClient = newApiClient
						isRetryRequired = true
					}
				} else if galasaError.IsRetryRequired() {
					log.Printf("Rate limit exceeded on attempt %v/%v", (attempt + 1), maxAttempts)
					isRetryRequired = true
				}

				if isRetryRequired && (attempt + 1) < maxAttempts {
					log.Printf("Retrying in %v second(s)", retryBackoffSeconds)
					timeService.Sleep(time.Duration(retryBackoffSeconds) * time.Second)
					isDone = false
					attempt++
				}
            }
        }
    }
    return err
}

// ExecuteCommandWithRateLimitRetries keeps trying until we've tried enough, it worked, 
// or it's failed too many times with rate limit issues.
func (retrier *CommsRetrierImpl) ExecuteCommandWithRateLimitRetries(
    commandExecutionFunc func() error,
) error {
    var err error
    isDone := false
	attempt := 0
	maxAttempts := retrier.maxAttempts
	retryBackoffSeconds := retrier.retryBackoffSeconds
	timeService := retrier.timeService

    for !isDone {

        err = commandExecutionFunc()

        isDone = true
        if err != nil {

            // Try to convert the error received from the command into an API error
            galasaError, isGalasaError := err.(galasaErrors.GalasaCommsError)
            if isGalasaError && galasaError.IsRetryRequired() {
				log.Printf("Rate limit exceeded on attempt %v/%v", (attempt + 1), maxAttempts)
				
				if (attempt + 1) < maxAttempts {
					log.Printf("Retrying in %v second(s)", retryBackoffSeconds)
					timeService.Sleep(time.Duration(retryBackoffSeconds) * time.Second)
					isDone = false
					attempt++
				}
            }
        }
    }
    return err
}
