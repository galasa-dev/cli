/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
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
	apiServerUrl string) error {

	var err error = nil
	var runs []galasaapi.Run
	var artifactPaths []string

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
	return err
}

// Retrieves the paths of all artifacts for a given test run using its runId.
func GetArtifactPathsFromRestApi(runId string, apiServerUrl string) ([]string, error) {

	var err error = nil
	var artifactPaths []string
	log.Println("Retrieving artifact paths for the given run")

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
	log.Printf("%v artifact path(s) found", len(artifactPaths))

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
	fileName := pathParts[len(pathParts) - 1]
	targetFilePath := filepath.Join(runDirectory, fileName)

	// Check if a new file should be created or if an existing one should be overwritten.
	fileExists, err := fileSystem.Exists(targetFilePath)
	if err == nil {
		if fileExists && !shouldOverwrite {

			// The --force flag was not provided, throw an error.
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CANNOT_OVERWRITE_FILE, targetFilePath)
			return err
		}

		// Set up the directory structure and artifact file on the host's file system
		var newFile io.Writer = nil
		newFile, err = CreateEmptyArtifactFile(fileSystem, targetFilePath)
		if err == nil {
			log.Printf("Writing artifact '%s' to '%s' on local file system", fileName, targetFilePath)

			// Set up a byte buffer to gradually read the downloaded file to avoid out-of-memory issues.
			reader := bufio.NewReader(fileDownloaded)
			buffer := make([]byte, bufferCapacity)

			// Read a chunk of the downloaded file into the buffer and write the buffer's contents to the
			// newly-created file.
			for {
				bytesRead, err := reader.Read(buffer)
				if err != nil {
					if err == io.EOF {
						// There was nothing else to read.
						log.Printf("Artifact '%s' written to '%s' OK", fileName, targetFilePath)
						break
					}
					return err
				}

				_, err = newFile.Write(buffer[:bytesRead])
				if err != nil {
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_WRITE_FILE, targetFilePath, err.Error())
					return err
				}
			}
		}
	}
	return err
}

// Creates an empty file representing an artifact that is being written. Any parent directories that do not exist
// will be created. Returns the empty file that was created or an error if any file creation operations failed. 
func CreateEmptyArtifactFile(fileSystem utils.FileSystem, targetFilePath string) (io.Writer, error) {
	
	var err error = nil
	var newFile io.Writer = nil

	targetDirectoryPath := filepath.Dir(targetFilePath)
	err = fileSystem.MkdirAll(targetDirectoryPath)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_CREATE_FOLDERS, targetDirectoryPath, err.Error())
	} else {
		newFile, err = fileSystem.Create(targetFilePath)
		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_WRITE_FILE, targetFilePath, err.Error())
		}
	}
	return newFile, err
}

// Retrieves an artifact for a given test run using its runId from the ecosystem API.
func GetFileFromRestApi(runId string, artifactPath string, apiServerUrl string) (*os.File, error) {

	var err error = nil
	log.Printf("Downloading artifact '%s' from API server", artifactPath)

	// A HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	fileDownloaded, httpResponse, err := restClient.ResultArchiveStoreAPIApi.
		GetRasRunArtifactByPath(context.Background(), runId, artifactPath).
		Execute()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DOWNLOADING_ARTIFACT_FAILED, artifactPath, err.Error())
	} else {
		log.Printf("Downloaded artifact '%s' from API server OK", artifactPath)
	}

	defer httpResponse.Body.Close()
	return fileDownloaded, err
}
