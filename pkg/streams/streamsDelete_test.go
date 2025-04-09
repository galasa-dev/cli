/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package streams

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func createMockStream(name string, description string) galasaapi.Stream {

	stream := *galasaapi.NewStream()
	streamMetadata := *galasaapi.NewStreamMetadata()

	streamMetadata.SetName(name)
	streamMetadata.SetDescription(description)

	stream.SetMetadata(streamMetadata)

	return stream

}

func WriteMockStreamResponse(
	t *testing.T,
	writer http.ResponseWriter,
	req *http.Request,
	name string,
	streamResultStrings []string) {

	writer.Header().Set("Content-Type", "application/json")
	values := req.URL.Path
	path := strings.Split(values, "/")
	streamPathVar := path[2]
	assert.Equal(t, streamPathVar, name)

	writer.Write([]byte(fmt.Sprintf(`
	{
		"metadata":{
			"name": "%s",
			"description": "This is a dummy stream"
		}
	}`, name)))

}

func TestStreamUserDeleteAUser(t *testing.T) {

	//Given...
	name := "mystream"
	description := "This is a dummy stream"

	streamToDelete := createMockStream(name, description)
	streamToDeleteBytes, _ := json.Marshal(streamToDelete)
	streamToDeleteJson := string(streamToDeleteBytes)

	getStreamInteraction := utils.NewHttpInteraction("/streams/"+name, http.MethodGet)
	getStreamInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
		WriteMockStreamResponse(t, writer, req, name, []string{streamToDeleteJson})
	}

	deleteStreamInteraction := utils.NewHttpInteraction("/streams/"+name, http.MethodDelete)
	deleteStreamInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	}

	interactions := []utils.HttpInteraction{
		getStreamInteraction,
		deleteStreamInteraction,
	}

	server := utils.NewMockHttpServer(t, interactions)
	defer server.Server.Close()

	console := utils.NewMockConsole()
	apiServerUrl := server.Server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockByteReader := utils.NewMockByteReader()

	// When...
	err := DeleteStream(
		name,
		apiClient,
		mockByteReader)

	// Then...
	assert.Nil(t, err, "DeleteStream returned an unexpected error")
	assert.Empty(t, console.ReadText(), "The console was written to on a successful deletion, it should be empty")
}

func TestStreamDeleteThrowsAnUnexpectedError(t *testing.T) {

	//Given...
	name := "mystream"
	description := "This is a dummy stream"

	streamToDelete := createMockStream(name, description)
	streamToDeleteBytes, _ := json.Marshal(streamToDelete)
	streamToDeleteJson := string(streamToDeleteBytes)

	getStreamInteraction := utils.NewHttpInteraction("/streams/"+name, http.MethodGet)
	getStreamInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
		WriteMockStreamResponse(t, writer, req, name, []string{streamToDeleteJson})
	}

	deleteStreamInteraction := utils.NewHttpInteraction("/streams/"+name, http.MethodDelete)
	deleteStreamInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{}`))
	}

	interactions := []utils.HttpInteraction{
		getStreamInteraction,
		deleteStreamInteraction,
	}

	server := utils.NewMockHttpServer(t, interactions)
	defer server.Server.Close()

	apiServerUrl := server.Server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockByteReader := utils.NewMockByteReader()

	// When...
	err := DeleteStream(
		name,
		apiClient,
		mockByteReader)

	// Then...
	assert.NotNil(t, err, "DeleteStream returned an unexpected error")
	assert.Contains(t, err.Error(), strconv.Itoa(http.StatusInternalServerError))
	assert.Contains(t, err.Error(), "GAL1239E")
}
