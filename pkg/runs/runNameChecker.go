/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"regexp"

	"github.com/galasa-dev/cli/pkg/errors"
)

var (
	// Expect the pattern letters + number
	// ^ matches the start of the string
	// $ matches the end of the string
	// If ^ and $ are not specified, then M2M2 would match (two instances of the match)!
	RUN_NAME_PATTERN *regexp.Regexp = regexp.MustCompile("^[a-zA-Z]+[0-9]+$")
)

// ---------------------------------------------------
// Functions called by other things

// ValidateFlagValue - Checks that a flag value is valid, as much as we can.
// Returns an error if it's invalid, nil if it looks valid.
// This function does not consult with an ecosystem, just checks the
// format of the runName.
func ValidateRunName(value string) error {

	var err error

	isMatching := RUN_NAME_PATTERN.MatchString(value)

	if !isMatching {
		err = errors.NewGalasaError(errors.GALASA_ERROR_INVALID_FLAG_VALUE, value)
	}

	return err
}
