/*
 * Copyright contributors to the Galasa project
 */
package formatters

import (
	"fmt"
	"strconv"
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

func NewDetailsFormatter() RunsFormatter {
	return new(DetailsFormatter)

}

func (*DetailsFormatter) GetName() string {
	return DETAILS_FORMATTER_NAME
}

func (*DetailsFormatter) FormatRuns(runs []galasaapi.Run, apiServerUrl string) (string, error) {
	var result string = ""
	var err error = nil

	if len(runs) < 1 {
		return result, err
	}

	buff := strings.Builder{}

	for i, run := range runs {
		var duration string
		var startTimeString string
		var endTimeString string
		startTimeStringRaw := run.TestStructure.GetStartTime()
		endTimeStringRaw := run.TestStructure.GetEndTime()

		if len(startTimeStringRaw) > 0 {
			startTimeString = formatTime(startTimeStringRaw)
			if len(endTimeStringRaw) > 0 {
				endTimeString = formatTime(endTimeStringRaw)
			}
		}

		startTime, err := time.Parse("2006-01-02 15:04:05", startTimeString)
		if err == nil {
			endTime, err := time.Parse("2006-01-02 15:04:05", endTimeString)
			if err == nil {
				duration = strconv.FormatInt(endTime.Sub(startTime).Milliseconds(), 10)
			}
		}

		var table = [][]string{
			{"name", ":  " + run.TestStructure.GetRunName()},
			{"status", ":  " + run.TestStructure.GetStatus()},
			{"result", ":  " + run.TestStructure.GetResult()},
			{"queued-time", ":  " + formatTime(run.TestStructure.GetQueued())},
			{"start-time", ":  " + startTimeString},
			{"end-time", ":  " + endTimeString},
			{"duration(ms)", ":  " + duration},
			{"test-name", ":  " + run.TestStructure.GetTestName()},
			{"requestor", ":  " + run.TestStructure.GetRequestor()},
			{"bundle", ":  " + run.TestStructure.GetBundle()},
			{"run-log", ":  " + apiServerUrl + "/ras/run/" + run.GetRunId() + "/runlog"},
		}

		writeTableToBuff(&buff, table)

		buff.WriteString("\n")
		var methodTable [][]string
		var headers = []string{"method", "type", "status", "result", "start-time", "end-time", "duration(ms)"}
		methodTable = append(methodTable, headers)

		for _, method := range run.TestStructure.Methods {
			var duration string
			startTimeStringRaw := method.GetStartTime()
			startTimeString := formatTime(startTimeStringRaw)

			endTimeStringRaw := method.GetEndTime()
			endTimeString := formatTime(endTimeStringRaw)

			startTime, err := time.Parse("2006-01-02 15:04:05", startTimeString)
			if err == nil {
				endTime, err := time.Parse("2006-01-02 15:04:05", endTimeString)
				if err == nil {
					duration = strconv.FormatInt(endTime.Sub(startTime).Milliseconds(), 10)
				}
			}
			var line []string
			line = append(line,
				method.GetMethodName(),
				method.GetType(),
				method.GetStatus(),
				method.GetResult(),
				startTimeString,
				endTimeString,
				duration,
			)
			methodTable = append(methodTable, line)
		}

		writeTableToBuff(&buff, methodTable)

		if i < len(runs)-1 {
			buff.WriteString("\n---\n\n")
		}

	}

	result = buff.String()

	return result, err
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

func formatTime(rawTime string) string {
	formattedTimeString := rawTime[0:10] + " " + rawTime[11:19]
	return formattedTimeString
}
