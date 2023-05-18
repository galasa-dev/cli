/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"fmt"
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

	RUN_U456_v2 = `{
		"runId": "xxx543xxx",
		"testStructure": {
			"runName": "U456",
			"bundle": "myBun2",	
			"testName": "myTestPackage.MyTest2",
			"testShortName": "MyTestName22",	
			"requestor": "unitTesting22",
			"status" : "Finished",
			"result" : "LongResultString",
			"queued" : null,	
			"startTime": "now",
			"endTime": "now",
			"methods": [{
				"className": "myTestPackage22.MyTestName2",
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
	runName := "U456"
	server := NewRunsGetServletMock(t, http.StatusOK, runName, RUN_U456)
	defer server.Close()

	outputFormat := "summary"
	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	// We expect
	if err != nil {
		assert.Fail(t, "Failed with an error when we expected it to pass. Error is "+err.Error())
	} else {
		textGotBack := mockConsole.ReadText()
		assert.Contains(t, textGotBack, runName)
		want :=
			"RunName Status   Result ShortTestName\n" +
				"U456    Finished Passed MyTestName\n"
		assert.Equal(t, textGotBack, want)
	}
}

func TestRunsGetOfRunNameWhichDoesNotExistProducesEmptyPage(t *testing.T) {
	// Given ...
	runName := "garbage"
	server := NewRunsGetServletMock(t, http.StatusOK, runName)
	defer server.Close()

	mockConsole := utils.NewMockConsole()

	outputFormatString := "summary"
	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, outputFormatString, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	// We expect

	if err != nil {
		assert.Fail(t, "Garbage runname value should not have failed "+err.Error())
	} else {
		textGotBack := mockConsole.ReadText()
		want := "RunName Status Result ShortTestName\n"
		assert.Equal(t, textGotBack, want)
	}

}

func NewRunsGetServletMock(t *testing.T, status int, runName string, runResultStrings ...string) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ras/runs" {
			t.Errorf("Expected to request '/ras/run', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		values := r.URL.Query()
		pageRequestedStr := values.Get("page")
		runNameQueryParameter := values.Get("runname")
		pageRequested, _ := strconv.Atoi(pageRequestedStr)
		assert.Equal(t, pageRequested, 1)

		assert.Equal(t, runNameQueryParameter, runName)

		combinedRunResultStrings := ""
		for index, runResult := range runResultStrings {
			if index > 0 {
				combinedRunResultStrings += ","
			}
			combinedRunResultStrings += runResult
		}

		w.Write([]byte(fmt.Sprintf(`
		{
			"pageNumber": 1,
			"pageSize": 1,
			"numPages": 1,
			"amountOfRuns": %d,
			"runs":[ %s ]
		}`, len(runResultStrings), combinedRunResultStrings)))
	}))

	return server
}

func TestRunsGetWhereRunNameExistsTwiceProducesTwoRunResultLines(t *testing.T) {
	// Given ...
	runName := "U456"
	server := NewRunsGetServletMock(t, http.StatusOK, runName, RUN_U456, RUN_U456_v2)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
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
		assert.Contains(t, textGotBack, runName)
		want :=
			"RunName Status   Result           ShortTestName\n" +
				"U456    Finished Passed           MyTestName\n" +
				"U456    Finished LongResultString MyTestName22\n"
		assert.Equal(t, textGotBack, want)
	}
}

func TestFailingGetRunsRequestReturnsError(t *testing.T) {

	// Given...
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	runName := "garbage"
	mockConsole := utils.NewMockConsole()
	outputFormatString := "summary"
	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, outputFormatString, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	assert.Contains(t, err.Error(), "GAL1068")
}
