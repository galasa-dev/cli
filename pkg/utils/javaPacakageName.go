/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"strings"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
)

const (
	LOWER_CASE_LETTERS = "abcdefghijklmnopqrstuvwxyz"
	DIGITS             = "0123456789"
	SEPARATOR          = "."
)

var (
	validStartingCharacters = LOWER_CASE_LETTERS + DIGITS
	validMiddleCharacters   = LOWER_CASE_LETTERS + DIGITS + SEPARATOR
	validLastCharacters     = LOWER_CASE_LETTERS + DIGITS
)

// To validate the string as a valid java package name before we start to use it.
func ValidateJavaPackageName(javaPackageName string) error {
	var err error = nil

	if javaPackageName == "" {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PACKAGE_NAME_BLANK)
	}

	if err == nil {
		// Check all the middle characters
		for _, c := range javaPackageName {
			charToLookFor := string(c)
			if !strings.Contains(validMiddleCharacters, charToLookFor) {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_CHAR_IN_PACKAGE_NAME, javaPackageName, charToLookFor)
				break
			}
		}
	}

	if err == nil {
		// Check the first character.
		firstChar := string(javaPackageName[0])
		if !strings.Contains(validStartingCharacters, firstChar) {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_FIRST_CHAR_IN_PKG_NAME, javaPackageName, firstChar)
		}
	}

	if err == nil {
		// Check the last character.
		length := len(javaPackageName)
		lastChar := string(javaPackageName[length-1])
		if !strings.Contains(validLastCharacters, lastChar) {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_LAST_CHAR_IN_PKG_NAME, javaPackageName, lastChar)
		}
	}
	return err
}
