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

	RUN_C123 = `{
		"runId": "xxx543xxx",
		"testStructure": {
			"runName": "C123",
			"bundle": "myBundleId",	
			"testName": "myTestPackage.MyTestName2",
			"testShortName": "MyTestName2",	
			"requestor": "unitTesting",
			"status" : "Submitted",
			"result" : "UNKNOWN",
			"queued" : null,	
			"startTime": "now",
			"endTime": "now",
			"methods": [{
				"className": "myTestPackage.MyTestName2",
				"methodName": "myTestMethodName",	
				"type": "test",	
				"status": "Done",	
				"result": "UNKNOWN",
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

func TestRunsGetOfRunNameWhichExistsProducesExpectedSummary(t *testing.T) {

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
		want := "" +
			"\nRunName Status   Result ShortTestName " +
			"\nU456    Finished Passed MyTestName    \n"
		assert.Equal(t, textGotBack, want)
	}
}

func TestRunsGetOfRunNameWhichDoesNotExistProducesError(t *testing.T) {
	// Given ...
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ras/run" {
			t.Errorf("Expected to request '/ras/run', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		values := r.URL.Query()
		pageRequestedStr := values.Get("page")
		pageRequested, _ := strconv.Atoi(pageRequestedStr)
		assert.Equal(t, pageRequested, 1)

		w.Write([]byte(`
		{
			"pageNumber": 1,
			"pageSize": 1,
			"numPages": 1,	
			"amountOfRuns": 0,
			"runs":[` + RUN_C123 + `]
		}`))
	}))
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	runName := "garbage"
	outputFormatString := "summary"
	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, outputFormatString, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	// We expect

	if err == nil {
		assert.Fail(t, "Garbage runname value should have given 404 error.")
	}
	assert.Contains(t, err.Error(), "GAL1068")
	assert.Contains(t, err.Error(), "404")

}

func TestRunsGetOfRunNameWhichExistsAmongstMultipleRunsProducesExpectedSummary(t *testing.T) {
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
		runnameURL := values.Get("runname")
		pageRequested, _ := strconv.Atoi(pageRequestedStr)
		assert.Equal(t, pageRequested, 1)

		if runnameURL == "C123" {
			w.Write([]byte(`
			{
				"pageNumber": 1,
				"pageSize": 1,
				"numPages": 1,	
				"amountOfRuns": 1,
				"runs":[` + RUN_C123 + `]
			}`))
		} else {
			w.Write([]byte(`
		{
			"pageNumber": 1,
			"pageSize": 1,
			"numPages": 1,	
			"amountOfRuns": 2,
			"runs":[` + RUN_U456 + "," + RUN_C123 + `]
		}`))
		}
	}))
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	runName := "C123"
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
		assert.Contains(t, textGotBack, "C123")
		want := "" +
			"\nRunName Status    Result  ShortTestName " +
			"\nC123    Submitted UNKNOWN MyTestName2   \n"
		assert.Equal(t, textGotBack, want)
	}
}
