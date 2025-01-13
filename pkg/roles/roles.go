/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package roles

import (
	"strings"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

func validateRoleName(nameToValidate string) (string, error) {
	var err error
	name := strings.TrimSpace(nameToValidate)

	if name == "" || !usesValidRoleCharacters(name) {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_ROLE_NAME)
	}
	return name, err
}

// Checks if a given string contains only characters in the Latin-1 character set (codepoints 0-255),
// returning true if so, and false otherwise
func usesValidRoleCharacters(str string) bool {
	isValidLatin1 := true
	for _, character := range str {
		if !((character >= 'a' && character <= 'z') ||
			(character >= 'A' && character <= 'Z') ||
			(character >= '0' && character <= '9')) {
			isValidLatin1 = false
			break
		}
	}
	return isValidLatin1
}
