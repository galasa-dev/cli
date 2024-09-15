/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
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

	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/runsformatter"
	"github.com/galasa-dev/cli/pkg/spi"
)

var (
	validFormatters = CreateFormatters()

	// When parsing the '--age' parameter value....
	// (^[\\D]*) - matches any leading garbage, which is non-digits. Should be empty.
	// ([0-9]+) - matches any digit sequence. Should be an integer.
	// (.*) - matches any time unit. Should be a valid time unit from our map above.
	agePartRegex *regexp.Regexp = regexp.MustCompile(`(^[\D]*)([0-9]+)(.*)`)

	// Comma separated string of valid status values
	// specifically, status values of an 'active' test
	activeStatusNames = "started,ending,generating,building,provstart,running,rundone,up"
)

// ---------------------------------------------------

// GetRuns - performs all the logic to implement the `galasactl runs get` command,
// but in a unit-testable manner.
func GetRuns(
	runName string,
	age string,
	requestorParameter string,
	resultParameter string,
	shouldGetActive bool,
	outputFormatString string,
	timeService spi.TimeService,
	console spi.Console,
	apiServerUrl string,
	apiClient *galasaapi.APIClient,
) error {
	var err error
	var fromAge int
	var toAge int

	log.Printf("GetRuns entered.")

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

	if (err == nil) && (resultParameter != "") {
		if shouldGetActive {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_ACTIVE_AND_RESULT_ARE_MUTUALLY_EXCLUSIVE)
		}
		if err == nil {
			resultParameter, err = ValidateResultParameter(resultParameter, apiClient)
		}
	}

	if err == nil {
		var chosenFormatter runsformatter.RunsFormatter
		chosenFormatter, err = validateOutputFormatFlagValue(outputFormatString, validFormatters)
		if err == nil {
			var runJson []galasaapi.Run
			runJson, err = GetRunsFromRestApi(runName, requestorParameter, resultParameter, fromAge, toAge, shouldGetActive, timeService, apiClient)
			if err == nil {
				// Some formatters need extra fields filled-in so they can be displayed.
				if chosenFormatter.IsNeedingMethodDetails() {
					runJson, err = GetRunDetailsFromRasSearchRuns(runJson, apiClient)
				}

				if err == nil {
					var outputText string

					//convert galsaapi.Runs tests into formattable data
					formattableTest := FormattableTestFromGalasaApi(runJson, apiServerUrl)
					outputText, err = chosenFormatter.FormatRuns(formattableTest)

					if err == nil {
						err = writeOutput(outputText, console)
					}
				}
			}
		}
	}
	log.Printf("GetRuns exiting. err is %v", err)
	return err
}

func CreateFormatters() map[string]runsformatter.RunsFormatter {
	validFormatters := make(map[string]runsformatter.RunsFormatter, 0)
	summaryFormatter := runsformatter.NewSummaryFormatter()
	validFormatters[summaryFormatter.GetName()] = summaryFormatter

	detailedFormatter := runsformatter.NewDetailsFormatter()
	validFormatters[detailedFormatter.GetName()] = detailedFormatter

	rawFormatter := runsformatter.NewRawFormatter()
	validFormatters[rawFormatter.GetName()] = rawFormatter

	return validFormatters
}

func writeOutput(outputText string, console spi.Console) error {
	err := console.WriteString(outputText)
	return err
}

// GetFormatterNamesString builds a string of comma separated, quoted formatter names
func GetFormatterNamesString(validFormatters map[string]runsformatter.RunsFormatter) string {
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
func validateOutputFormatFlagValue(outputFormatString string, validFormatters map[string]runsformatter.RunsFormatter) (runsformatter.RunsFormatter, error) {
	var err error

	chosenFormatter, isPresent := validFormatters[outputFormatString]

	if !isPresent {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OUTPUT_FORMAT, outputFormatString, GetFormatterNamesString(validFormatters))
	}

	return chosenFormatter, err
}

func GetRunDetailsFromRasSearchRuns(runs []galasaapi.Run, apiClient *galasaapi.APIClient) ([]galasaapi.Run, error) {
	var err error
	var runsDetails []galasaapi.Run = make([]galasaapi.Run, 0)
	var context context.Context = nil
	var details *galasaapi.Run
	var httpResponse *http.Response

	var restApiVersion string
	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {

		for _, run := range runs {
			runid := run.GetRunId()
			details, httpResponse, err = apiClient.ResultArchiveStoreAPIApi.GetRasRunById(context, runid).ClientApiVersion(restApiVersion).Execute()
			if err != nil {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
			} else {
				if httpResponse.StatusCode != http.StatusOK {

					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_NON_OK_STATUS, strconv.Itoa(httpResponse.StatusCode))
				} else {
					runsDetails = append(runsDetails, *details)
				}
			}
		}
	}

	return runsDetails, err
}

