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

	"github.com/galasa.dev/cli/pkg/formatters"
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
			 "queued" : "2023-05-10T06:00:13.043037Z",	
			 "startTime": "2023-05-10T06:00:36.159003Z",
			 "endTime": "2023-05-10T06:02:53.823338Z",
			 "methods": [{
				 "className": "myTestPackage.MyTestName",
				 "methodName": "myTestMethodName",	
				 "type": "test",	
				 "status": "Done",	
				 "result": "Success",
				 "startTime": "2023-05-10T06:00:13.254335Z",
				 "endTime": "2023-05-10T06:03:11.882739Z",	
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
			 "queued" : "2023-05-10T06:00:13.405966Z",	
			 "startTime": "2023-05-10T06:02:26.801649Z",
			 "endTime": "2023-05-10T06:04:04.448826Z",
			 "methods": [{
				 "className": "myTestPackage22.MyTestName2",
				 "methodName": "myTestMethodName",	
				 "type": "test",	
				 "status": "Done",	
				 "result": "UNKNOWN",
				 "startTime": "2023-05-10T06:02:28.457784Z",
				 "endTime": "2023-05-10T06:04:28.585024Z",	
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
	validFormatters := createFormatters()
	outputFormatter, err := validateOutputFormatFlagValue("summary", validFormatters)
	if err != nil {
		assert.Fail(t, "Summary validate gave unexpected error "+err.Error())
	}
	assert.NotNil(t, outputFormatter)
}

func TestOutputFormatGarbageStringValidationGivesError(t *testing.T) {
	validFormatters := createFormatters()
	_, err := validateOutputFormatFlagValue("garbage", validFormatters)
	if err == nil {
		assert.Fail(t, "Garbage output format flag value should have given validation error.")
	}
	assert.Contains(t, err.Error(), "GAL1067")
	assert.Contains(t, err.Error(), "'garbage'")
	assert.Contains(t, err.Error(), "'summary'")
	assert.Contains(t, err.Error(), "'details'")
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
			"name status   result test-name\n" +
				"U456 Finished Passed myTestPackage.MyTestName\n"
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
		want := ""
		assert.Equal(t, textGotBack, want)
	}

}

func NewRunsGetServletMock(t *testing.T, status int, runName string, runResultStrings ...string) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ras/run" {
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
			"name status   result           test-name\n" +
				"U456 Finished Passed           myTestPackage.MyTestName\n" +
				"U456 Finished LongResultString myTestPackage.MyTest2\n"
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

func TestOutputFormatDetailsValidatesOk(t *testing.T) {
	validFormatters := createFormatters()
	outputFormatter, err := validateOutputFormatFlagValue("details", validFormatters)
	if err != nil {
		assert.Fail(t, "Details validate gave unexpected error "+err.Error())
	}
	assert.NotNil(t, outputFormatter)
}

func TestRunsGetOfRunNameWhichExistsProducesExpectedDetails(t *testing.T) {

	// Given ...
	runName := "U456"
	server := NewRunsGetServletMock(t, http.StatusOK, runName, RUN_U456)
	defer server.Close()

	outputFormat := "details"
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
			"name         :  U456\n" +
				"status       :  Finished\n" +
				"result       :  Passed\n" +
				"queued-time  :  2023-05-10 06:00:13\n" +
				"start-time   :  2023-05-10 06:00:36\n" +
				"end-time     :  2023-05-10 06:02:53\n" +
				"duration(ms) :  137000\n" +
				"test-name    :  myTestPackage.MyTestName\n" +
				"requestor    :  unitTesting\n" +
				"bundle       :  myBundleId\n" +
				"run-log      :  " + apiServerUrl + "/ras/run/xxx876xxx/runlog\n" +
				"\n" +
				"method           type status result  start-time          end-time            duration(ms)\n" +
				"myTestMethodName test Done   Success 2023-05-10 06:00:13 2023-05-10 06:03:11 178000\n"
		assert.Equal(t, textGotBack, want)
	}
}

func TestGetFormatterNamesStringMultipleFormattersFormatsOk(t *testing.T) {
	validFormatters := make(map[string]formatters.RunsFormatter, 0)
	validFormatters["first"] = nil
	validFormatters["second"] = nil

	result := getFormatterNamesString(validFormatters)

	assert.NotNil(t, result)
	assert.Equal(t, result, "'first', 'second'")
}

func TestRunsGetOfRunNameWhichExistsProducesExpectedRaw(t *testing.T) {

	// Given ...
	runName := "U456"
	server := NewRunsGetServletMock(t, http.StatusOK, runName, RUN_U456)
	defer server.Close()

	outputFormat := "raw"
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
		want := "U456|Finished|Passed|2023-05-10 06:00:13|2023-05-10 06:00:36|2023-05-10 06:02:53|137000|myTestPackage.MyTestName|unitTesting|myBundleId|" + apiServerUrl + "/ras/run/xxx876xxx/runlog"

		assert.Equal(t, textGotBack, want)
	}
}
