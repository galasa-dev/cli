/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

// DownloadArtifacts - performs all the logic to implement the `galasactl runs download` command,
// but in a unit-testable manner.
func DownloadArtifacts(
	runName string,
	forceDownload bool,
	fileSystem files.FileSystem,
	timeService utils.TimeService,
	console utils.Console,
	apiServerUrl string,
	runDownloadTargetFolder string,
) error {

	var err error = nil
	var runs []galasaapi.Run

	if (err == nil) && (runName != "") {
		err = ValidateRunName(runName)
	}
	if err == nil {
		requestorParameter := ""
		resultParameter := ""
		fromAgeHours := 0
		toAgeHours := 0
		shouldGetActive := false
		runs, err = GetRunsFromRestApi(runName, requestorParameter, resultParameter, fromAgeHours, toAgeHours, shouldGetActive, timeService, apiServerUrl)
		if err == nil {
			if len(runs) > 1 {
				// get list of runs that are reRuns - get list of runs that are reRuns of each other
				// create a map of lists of reRuns - key is queued time, value is the run
				reRunsByQueuedTime := createMapOfReRuns(runs)

				err = downloadReRunArtfifacts(
					reRunsByQueuedTime,
					forceDownload,
					fileSystem,
					apiServerUrl,
					console,
					timeService,
					runDownloadTargetFolder,
				)

			} else if len(runs) == 1 {
				var folderName string
				folderName, err = nameDownloadFolder(runs[0], runName, timeService)
				if err == nil {
					err = downloadArtifactsToDirectory(apiServerUrl, folderName, runs[0], fileSystem, forceDownload, console, runDownloadTargetFolder)
				}
			} else {
				log.Printf("No artifacts to download for run: '%s'", runName)
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_NO_ARTIFACTS_TO_DOWNLOAD, runName)
			}
		}
	}

	return err
}

func downloadReRunArtfifacts(
	reRunsByQueuedTime map[string][]galasaapi.Run,
	forceDownload bool,
	fileSystem files.FileSystem,
	apiServerUrl string,
	console utils.Console,
	timeService utils.TimeService,
	runDownloadTargetFolder string,
) error {
	var err error = nil
	for _, reRunsList := range reRunsByQueuedTime {
		if err == nil {
			for reRunIndex, reRun := range reRunsList {
				if err == nil {
					directoryName := nameReRunArtifactDownloadDirectory(reRun, reRunIndex, timeService)
					err = downloadArtifactsToDirectory(
						apiServerUrl,
						directoryName,
						reRun,
						fileSystem,
						forceDownload,
						console,
						runDownloadTargetFolder,
					)
				}
			}
		}
	}
	return err
}

func createMapOfReRuns(runs []galasaapi.Run) map[string][]galasaapi.Run {
	// key = queued time
	// value = []runs
	// This function creates a map of queued times. Each time has a list of runs with that queued time.
	// helpful if we had two unrelated tests with the same runname, which both had a set of reruns.
	// e.g.
	// "2023-05-10T06:00:13.043037Z" = [reRun1, reRun2, reRun3]
	// "2023-01-10T16:10:13.043037Z" = [reRunA, reRunB, reRunC]
	reRunsByQueuedTime := make(map[string][]galasaapi.Run, 0)
	queuedTime := ""
	for _, run := range runs {
		queuedTime = run.TestStructure.GetQueued()
		runsWithSameQueuedTime := reRunsByQueuedTime[queuedTime]
		runsWithSameQueuedTime = append(runsWithSameQueuedTime, run)
		reRunsByQueuedTime[queuedTime] = runsWithSameQueuedTime
	}
	return reRunsByQueuedTime
}

func nameReRunArtifactDownloadDirectory(reRun galasaapi.Run, reRunIndex int, timeService utils.TimeService) string {
	result := reRun.TestStructure.GetResult()
	runName := reRun.TestStructure.GetRunName()
	directoryName := runName
	if result == "" {
		// Add timestamp of download to folder name
		downloadedTime := timeService.Now().Format("2006-01-02_15:04:05")

		directoryName = runName + "-" + strconv.Itoa(reRunIndex+1) + "-" + downloadedTime
	} else {
		directoryName = runName + "-" + strconv.Itoa(reRunIndex+1)
	}
	return directoryName
}

func nameDownloadFolder(run galasaapi.Run, runName string, timeService utils.TimeService) (string, error) {
	directoryName := runName
	var err error = nil
	result := run.TestStructure.GetResult()
	if result == "" {
		downloadedTime := timeService.Now().Format("2006-01-02_15:04:05")
		directoryName = runName + "-" + downloadedTime
	}
	return directoryName, err
}

