/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"log"
	"strings"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/launcher"
)

func GetStreams(launcher launcher.Launcher) ([]string, error) {
	log.Println("Getting streams list.")
	streams, err := launcher.GetStreams()
	return streams, err
}

func ValidateStream(streams []string, stream string) error {
	log.Println("Validating streams list.")
	for _, s := range streams {
		if s == stream {
			log.Println("Stream is found in the list of valid streams.")
			return nil
		}
	}

	// Build the error message.
	var error *galasaErrors.GalasaError
	if len(streams) < 1 {
		// No streams configured.
		error = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_NO_STREAMS_CONFIGURED, stream)
	} else {

		var buffer strings.Builder
		var availableStreamsList string
		for _, s := range streams {
			buffer.WriteString(" '")
			buffer.WriteString(s)
			buffer.WriteString("'")
		}
		availableStreamsList = buffer.String()
		error = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_STREAM, stream, availableStreamsList)
	}

	return error
}
