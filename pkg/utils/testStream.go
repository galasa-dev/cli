//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package utils

import (
	"errors"
	"strings"

	"github.com/galasa.dev/cli/pkg/galasaapi"
)

func FetchTestStreams(apiClient *galasaapi.APIClient) []string {
	cpsProperty, _, err := apiClient.ConfigurationPropertyStoreAPIApi.GetCpsNamespaceCascadeProperty(nil, "framework", "test", "streams").Execute()
	if (err != nil) {
		panic(err)
	}

	if cpsProperty.Value == nil {
		return make([]string, 0)
	}

	return strings.Split(*cpsProperty.Value, ",")
}

func ValidateStream(streams []string, stream string) (error) {
    for _, s := range streams {
        if s == stream {
            return nil
        }
    }

    return errors.New("Stream \"" + stream + "\" is missing from ecosystem")
}