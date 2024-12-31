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

func validateRoleName(name string) (string, error) {
	var err error
	name = strings.TrimSpace(name)

	if name == "" || strings.ContainsAny(name, " .\n\t") || !isLatin1(name) {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_ROLE_NAME)
	}
	return name, err
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
