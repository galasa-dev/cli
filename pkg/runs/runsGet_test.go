/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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

func NewRunsGetServletMock(t *testing.T, status int, runName string, runResultStrings ...string) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/ras/runs/") {
			ConfigureServerForDetailsEndpoint(t, w, r, status, runResultStrings...)
		} else {
			ConfigureServerForRasRunsEndpoint(t, w, r, runName, status, runResultStrings...)
		}

	}))

	return server
}

func ConfigureServerForDetailsEndpoint(t *testing.T, w http.ResponseWriter, r *http.Request, status int, runResultStrings ...string) {
	if r.Header.Get("Accept") != "application/json" {
		t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
	}
	urlParts := strings.Split(r.URL.Path, "/")
	runid := urlParts[3]
	for _, runResult := range runResultStrings {
		assert.Contains(t, runResult, runid)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	combinedRunResultStrings := ""
	for index, runResult := range runResultStrings {
		if index > 0 {
			combinedRunResultStrings += ","
		}
		combinedRunResultStrings += runResult
	}

	w.Write([]byte(fmt.Sprintf(`
			%s 
		`, combinedRunResultStrings)))
}

func ConfigureServerForRasRunsEndpoint(t *testing.T, w http.ResponseWriter, r *http.Request, runName string, status int, runResultStrings ...string) {
	if r.URL.Path != "/ras/runs" {
		t.Errorf("Expected to request '/ras/runs', got: %s", r.URL.Path)
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
}

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
	assert.Contains(t, err.Error(), "'raw'")
}

func TestRunsGetOfRunNameWhichExistsProducesExpectedSummary(t *testing.T) {

	// Given ...
	runName := "U456"
	age := "2d:24h"
	server := NewRunsGetServletMock(t, http.StatusOK, runName, RUN_U456)
	defer server.Close()

	outputFormat := "summary"
	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	// We expect
	if err != nil {
		assert.Fail(t, "Failed with an error when we expected it to pass. Error is "+err.Error())
	} else {
		textGotBack := mockConsole.ReadText()
		assert.Contains(t, textGotBack, runName)
		want :=
			"submitted-time      name status   result test-name\n" +
				"2023-05-10 06:00:13 U456 Finished Passed myTestPackage.MyTestName\n" +
				"\n" +
				"Total:1 Passed:1 PassedWithDefects:0 Failed:0 EnvFail:0 FailedWithDefects:0\n"
		assert.Equal(t, textGotBack, want)
	}
}

func TestRunsGetOfRunNameWhichDoesNotExistProducesError(t *testing.T) {
	// Given ...
	age := "2d:24h"
	runName := "garbage"
	server := NewRunsGetServletMock(t, http.StatusOK, runName)
	defer server.Close()

	mockConsole := utils.NewMockConsole()

	outputFormat := "summary"
	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	// We expect

	assert.NotNil(t, err, "Garbage runname value should not have failed.")
	if err != nil {
		assert.ErrorContains(t, err, "GAL1075E")
		assert.ErrorContains(t, err, runName)
	}
}

func TestRunsGetWhereRunNameExistsTwiceProducesTwoRunResultLines(t *testing.T) {
	// Given ...
	age := ""
	runName := "U456"
	server := NewRunsGetServletMock(t, http.StatusOK, runName, RUN_U456, RUN_U456_v2)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	outputFormat := "summary"
	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	// We expect
	if err != nil {
		assert.Fail(t, "Failed with an error when we expected it to pass. Error is "+err.Error())
	} else {
		textGotBack := mockConsole.ReadText()
		assert.Contains(t, textGotBack, runName)
		want :=
			"submitted-time      name status   result           test-name\n" +
				"2023-05-10 06:00:13 U456 Finished Passed           myTestPackage.MyTestName\n" +
				"2023-05-10 06:00:13 U456 Finished LongResultString myTestPackage.MyTest2\n" +
				"\n" +
				"Total:2 Passed:1 PassedWithDefects:0 Failed:0 EnvFail:0 FailedWithDefects:0\n"
		assert.Equal(t, textGotBack, want)
	}
}

func TestFailingGetRunsRequestReturnsError(t *testing.T) {

	// Given...
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	age := ""
	runName := "garbage"
	mockConsole := utils.NewMockConsole()
	outputFormat := "summary"
	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	assert.Contains(t, err.Error(), "GAL1075")
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
	age := ""
	runName := "U456"
	server := NewRunsGetServletMock(t, http.StatusOK, runName, RUN_U456)
	defer server.Close()

	outputFormat := "details"
	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	// We expect
	if err != nil {
		assert.Fail(t, "Failed with an error when we expected it to pass. Error is "+err.Error())
	} else {
		textGotBack := mockConsole.ReadText()
		assert.Contains(t, textGotBack, runName)
		want :=
			"name           : U456\n" +
				"status         : Finished\n" +
				"result         : Passed\n" +
				"submitted-time : 2023-05-10 06:00:13\n" +
				"start-time     : 2023-05-10 06:00:36\n" +
				"end-time       : 2023-05-10 06:02:53\n" +
				"duration(ms)   : 137664\n" +
				"test-name      : myTestPackage.MyTestName\n" +
				"requestor      : unitTesting\n" +
				"bundle         : myBundleId\n" +
				"run-log        : " + apiServerUrl + "/ras/runs/xxx876xxx/runlog\n" +
				"\n" +
				"method           type status result  start-time          end-time            duration(ms)\n" +
				"myTestMethodName test Done   Success 2023-05-10 06:00:13 2023-05-10 06:03:11 178628\n" +
				"\n" +
				"Total:1 Passed:1 PassedWithDefects:0 Failed:0 EnvFail:0 FailedWithDefects:0\n"
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

func TestAPIInternalErrorIsHandledOk(t *testing.T) {
	// Given ...
	age := ""
	runName := "U456"
	server := NewRunsGetServletMock(t, http.StatusInternalServerError, runName, RUN_U456)
	defer server.Close()

	outputFormat := "details"
	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	// We expect
	assert.Error(t, err)
	assert.ErrorContains(t, err, "500")
	assert.ErrorContains(t, err, "GAL1068")
}

func TestRunsGetOfRunNameWhichExistsProducesExpectedRaw(t *testing.T) {

	// Given ...
	age := ""
	runName := "U456"
	server := NewRunsGetServletMock(t, http.StatusOK, runName, RUN_U456)
	defer server.Close()

	outputFormat := "raw"
	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	// We expect
	assert.Nil(t, err)
	textGotBack := mockConsole.ReadText()
	assert.Contains(t, textGotBack, runName)
	want := "U456|Finished|Passed|2023-05-10T06:00:13.043037Z|2023-05-10T06:00:36.159003Z|2023-05-10T06:02:53.823338Z|137664|myTestPackage.MyTestName|unitTesting|myBundleId|" + apiServerUrl + "/ras/run/xxx876xxx/runlog\n"
	assert.Equal(t, textGotBack, want)
}

func TestRunsGetWithFromAndToAge(t *testing.T) {

	// Given ...
	age := "5d:12h"

	//When ...
	from, to, err := getTimesFromAge(age)

	// Then...
	// We expect
	// from = 5*24 = 120
	// to   = 12*1 = 12
	assert.Nil(t, err)
	assert.NotNil(t, from)
	assert.NotNil(t, to)
	assert.EqualValues(t, 120, from)
	assert.EqualValues(t, 12, to)
}

func TestRunsGetWithJustFromAge(t *testing.T) {

	// Given
	age := "20d"

	// When
	from, to, err := getTimesFromAge(age)

	// Then...
	// We expect
	// from = 20*24    = 480
	// to not provided = 0
	assert.Nil(t, err)
	assert.NotNil(t, from)
	assert.NotNil(t, to)
	assert.EqualValues(t, 480, from)
	assert.EqualValues(t, 0, to)
}

func TestRunsGetWithNoRunNameAndNoFromAgeReturnsError(t *testing.T) {

	// Given
	age := ""

	// When
	_, _, err := getTimesFromAge(age)

	// Then...
	// We expect
	assert.Error(t, err)
	assert.ErrorContains(t, err, "GAL1079")
}

func TestRunsGetWithBadlyFormedFromAndToParameter(t *testing.T) {

	// Given
	age := "12m:1y"

	// When
	_, _, err := getTimesFromAge(age)

	// Then...
	// We expect
	assert.Error(t, err)
	assert.ErrorContains(t, err, "GAL1079")
}

func TestRunsGetWithOlderToAgeThanFromAge(t *testing.T) {

	// Given
	age := "1d:3d"

	// When
	_, _, err := getTimesFromAge(age)

	// Then...
	// We expect
	assert.Error(t, err)
	assert.ErrorContains(t, err, "GAL1077")
}

func TestRunsGetURLQueryWithFromAndToDate(t *testing.T) {
	// Given ...
	age := "5d:12h"
	runName := "U456"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.NotNil(t, query.Get("from"))
		assert.NotEqualValues(t, query.Get("from"), "")
		assert.NotNil(t, query.Get("to"))
		assert.NotEqualValues(t, query.Get("to"), "")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`
		 {
			 "pageNumber": 1,
			 "pageSize": 1,
			 "numPages": 1,
			 "amountOfRuns": 0,
			 "runs":[]
		 }`))
	}))
	defer server.Close()

	outputFormat := "summary"
	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then ...
	assert.Nil(t, err)
}

