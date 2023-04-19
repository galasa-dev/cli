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

	// TODO: Should we validate the runname? Can we ?
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

// func renderRuns(format OutputFormat, runs []galasaapi.Run) (string, error) {
// 	var err error = nil
// 	buff := strings.Builder{}
// 	for _, run := range runs {
// 		line := fmt.Sprintf("%s %s %s\n", run.TestStructure.GetRunName(), run.TestStructure.GetStatus(), run.TestStructure.GetResult())
// 		buff.WriteString(line)
// 	}
// 	result := buff.String()
// 	return result, err
// }

func renderRuns(outputFormat OutputFormat, runs []galasaapi.Run) (string, error) {
	var err error = nil
	var formattedOutput string
	//can switch on the output format in the future. Currently this is all for outputFormat = 'summary'
	switch outputFormat {
	case 0:
		//outputFormat = 'summary'
		var output [][]string
		var headers = []string{"RunName", "Status", "Result", "ShortTestName"}
		output = append(output, headers)
		for _, run := range runs {
			var line []string
			line = append(line, run.TestStructure.GetRunName(), run.TestStructure.GetStatus(), run.TestStructure.GetResult(), run.TestStructure.GetTestShortName())
			output = append(output, line)
		}
		formattedOutput, err = formatSummaryConsole(output, runs)

	// case 1:
	// 	//outputFormat = 'raw'
	// case 2:
	// 	//outputFormat = 'detailed'
	default:
		//outputFormat = 'summary'
	}

	return formattedOutput, err

}

func formatSummaryConsole(table [][]string, runs []galasaapi.Run) (string, error) {
	var err error = nil
	buff := strings.Builder{}
	columnLengths := make([]int, len(table[0]))
	for _, row := range table {
		for i, val := range row {
			if len(val) > columnLengths[i] {
				columnLengths[i] = len(val)
			}
		}
	}

	for _, row := range table {
		var line string = "\n"
		for j, val := range row {
			line += fmt.Sprintf("%-*s ", columnLengths[j], val)
		}
		buff.WriteString(line)
	}
	buff.WriteString("\n")
	result := buff.String()

	return result, err
}

func validateOutputFormatFlagValue(outputFormatString string) (OutputFormat, error) {
	var err error
	var outputFormat OutputFormat

	switch outputFormatString {
	case "summary":
		outputFormat = OUTPUT_FORMAT_SUMMARY

	default:
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OUTPUT_FORMAT, outputFormatString)
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
				// Add all the runs into our set of results.
				// Note: The ... syntax means 'all of the array', so they all get appended at once.
				results = append(results, runsOnThisPage...)

				// Have we processed the last page ?
				if pageNumberWanted == runData.GetNumPages() {
					gotAllResults = true
				}
			}
		}
	}

	return results, err
}
