/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const (
	// An active run that should be finished now.
	RUN_U123_FIRST_RUN = `{
		"runId": "xxx122xxx",
		"artifacts": [],
		"testStructure": {
			"runName": "U123",
			"bundle": "myBundleId",
			"testName": "myTestPackage.MyTestName",
			"testShortName": "MyTestName",
			"requestor": "unitTesting",
			"status" : "building",
			"queued" : "2023-05-10T06:00:00.000000Z",
			"startTime": "2023-05-10T06:00:10.000000Z",
			"methods": [{
				"className": "myTestPackage.MyTestName",
				"methodName": "myTestMethodName",
				"type": "test",
				"runLogStart":null,
				"runLogEnd":null,
				"befores":[], 
				"afters": []
			}]
		}
	}`
	// An active run
	RUN_U123_RE_RUN = `{
		"runId": "xxx123xxx",
		"artifacts": [],
		"testStructure": {
			"runName": "U123",
			"bundle": "myBundleId",
			"testName": "myTestPackage.MyTestName",
			"testShortName": "MyTestName",
			"requestor": "unitTesting",
			"status" : "building",
			"queued" : "2023-05-10T06:00:13.043037Z",
			"startTime": "2023-05-10T06:00:36.159003Z",
			"methods": [{
				"className": "myTestPackage.MyTestName",
				"methodName": "myTestMethodName",
				"type": "test",
				"runLogStart":null,
				"runLogEnd":null,
				"befores":[], 
				"afters": []
			}]
		}
	}`
	// A finished run
	RUN_U120 = `{
		"runId": "xxx120xxx",
		"artifacts": [],
		"testStructure": {
			"runName": "U120",
			"bundle": "myBundleId",
			"testName": "myTestPackage.MyTestName",
			"testShortName": "MyTestName",
			"requestor": "unitTesting",
			"status" : "finished",
			"result": "Passed",
			"queued" : "2023-05-10T06:00:13.043037Z",
			"startTime": "2023-05-10T06:00:36.159003Z",
			"endTime": "2023-05-10T06:01:36.159003Z",
			"methods": [{
				"className": "myTestPackage.MyTestName",
				"methodName": "myTestMethodName",
				"type": "test",
				"status": "finished",
        		"result": "Passed",
				"startTime": "2023-05-10T06:00:36.159003Z",
				"endTime": "2023-05-10T06:01:36.159003Z",
				"runLogStart":0,
				"runLogEnd":0,
				"befores":[], 
				"afters": []
			}]
		}
	}`
	// Another finished run
	RUN_U121 = `{
		"runId": "xxx121xxx",
		"artifacts": [],
		"testStructure": {
			"runName": "U121",
			"bundle": "myBundleId",
			"testName": "myTestPackage.MyTestName",
			"testShortName": "MyTestName",
			"requestor": "unitTesting",
			"status" : "finished",
			"result": "Passed",
			"queued" : "2023-05-10T06:00:13.043037Z",
			"startTime": "2023-05-10T06:00:36.159003Z",
			"endTime": "2023-05-10T06:01:36.159003Z",
			"methods": [{
				"className": "myTestPackage.MyTestName",
				"methodName": "myTestMethodName",
				"type": "test",
				"status": "finished",
        		"result": "Passed",
				"startTime": "2023-05-10T06:00:36.159003Z",
				"endTime": "2023-05-10T06:01:36.159003Z",
				"runLogStart":0,
				"runLogEnd":0,
				"befores":[], 
				"afters": []
			}]
		}
	}`
)

