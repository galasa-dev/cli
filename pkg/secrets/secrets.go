/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package secrets

import (
    "strings"

    galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

func validateSecretName(secretName string) (string, error) {
    var err error
    secretName = strings.TrimSpace(secretName)

    if secretName == "" || strings.ContainsAny(secretName, " .\n\t") || !isLatin1(secretName) {
        err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_SECRET_NAME)
    }
    return secretName, err
}

func validateDescription(description string) (string, error) {
    var err error
    description = strings.TrimSpace(description)

    if description == "" || !isLatin1(description) {
        err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_SECRET_DESCRIPTION)
    }
    return description, err
}

// Checks if a given string contains only characters in the Latin-1 character set (codepoints 0-255),
// returning true if so, and false otherwise
func isLatin1(str string) bool {
	isValidLatin1 := true
	for _, character := range str {
		if character > 255 {
			isValidLatin1 = false
			break
		}
	}
	return isValidLatin1
}
