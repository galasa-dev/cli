/*
 * Copyright contributors to the Galasa project
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
	errorText, exitCode := extractErrorDetails(err)
	assert.Equal(t, errorText, "my text", "Failed to extract the exit text from an error!")
	assert.Equal(t, 1, exitCode, "Wrong default exit code")
}

func TestCanGetNormalExitCodeAndErrorTextFromAGalasaError(t *testing.T) {
	err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_MISSING_ACTION_FLAGS)
	errorText, exitCode := extractErrorDetails(err)
	assert.Contains(t, errorText, "GAL1009E", "Failed to extract the exit text from a galasa error!")
	assert.Equal(t, 1, exitCode, "Wrong default exit code")
}

func TestCanGetTestsFailedExitCodeAndErrorTextFromATestFailedGalasaErrorPointer(t *testing.T) {
	err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_TESTS_FAILED, 14)
	errorText, exitCode := extractErrorDetails(err)
	assert.Contains(t, errorText, "GAL1017E", "Failed to extract the exit text from a galasa error!")
	assert.Equal(t, 2, exitCode, "Wrong default exit code")
}
