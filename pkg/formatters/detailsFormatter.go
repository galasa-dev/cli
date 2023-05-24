/*
 * Copyright contributors to the Galasa project
 */
package formatters

import (
	"fmt"
	"strings"
	"time"

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

	if len(runs) < 1 {
		return result, err
	}

	buff := strings.Builder{}

	for i, run := range runs {

		coreDetailsTable := tabulateCoreRunDetails(run, apiServerUrl)
		writeTableToBuff(&buff, coreDetailsTable)

		buff.WriteString("\n")

		methodTable := initialiseMethodTable()
		methodTable = tabulateRunMethodsToTable(run.TestStructure.GetMethods(), methodTable)
		writeTableToBuff(&buff, methodTable)

		if i < len(runs)-1 {
			buff.WriteString("\n---\n\n")
		}

	}

	result = buff.String()

	return result, err
}

// -----------------------------------------------------
// Internal functions
func tabulateCoreRunDetails(run galasaapi.Run, apiServerUrl string) [][]string {
	var duration string = ""
	var startTimeStringReadable string = ""
	var endTimeStringReadable string = ""

	var startTimeForDuration time.Time
	var endTimeForDuration time.Time

	startTimeStringRaw := run.TestStructure.GetStartTime()
	endTimeStringRaw := run.TestStructure.GetEndTime()

	if len(startTimeStringRaw) > 0 {
		startTimeStringReadable = formatTimeReadable(startTimeStringRaw)
		startTimeForDuration = formatTimeForDurationCalculation(startTimeStringRaw)

		if len(endTimeStringRaw) > 0 {
			endTimeStringReadable = formatTimeReadable(endTimeStringRaw)
			endTimeForDuration = formatTimeForDurationCalculation(endTimeStringRaw)
			duration = calculateDurationMilliseconds(startTimeForDuration, endTimeForDuration)
		}

	}

	var table = [][]string{
		{"name", ":  " + run.TestStructure.GetRunName()},
		{"status", ":  " + run.TestStructure.GetStatus()},
		{"result", ":  " + run.TestStructure.GetResult()},
		{"queued-time", ":  " + formatTimeReadable(run.TestStructure.GetQueued())},
		{"start-time", ":  " + startTimeStringReadable},
		{"end-time", ":  " + endTimeStringReadable},
		{"duration(ms)", ":  " + duration},
		{"test-name", ":  " + run.TestStructure.GetTestName()},
		{"requestor", ":  " + run.TestStructure.GetRequestor()},
		{"bundle", ":  " + run.TestStructure.GetBundle()},
		{"run-log", ":  " + apiServerUrl + "/ras/runs/" + run.GetRunId() + "/runlog"},
	}
	return table
}

func initialiseMethodTable() [][]string {
	var methodTable [][]string
	var headers = []string{"method", "type", "status", "result", "start-time", "end-time", "duration(ms)"}
	methodTable = append(methodTable, headers)

	return methodTable
}

func tabulateRunMethodsToTable(methods []galasaapi.TestMethod, methodTable [][]string) [][]string {
	for _, method := range methods {
		var duration string = ""
		var startTimeStringReadable string = ""
		var endTimeStringReadable string = ""

		var startTimeForDuration time.Time
		var endTimeForDuration time.Time

		startTimeStringRaw := method.GetStartTime()
		endTimeStringRaw := method.GetEndTime()

		if len(startTimeStringRaw) > 0 {
			startTimeStringReadable = formatTimeReadable(startTimeStringRaw)
			startTimeForDuration = formatTimeForDurationCalculation(startTimeStringRaw)
			if len(endTimeStringRaw) > 0 {
				endTimeStringReadable = formatTimeReadable(endTimeStringRaw)
				endTimeForDuration = formatTimeForDurationCalculation(endTimeStringRaw)
				duration = calculateDurationMilliseconds(startTimeForDuration, endTimeForDuration)
			}
		}

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

func writeTableToBuff(buff *strings.Builder, table [][]string) {
	columnLengths := calculateMaxLengthOfEachColumn(table)

	for _, row := range table {
		for column, val := range row {

			// For every column except the last one, add spacing.
			if column < len(row)-1 {
				// %-*s : variable space-padding length, padding is on the right.
				buff.WriteString(fmt.Sprintf("%-*s", columnLengths[column], val))
				buff.WriteString(" ")
			} else {
				buff.WriteString(val)
			}
		}
		buff.WriteString("\n")
	}
}
