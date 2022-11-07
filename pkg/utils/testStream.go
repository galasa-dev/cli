/*
 * Copyright contributors to the Galasa project
 */

package utils

import (
	"errors"
	"fmt"
	"strings"

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
	var errorMsg = ""
	if len(streams) < 1 {
		template := "Stream \"%s\" is not found in the ecosystem. There are no streams set up. " +
			"Ask your Galasa system administrator to add a new stream with the desired name."
		errorMsg = fmt.Sprintf(template, stream)
	} else {
		template := "Stream \"%s\" is not found in the ecosystem. Valid streams are:%s. " +
			"Try again using a valid stream, or ask your Galasa system administrator to " +
			"add a new stream with the desired name."
		var buffer strings.Builder
		for _, s := range streams {
			buffer.WriteString(" '")
			buffer.WriteString(s)
			buffer.WriteString("'")
		}
		errorMsg = fmt.Sprintf(template, stream, buffer.String())
	}

	return errors.New(errorMsg)
}