// Retrieves test runs from the ecosystem API that match a given runName.
// Multiple test runs can be returned as the runName is not unique.
func GetRunsFromRestApi(
	runName string,
	requestorParameter string,
	resultParameter string,
	fromAgeMins int,
	toAgeMins int,
	shouldGetActive bool,
	timeService spi.TimeService,
	apiClient *galasaapi.APIClient,
) ([]galasaapi.Run, error) {

	var err error
	var results []galasaapi.Run = make([]galasaapi.Run, 0)

	var context context.Context = nil

	now := timeService.Now()
	fromTime := now.Add(-(time.Duration(fromAgeMins) * time.Minute)).UTC() // Add a minus, so subtract

	toTime := now.Add(-(time.Duration(toAgeMins) * time.Minute)).UTC() // Add a minus, so subtract

	var pageNumberWanted int32 = 1
	gotAllResults := false
	var restApiVersion string
	var pageCursor string

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {

		for (!gotAllResults) && (err == nil) {

			var runData *galasaapi.RunResults
			var httpResponse *http.Response
			log.Printf("Requesting page '%d' ", pageNumberWanted)
			apicall := apiClient.ResultArchiveStoreAPIApi.GetRasSearchRuns(context).ClientApiVersion(restApiVersion).IncludeCursor("true")
			if fromAgeMins != 0 {
				apicall = apicall.From(fromTime)
			}
			if toAgeMins != 0 {
				apicall = apicall.To(toTime)
			}
			if runName != "" {
				apicall = apicall.Runname(runName)
			}
			if requestorParameter != "" {
				apicall = apicall.Requestor(requestorParameter)
			}
			if resultParameter != "" {
				apicall = apicall.Result(resultParameter)
			}
			if shouldGetActive {
				apicall = apicall.Status(activeStatusNames)
			}
			if pageCursor != "" {
				apicall = apicall.Cursor(pageCursor)
			}
			apicall = apicall.Sort("from:desc")
			runData, httpResponse, err = apicall.Execute()

			if err != nil {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, err.Error())
			} else {
				if httpResponse.StatusCode != http.StatusOK {
					httpError := "\nhttp response status code: " + strconv.Itoa(httpResponse.StatusCode)
					errString := httpError
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_RUNS_FAILED, errString)
				} else {

					// Copy the results from this page into our bigger list of results.
					runsOnThisPage := runData.GetRuns()
					// Add all the runs into our set of results.
					// Note: The ... syntax means 'all of the array', so they all get appended at once.
					results = append(results, runsOnThisPage...)

					// Have we processed the last page ?
					if !runData.HasNextCursor() || len(runsOnThisPage) < int(runData.GetPageSize()) {
						gotAllResults = true
					} else {
						pageCursor = runData.GetNextCursor()
						pageNumberWanted++
					}
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

	var err error
	var fromAge int
	var toAge int = 0

	ageParts := strings.Split(age, ":")

	if len(ageParts) > 2 {
		// Too many colons.
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_AGE_PARAMETER, age, GetTimeUnitsForErrorMessage())
	} else {
		// No colons !... only 'from' time specified.
		fromPart := ageParts[0]
		if !agePartRegex.MatchString(fromPart) {
			// Invalid from part.
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_FROM_AGE_SPECIFIED, age, GetTimeUnitsForErrorMessage())
		} else {
			fromAge, err = getMinutesFromAgePart(fromPart, age)

			if fromAge == 0 {
				// 'from' can't be 0 hours.
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_AGE_PARAMETER, age, GetTimeUnitsForErrorMessage())
			}
		}

		if (err == nil) && (len(ageParts) > 1) {
			// One colon, indicates there is a 'to' part.
			toPart := ageParts[1]
			toAge, err = getMinutesFromAgePart(toPart, age)
			if err == nil {
				// From value must be bigger than to value
				if toAge > 0 && fromAge <= toAge {
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FROM_AGE_SMALLER_THAN_TO_AGE, age)
				}
			}
		}
	}

	return fromAge, toAge, err
}

// getHoursFromAgePart - Input value is '15d' or '14h' for example.
func getMinutesFromAgePart(agePart string, errorMessageValue string) (int, error) {
	var err error
	var minutes int = 0
	var duration int

	// Separate the integer part from time unit part.
	// Expecting something like 14d which will then break down into further matches: Index 0 is 14d, index 1 is 14, index 2 is d
	durationPart := agePartRegex.FindStringSubmatch(agePart)
	leadingGarbage := durationPart[1]
	durationNumber := durationPart[2]
	durationUnitStr := durationPart[3]

	if leadingGarbage != "" {
		// Some leading garbage prior to the 'FROM' field. It must be empty.
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_AGE_PARAMETER, errorMessageValue, GetTimeUnitsForErrorMessage())
	} else {

		if len(durationPart) == 0 {
			// Invalid from. It must be some time in the past.
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_FROM_AGE_SPECIFIED, errorMessageValue, GetTimeUnitsForErrorMessage())
		} else {
			// we can extract the integer part now

			duration, err = getValueAsInt(durationNumber)
			if err == nil {
				if duration < 0 {
					// Number part of the duration can't be negative.
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_NEGATIVE_AGE_SPECIFIED, errorMessageValue, GetTimeUnitsForErrorMessage())
				} else {

					timeUnit, isRecognisedTimeUnit := GetTimeUnitFromShortName(durationUnitStr)
					if !isRecognisedTimeUnit {
						// Bad time unit.
						err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_BAD_TIME_UNIT_AGE_SPECIFIED, errorMessageValue, GetTimeUnitsForErrorMessage())
					} else {
						minutes = duration * timeUnit.GetMinuteMultiplier()
					}
				}
			}
		}
	}
	return minutes, err
}

func getValueAsInt(value string) (int, error) {
	var age int
	var err error
	if age, err = strconv.Atoi(value); err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_FROM_OR_TO_PARAMETER, value)
	}
	return age, err
}