func TestRunsGetURLQueryJustFromAge(t *testing.T) {
	// Given ...
	age := "2d"
	runName := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.NotNil(t, query.Get("from"))
		assert.NotEqualValues(t, query.Get("from"), "")
		assert.EqualValues(t, query.Get("to"), "")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`
		 {
			 "pageNumber": 1,
			 "pageSize": 1,
			 "numPages": 1,
			 "amountOfRuns": 0,
			 "runs":[]
		 }`))
	}))
	defer server.Close()

	outputFormat := "summary"
	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then ...
	assert.Nil(t, err)
}

func TestRunsGetURLQueryWithNoRunNameAndNoFromAgeReturnsError(t *testing.T) {
	// Given ...
	age := ""
	runName := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.EqualValues(t, query.Get("from"), "")
		assert.EqualValues(t, query.Get("to"), "")
		assert.EqualValues(t, query.Get("runname"), "")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`
		 {
			 "pageNumber": 1,
			 "pageSize": 1,
			 "numPages": 1,
			 "amountOfRuns": 0,
			 "runs":[]
		 }`))
	}))
	defer server.Close()

	outputFormat := "summary"
	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then ...
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "1080")
}

func TestRunsGetURLQueryWithOlderToAgeThanFromAgeReturnsError(t *testing.T) {
	// Given ...
	age := "1d:1w"
	runName := "U456"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.EqualValues(t, query.Get("from"), "")
		assert.EqualValues(t, query.Get("to"), "")
		assert.EqualValues(t, query.Get("runname"), "U456")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`
		 {
			 "pageNumber": 1,
			 "pageSize": 1,
			 "numPages": 1,
			 "amountOfRuns": 0,
			 "runs":[]
		 }`))
	}))
	defer server.Close()

	outputFormat := "summary"
	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then ...
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "1077")
}

func TestRunsGetURLQueryWithBadlyFormedFromAndToParameterReturnsError(t *testing.T) {
	// Given ...
	age := "12m:1y"
	runName := "U456"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.EqualValues(t, query.Get("from"), "")
		assert.EqualValues(t, query.Get("to"), "")
		assert.EqualValues(t, query.Get("runname"), "U456")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`
		 {
			 "pageNumber": 1,
			 "pageSize": 1,
			 "numPages": 1,
			 "amountOfRuns": 0,
			 "runs":[]
		 }`))
	}))
	defer server.Close()

	outputFormat := "summary"
	mockConsole := utils.NewMockConsole()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := GetRuns(runName, age, outputFormat, mockTimeService, mockConsole, apiServerUrl)

	// Then ...
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "1079")
}

