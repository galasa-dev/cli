/*
 * Copyright contributors to the Galasa project
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
			 "result" : "LongResultString",
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
			"result" : "Fassed",
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
	runId string,
	runName string,
	artifactList []MockArtifact,
	runResultStrings []string) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

		acceptHeader := req.Header.Get("Accept")
		switch req.URL.Path {
		case "/ras/runs":
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsResponse(t, writer, req, runName, runResultStrings)

		case fmt.Sprintf(`/ras/runs/%s/artifacts`, runId):
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsArtifactsResponse(t, writer, req, artifactList)
		}

		runsFilesEndpoint := fmt.Sprintf(`/ras/runs/%s/files`, runId)
		if strings.HasPrefix(req.URL.Path, runsFilesEndpoint) {
			for _, artifact := range artifactList {
				if req.URL.Path == (runsFilesEndpoint + artifact.path) {
					WriteMockRasRunsFilesResponse(t, writer, req, artifact.path)
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
	server := NewRunsDownloadServletMock(t, http.StatusOK, runId, runName, []MockArtifact{*dummyTxtArtifact}, []string{RUN_U27})
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := utils.NewOverridableMockFileSystem()

	mockFile := utils.MockFile{}
	mockFile.VirtualFunction_Write = func(content []byte) (int, error) {
		return 0, errors.New("simulating failed file write")
	}

	mockFileSystem.VirtualFunction_Create = func(path string) (io.Writer, error) {
		return &mockFile, nil
	}

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	assert.Contains(t, err.Error(), "GAL1042")
}

func TestRunsDownloadFailingFileCreationReturnsError(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx543xxx"
	forceDownload := false

	dummyTxtArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)
	server := NewRunsDownloadServletMock(t, http.StatusOK, runId, runName, []MockArtifact{*dummyTxtArtifact}, []string{RUN_U27})
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := utils.NewOverridableMockFileSystem()

	mockFileSystem.VirtualFunction_Create = func(path string) (io.Writer, error) {
		return nil, errors.New("simulating failed folder creation")
	}

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	assert.Contains(t, err.Error(), "GAL1042")
}

func TestRunsDownloadFailingFolderCreationReturnsError(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx543xxx"
	forceDownload := false

	dummyTxtArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)
	server := NewRunsDownloadServletMock(t, http.StatusOK, runId, runName, []MockArtifact{*dummyTxtArtifact}, []string{RUN_U27})
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := utils.NewOverridableMockFileSystem()

	mockFileSystem.VirtualFunction_MkdirAll = func(path string) error {
		return errors.New("simulating failed folder creation")
	}

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	assert.Contains(t, err.Error(), "GAL1041")
}

func TestRunsDownloadExistingFileForceOverwritesMultipleArtifactsToFileSystem(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx543xxx"
	forceDownload := true

	dummyTxtArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)
	dummyRunLog := NewMockArtifact("/run.log", "text/plain", 203)
	mockArtifacts := []MockArtifact{
		*dummyTxtArtifact,
		*dummyRunLog,
	}
	server := NewRunsDownloadServletMock(t, http.StatusOK, runId, runName, mockArtifacts, []string{RUN_U27})
	defer server.Close()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()
	mockConsole := utils.NewMockConsole()

	mockFileSystem := utils.NewMockFileSystem()
	mockFileSystem.WriteTextFile(runName+dummyTxtArtifact.path, "dummy text file")
	mockFileSystem.WriteTextFile(runName+dummyRunLog.path, "dummy log")

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	assert.Nil(t, err)

	downloadedTxtArtifactContents, err := mockFileSystem.ReadTextFile(runName + dummyTxtArtifact.path)
	assert.Nil(t, err)

	downloadedRunLogContents, err := mockFileSystem.ReadTextFile(runName + dummyRunLog.path)
	assert.Nil(t, err)

	assert.Equal(t, dummyTxtArtifact.path, downloadedTxtArtifactContents)
	assert.Equal(t, dummyRunLog.path, downloadedRunLogContents)
}

func TestRunsDownloadExistingFileNoForceReturnsError(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx543xxx"
	forceDownload := false

	mockArtifacts := []MockArtifact{
		*NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024),
		*NewMockArtifact("/run.log", "text/plain", 203),
	}
	server := NewRunsDownloadServletMock(t, http.StatusOK, runId, runName, mockArtifacts, []string{RUN_U27})
	defer server.Close()

	apiServerUrl := server.URL
	mockConsole := utils.NewMockConsole()
	mockTimeService := utils.NewMockTimeService()

	mockFileSystem := utils.NewMockFileSystem()
	mockFileSystem.WriteTextFile(runName+"/dummy.txt", "dummy text file")
	mockFileSystem.WriteTextFile(runName+"/run.log", "dummy log")

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	assert.Contains(t, err.Error(), "GAL1036E")
}

func TestRunsDownloadWritesMultipleArtifactsToFileSystem(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx543xxx"
	forceDownload := false

	dummyTxtArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)
	dummyGzArtifact := NewMockArtifact("/artifacts/dummy.gz", "application/x-gzip", 342)
	dummyRunLogArtifact := NewMockArtifact("/run.log", "text/plain", 203)
	mockArtifacts := []MockArtifact{
		*dummyTxtArtifact,
		*dummyGzArtifact,
		*dummyRunLogArtifact,
	}
	server := NewRunsDownloadServletMock(t, http.StatusOK, runId, runName, mockArtifacts, []string{RUN_U27})
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := utils.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	downloadedTxtArtifactExists, _ := mockFileSystem.Exists(runName + dummyTxtArtifact.path)
	downloadedGzArtifactExists, _ := mockFileSystem.Exists(runName + dummyGzArtifact.path)
	downloadedRunLogArtifactExists, _ := mockFileSystem.Exists(runName + dummyRunLogArtifact.path)

	assert.Nil(t, err)
	assert.True(t, downloadedTxtArtifactExists)
	assert.True(t, downloadedGzArtifactExists)
	assert.True(t, downloadedRunLogArtifactExists)
}

func TestRunsDownloadWritesSingleArtifactToFileSystem(t *testing.T) {
	// Given ...
	runName := "U27"
	runId := "xxx543xxx"
	forceDownload := false

	dummyArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)
	server := NewRunsDownloadServletMock(t, http.StatusOK, runId, runName, []MockArtifact{*dummyArtifact}, []string{RUN_U27})
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := utils.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	downloadedArtifactExists, _ := mockFileSystem.Exists(runName + dummyArtifact.path)

	assert.Nil(t, err)
	assert.True(t, downloadedArtifactExists)
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
	mockFileSystem := utils.NewMockFileSystem()
	forceDownload := false

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

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
	mockFileSystem := utils.NewMockFileSystem()
	forceDownload := false

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	assert.Contains(t, err.Error(), "GAL1073")
}

func TestRunsDownloadWritesSingleArtifactToFileSystemMultipleReRuns(t *testing.T) {
	// Given ...
	runName := "U27"
	runIds := []string{"xxx543xxx", "xxx987xxx"}
	queuedTime := "2023-05-10_06:00:13"

	runResultStrings := []string{RUN_U27, RUN_U27V2}

	forceDownload := true
	dummyTxtArtifact := NewMockArtifact("/artifacts/dummy.txt", "text/plain", 1024)
	dummyGzArtifact := NewMockArtifact("/artifacts/dummy.gz", "application/x-gzip", 342)
	dummyRunLogArtifact := NewMockArtifact("/run.log", "text/plain", 203)
	mockArtifacts := []MockArtifact{
		*dummyTxtArtifact,
		*dummyGzArtifact,
		*dummyRunLogArtifact,
	}

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

		acceptHeader := req.Header.Get("Accept")
		if req.URL.Path == "/ras/runs" {
			assert.Equal(t, "application/json", acceptHeader, "Expected Accept: application/json header, got: %s", acceptHeader)
			WriteMockRasRunsResponse(t, writer, req, runName, runResultStrings)
		} else {
			for _, runId := range runIds {
				runsFilesEndpoint := fmt.Sprintf(`/ras/runs/%s/files`, runId)
				if strings.HasPrefix(req.URL.Path, runsFilesEndpoint) {
					for _, artifact := range mockArtifacts {
						if req.URL.Path == (runsFilesEndpoint + artifact.path) {
							WriteMockRasRunsFilesResponse(t, writer, req, artifact.path)
						}
					}
				}
			}
		}
	}))
	defer server.Close()

	mockConsole := utils.NewMockConsole()
	mockFileSystem := utils.NewMockFileSystem()

	apiServerUrl := server.URL
	mockTimeService := utils.NewMockTimeService()

	// When...
	err := DownloadArtifacts(runName, forceDownload, mockFileSystem, mockTimeService, mockConsole, apiServerUrl)

	// Then...
	// U27-1-2023-06-01_00:00:00 test did not finish
	// U27-2 					 test finished
	downloadedTxtArtifactExists1, _ := mockFileSystem.Exists(fmt.Sprintf("%s-1-%s%s", runName, queuedTime, dummyTxtArtifact.path))
	downloadedGzArtifactExists1, _ := mockFileSystem.Exists(fmt.Sprintf("%s-1-%s%s", runName, queuedTime, dummyGzArtifact.path))
	downloadedRunLogArtifactExists1, _ := mockFileSystem.Exists(fmt.Sprintf("%s-1-%s%s", runName, queuedTime, dummyRunLogArtifact.path))

	downloadedTxtArtifactExists2, _ := mockFileSystem.Exists(fmt.Sprintf("%s-3%s", runName, dummyTxtArtifact.path))
	downloadedGzArtifactExists2, _ := mockFileSystem.Exists(fmt.Sprintf("%s-3%s", runName, dummyGzArtifact.path))
	downloadedRunLogArtifactExists2, _ := mockFileSystem.Exists(fmt.Sprintf("%s-3%s", runName, dummyRunLogArtifact.path))

	assert.Nil(t, err)

	assert.True(t, downloadedTxtArtifactExists1)
	assert.True(t, downloadedGzArtifactExists1)
	assert.True(t, downloadedRunLogArtifactExists1)

	assert.True(t, downloadedTxtArtifactExists2)
	assert.True(t, downloadedGzArtifactExists2)
	assert.True(t, downloadedRunLogArtifactExists2)
}
