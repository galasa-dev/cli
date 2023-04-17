/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

// ---------------------------------------------------
// As close as we can get to an enum in Go...
type OutputFormat int64

const (
	OUTPUT_FORMAT_SUMMARY = OutputFormat(0)
)

const (
	HTTP_STATUS_CODE_OK = 200
)

// GetRuns - performs all the logic to implement the `galasactl runs get` command,
// but in a unit-testable manner.
func GetRuns(
	runName string,
	outputFormatString string,
	timeService utils.TimeService,
	console utils.Console,
	apiServerUrl string,
) error {

	var err error
	var outputFormat OutputFormat

	outputFormat, err = validateOutputFormatFlagValue(outputFormatString)
	if err == nil {
		var runJson []galasaapi.Run
		runJson, err = GetRunsFromRestApi(runName, outputFormat, timeService, apiServerUrl)
		if err == nil {
			var outputText string
			outputText, err = renderRuns(outputFormat, runJson)
			if err == nil {
				err = writeOutput(outputText, console)
			}
		}
	}

	return err
}

func writeOutput(outputText string, console utils.Console) error {
	err := console.WriteString(outputText)
	return err
}

func renderRuns(format OutputFormat, runs []galasaapi.Run) (string, error) {
	var err error = nil
	buff := strings.Builder{}
	for _, run := range runs {
		line := fmt.Sprintf("%s %s %s\n", run.TestStructure.GetRunName(), run.TestStructure.GetStatus(), run.TestStructure.GetResult())
		buff.WriteString(line)
	}
	result := buff.String()
	return result, err
}

func validateOutputFormatFlagValue(outputFormatString string) (OutputFormat, error) {
	var err error
	var outputFormat OutputFormat

	switch outputFormatString {
	case "summary":

	default:
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OUTPUT_FORMAT, outputFormat)
	}

	return outputFormat, err
}

func GetRunsFromRestApi(
	runName string,
	outputFormat OutputFormat,
	timeService utils.TimeService,
	apiServerUrl string,
) ([]galasaapi.Run, error) {

	var err error = nil
	var results []galasaapi.Run = make([]galasaapi.Run, 0)

	var context context.Context = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	toTime := timeService.Now()
	var pageNumberWanted int32 = 1
	gotAllResults := false

	for (!gotAllResults) && (err == nil) {

		var runData *galasaapi.RunResults
		var httpResponse *http.Response
		log.Printf("Requesting page '%d' ", pageNumberWanted)
		runData, httpResponse, err = restClient.ResultArchiveStoreAPIApi.
			GetRasSearchRuns(context).
			To(toTime).
			Runname(runName).
			Page(pageNumberWanted).
			Sort("to:desc").
			Execute()

		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
		} else {
			if httpResponse.StatusCode != HTTP_STATUS_CODE_OK {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
			} else {

				// Copy the results from this page into our bigger list of results.
				runsOnThisPage := runData.GetRuns()
				for _, run := range runsOnThisPage {
					results = append(results, run)
				}

				// Have we processed the last page ?
				if pageNumberWanted == runData.GetNumPages() {
					gotAllResults = true
				}
			}
		}
	}

	return results, err
}

func SummaryOutput(console utils.Console, text string) {
	console.WriteString(text)
}
