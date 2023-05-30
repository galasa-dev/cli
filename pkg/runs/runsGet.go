/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/formatters"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

var (
	validFormatters = createFormatters()
)

// ---------------------------------------------------

// GetRuns - performs all the logic to implement the `galasactl runs get` command,
// but in a unit-testable manner.
func GetRuns(
	runName string,
	age string,
	outputFormatString string,
	timeService utils.TimeService,
	console utils.Console,
	apiServerUrl string,
) error {

	var err error
	var fromAge int
	var toAge int

	if (runName == "") && (age == "") {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_NO_RUNNAME_OR_AGE_SPECIFIED)
	}

	if (err == nil) && (runName != "") {
		// Validate the runName as best we can without contacting the ecosystem.
		err = ValidateRunName(runName)
	}

	if (err == nil) && (age != "") {
		fromAge, toAge, err = getTimesFromAge(age)
	}

	if err == nil {
		var chosenFormatter formatters.RunsFormatter
		chosenFormatter, err = validateOutputFormatFlagValue(outputFormatString, validFormatters)
		if err == nil {
			var runJson []galasaapi.Run
			runJson, err = GetRunsFromRestApi(runName, fromAge, toAge, timeService, apiServerUrl)
			if err == nil {
				// Some formatters need extra fields filled-in so they can be displayed.
				if chosenFormatter.IsNeedingMethodDetails() {
					runJson, err = GetRunDetailsFromRasSearchRuns(runJson, apiServerUrl)
				}

				if err == nil {
					var outputText string
					outputText, err = chosenFormatter.FormatRuns(runJson, apiServerUrl)
					if err == nil {
						err = writeOutput(outputText, console)
					}
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

	rawFormatter := formatters.NewRawFormatter()
	validFormatters[rawFormatter.GetName()] = rawFormatter

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
	fromAgeHours int,
	toAgeHours int,
	timeService utils.TimeService,
	apiServerUrl string,
) ([]galasaapi.Run, error) {

	var err error = nil
	var results []galasaapi.Run = make([]galasaapi.Run, 0)

	var context context.Context = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	now := timeService.Now()
	fromTime := now.Add(-(time.Duration(fromAgeHours) * time.Hour)).UTC() // Add a minus, so subtract

	toTime := now.Add(-(time.Duration(toAgeHours) * time.Hour)).UTC() // Add a minus, so subtract

	var pageNumberWanted int32 = 1
	gotAllResults := false

	for (!gotAllResults) && (err == nil) {

		var runData *galasaapi.RunResults
		var httpResponse *http.Response
		log.Printf("Requesting page '%d' ", pageNumberWanted)
		apicall := restClient.ResultArchiveStoreAPIApi.GetRasSearchRuns(context)
		if fromAgeHours != 0 {
			apicall = apicall.From(fromTime)
		}
		if toAgeHours != 0 {
			apicall = apicall.To(toTime)
		}
		if runName != "" {
			apicall = apicall.Runname(runName)
		}
		apicall = apicall.Page(pageNumberWanted)
		apicall = apicall.Sort("to:desc")
		runData, httpResponse, err = apicall.Execute()

		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
		} else {
			if httpResponse.StatusCode != http.StatusOK {
				httpError := "\nhttp response status code: " + strconv.Itoa(httpResponse.StatusCode)
				errString := err.Error() + httpError
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, errString)
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

func getTimesFromAge(age string) (int, int, error) {
	// Validate the age parameter
	// Validate that the time unit is either 'w', 'd', 'h'
	// Validate that if both FROM and TO are specified, FROM is older than TO

	// Make a map of how many hours for each unit so can compare from and to values consistently
	// Can be extended to support other units
	var timeUnits = make(map[string]int)
	timeUnits["w"] = 168
	timeUnits["d"] = 24
	timeUnits["h"] = 1

	regex := "([0-9]+)([dhw])"
	re := regexp.MustCompile(regex)

	submatches := re.FindAllStringSubmatch(age, -1)

	var err error
	var fromAge int
	var toAge int

	if len(submatches) != 0 {
		from := submatches[0] // Expecting something like 14d which will then break down into further matches: Index 0 is 14d, index 1 is 14, index 2 is d
		fromAge, err = getValueAsInt(from[1])
		if err == nil {
			if fromAge != 0 {
				fromAge = fromAge * timeUnits[from[2]]
			} else {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FROM_AGE_NOT_SPECIFIED, fromAge)
			}

			// The user has also specified a to age
			if len(submatches) > 1 {
				to := submatches[1]
				toAge, err = getValueAsInt(to[1])
				if err == nil {
					if toAge != 0 {
						toAge = toAge * timeUnits[to[2]]

						// From value has to be bigger than to value
						if fromAge <= toAge {
							err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FROM_AGE_SMALLER_THAN_TO_AGE, age)
						}
					}
				}
			}
		}
	} else {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_AGE_PARAMETER, age)
	}

	return fromAge, toAge, err
}

func getValueAsInt(value string) (int, error) {
	var age int
	var err error
	if age, err = strconv.Atoi(value); err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_FROM_OR_TO_PARAMETER, value)
	}
	return age, err
}
