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

	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const (
	RUN_U1 = `{
		 "runId": "xxx876xxx",
		 "testStructure": {
			 "runName": "U1",
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

	RUN_U2 = `{
		 "runId": "xxx543xxx",
		 "testStructure": {
			 "runName": "U2",
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

func WriteMockRasRunsResponse(
	t *testing.T,
	w http.ResponseWriter,
	r *http.Request,
	status int,
	runName string,
	runResultStrings ...string) {

	if r.Header.Get("Accept") != "application/json" {
		t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
	}
	w.Header().Set("Content-Type", "application/json")

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

func NewRunsDownloadServletMock(
	t *testing.T,
	status int,
	runId string,
	runName string,
	artifactId string,
	runResultStrings ...string) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.URL.Path {
		case "/ras/run":
			WriteMockRasRunsResponse(t, w, r, status, runName, runResultStrings...)

		case fmt.Sprintf(`/ras/run/%s/artifacts`, runId):
			// mock response for listing artifacts

		case fmt.Sprintf(`/ras/run/%s/artifacts/%s`, runId, artifactId):
			// mock response for downloading a single artifact
		}

		w.WriteHeader(status)
	}))

	return server
}

// ------------------------------------------------------------------

func TestRunsDownloadOfRunNameWhichDoesNotExistDisplaysNoArtifactsToDownload(t *testing.T) {
	// Given ...
	runName := "garbage"
	runId := "U1"
	artifactId := "artifact1"
	server := NewRunsDownloadServletMock(t, http.StatusOK, runId, runName, artifactId)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := utils.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	if err != nil {
		assert.Fail(t, "Garbage runname value should not have failed "+err.Error())
	} else {
		textGotBack := mockConsole.ReadText()
		want := "No artifacts to download.\n"
		assert.Equal(t, textGotBack, want)
	}

}

func TestFailingGetArtifactsRequestReturnsError(t *testing.T) {

	// Given...
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	runName := "garbage"
	mockConsole := utils.NewMockConsole()
	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()
	mockFileSystem := utils.NewMockFileSystem()


	// When...
	err := DownloadArtifacts(runName, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	assert.Contains(t, err.Error(), "GAL1068")
}
