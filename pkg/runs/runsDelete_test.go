/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewRunsDeleteServletMock(
	t *testing.T,
	runName string,
	runId string,
	runResultStrings []string,
) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

		assert.NotEmpty(t, req.Header.Get("ClientApiVersion"))
		acceptHeader := req.Header.Get("Accept")
		if req.URL.Path == "/ras/runs" {
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsResponse(t, writer, req, runName, runResultStrings)
		} else if req.URL.Path == "/ras/runs/"+runId {
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsPutStatusQueuedResponse(t, writer, req, runName)
		}
	}))

	return server
}

// func TestCanDeleteARun(t *testing.T) {
// 	//Given
// 	runs := make([]galasaapi.Run, 0)

// 	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
// 		assert.NotEmpty(t, req.Header.Get("ClientApiVersion"))
// 	}))

// 	runName := "J20"

// 	//When
// 	console := NewMockConsole()
// 	apiServerUrl := server.URL
// 	apiClient := api.InitialiseAPI(apiServerUrl)
// 	mockTimeService := utils.NewMockTimeService()

// 	err := RunsDelete(
// 		runName,
// 		console,
// 		apiServerUrl,
// 		apiClient,
// 		mockTimeService)

// 	//Then
// 	assert.NotNil(t, err, "RunsDelete returned an unexpected error %s", err.Error)

// }
