/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package streams

import (
	"regexp"
	"strings"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

func validateStreamName(streamName string) (string, error) {

	var err error
	streamName = strings.TrimSpace(streamName)

	if streamName == "" {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_MISSING_STREAM_NAME_FLAG)
	} else {
		// Check that stream name only contains alphanumeric characters, underscores, and hyphens
        validChars := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
        if !validChars.MatchString(streamName) {
            err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_STREAM_NAME)
        }
	}

	return streamName, err

}
