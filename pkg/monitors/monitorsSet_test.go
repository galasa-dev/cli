/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package monitors

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func readMonitorRequestBody(req *http.Request) galasaapi.UpdateGalasaMonitorRequest {
    var monitorUpdateRequest galasaapi.UpdateGalasaMonitorRequest
    requestBodyBytes, _ := io.ReadAll(req.Body)
    defer req.Body.Close()

    _ = json.Unmarshal(requestBodyBytes, &monitorUpdateRequest)
    return monitorUpdateRequest
}

func TestCanEnableAMonitor(t *testing.T) {
    // Given...
    monitorName := "customManagerCleanup"
    isEnabled := "true"

    // Create the expected HTTP interactions with the API server
    putMonitorInteraction := utils.NewHttpInteraction("/monitors/" + monitorName, http.MethodPut)
    putMonitorInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        requestBody := readMonitorRequestBody(req)
        assert.Equal(t, requestBody.Data.GetIsEnabled(), true)

        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusNoContent)
    }

    interactions := []utils.HttpInteraction{
        putMonitorInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetMonitor(
        monitorName,
        isEnabled,
        apiClient,
        mockByteReader)

    // Then...
    assert.Nil(t, err, "EnableMonitor returned an unexpected error")
}

func TestEnableMonitorWithBlankNameDisplaysError(t *testing.T) {
    // Given...
    monitorName := "    "
    isEnabled := "true"

    // The client-side validation should fail, so no HTTP interactions will be performed
    interactions := []utils.HttpInteraction{}

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetMonitor(
        monitorName,
        isEnabled,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "EnableMonitor did not return an error as expected")
    consoleOutputText := err.Error()
    assert.Contains(t, consoleOutputText, "GAL1225E")
    assert.Contains(t, consoleOutputText, " Invalid monitor name provided")
}

func TestEnableNonExistantMonitorDisplaysError(t *testing.T) {
    // Given...
    nonExistantMonitor := "monitorDoesNotExist123"
    isEnabled := "true"

    // Create the expected HTTP interactions with the API server
    enableMonitorInteraction := utils.NewHttpInteraction("/monitors/" + nonExistantMonitor, http.MethodPut)
    enableMonitorInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusNotFound)
        writer.Write([]byte(`{ "error_message": "No such monitor exists" }`))
    }

    interactions := []utils.HttpInteraction{ enableMonitorInteraction }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetMonitor(
        nonExistantMonitor,
        isEnabled,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "EnableMonitor did not return an error but it should have")
    assert.ErrorContains(t, err, "GAL1222E")
}

func TestEnableMonitorFailsWithNoExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    monitorName := "myMonitor"
    isEnabled := "true"

    // Create the expected HTTP interactions with the API server
    enableMonitorInteraction := utils.NewHttpInteraction("/monitors/" + monitorName, http.MethodPut)
    enableMonitorInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusInternalServerError)
    }

    interactions := []utils.HttpInteraction{
        enableMonitorInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetMonitor(
        monitorName,
        isEnabled,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "EnableMonitor did not return an error but it should have")
    assert.ErrorContains(t, err, "GAL1219E")
}

func TestEnableMonitorFailsWithNonJsonContentTypeExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    monitorName := "myMonitor"
    isEnabled := "true"

    // Create the expected HTTP interactions with the API server
    enableMonitorInteraction := utils.NewHttpInteraction("/monitors/" + monitorName, http.MethodPut)
    enableMonitorInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Header().Set("Content-Type", "application/notJsonOnPurpose")
        writer.Write([]byte("something not json but non-zero-length."))
    }

    interactions := []utils.HttpInteraction{
        enableMonitorInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetMonitor(
        monitorName,
        isEnabled,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "EnableMonitor did not return an error but it should have")
    assert.ErrorContains(t, err, strconv.Itoa(http.StatusInternalServerError))
    assert.ErrorContains(t, err, "GAL1223E")
    assert.ErrorContains(t, err, "Error details from the server are not in the json format")
}

func TestCanDisableAMonitor(t *testing.T) {
    // Given...
    monitorName := "customManagerCleanup"
    isEnabled := "false"

    // Create the expected HTTP interactions with the API server
    putMonitorInteraction := utils.NewHttpInteraction("/monitors/" + monitorName, http.MethodPut)
    putMonitorInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        requestBody := readMonitorRequestBody(req)
        assert.Equal(t, requestBody.Data.GetIsEnabled(), false)

        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusNoContent)
    }

    interactions := []utils.HttpInteraction{
        putMonitorInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetMonitor(
        monitorName,
        isEnabled,
        apiClient,
        mockByteReader)

    // Then...
    assert.Nil(t, err, "DisableMonitor returned an unexpected error")
}
