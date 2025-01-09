/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

// Checks if a given string contains only characters in the Latin-1 character set (codepoints 0-255),
// returning true if so, and false otherwise
func IsLatin1(str string) bool {
	isValidLatin1 := true
	for _, character := range str {
		if character > 255 {
			isValidLatin1 = false
			break
		}
	}
	return isValidLatin1
}
