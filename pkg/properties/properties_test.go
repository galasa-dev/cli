/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidNamespaceFormatReturnsOk(t *testing.T) {
	//Given
	validNamespaceFormat := "framework"

	//When
	err := validateNamespaceFormat(validNamespaceFormat)

	//Then
	assert.Nil(t, err)
}

func TestValidNamespaceFormatWithNumbersReturnsOk(t *testing.T) {
	//Given
	validNamespaceFormat := "fra4mework5"

	//When
	err := validateNamespaceFormat(validNamespaceFormat)

	//Then
	assert.Nil(t, err)
}

func TestInvalidNamespaceFormatStartingWithNumberReturnsError(t *testing.T) {
	//Given
	invalidNamespaceFormat := "1framework"

	//When
	err := validateNamespaceFormat(invalidNamespaceFormat)

	//Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1138")
}

func TestInvalidNamespaceFormatWithSpecialCharacterReturnsError(t *testing.T) {
	//Given
	invalidNamespaceFormat := "frame-work"

	//When
	err := validateNamespaceFormat(invalidNamespaceFormat)

	//Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1138")
}

func TestInvalidNamespaceFormatWithNumbersAndSpecialCharacterReturnsError(t *testing.T) {
	//Given
	invalidNamespaceFormat := "fr8amework-"

	//When
	err := validateNamespaceFormat(invalidNamespaceFormat)

	//Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1138")
}

func TestInvalidNamespaceFormatWithStartingCapitalLetterReturnsError(t *testing.T) {
	//Given
	invalidNamespaceFormat := "Framework"

	//When
	err := validateNamespaceFormat(invalidNamespaceFormat)

	//Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1138")
}

func TestInvalidNamespaceFormatWithCapitalLetterReturnsError(t *testing.T) {
	//Given
	invalidNamespaceFormat := "frameWork"

	//When
	err := validateNamespaceFormat(invalidNamespaceFormat)

	//Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1138")
}
