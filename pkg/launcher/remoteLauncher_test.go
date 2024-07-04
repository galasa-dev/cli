/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package launcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/stretchr/testify/assert"
)

func TestProcessingGoodPropertiesExtractsStreamsOk(t *testing.T) {

	var inputProperties []galasaapi.GalasaProperty = make([]galasaapi.GalasaProperty, 0)

	name1 := "thames"
	name1full := "test.stream." + name1 + ".repo"
	name2 := "avon"
	name2full := "test.stream." + name2 + ".repo"

	inputProperties = append(inputProperties, galasaapi.GalasaProperty{
		Metadata: &galasaapi.GalasaPropertyMetadata{
			Name: &name1full,
		},
	})

	inputProperties = append(inputProperties, galasaapi.GalasaProperty{
		Metadata: &galasaapi.GalasaPropertyMetadata{
			Name: &name2full,
		},
	})

	streams, err := getStreamNamesFromProperties(inputProperties)
	assert.Nil(t, err)
	assert.NotNil(t, streams)
	assert.Equal(t, 2, len(streams))

	assert.Equal(t, streams[0], name1)
	assert.Equal(t, streams[1], name2)
}

func TestProcessingEmptyPropertiesListExtractsZeroStreamsOk(t *testing.T) {

	var inputProperties []galasaapi.GalasaProperty = make([]galasaapi.GalasaProperty, 0)

	streams, err := getStreamNamesFromProperties(inputProperties)

	assert.Nil(t, err)
	assert.NotNil(t, streams)
	assert.Equal(t, 0, len(streams))
}

func TestGetTestCatalogHttpErrorGetsReported(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("URL arrived at the mock test server: %s\n", r.RequestURI)
		switch r.RequestURI {
		case "/cps/framework/properties?prefix=test.stream.myStream&suffix=location":

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			name := "mycpsPropName"
			value := "a duff value" // This is intentionally duff, which will cause an HTTP error when the production code tries to GET using this as a URL.
			payload := []galasaapi.GalasaProperty{
				{
					Metadata: &galasaapi.GalasaPropertyMetadata{
						Name: &name,
					},
					Data: &galasaapi.GalasaPropertyData{
						Value: &value,
					},
				},
			}
			payloadBytes, _ := json.Marshal(payload)
			w.Write(payloadBytes)

			fmt.Printf("mock server sending payload: %s\n", string(payloadBytes))

		}
	}))
	defer server.Close()

	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)

	launcher := NewRemoteLauncher(apiServerUrl, apiClient)

	_, err := launcher.GetTestCatalog("myStream")

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "GAL1144E") // Failed to get the test catalog.
}
