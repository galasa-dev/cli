/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package formatters

import (
	"strings"

	"github.com/galasa.dev/cli/pkg/galasaapi"
)

// -----------------------------------------------------
// Detailed format.
const (
	DETAILS_FORMATTER_NAME = "details"
)

type DetailsFormatter struct {
}

// -----------------------------------------------------
// Constructors
func NewDetailsFormatter() RunsFormatter {
	return new(DetailsFormatter)

}

// -----------------------------------------------------
// Functions in the RunsFormatter interface
func (*DetailsFormatter) GetName() string {
	return DETAILS_FORMATTER_NAME
}

func (*DetailsFormatter) IsNeedingMethodDetails() bool {
	return true
}

func (*DetailsFormatter) FormatRuns(runs []galasaapi.Run, apiServerUrl string) (string, error) {
	var result string = ""
	var err error = nil

	totalResults := len(runs)
	resultCountsMap := initialiseResultMap()

	buff := strings.Builder{}

	if len(runs) > 0 {

		for i, run := range runs {
			accumulateResults(resultCountsMap, run)
			coreDetailsTable := tabulateCoreRunDetails(run, apiServerUrl)
			coreDetailsColumnLengths := calculateMaxLengthOfEachColumn(coreDetailsTable)
			writeFormattedTableToStringBuilder(coreDetailsTable, &buff, coreDetailsColumnLengths)

			buff.WriteString("\n")

			methodTable := initialiseMethodTable()
			methodTable = tabulateRunMethodsToTable(run.TestStructure.GetMethods(), methodTable)
			methodColumnLengths := calculateMaxLengthOfEachColumn(methodTable)
			writeFormattedTableToStringBuilder(methodTable, &buff, methodColumnLengths)

			if i < len(runs)-1 {
				buff.WriteString("\n---\n\n")
			}

		}

		buff.WriteString("\n")
	}
	totalReportString := generateResultTotalsReport(totalResults, resultCountsMap)
	buff.WriteString(totalReportString + "\n")

	result = buff.String()

	return result, err
}

// -----------------------------------------------------
// Internal functions
func tabulateCoreRunDetails(run galasaapi.Run, apiServerUrl string) [][]string {
	startTimeStringRaw := run.TestStructure.GetStartTime()
	endTimeStringRaw := run.TestStructure.GetEndTime()

	startTimeStringReadable := getReadableTime(startTimeStringRaw)
	endTimeStringReadable := getReadableTime(endTimeStringRaw)

	duration := getDuration(startTimeStringRaw, endTimeStringRaw)

	var table = [][]string{
		{HEADER_RUNNAME, ": " + run.TestStructure.GetRunName()},
		{HEADER_STATUS, ": " + run.TestStructure.GetStatus()},
		{HEADER_RESULT, ": " + run.TestStructure.GetResult()},
		{HEADER_SUBMITTED_TIME, ": " + formatTimeReadable(run.TestStructure.GetQueued())},
		{HEADER_START_TIME, ": " + startTimeStringReadable},
		{HEADER_END_TIME, ": " + endTimeStringReadable},
		{HEADER_DURATION, ": " + duration},
		{HEADER_TEST_NAME, ": " + run.TestStructure.GetTestName()},
		{HEADER_REQUESTOR, ": " + run.TestStructure.GetRequestor()},
		{HEADER_BUNDLE, ": " + run.TestStructure.GetBundle()},
		{HEADER_RUN_LOG, ": " + apiServerUrl + RAS_RUNS_URL + run.GetRunId() + "/runlog"},
	}
	return table
}

func initialiseMethodTable() [][]string {
	var methodTable [][]string
	var headers = []string{HEADER_METHOD_NAME, HEADER_METHOD_TYPE, HEADER_STATUS, HEADER_RESULT, HEADER_START_TIME, HEADER_END_TIME, HEADER_DURATION}
	methodTable = append(methodTable, headers)

	return methodTable
}

func tabulateRunMethodsToTable(methods []galasaapi.TestMethod, methodTable [][]string) [][]string {
	for _, method := range methods {
		startTimeStringRaw := method.GetStartTime()
		endTimeStringRaw := method.GetEndTime()

		startTimeStringReadable := getReadableTime(startTimeStringRaw)
		endTimeStringReadable := getReadableTime(endTimeStringRaw)

		duration := getDuration(startTimeStringRaw, endTimeStringRaw)

		var line []string
		line = append(line,
			method.GetMethodName(),
			method.GetType(),
			method.GetStatus(),
			method.GetResult(),
			startTimeStringReadable,
			endTimeStringReadable,
			duration,
		)
		methodTable = append(methodTable, line)
	}

	return methodTable
}
