/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package api

import (
	"log"
	"net/http"
	"time"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/spi"
)

var (
	RATE_LIMIT_STATUS_CODES_MAP = map[int]struct{}{
		http.StatusServiceUnavailable: {},
		http.StatusTooManyRequests: {},
	}
)

type CommsRetrier struct {
	maxAttempts int
	retryBackoffSeconds float64
	timeService spi.TimeService
}

func NewCommsRetrier(maxAttempts int, retryBackoffSeconds float64, timeService spi.TimeService) *CommsRetrier {
	return &CommsRetrier{
		maxAttempts: maxAttempts,
		retryBackoffSeconds: retryBackoffSeconds,
		timeService: timeService,
	}
}

// ExecuteCommandWithRateLimitRetries keeps trying until we've tried enough, it worked, 
// or it's failed too many times with rate limit issues.
func (retrier *CommsRetrier) ExecuteCommandWithRateLimitRetries(
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
