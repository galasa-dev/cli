/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/galasa.dev/cli/pkg/files"
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
			 "queued" : "2023-05-10T06:00:13.043037Z",	
			 "startTime": "2023-05-10T06:00:36.159003Z",
			 "endTime": "2023-05-10T06:02:53.823338Z",
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

	RUN_U27 = `{
		 "runId": "xxx543xxx",
		 "testStructure": {
			 "runName": "U27",
			 "bundle": "myBun27",	
			 "testName": "myTestPackage.MyTest27",
			 "testShortName": "MyTestName27",	
			 "requestor": "unitTesting27",
			 "status" : "building",
			 "queued" : "2023-05-10T06:00:13.043037Z",	
			 "startTime": "2023-05-10T06:00:36.159003Z",
			 "methods": [{
				 "className": "myTestPackage27.MyTestName27",
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

	RUN_U27V2 = `{
		"runId": "xxx987xxx",
		"testStructure": {
			"runName": "U27",
			"bundle": "myBun27",	
			"testName": "myTestPackage.MyTest27",
			"testShortName": "MyTestName27",	
			"requestor": "unitTesting27",
			"status" : "finished",
			"result" : "Passed",
			"queued" : "2023-05-10T06:00:13.043037Z",	
			"startTime": "2023-05-10T06:01:36.159003Z",
			"endTime": "2023-05-10T06:02:53.823338Z",
			"methods": [{
				"className": "myTestPackage27.MyTestName27",
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
			"artifactPath": "myPathToArtifact2",	
			"contentType":	"application/json"
		}]
   }`

	RUN_U27V3 = `{
		"runId": "xxx1234xxx",
		"testStructure": {
			"runName": "U27",
			"bundle": "myBun27",	
			"testName": "myTestPackage.MyTest27",
			"testShortName": "MyTestName27",	
			"requestor": "unitTesting27",
			"status" : "finished",
			"result" : "EnvFail",
			"queued" : "2022-05-10T04:00:13.043037Z",	
			"startTime": "2022-05-10T04:01:36.159003Z",
			"methods": [{
				"className": "myTestPackage27.MyTestName27",
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
			"artifactPath": "myPathToArtifact2",	
			"contentType":	"application/json"
		}]
	}`

	RUN_U27V4 = `{
		"runId": "xxx4321xxx",
		"testStructure": {
			"runName": "U27",
			"bundle": "myBun27",	
			"testName": "myTestPackage.MyTest27",
			"testShortName": "MyTestName27",	
			"requestor": "unitTesting27",
			"status" : "finished",
			"result" : "EnvFail",
			"queued" : "2022-05-10T04:00:13.043037Z",	
			"startTime": "2022-05-10T04:10:36.159003Z",
			"methods": [{
				"className": "myTestPackage27.MyTestName27",
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
			"artifactPath": "myPathToArtifact2",	
			"contentType":	"application/json"
		}]
	}`

	RUN_U27V5 = `{
		"runId": "xxx2468xxx",
		"testStructure": {
			"runName": "U27",
			"bundle": "myBun27",	
			"testName": "myTestPackage.MyTest27",
			"testShortName": "MyTestName27",	
			"requestor": "unitTesting27",
			"status" : "Building",
			"queued" : "2022-05-10T04:00:13.043037Z",	
			"startTime": "2022-05-10T04:10:36.159003Z",
			"methods": [{
				"className": "myTestPackage27.MyTestName27",
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
			"artifactPath": "myPathToArtifact2",	
			"contentType":	"application/json"
		}]
	}`
)

type MockArtifact struct {
	path        string
	contentType string
	size        int
}

func NewMockArtifact(mockPath string, mockContentType string, mockSize int) *MockArtifact {
	mockArtifact := MockArtifact{path: mockPath, contentType: mockContentType, size: mockSize}
	return &mockArtifact
}

//------------------------------------------------------------------
// Helper methods
//------------------------------------------------------------------

// Sets a response for requests to /ras/runs
func WriteMockRasRunsResponse(
	t *testing.T,
	writer http.ResponseWriter,
	req *http.Request,
	runName string,
	runResultStrings []string) {

	writer.Header().Set("Content-Type", "application/json")

	values := req.URL.Query()
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

	writer.Write([]byte(fmt.Sprintf(`
	{
		"pageNumber": 1,
		"pageSize": 1,
		"numPages": 1,
		"amountOfRuns": %d,
		"runs":[ %s ]
	}`, len(runResultStrings), combinedRunResultStrings)))
}

// Sets a response for requests to /ras/runs/{runId}/artifacts
func WriteMockRasRunsArtifactsResponse(
	t *testing.T,
	writer http.ResponseWriter,
	req *http.Request,
	artifactsList []MockArtifact) {

	writer.Header().Set("Content-Type", "application/json")

	artifactsListJsonString := ""
	for index, artifact := range artifactsList {
		if index > 0 {
			artifactsListJsonString += ","
		}
		artifactsListJsonString += fmt.Sprintf(`{
			"path": "%s",
			"contentType": "%s",
			"size": "%d"
		}`, artifact.path, artifact.contentType, artifact.size)
	}

	writer.Write([]byte(fmt.Sprintf(`
	[ %s ]
	`, artifactsListJsonString)))

}

// Sets a response for requests to /ras/runs/{runId}/files/{artifactPath}
func WriteMockRasRunsFilesResponse(
	t *testing.T,
	writer http.ResponseWriter,
	req *http.Request,
	desiredContents string) {

	writer.Header().Set("Content-Disposition", "attachment")
	writer.Write([]byte(desiredContents))
}

// Creates a new mock server to handle requests from test methods
func NewRunsDownloadServletMock(
	t *testing.T,
	status int,
	runName string,
	runResultStrings []string,
	runs map[string][]MockArtifact) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

		acceptHeader := req.Header.Get("Accept")
		if req.URL.Path == "/ras/runs" {
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsResponse(t, writer, req, runName, runResultStrings)

		} else {
			for runId, artifacts := range runs {
				if req.URL.Path == "/ras/runs/"+runId+"/artifacts" {
					assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
					WriteMockRasRunsArtifactsResponse(t, writer, req, artifacts)
				} else {
					runsFilesEndpoint := fmt.Sprintf(`/ras/runs/%s/files`, runId)
					if strings.HasPrefix(req.URL.Path, runsFilesEndpoint) {
						for _, artifact := range artifacts {
							if req.URL.Path == (runsFilesEndpoint + artifact.path) {
								WriteMockRasRunsFilesResponse(t, writer, req, artifact.path)
							}
						}
					}
				}
			}
		}

		writer.WriteHeader(status)
	}))

	return server
}

//------------------------------------------------------------------
// Test methods
//------------------------------------------------------------------

func TestRunsDownloadFailingFileWriteReturnsError(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx543xxx"
	forceDownload := false

	dummyTxtArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)

	runs := make(map[string][]MockArtifact, 0)
	runs[runId] = []MockArtifact{*dummyTxtArtifact}

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, []string{RUN_U27}, runs)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := files.NewOverridableMockFileSystem()

	mockFile := files.MockFile{}
	mockFile.VirtualFunction_Write = func(content []byte) (int, error) {
		return 0, errors.New("simulating failed file write")
	}

	mockFileSystem.VirtualFunction_Create = func(path string) (io.Writer, error) {
		return &mockFile, nil
	}

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	assert.Contains(t, err.Error(), "GAL1042")
}

func TestRunsDownloadFailingFileCreationReturnsError(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx543xxx"
	forceDownload := false

	dummyTxtArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)

	runs := make(map[string][]MockArtifact, 0)
	runs[runId] = []MockArtifact{*dummyTxtArtifact}

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, []string{RUN_U27}, runs)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := files.NewOverridableMockFileSystem()

	mockFileSystem.VirtualFunction_Create = func(path string) (io.Writer, error) {
		return nil, errors.New("simulating failed folder creation")
	}

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	assert.Contains(t, err.Error(), "GAL1042")
}

func TestRunsDownloadFailingFolderCreationReturnsError(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx543xxx"
	forceDownload := false

	dummyTxtArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)

	runs := make(map[string][]MockArtifact, 0)
	runs[runId] = []MockArtifact{*dummyTxtArtifact}

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, []string{RUN_U27}, runs)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := files.NewOverridableMockFileSystem()

	mockFileSystem.VirtualFunction_MkdirAll = func(path string) error {
		return errors.New("simulating failed folder creation")
	}

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	assert.Contains(t, err.Error(), "GAL1041")
}

func TestRunsDownloadExistingFileForceOverwritesMultipleArtifactsToFileSystem(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx987xxx"
	forceDownload := true

	dummyTxtArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)
	dummyRunLog := NewMockArtifact("/run.log", "text/plain", 203)
	mockArtifacts := []MockArtifact{
		*dummyTxtArtifact,
		*dummyRunLog,
	}

	runs := make(map[string][]MockArtifact, 0)
	runs[runId] = mockArtifacts

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, []string{RUN_U27V2}, runs)
	defer server.Close()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()
	mockConsole := utils.NewMockConsole()

	mockFileSystem := files.NewMockFileSystem()
	mockFileSystem.WriteTextFile(runName+dummyTxtArtifact.path, "dummy text file")
	mockFileSystem.WriteTextFile(runName+dummyRunLog.path, "dummy log")

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	assert.Nil(t, err)

	downloadedTxtArtifactContents, err := mockFileSystem.ReadTextFile(runName + dummyTxtArtifact.path)
	assert.Nil(t, err)

	downloadedRunLogContents, err := mockFileSystem.ReadTextFile(runName + dummyRunLog.path)
	assert.Nil(t, err)

	assert.Equal(t, dummyTxtArtifact.path, downloadedTxtArtifactContents)
	assert.Equal(t, dummyRunLog.path, downloadedRunLogContents)

	textGotBack := mockConsole.ReadText()
	assert.Contains(t, textGotBack, runName)
}

func TestRunsDownloadExistingFileNoForceReturnsError(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx987xxx"
	forceDownload := false

	mockArtifacts := []MockArtifact{
		*NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024),
		*NewMockArtifact("/run.log", "text/plain", 203),
	}

	runs := make(map[string][]MockArtifact, 0)
	runs[runId] = mockArtifacts

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, []string{RUN_U27V2}, runs)
	defer server.Close()

	apiServerUrl := server.URL
	mockConsole := utils.NewMockConsole()
	mockTimeService := utils.NewMockTimeService()

	mockFileSystem := files.NewMockFileSystem()
	mockFileSystem.WriteTextFile(runName+"/dummy.txt", "dummy text file")
	mockFileSystem.WriteTextFile(runName+"/run.log", "dummy log")

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1036E")
}

func TestRunsDownloadWritesMultipleArtifactsToFileSystem(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx987xxx"
	forceDownload := false

	dummyTxtArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)
	dummyGzArtifact := NewMockArtifact("/artifacts/dummy.gz", "application/x-gzip", 342)
	dummyRunLogArtifact := NewMockArtifact("/run.log", "text/plain", 203)
	mockArtifacts := []MockArtifact{
		*dummyTxtArtifact,
		*dummyGzArtifact,
		*dummyRunLogArtifact,
	}

	runs := make(map[string][]MockArtifact, 0)
	runs[runId] = mockArtifacts

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, []string{RUN_U27V2}, runs)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := files.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	downloadedTxtArtifactExists, _ := mockFileSystem.Exists(runName + dummyTxtArtifact.path)
	downloadedGzArtifactExists, _ := mockFileSystem.Exists(runName + dummyGzArtifact.path)
	downloadedRunLogArtifactExists, _ := mockFileSystem.Exists(runName + dummyRunLogArtifact.path)

	assert.Nil(t, err)
	assert.True(t, downloadedTxtArtifactExists)
	assert.True(t, downloadedGzArtifactExists)
	assert.True(t, downloadedRunLogArtifactExists)

	textGotBack := mockConsole.ReadText()
	assert.Contains(t, textGotBack, runName)
}

func TestRunsDownloadWritesSingleArtifactToFileSystem(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx987xxx"
	forceDownload := false

	dummyArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)

	runs := make(map[string][]MockArtifact, 0)
	runs[runId] = []MockArtifact{*dummyArtifact}

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, []string{RUN_U27V2}, runs)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := files.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	downloadedArtifactExists, _ := mockFileSystem.Exists(runName + dummyArtifact.path)

	assert.Nil(t, err)
	assert.True(t, downloadedArtifactExists)

	textGotBack := mockConsole.ReadText()
	assert.Contains(t, textGotBack, runName)
}

func TestFailingGetFileRequestReturnsError(t *testing.T) {

	// Given...
	runName := "U1"
	runId := "xxx876xxx"
	dummyArtifact := NewMockArtifact("/artifacts/dummy.gz", "application/x-gzip", 30)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/ras/runs":
			WriteMockRasRunsResponse(t, writer, req, runName, []string{RUN_U1})

		case fmt.Sprintf(`/ras/runs/%s/artifacts`, runId):
			WriteMockRasRunsArtifactsResponse(t, writer, req, []MockArtifact{*dummyArtifact})

		case fmt.Sprintf(`/ras/runs/%s/files%s`, runId, dummyArtifact.path):
			// Make the request to download an artifact fail
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()
	mockFileSystem := files.NewMockFileSystem()
	forceDownload := false

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	assert.Contains(t, err.Error(), "GAL1074")
}

func TestFailingGetArtifactsRequestReturnsError(t *testing.T) {

	// Given...
	runName := "U1"
	runId := "xxx876xxx"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/ras/runs":
			WriteMockRasRunsResponse(t, writer, req, runName, []string{RUN_U1})

		case fmt.Sprintf(`/ras/runs/%s/artifacts`, runId):
			// Make the request to list artifacts fail
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()
	mockFileSystem := files.NewMockFileSystem()
	forceDownload := false

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	assert.Contains(t, err.Error(), "GAL1073")
}

func TestRunsDownloadMultipleReRunsWithCorrectOrderFolders(t *testing.T) {
	// Given ...
	runName := "U27"
	runId1 := "xxx543xxx"
	runId2 := "xxx987xxx"

	runResultStrings := []string{RUN_U27, RUN_U27V2}

	forceDownload := true
	dummyTxtArtifactRunId1 := NewMockArtifact("/artifacts/dummy1.txt", "text/plain", 1024)
	mockArtifactsRunId1 := []MockArtifact{
		*dummyTxtArtifactRunId1,
	}

	dummyTxtArtifactRunId2 := NewMockArtifact("/artifacts/dummy2.txt", "text/plain", 1024)
	mockArtifactsRunId2 := []MockArtifact{
		*dummyTxtArtifactRunId2,
	}

	runs := make(map[string][]MockArtifact, 0)
	runs[runId1] = mockArtifactsRunId1
	runs[runId2] = mockArtifactsRunId2

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, runResultStrings, runs)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := files.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()
	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	// U27-1-2023-2023-05-10T06:00:13 	(test did not finish)
	// U27-2 					 		(test finished)
	run1FolderName := runName + "-1-" + mockTimeService.Now().Format("2006-01-02_15:04:05")
	run2FolderName := runName + "-2"

	downloadedTxtArtifactExists1, _ := mockFileSystem.Exists(run1FolderName + dummyTxtArtifactRunId1.path)
	downloadedTxtArtifactExists2, _ := mockFileSystem.Exists(run2FolderName + dummyTxtArtifactRunId2.path)
	assert.Nil(t, err)

	assert.True(t, downloadedTxtArtifactExists1)
	assert.True(t, downloadedTxtArtifactExists2)

	textGotBack := mockConsole.ReadText()
	assert.Contains(t, textGotBack, run1FolderName)
	assert.Contains(t, textGotBack, run2FolderName)
}

func TestRunsDownloadMultipleSetsOfUnrelatedReRunsWithCorrectOrderFolders(t *testing.T) {
	// Given ...
	runName := "U27"
	runId1a := "xxx543xxx"
	runId1b := "xxx987xxx"

	runId2a := "xxx1234xxx"
	runId2b := "xxx4321xxx"

	runResultStrings := []string{RUN_U27, RUN_U27V2, RUN_U27V3, RUN_U27V4}

	forceDownload := true

	dummyTxtArtifactRunId1a := NewMockArtifact("/artifacts/dummy1a.txt", "text/plain", 1024)
	mockArtifactsRunId1a := []MockArtifact{
		*dummyTxtArtifactRunId1a,
	}

	dummyTxtArtifactRunId1b := NewMockArtifact("/artifacts/dummy1b.txt", "text/plain", 1024)
	mockArtifactsRunId1b := []MockArtifact{
		*dummyTxtArtifactRunId1b,
	}

	dummyTxtArtifactRunId2a := NewMockArtifact("/artifacts/dummy2a.txt", "text/plain", 1024)
	mockArtifactsRunId2a := []MockArtifact{
		*dummyTxtArtifactRunId2a,
	}

	dummyTxtArtifactRunId2b := NewMockArtifact("/artifacts/dummy2b.txt", "text/plain", 1024)
	mockArtifactsRunId2b := []MockArtifact{
		*dummyTxtArtifactRunId2b,
	}

	runs := make(map[string][]MockArtifact, 0)
	runs[runId1a] = mockArtifactsRunId1a
	runs[runId1b] = mockArtifactsRunId1b
	runs[runId2a] = mockArtifactsRunId2a
	runs[runId2b] = mockArtifactsRunId2b

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, runResultStrings, runs)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := files.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()
	mockTimeService.Sleep(time.Second)

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	// U27-1-2023-05-10T06:00:13 	(test did not finish)
	// U27-2 					 	(test finished)
	// U27-1-2022-05-10T04:00:13 	(test did not finish)
	// U27-2-2022-05-10T04:00:13    (test did not finish)
	run1aFolderName := runName + "-1-" + mockTimeService.Now().Format("2006-01-02_15:04:05")
	run1bFolderName := runName + "-2"
	run2aFolderName := runName + "-1"
	run2bFolderName := runName + "-2"

	downloadedTxtArtifactExists1a, _ := mockFileSystem.Exists(runName + "-1-" + mockTimeService.Now().Format("2006-01-02_15:04:05") + dummyTxtArtifactRunId1a.path)
	downloadedTxtArtifactExists1b, _ := mockFileSystem.Exists(runName + "-2" + dummyTxtArtifactRunId1b.path)

	downloadedTxtArtifactExists2a, _ := mockFileSystem.Exists(runName + "-1" + dummyTxtArtifactRunId2a.path)
	downloadedTxtArtifactExists2b, _ := mockFileSystem.Exists(runName + "-2" + dummyTxtArtifactRunId2b.path)

	assert.Nil(t, err)

	assert.True(t, downloadedTxtArtifactExists1a)
	assert.True(t, downloadedTxtArtifactExists1b)
	assert.True(t, downloadedTxtArtifactExists2a)
	assert.True(t, downloadedTxtArtifactExists2b)

	textGotBack := mockConsole.ReadText()
	assert.Contains(t, textGotBack, run1aFolderName)
	assert.Contains(t, textGotBack, run1bFolderName)
	assert.Contains(t, textGotBack, run2aFolderName)
	assert.Contains(t, textGotBack, run2bFolderName)
}

func TestRunsDownloadWithValidRunNameNoArtifacts(t *testing.T) {
	// Given ...
	runName := "U27"

	runResultStrings := []string{}

	forceDownload := true

	runs := make(map[string][]MockArtifact, 0)

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, runResultStrings, runs)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := files.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")
	// Then...

	assert.Contains(t, err.Error(), "GAL1083E")
	assert.Contains(t, err.Error(), runName)
	assert.Contains(t, err.Error(), "No artifacts")

}

func TestRunsDownloadWithInvalidRunName(t *testing.T) {
	// Given ...
	runName := "garbage"

	runResultStrings := []string{}

	forceDownload := true

	runs := make(map[string][]MockArtifact, 0)

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, runResultStrings, runs)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := files.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...

	assert.Contains(t, err.Error(), "GAL1075E")
	assert.Contains(t, err.Error(), "garbage")

}

func TestRunsDownloadAddsTimestampToFolderIfRunNotFinished(t *testing.T) {
	// Given ...
	runName := "U27"
	runId1 := "xxx2468xxx"

	runResultStrings := []string{RUN_U27V5}

	forceDownload := true
	dummyTxtArtifactRunId1 := NewMockArtifact("/artifacts/dummy1.txt", "text/plain", 1024)
	mockArtifactsRunId1 := []MockArtifact{
		*dummyTxtArtifactRunId1,
	}

	runs := make(map[string][]MockArtifact, 0)
	runs[runId1] = mockArtifactsRunId1

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, runResultStrings, runs)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := files.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()
	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, ".")

	// Then...
	run1FolderName := runName + "-" + mockTimeService.Now().Format("2006-01-02_15:04:05")

	downloadedTxtArtifactExists1, _ := mockFileSystem.Exists(run1FolderName + dummyTxtArtifactRunId1.path)

	assert.Nil(t, err)

	assert.True(t, downloadedTxtArtifactExists1)

	textGotBack := mockConsole.ReadText()
	assert.Contains(t, textGotBack, run1FolderName)
}

func TestRunsDownloadWritesSingleArtifactToDestinationFolder(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx987xxx"
	forceDownload := false

	dummyArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)

	runs := make(map[string][]MockArtifact, 0)
	runs[runId] = []MockArtifact{*dummyArtifact}

	server := NewRunsDownloadServletMock(t, http.StatusOK, runName, []string{RUN_U27V2}, runs)
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := files.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl, "/myfolder")

	// Then...
	downloadedArtifactExists, _ := mockFileSystem.Exists("/myfolder/" + runName + dummyArtifact.path)

	assert.Nil(t, err)
	assert.True(t, downloadedArtifactExists)

	textGotBack := mockConsole.ReadText()
	assert.Contains(t, textGotBack, "GAL2501I")
	assert.Contains(t, textGotBack, "/myfolder/"+runName)
}
