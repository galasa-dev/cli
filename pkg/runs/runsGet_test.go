/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	// "github.com/galasa.dev/cli/pkg/utils"

	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const (
	RUN_U456 = `{
		"runId": "xxx876xxx",
		"testStructure": {
			"runName": "U456",
			"bundle": "myBundleId",	
			"testName": "myTestPackage.MyTestName",
			"testShortName": "MyTestName",	
			"requestor": "unitTesting",
			"status" : "Finished",
			"result" : "Passed",
			"queued" : null,	
			"startTime": "now",
			"endTime": "now",
			"methods": [{
				"className": "myTestPackage.MyTestName",
				"methodName": "myTestMethodName",	
				"type": "test",	
				"status": "Done",	
				"result": "Success",
				"startTime": null,
				"endTime": null,	
				"runLogStart":null,	
				"runLogEnd":null,	
				"befores":[]
			}]
		},
		"artifacts": [{
			"artifactPath": "myPathToArtifact1",	
			"contentType":	"application/json"
		}]
	}`
)

// ------------------------------------------------------------------
// Testing that the output format string passed by the user on the command-line
// is valid and supported.
func TestOutputFormatSummaryValidatesOk(t *testing.T) {

	outputFormat, err := validateOutputFormatFlagValue("summary")
	if err != nil {
		assert.Fail(t, "Summary validate gave unexpected error "+err.Error())
	}
	assert.Equal(t, outputFormat, OUTPUT_FORMAT_SUMMARY)
}

func TestOutputFormatGarbageStringValidationGivesError(t *testing.T) {

	_, err := validateOutputFormatFlagValue("garbage")
	if err == nil {
		assert.Fail(t, "Garbage output format flag value should have given validation error.")
	}
	assert.Contains(t, err.Error(), "GAL1067")
}

func TestRunsGetOfRunIdWhichExistsProducesExpectedSummary(t *testing.T) {

	// Given ...
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ras/run" {
			t.Errorf("Expected to request '/ras/run', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		values := r.URL.Query()
		pageRequestedStr := values.Get("page")
		pageRequested, _ := strconv.Atoi(pageRequestedStr)
		assert.Equal(t, pageRequested, 1)

		w.Write([]byte(`
		{
			"pageNumber": 1,
			"pageSize": 1,
			"numPages": 1,	
			"amountOfRuns": 1,
			"runs":[` + RUN_U456 + `]
		}`))
	}))
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	runName := "U456"
	outputFormatString := "summary"
	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, outputFormatString, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	// We expect
	if err != nil {
		assert.Fail(t, "Failed with an error when we expected it to pass. Error is "+err.Error())
	} else {
		textGotBack := mockConsole.ReadText()
		assert.Contains(t, textGotBack, "U456")
	}
}