func downloadArtifactsToDirectory(apiServerUrl string,
	directoryName string,
	run galasaapi.Run,
	fileSystem files.FileSystem,
	forceDownload bool,
	console utils.Console,
	runDownloadTargetFolder string,
) error {

	runId := run.GetRunId()

	// We want to base the directory we download to on the destination folder.
	// If the destination folder is "." (current folder/relative)
	// then ignore, as prefixing with a "." just adds noise when we log and
	// print out the path.
	if runDownloadTargetFolder != "." {
		directoryName = filepath.Join(runDownloadTargetFolder, directoryName)
	}

	filesWrittenOkCount := 0

	artifactPaths, err := GetArtifactPathsFromRestApi(runId, apiServerUrl)
	if err == nil {
		for _, artifactPath := range artifactPaths {
			if err == nil {
				var artifactData io.Reader
				var httpResponse *http.Response
				var isArtifactDataEmpty bool
				artifactData, isArtifactDataEmpty, httpResponse, err = GetFileFromRestApi(runId, strings.TrimPrefix(artifactPath, "/"), apiServerUrl)
				if err == nil {
					if !isArtifactDataEmpty {
						err = WriteArtifactToFileSystem(fileSystem, directoryName, artifactPath, artifactData, forceDownload, console)
						if err == nil {
							filesWrittenOkCount += 1
						}
					}
				}

				closeErr := httpResponse.Body.Close()
				// The first error is most important so needs preserving...
				if closeErr != nil && err == nil {
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_HTTP_RESPONSE_CLOSE_FAILED, closeErr.Error())
				}
			}
		}
	}

	// Write out the number of files downloaded to the folder xxx
	if filesWrittenOkCount > 0 {
		msg := fmt.Sprintf(
			galasaErrors.GALASA_INFO_FOLDER_DOWNLOADED_TO.Template,
			filesWrittenOkCount,
			directoryName,
		)
		consoleErr := console.WriteString(msg)
		// Console error is not as important to report as the original error if there was one.
		if consoleErr != nil && err == nil {
			err = consoleErr
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
	fileSystem files.FileSystem,
	runDirectory string,
	artifactPath string,
	fileDownloaded io.Reader,
	shouldOverwrite bool,
	console utils.Console) error {

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
				} else {
					log.Printf("Artifact '%s' write to '%s' failed. '%s'", fileName, targetFilePath, err.Error())
				}
			}
		}
	}
	return err
}

// Writes the contents of a given source file into a given target file using a buffer
// to read and write the contents in chunks.
func TransferContent(sourceFile io.Reader, targetFile io.Writer, targetFilePath string) error {

	log.Printf("TransferContent: Entered. targetFilePath: %s sourceFile type:%s", targetFilePath, reflect.TypeOf(sourceFile))

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
			log.Printf("TransferContent: Read(...) returned %s", err.Error())

			// Suppress end of file error, as we want to read until the end of the file.
			if err == io.EOF {
				// There was nothing else to read.
				err = nil
			}
			break
		}

		log.Printf("TransferContent: Read(...) returned %d bytes of data", bytesRead)

		_, err = targetFile.Write(buffer[:bytesRead])
		if err != nil {
			log.Printf("TransferContent: Write() failed. err=%s", err.Error())
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_WRITE_FILE, targetFilePath, err.Error())
			break
		}
	}

	if err != nil {
		log.Printf("TransferContent problem. File: %s. err: %s", targetFilePath, err.Error())
	}

	return err
}

// Creates an empty file representing an artifact that is being written. Any parent directories that do not exist
// will be created. Returns the empty file that was created or an error if any file creation operations failed.
func CreateEmptyArtifactFile(fileSystem files.FileSystem, targetFilePath string) (io.Writer, error) {

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

// GetFileFromRestApi Retrieves an artifact for a given test run using its runId from the ecosystem API.
// Note: The call leaves closing the http request as a responsibility of the caller.
func GetFileFromRestApi(runId string, artifactPath string, apiServerUrl string) (io.Reader, bool, *http.Response, error) {

	var err error = nil
	isFileEmpty := false
	log.Printf("Downloading artifact '%s' from API server", artifactPath)

	// A HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	fileDownloaded, httpResponse, err := restClient.ResultArchiveStoreAPIApi.
		GetRasRunArtifactByPath(context.Background(), runId, artifactPath).
		Execute()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DOWNLOADING_ARTIFACT_FAILED, artifactPath, err.Error())
		log.Printf("Failed to download artifact. %s", err.Error())
	} else {
		log.Printf("Artifact '%s' http response from API server OK", artifactPath)

		if fileDownloaded == nil {
			log.Printf("Artifact '%s' http response returned nil file content.", artifactPath)
			isFileEmpty = true
		}
	}

	return fileDownloaded, isFileEmpty, httpResponse, err
}
