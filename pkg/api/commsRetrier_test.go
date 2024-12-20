/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package api

import (
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)


func createMockErrorMessageType() *galasaErrors.MessageType {
	return galasaErrors.NewMessageType("TEST123: simulating a failure on attempt %v", 123, false)
}

func TestExecuteCommandWithRateLimitRetriesOnlyRunsOnceOnSuccess(t *testing.T) {
    // Given...
    maxAttempts := 3
    retryBackoffSeconds := 1
	
	now := time.Now()
    timeService := utils.NewOverridableMockTimeService(now)
	commsRetrier := NewCommsRetrier(maxAttempts, float64(retryBackoffSeconds), timeService)
    runCounter := 0
    executionFunc := func() error {
        runCounter++
        return nil
    }

    // When...
    err := commsRetrier.ExecuteCommandWithRateLimitRetries(executionFunc)

    // Then...
    assert.Nil(t, err)
    assert.Equal(t, 1, runCounter, "The execution function should only have been run once")
	assert.Equal(t, now, timeService.Now(), "Time should not have advanced")
}

func TestExecuteCommandWithRateLimitRetriesTriesAgainOnRateLimitFailure(t *testing.T) {
    // Given...
    maxAttempts := 3
    retryBackoffSeconds := 10

	now := time.Now()
    timeService := utils.NewOverridableMockTimeService(now)
	commsRetrier := NewCommsRetrier(maxAttempts, float64(retryBackoffSeconds), timeService)
    attemptCounter := 0
    executionFunc := func() error {
        var err error
        attemptCounter++
        if attemptCounter != 2 {
            err = galasaErrors.NewGalasaErrorWithHttpStatusCode(http.StatusTooManyRequests, createMockErrorMessageType(), attemptCounter)
        }
        return err
    }

    // When...
    err := commsRetrier.ExecuteCommandWithRateLimitRetries(executionFunc)

    // Then...
    assert.Nil(t, err)
    assert.Equal(t, 2, attemptCounter, "The execution function should have been run twice")
	assert.Equal(t, now.Add(10 * time.Second), timeService.Now(), "Time should have advanced after each attempt")
}

func TestExecuteCommandWithRateLimitRetriesGivesUpAfterMaxAttempts(t *testing.T) {
    // Given...
    maxAttempts := 4
    retryBackoffSeconds := 10

	now := time.Now()
    timeService := utils.NewOverridableMockTimeService(now)
	commsRetrier := NewCommsRetrier(maxAttempts, float64(retryBackoffSeconds), timeService)
    attemptCounter := 0
    executionFunc := func() error {
        attemptCounter++
        return galasaErrors.NewGalasaErrorWithHttpStatusCode(http.StatusTooManyRequests, createMockErrorMessageType(), attemptCounter)
    }

    // When...
    err := commsRetrier.ExecuteCommandWithRateLimitRetries(executionFunc)

    // Then...
    assert.NotNil(t, err)
    assert.ErrorContains(t, err, strconv.Itoa(maxAttempts), "The last error should have been returned")
    assert.Equal(t, maxAttempts, attemptCounter, "The execution function should have been run the maximum number of times")
	assert.Equal(t, now.Add(30 * time.Second), timeService.Now(), "Time should have advanced after each attempt")
}

func TestExecuteCommandWithRateLimitRetriesRunsOnceOnNonRateLimitedFailure(t *testing.T) {
    // Given...
    maxAttempts := 3
    retryBackoffSeconds := 1

	now := time.Now()
    timeService := utils.NewOverridableMockTimeService(now)
	commsRetrier := NewCommsRetrier(maxAttempts, float64(retryBackoffSeconds), timeService)
    attemptCounter := 0
    errorMsg := "simulating an error that is not related to the API server response"
    executionFunc := func() error {
        attemptCounter++
        return errors.New(errorMsg)
    }

    // When...
    err := commsRetrier.ExecuteCommandWithRateLimitRetries(executionFunc)

    // Then...
    assert.NotNil(t, err)
    assert.ErrorContains(t, err, errorMsg)
    assert.Equal(t, 1, attemptCounter, "The execution function should only have been run once")
	assert.Equal(t, now, timeService.Now(), "Time should not have advanced")
}