func WriteMockRasRunsPutStatusQueuedResponse(
	t *testing.T,
	writer http.ResponseWriter,
	req *http.Request,
	runName string) {

	var statusCode int
	var response string

	if runName == "U123" {
		statusCode = 200
		response = fmt.Sprintf("Successfully reset run %s", runName)
		writer.Header().Set("Content-Type", "text/plain")
	} else if runName == "U120" {
		statusCode = 400
		response = `{
			"error_code": 5049, 
			"error_message": "GAL5049E: Error occured when trying to reset the run 'U120'. The run has already completed."
		}`
		writer.Header().Set("Content-Type", "application/json")
	} else if runName == "U121" {
		statusCode = 400
		response = `{{
			not for parsing
		}`
		writer.Header().Set("Content-Type", "application/json")
	}

	writer.WriteHeader(statusCode)
	writer.Write([]byte(response))
}

func NewRunsResetServletMock(
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

//------------------------------------------------------------------
// Test methods
//------------------------------------------------------------------

func TestRunsResetWithActiveRunReturnsOK(t *testing.T) {
	// Given ...
	runName := "U123"
	runId := "xxx123xxx"

	runResultStrings := []string{RUN_U123_RE_RUN}

	server := NewRunsResetServletMock(t, runName, runId, runResultStrings)
	defer server.Close()

	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := ResetRun(runName, mockTimeService, mockConsole, apiServerUrl, apiClient)

	// Then...
	assert.Nil(t, err)
	textGotBack := mockConsole.ReadText()
	assert.Contains(t, textGotBack, "GAL2503I")
	assert.Contains(t, textGotBack, runName)
}

func TestRunsResetWithMultipleActiveRunsReturnsError(t *testing.T) {
	// Given ...
	runName := "U123"
	runId := "xxx123xxx"

	runResultStrings := []string{RUN_U123_FIRST_RUN, RUN_U123_RE_RUN}

	server := NewRunsResetServletMock(t, runName, runId, runResultStrings)
	defer server.Close()

	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := ResetRun(runName, mockTimeService, mockConsole, apiServerUrl, apiClient)

	// Then...
	assert.Contains(t, err.Error(), "GAL1131")
	assert.Contains(t, err.Error(), runName)
}

func TestRunsResetWithNoActiveRunReturnsError(t *testing.T) {
	// Given ...
	runName := "U123"
	runId := "xxx123xxx"

	runResultStrings := []string{}

	server := NewRunsResetServletMock(t, runName, runId, runResultStrings)
	defer server.Close()

	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := ResetRun(runName, mockTimeService, mockConsole, apiServerUrl, apiClient)

	// Then...
	assert.Contains(t, err.Error(), "GAL1132")
	assert.Contains(t, err.Error(), runName)
}

func TestRunsResetWithInvalidRunNameReturnsError(t *testing.T) {
	// Given ...
	runName := "garbage"
	runId := "xxx123xxx"

	runResultStrings := []string{RUN_U123_RE_RUN}

	server := NewRunsResetServletMock(t, runName, runId, runResultStrings)
	defer server.Close()

	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := ResetRun(runName, mockTimeService, mockConsole, apiServerUrl, apiClient)

	// Then...
	assert.Contains(t, err.Error(), "GAL1075")
	assert.Contains(t, err.Error(), runName)
}

func TestRunsResetWhereOperationFailedServerSideReturnsError(t *testing.T) {
	// Given ...
	runName := "U120"
	runId := "xxx120xxx"

	runResultStrings := []string{RUN_U120}

	server := NewRunsResetServletMock(t, runName, runId, runResultStrings)
	defer server.Close()

	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := ResetRun(runName, mockTimeService, mockConsole, apiServerUrl, apiClient)

	// Then...
	assert.Error(t, err)
	assert.ErrorContains(t, err, "GAL1133")
}

func TestRunsResetWhereServerSideResponseCannotBeParsedReturnsError(t *testing.T) {
	// Given ...
	runName := "U121"
	runId := "xxx121xxx"

	runResultStrings := []string{RUN_U121}

	server := NewRunsResetServletMock(t, runName, runId, runResultStrings)
	defer server.Close()

	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := ResetRun(runName, mockTimeService, mockConsole, apiServerUrl, apiClient)

	// Then...
	assert.Error(t, err)
	assert.ErrorContains(t, err, "GAL1134")
}
