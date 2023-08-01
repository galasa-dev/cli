/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"errors"
	"testing"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCanGetNormalExitCodeAndErrorTextFromAnError(t *testing.T) {
	var err = errors.New("my text")
	errorText, exitCode, isStackTraceWanted := extractErrorDetails(err)
	assert.Equal(t, errorText, "my text", "Failed to extract the exit text from an error!")
	assert.Equal(t, 1, exitCode, "Wrong default exit code")
	assert.True(t, isStackTraceWanted, "We want stack trace from non-galasa errors")
}

func TestCanGetNormalExitCodeAndErrorTextFromAGalasaError(t *testing.T) {
	err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_MISSING_ACTION_FLAGS)
	errorText, exitCode, isStackTraceWanted := extractErrorDetails(err)
	assert.Contains(t, errorText, "GAL1009E", "Failed to extract the exit text from a galasa error!")
	assert.Equal(t, 1, exitCode, "Wrong default exit code")
	assert.False(t, isStackTraceWanted, "We don't want stack trace from galasa errors")
}

func TestCanGetTestsFailedExitCodeAndErrorTextFromATestFailedGalasaErrorPointer(t *testing.T) {
	err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_TESTS_FAILED, 14)
	errorText, exitCode, isStackTraceWanted := extractErrorDetails(err)
	assert.Contains(t, errorText, "GAL1017E", "Failed to extract the exit text from a galasa error!")
	assert.Equal(t, 2, exitCode, "Wrong default exit code")
	assert.False(t, isStackTraceWanted, "We don't want stack trace from galasa errors")
}
