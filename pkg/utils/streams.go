/*
 * Copyright contributors to the Galasa project
 */

package utils

import (
	"strings"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/galasaapi"
)

func FetchTestStreams(apiClient *galasaapi.APIClient) []string {
	cpsProperty, _, err := apiClient.ConfigurationPropertyStoreAPIApi.GetCpsNamespaceCascadeProperty(nil, "framework", "test", "streams").Execute()
	if err != nil {
		panic(err)
	}

	if cpsProperty.Value == nil {
		return make([]string, 0)
	}

	return strings.Split(*cpsProperty.Value, ",")
}

func ValidateStream(streams []string, stream string) error {
	for _, s := range streams {
		if s == stream {
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
