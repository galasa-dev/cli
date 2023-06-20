/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------
// Functions
func TestValidResultNamePassesValidation(t *testing.T) {
	result, err := ValidateResultName("passed")
	assert.Nil(t, err)
	assert.Equal(t, "Passed", result)
}

func TestValidResultNameWithMixedCasePassesValidation(t *testing.T) {
	result, err := ValidateResultName("fAiLeD")
	assert.Nil(t, err)
	assert.Equal(t, "Failed", result)
}

func TestInvalidResultNameFailsWithError(t *testing.T) {
	checkInvalidResultNameFails(t, "garbage")
}

func checkInvalidResultNameFails(t *testing.T, resultName string) {
	_, err := ValidateResultName(resultName)
	assert.NotNil(t, err, "Should not have validated OK.")
	if err != nil {
		assert.ErrorContains(t, err, "GAL1085E")
		assert.ErrorContains(t, err, resultName)
		assert.Contains(t, err.Error(), "'Passed'")
		assert.Contains(t, err.Error(), "'Failed'")
		assert.Contains(t, err.Error(), "'EnvFail'")
		assert.Contains(t, err.Error(), "'UNKNOWN'")
	}
}
