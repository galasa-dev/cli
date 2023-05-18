/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"context"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/formatters"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

// ---------------------------------------------------

// GetRuns - performs all the logic to implement the `galasactl runs get` command,
// but in a unit-testable manner.
func GetRuns(
	runName string,
	outputFormatString string,
	timeService utils.TimeService,
	console utils.Console,
	apiServerUrl string,
) error {

	// TODO: Should we validate the runname? Can we ?
	validFormatters := createFormatters()
	chosenFormatter, err := validateOutputFormatFlagValue(outputFormatString, validFormatters)
	if err == nil {
		var runJsonArray []galasaapi.Run
		runJsonArray, err = GetRunsFromRestApi(runName, timeService, apiServerUrl)
		if err == nil {
			var outputText string
			if chosenFormatter.GetName() != "summary" {
				runJsonArray, err = GetRunDetailsFromRasSearchRuns(runJsonArray, apiServerUrl)
			}
			if err == nil {
				outputText, err = chosenFormatter.FormatRuns(runJsonArray, apiServerUrl)
				if err == nil {
					err = writeOutput(outputText, console)
				}
			}

		}
	}

	return err
}

func createFormatters() map[string]formatters.RunsFormatter {
	validFormatters := make(map[string]formatters.RunsFormatter, 0)
	summaryFormatter := formatters.NewSummaryFormatter()
	validFormatters[summaryFormatter.GetName()] = summaryFormatter

	detailedFormatter := formatters.NewDetailsFormatter()
	validFormatters[detailedFormatter.GetName()] = detailedFormatter

	return validFormatters
}

func writeOutput(outputText string, console utils.Console) error {
	err := console.WriteString(outputText)
	return err
}

// getFormatterNamesString builds a string of comma separated, quoted formatter names
func getFormatterNamesString(validFormatters map[string]formatters.RunsFormatter) string {
	// extract names into a sorted slice
	names := make([]string, 0, len(validFormatters))
	for name := range validFormatters {
		names = append(names, name)
	}
	sort.Strings(names)

	// render list of sorted names into string
	formatterNames := strings.Builder{}

	for count, formatterName := range names {

		if count != 0 {
			formatterNames.WriteString(", ")
		}
		formatterNames.WriteString("'" + formatterName + "'")
	}

	return formatterNames.String()
}

// Ensures the user has provided a valid output format as part of the "runs get" command.
func validateOutputFormatFlagValue(outputFormatString string, validFormatters map[string]formatters.RunsFormatter) (formatters.RunsFormatter, error) {
	var err error

	chosenFormatter, isPresent := validFormatters[outputFormatString]

	if !isPresent {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OUTPUT_FORMAT, outputFormatString, getFormatterNamesString(validFormatters))
	}

	return chosenFormatter, err
}

func GetRunDetailsFromRasSearchRuns(runs []galasaapi.Run, apiServerUrl string) ([]galasaapi.Run, error) {
	var err error = nil
	var runsDetails []galasaapi.Run = make([]galasaapi.Run, 0)
	restClient := api.InitialiseAPI(apiServerUrl)
	var context context.Context = nil
	var details *galasaapi.Run
	var httpResponse *http.Response

	for _, run := range runs {
		runid := run.GetRunId()
		details, httpResponse, err = restClient.ResultArchiveStoreAPIApi.GetRasRunById(context, runid).Execute()
		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
		} else {
			if httpResponse.StatusCode != http.StatusOK {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
			} else {
				runsDetails = append(runsDetails, *details)
			}
		}
	}

	return runsDetails, err
}

// Retrieves test runs from the ecosystem API that match a given runName.
// Multiple test runs can be returned as the runName is not unique.
func GetRunsFromRestApi(
	runName string,
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
			if httpResponse.StatusCode != http.StatusOK {
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
