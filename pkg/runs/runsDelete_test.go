/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func createMockRun(runName string, runId string) galasaapi.Run {
	run := *galasaapi.NewRun()
	run.SetRunId(runId)
	testStructure := *galasaapi.NewTestStructure()
	testStructure.SetRunName(runName)

	run.SetTestStructure(testStructure)
	return run
}

func TestCanDeleteARun(t *testing.T) {
	// Given...
	runName := "J20"
	runId := "J234567890"

	// Create the mock run to be deleted
	runToDelete := createMockRun(runName, runId)
	runToDeleteBytes, _ := json.Marshal(runToDelete)
	runToDeleteJson := string(runToDeleteBytes)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

		assert.NotEmpty(t, req.Header.Get("ClientApiVersion"))
		acceptHeader := req.Header.Get("Accept")
		if req.URL.Path == "/ras/runs" {
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsResponse(t, writer, req, runName, []string{ runToDeleteJson })
		} else if req.URL.Path == "/ras/runs/"+runId {
			writer.WriteHeader(http.StatusNoContent)
		}
	}))

	console := utils.NewMockConsole()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := RunsDelete(
		runName,
		console,
		apiServerUrl,
		apiClient,
		mockTimeService)

	// Then...
	assert.Nil(t, err, "RunsDelete returned an unexpected error")
	assert.Empty(t, console.ReadText(), "The console was written to on a successful deletion, it should be empty")
}

func TestDeleteNonExistantRunDisplaysError(t *testing.T) {
	// Given...
	nonExistantRunName := "runDoesNotExist123"

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		assert.NotEmpty(t, req.Header.Get("ClientApiVersion"))
		acceptHeader := req.Header.Get("Accept")
		if req.URL.Path == "/ras/runs" {
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsResponse(t, writer, req, nonExistantRunName, []string{})
		} else {
			assert.Fail(t, "An unexpected http request was issued to the test case.")
		}
	}))

	console := utils.NewMockConsole()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := RunsDelete(
		nonExistantRunName,
		console,
		apiServerUrl,
		apiClient,
		mockTimeService)

	// Then...
	assert.NotNil(t, err, "RunsDelete did not return an error but it should have")
	consoleOutputText := console.ReadText() 
	assert.Contains(t, consoleOutputText, nonExistantRunName)
	assert.Contains(t, consoleOutputText, "GAL1163E")
	assert.Contains(t, consoleOutputText, "The run named 'runDoesNotExist123' could not be deleted")
}

func TestRunsDeleteFailsWithNoExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
	// Given...
	runName := "J20"
	runId := "J234567890"

	// Create the mock run to be deleted
	runToDelete := createMockRun(runName, runId)
	runToDeleteBytes, _ := json.Marshal(runToDelete)
	runToDeleteJson := string(runToDeleteBytes)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

		assert.NotEmpty(t, req.Header.Get("ClientApiVersion"))
		acceptHeader := req.Header.Get("Accept")
		if req.URL.Path == "/ras/runs" {
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsResponse(t, writer, req, runName, []string{ runToDeleteJson })
		} else if req.URL.Path == "/ras/runs/"+runId {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}))

	console := utils.NewMockConsole()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := RunsDelete(
		runName,
		console,
		apiServerUrl,
		apiClient,
		mockTimeService)

	// Then...
	assert.NotNil(t, err, "RunsDelete returned an unexpected error")
	consoleText := console.ReadText()
	assert.Contains(t, consoleText , runName)
	assert.Contains(t, consoleText , "GAL1159E")
}

func TestRunsDeleteFailsWithNonJsonContentTypeExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
	// Given...
	runName := "J20"
	runId := "J234567890"

	// Create the mock run to be deleted
	runToDelete := createMockRun(runName, runId)
	runToDeleteBytes, _ := json.Marshal(runToDelete)
	runToDeleteJson := string(runToDeleteBytes)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

		assert.NotEmpty(t, req.Header.Get("ClientApiVersion"))
		acceptHeader := req.Header.Get("Accept")
		if req.URL.Path == "/ras/runs" {
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsResponse(t, writer, req, runName, []string{ runToDeleteJson })
		} else if req.URL.Path == "/ras/runs/"+runId {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Header().Set("Content-Type", "application/notJsonOnPurpose")
			writer.Write([]byte("something not json but non-zero-length."))
		}
	}))

	console := utils.NewMockConsole()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := RunsDelete(
		runName,
		console,
		apiServerUrl,
		apiClient,
		mockTimeService)

	// Then...
	assert.NotNil(t, err, "RunsDelete returned an unexpected error")
	consoleText := console.ReadText()
	assert.Contains(t, consoleText, runName)
	assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
	assert.Contains(t, consoleText, "GAL1164E")
	assert.Contains(t, consoleText, "Error details from the server are not in the json format")
}


// func TestRunsDeleteFailsWithBadlyFormedJsonContentExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
// 	// Given...
// 	runName := "J20"
// 	runId := "J234567890"

// 	// Create the mock run to be deleted
// 	runToDelete := createMockRun(runName, runId)
// 	runToDeleteBytes, _ := json.Marshal(runToDelete)
// 	runToDeleteJson := string(runToDeleteBytes)

// 	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

// 		assert.NotEmpty(t, req.Header.Get("ClientApiVersion"))
// 		acceptHeader := req.Header.Get("Accept")
// 		if req.URL.Path == "/ras/runs" {
// 			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
// 			WriteMockRasRunsResponse(t, writer, req, runName, []string{ runToDeleteJson })
// 		} else if req.URL.Path == "/ras/runs/"+runId {
// 			writer.WriteHeader(http.StatusInternalServerError)
// 			writer.Header().Add("Content-Type", "application/json")
// 			writer.Write([]byte(`{ "this", "isBadJson because it doesnt end in a close braces" `))
// 		}
// 	}))

// 	console := utils.NewMockConsole()
// 	apiServerUrl := server.URL
// 	apiClient := api.InitialiseAPI(apiServerUrl)
// 	mockTimeService := utils.NewMockTimeService()

// 	// When...
// 	err := RunsDelete(
// 		runName,
// 		console,
// 		apiServerUrl,
// 		apiClient,
// 		mockTimeService)

// 	// Then...
// 	assert.NotNil(t, err, "RunsDelete returned an unexpected error")
// 	consoleText := console.ReadText()
// 	assert.Contains(t, consoleText, runName)
// 	assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
// 	assert.Contains(t, consoleText, "Gxxxx")
// 	assert.Contains(t, consoleText, "Error details from the server are not in the json format")
// }
