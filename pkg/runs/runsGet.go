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
	validFormatters = CreateFormatters()

	// Make a map of how many hours for each unit so can compare from and to values consistently
	// Can be extended to support other units

	timeUnits = CreateTimeUnits()

	// When parsing the '--age' parameter value....
	// (^[\\D]*) - matches any leading garbage, which is non-digits. Should be empty.
	// ([0-9]+) - matches any digit sequence. Should be an integer.
	// (.*) - matches any time unit. Should be a valid time unit from our map above.
	agePartRegex *regexp.Regexp = regexp.MustCompile(`(^[\D]*)([0-9]+)(.*)`)
)

// ---------------------------------------------------

// GetRuns - performs all the logic to implement the `galasactl runs get` command,
// but in a unit-testable manner.
func GetRuns(
	runName string,
	age string,
	requestorParameter string,
	resultParameter string,
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

	if (err == nil) && (resultParameter != "") {
		resultParameter, err = ValidateResultParameter(resultParameter, apiServerUrl)
	}

	if err == nil {
		var chosenFormatter formatters.RunsFormatter
		chosenFormatter, err = validateOutputFormatFlagValue(outputFormatString, validFormatters)
		if err == nil {
			var runJson []galasaapi.Run
			runJson, err = GetRunsFromRestApi(runName, requestorParameter, resultParameter, fromAge, toAge, timeService, apiServerUrl)
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

func CreateTimeUnits() map[string]TimeUnit {
	timeUnits := make(map[string]TimeUnit, 0)

	unitWeeks := newTimeUnit(TIME_UNIT_WEEKS_LONG, 10080)
	timeUnits[TIME_UNIT_WEEKS_SHORT] = *unitWeeks

	unitDays := newTimeUnit(TIME_UNIT_DAYS_LONG, 1440)
	timeUnits[TIME_UNIT_DAYS_SHORT] = *unitDays

	unitHours := newTimeUnit(TIME_UNIT_HOURS_LONG, 60)
	timeUnits[TIME_UNIT_HOURS_SHORT] = *unitHours

	unitMinutes := newTimeUnit(TIME_UNIT_MINUTES_LONG, 1)
	timeUnits[TIME_UNIT_MINUTES_SHORT] = *unitMinutes

	return timeUnits
}

func CreateFormatters() map[string]formatters.RunsFormatter {
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

// GetFormatterNamesString builds a string of comma separated, quoted formatter names
func GetFormatterNamesString(validFormatters map[string]formatters.RunsFormatter) string {
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
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OUTPUT_FORMAT, outputFormatString, GetFormatterNamesString(validFormatters))
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
	requestorParameter string,
	resultParameter string,
	fromAgeMins int,
	toAgeMins int,
	timeService utils.TimeService,
	apiServerUrl string,
) ([]galasaapi.Run, error) {

	var err error = nil
	var results []galasaapi.Run = make([]galasaapi.Run, 0)

	var context context.Context = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	now := timeService.Now()
	fromTime := now.Add(-(time.Duration(fromAgeMins) * time.Minute)).UTC() // Add a minus, so subtract

	toTime := now.Add(-(time.Duration(toAgeMins) * time.Minute)).UTC() // Add a minus, so subtract

	var pageNumberWanted int32 = 1
	gotAllResults := false

	for (!gotAllResults) && (err == nil) {

		var runData *galasaapi.RunResults
		var httpResponse *http.Response
		log.Printf("Requesting page '%d' ", pageNumberWanted)
		apicall := restClient.ResultArchiveStoreAPIApi.GetRasSearchRuns(context)
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
				} else {
					pageNumberWanted++
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
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_AGE_PARAMETER, age, GetTimeUnitsForErrorMessage(timeUnits))
	} else {
		// No colons !... only 'from' time specified.
		fromPart := ageParts[0]
		if !agePartRegex.MatchString(fromPart) {
			// Invalid from part.
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_FROM_AGE_SPECIFIED, age, GetTimeUnitsForErrorMessage(timeUnits))
		} else {
			fromAge, err = getMinutesFromAgePart(fromPart, age)

			if fromAge == 0 {
				// 'from' can't be 0 hours.
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_AGE_PARAMETER, age, GetTimeUnitsForErrorMessage(timeUnits))
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
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_AGE_PARAMETER, errorMessageValue, GetTimeUnitsForErrorMessage(timeUnits))
	} else {

		if len(durationPart) == 0 {
			// Invalid from. It must be some time in the past.
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_FROM_AGE_SPECIFIED, errorMessageValue, GetTimeUnitsForErrorMessage(timeUnits))
		} else {
			// we can extract the integer part now

			duration, err = getValueAsInt(durationNumber)
			if err == nil {
				if duration < 0 {
					// Number part of the duration can't be negative.
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_NEGATIVE_AGE_SPECIFIED, errorMessageValue, GetTimeUnitsForErrorMessage(timeUnits))
				} else {

					timeUnit, isRecognisedTimeUnit := timeUnits[durationUnitStr]
					if !isRecognisedTimeUnit {
						// Bad time unit.
						err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_BAD_TIME_UNIT_AGE_SPECIFIED, errorMessageValue, GetTimeUnitsForErrorMessage(timeUnits))
					} else {
						minutes = duration * timeUnit.getMinuteMultiplier()
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

func GetTimeUnitsForErrorMessage(timeUnits map[string]TimeUnit) string {
	outputString := strings.Builder{}
	count := 0
	for initial, unit := range timeUnits {

		if count != 0 {
			outputString.WriteString(", ")
		}
		outputString.WriteString("'" + initial + "' (" + unit.getName() + ")")
		count++
	}

	return outputString.String()
}
