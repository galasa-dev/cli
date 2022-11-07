/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var streams = []string{"bsf", "prod"}

func TestValidateStreamFound(t *testing.T) {
	actual := ValidateStream(streams, "prod")
	assert.Nil(t, actual, "Validate didn't find the correct stream!")
}

func TestValidateStreamNotFoundNormal(t *testing.T) {
	err := ValidateStream(streams, "apple")
	expected := "Stream \"apple\" is not found in the ecosystem. Valid streams are: bsf prod"
	if assert.NotNil(t, err, "ValidateStream found a stream where it should have returned an error message.") {
		errorString := err.Error()
		assert.Equal(t, expected, errorString, "Validate didn't fail with the correct response.")
	}
}

func TestValidateStreamNotFoundNoStreamsToFind(t *testing.T) {
	emptyStreams := []string{}
	err := ValidateStream(emptyStreams, "orange")
	expected := "Stream \"orange\" is not found in the ecosystem. There are no streams set up."
	if assert.NotNil(t, err) {
		errorString := err.Error()
		assert.Equal(t, expected, errorString, "Validate didn't fail with the correct response.")
	}
}
