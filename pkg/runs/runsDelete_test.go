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
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func NewRunsDeleteServletMock(
	t *testing.T,
	runName string,
	runId string,
	runResultJsonStrings []string,
	deleteRunStatusCode int,
) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

		assert.NotEmpty(t, req.Header.Get("ClientApiVersion"))
		acceptHeader := req.Header.Get("Accept")
		if req.URL.Path == "/ras/runs" {
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsResponse(t, writer, req, runName, runResultJsonStrings)
		} else if req.URL.Path == "/ras/runs/"+runId {
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsDeleteResponse(t, writer, req, runName, deleteRunStatusCode)
		}
	}))

	return server
}

func WriteMockRasRunsDeleteResponse(
	t *testing.T,
	writer http.ResponseWriter,
	req *http.Request,
	runName string,
	statusCode int) {

	writer.WriteHeader(statusCode)
	
}

func createMockRun(runName string) galasaapi.Run {
	run := *galasaapi.NewRun()
	run.SetRunId(runName)
	testStructure := *galasaapi.NewTestStructure()
	testStructure.SetRunName(runName)

	run.SetTestStructure(testStructure)
	return run
}

func TestCanDeleteARun(t *testing.T) {
	// Given...
	runName := "J20"

	// Create the mock run to be deleted
	runToDelete := createMockRun(runName)
	runToDeleteBytes, _ := json.Marshal(runToDelete)
	runToDeleteJson := string(runToDeleteBytes)

	server := NewRunsDeleteServletMock(t, runName, runName, []string{ runToDeleteJson }, 204)

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
	nonExistantRunName := "run-does-not-exist"
	
	existingRunName := "J20"
	existingRun := createMockRun(existingRunName)
	existingRunBytes, _ := json.Marshal(existingRun)
	existingRunJson := string(existingRunBytes)

	server := NewRunsDeleteServletMock(t, nonExistantRunName, nonExistantRunName, []string{ existingRunJson }, 404)

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
}
