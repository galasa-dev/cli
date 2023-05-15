/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"io/ioutil"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

// DownloadArtifacts - performs all the logic to implement the `galasactl runs download` command,
// but in a unit-testable manner.
func DownloadArtifacts(
	runName string,
	fileSystem utils.FileSystem,
	timeService utils.TimeService,
	console utils.Console,
	apiServerUrl string,
) error {

	var err error
	var runJson []galasaapi.Run
	var artifacts []string

	if err == nil {

		runJson, err = GetRunsFromRestApi(runName, timeService, apiServerUrl)
		if err == nil {
			for _, run := range runJson {
				runId := run.GetRunId()
				artifacts, err = GetArtifactPathsFromRestApi(runId, apiServerUrl)
				if err == nil {
					for _, artifact := range artifacts {
						data, err := GetFileFromRestApi(runId, artifact, apiServerUrl)
						if err == nil {
							err = WriteFileToFileSystem(fileSystem, strings.Split(artifact, "/")[len(strings.Split(artifact, "/"))-1], data)
						}
					}
				}
			}
		}
	}
	// href="/ras/runs/0cff27025c9666fda73b82ab710011e9/files/artifacts/framework/cps_record.properties"
	return err
}

func GetArtifactPathsFromRestApi(
	runId string,
	apiServerUrl string,
) ([]string, error) {

	var err error = nil
	var results []string

	var context context.Context = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	var runData []galasaapi.ArtifactIndexEntry
	var httpResponse *http.Response
	runData, httpResponse, err = restClient.ResultArchiveStoreAPIApi.
	GetRasRunArtifactList(context, runId).
		Execute()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
	} else {
		if httpResponse.StatusCode != HTTP_STATUS_CODE_OK {
			body, _  := io.ReadAll(httpResponse.Body)
			bodyStr := string(body)
			err = galasaErrors.ThrowAPIError(bodyStr)
		} else {
			for _, artifact := range runData {
				results = append(results, artifact.GetPath())
			}
		}
	}

	return results, err
}

func WriteFileToFileSystem(fileSystem utils.FileSystem, filename string, byteStream []byte) error {
	var err error = nil
	err = fileSystem.WriteBinaryFile(filename, byteStream)
	return err
}

// Retrieves test runs from the ecosystem API that match a given runName.
// Multiple test runs can be returned as the runName is not unique.
func GetFileFromRestApi(
	runId string,
	artifactPath string,
	apiServerUrl string,
) ([]byte, error) {

	var err error = nil
	var result os.File 

	var context context.Context = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)


	var httpResponse *http.Response
	result, httpResponse, err = restClient.ResultArchiveStoreAPIApi.
		GetRasRunArtifactByPath(context, runId, artifactPath).
		Execute()
	size , _ := result.Stat()

	var resultbyte []byte = make([]byte, size.Size()) 
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
	} else {
		if httpResponse.StatusCode != HTTP_STATUS_CODE_OK {
			body, _  := io.ReadAll(httpResponse.Body)
			bodyStr := string(body)
			err = galasaErrors.ThrowAPIError(bodyStr)
		}else{
			readbytes , err := result.Read(resultbyte)
			if readbytes != int(size.Size()) || err !=nil{
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
			}
		}
	}

	return resultbyte, err
}
