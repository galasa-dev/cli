/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"bufio"
	"context"
	"io"
	"log"
	"path/filepath"
	"strconv"
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

	runs, err = GetRunsFromRestApi(runName, 0, 0, timeService, apiServerUrl)
	if err == nil {
		if len(runs) > 1 {
			// get list of runs that are reRuns - get list of runs that are reRuns of each other
			// create a map of lists of reRuns - key is queued time, value is the run

			reRunsByQueuedTime := make(map[string][]galasaapi.Run, 0)
			queuedTime := ""
			for _, run := range runs {
				queuedTime = run.TestStructure.GetQueued()
				runsWithSameQueuedTime := reRunsByQueuedTime[queuedTime]
				runsWithSameQueuedTime = append(runsWithSameQueuedTime, run)
				reRunsByQueuedTime[queuedTime] = runsWithSameQueuedTime
			}

			// SORTING BY START TIME TO DO - on api side not in cli

			for _, reRunsList := range reRunsByQueuedTime {
				for count, reRun := range reRunsList {
					directoryName := nameArtifactDownloadDirectory(reRun, count)
					err = downloadArtifactsToDirectory(apiServerUrl, directoryName, reRun, fileSystem, forceDownload)
				}
			}

		} else {
			for _, run := range runs {
				err = downloadArtifactsToDirectory(apiServerUrl, runName, run, fileSystem, forceDownload)
			}
			//err = downloadArtifactsToDirectory(apiServerUrl, runName, runs[0], fileSystem, forceDownload)
		}

	}
	return err
}

func nameArtifactDownloadDirectory(reRun galasaapi.Run, count int) string {
	endTime := reRun.TestStructure.GetEndTime()
	runName := reRun.TestStructure.GetRunName()
	directoryName := runName
	if endTime == "" {
		queuedTimeString := strings.Replace(reRun.TestStructure.GetQueued(), "T", "_", -1)
		queuedTimeString = strings.Split(queuedTimeString, ".")[0]

		directoryName = runName + "-" + strconv.Itoa(count+1) + "-" + queuedTimeString
	} else {
		directoryName = runName + "-" + strconv.Itoa(count+1)
	}
	return directoryName
}

func downloadArtifactsToDirectory(apiServerUrl string,
	directoryName string,
	reRun galasaapi.Run,
	fileSystem utils.FileSystem,
	forceDownload bool) error {
	runId := reRun.GetRunId()
	artifactPaths, err := GetArtifactPathsFromRestApi(runId, apiServerUrl)
	if err == nil {
		log.Printf("Creating folder %s", directoryName)
		for _, artifactPath := range artifactPaths {
			var artifactData io.Reader
			artifactData, err = GetFileFromRestApi(runId, strings.TrimPrefix(artifactPath, "/"), apiServerUrl)
			if err == nil {
				err = WriteArtifactToFileSystem(fileSystem, directoryName, artifactPath, artifactData, forceDownload)
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
	fileDownloaded io.Reader,
	shouldOverwrite bool) error {

	var err error = nil

	pathParts := strings.Split(artifactPath, "/")
	fileName := pathParts[len(pathParts)-1]
	targetFilePath := filepath.Join(runDirectory, artifactPath)

	// Check if a new file should be created or if an existing one should be overwritten.
	fileExists, err := fileSystem.Exists(targetFilePath)
	if err == nil {
		if fileExists && !shouldOverwrite {

			// The --force flag was not provided, throw an error.
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CANNOT_OVERWRITE_FILE, targetFilePath)
		} else {
			// Set up the directory structure and artifact file on the host's file system
			var newFile io.Writer = nil
			newFile, err = CreateEmptyArtifactFile(fileSystem, targetFilePath)
			if err == nil {
				log.Printf("Writing artifact '%s' to '%s' on local file system", fileName, targetFilePath)

				err = TransferContent(fileDownloaded, newFile, targetFilePath)
				if err == nil {
					log.Printf("Artifact '%s' written to '%s' OK", fileName, targetFilePath)
				}
			}
		}
	}
	return err
}

// Writes the contents of a given source file into a given target file using a buffer
// to read and write the contents in chunks.
func TransferContent(sourceFile io.Reader, targetFile io.Writer, targetFilePath string) error {
	var err error = nil

	// Set buffer capacity to 1KB
	bufferCapacity := 1024

	// Set up a byte buffer to gradually read the downloaded file to avoid out-of-memory issues.
	reader := bufio.NewReader(sourceFile)
	buffer := make([]byte, bufferCapacity)
	for {
		var bytesRead int
		bytesRead, err = reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				// There was nothing else to read.
				err = nil
			}
			break
		}

		_, err = targetFile.Write(buffer[:bytesRead])
		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_WRITE_FILE, targetFilePath, err.Error())
			break
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
func GetFileFromRestApi(runId string, artifactPath string, apiServerUrl string) (io.Reader, error) {

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

	httpResponse.Body.Close()
	return fileDownloaded, err
}
