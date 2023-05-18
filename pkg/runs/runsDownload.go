/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"bufio"
	"context"
	"io"
	"os"
	"strings"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

// DownloadArtifacts - performs all the logic to implement the `galasactl runs download` command,
// but in a unit-testable manner.
func DownloadArtifacts(
	runName string,
	forceDownload bool,
	fileSystem utils.FileSystem,
	timeService utils.TimeService,
	console utils.Console,
	apiServerUrl string,
) error {

	var err error = nil
	var runs []galasaapi.Run
	var artifactPaths []string

	if err == nil {
		runs, err = GetRunsFromRestApi(runName, timeService, apiServerUrl)
		if err == nil {
			for _, run := range runs {
				runId := run.GetRunId()
				artifactPaths, err = GetArtifactPathsFromRestApi(runId, apiServerUrl)
				if err == nil {
					for _, artifactPath := range artifactPaths {
						artifactData, err := GetFileFromRestApi(runId, strings.TrimPrefix(artifactPath, "/"), apiServerUrl)
						if err != nil {
							return err
						} else {
							err = WriteArtifactToFileSystem(fileSystem, runName, artifactPath, artifactData, forceDownload)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}
	return err
}

// Retrieves the paths of all artifacts for a given test run using its runId.
func GetArtifactPathsFromRestApi(
	runId string,
	apiServerUrl string,
) ([]string, error) {

	var err error = nil
	var artifactPaths []string

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	var artifactsList []galasaapi.ArtifactIndexEntry
	artifactsList, httpResponse, err := restClient.ResultArchiveStoreAPIApi.
		GetRasRunArtifactList(context.Background(), runId).
		Execute()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_ARTIFACTS_FAILED, err.Error())
	} else {
		for _, artifact := range artifactsList {
			artifactPaths = append(artifactPaths, artifact.GetPath())
		}
	}

	defer httpResponse.Body.Close()
	return artifactPaths, err
}

// Writes an artifact to the host's file system, creating a new directory for the run's artifacts
// to be written to. Existing files are only overwritten if the --force option is used as part of
// the "runs download" command.
func WriteArtifactToFileSystem(
	fileSystem utils.FileSystem,
	runDirectory string,
	artifactPath string,
	fileDownloaded *os.File,
	shouldOverwrite bool) error {

	var err error = nil
	bufferCapacity := 1024

	pathParts := strings.Split(artifactPath, "/")
	fileName := pathParts[len(pathParts)-1]
	filePath := runDirectory + "/" + fileName

	// Check if a new file should be created or if an existing one should be overwritten.
	fileExists, err := fileSystem.Exists(filePath)
	if err == nil {
		if fileExists && !shouldOverwrite {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CANNOT_OVERWRITE_FILE, filePath)

		} else {
			err = fileSystem.MkdirAll(runDirectory)
			if err == nil {
				newFile, err := fileSystem.Create(filePath)
				if err == nil {

					// Set up a byte buffer to gradually read the downloaded file to avoid out-of-memory issues.
					reader := bufio.NewReader(fileDownloaded)
					buffer := make([]byte, bufferCapacity)

					// Read a chunk of the downloaded file into the buffer and write the buffer's contents to the
					// newly-created file.
					for {
						bytesRead, err := reader.Read(buffer)
						if err != nil {
							if err == io.EOF {
								// There was nothing to read.
								break
							}
							return err
						}

						_, err = newFile.Write(buffer[:bytesRead])
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return err
}

// Retrieves an artifact for a given test run using its runId from the ecosystem API.
func GetFileFromRestApi(
	runId string,
	artifactPath string,
	apiServerUrl string) (*os.File, error) {

	var err error = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	fileDownloaded, httpResponse, err := restClient.ResultArchiveStoreAPIApi.
		GetRasRunArtifactByPath(context.Background(), runId, artifactPath).
		Execute()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DOWNLOADING_ARTIFACT_FAILED, artifactPath, err.Error())
	}

	defer httpResponse.Body.Close()
	return fileDownloaded, err
}
