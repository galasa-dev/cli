/*
 * Copyright contributors to the Galasa project
 */
package formatters

import (
	"fmt"
	"log"
	"strings"

	"github.com/galasa.dev/cli/pkg/galasaapi"
)

// -----------------------------------------------------
// Detailed format.
type DetailedFormatter struct {
}

func NewDetailedFormatter() RunsFormatter {
	return new(DetailedFormatter)
}

func (*DetailedFormatter) FormatRuns(runs []galasaapi.Run) (string, error) {
	var result string = ""
	var err error = nil

	if len(runs) < 1 {
		return result, err
	}

	buff := strings.Builder{}

	for _, run := range runs {

		var table = [][]string{
			{"name", ":  " + run.TestStructure.GetRunName()},
			{"status", ":  " + run.TestStructure.GetStatus()},
			{"result", ":  " + run.TestStructure.GetResult()},
			{"queued-time", ":  " + run.TestStructure.GetQueued()},
			{"start-time", ":  " + run.TestStructure.GetStartTime()},
			{"end-time", ":  " + run.TestStructure.GetEndTime()},
			{"test-name", ":  " + run.TestStructure.GetTestName()},
			{"requestor", ":  " + run.TestStructure.GetRequestor()},
			{"bundle", ":  " + run.TestStructure.GetBundle()},
		}

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
		buff.WriteString("\n")
		var methodTable [][]string
		var headers = []string{"method", "type", "status", "result", "start-time", "end-time", "duration(ms)"}
		methodTable = append(methodTable, headers)

		for _, method := range run.TestStructure.Methods {
			var line []string
			//duration := run.TestStructure.Methods[i].GetEndTime() - run.TestStructure.Methods[i].GetStartTime()
			line = append(line,
				method.GetMethodName(),
				method.GetType(),
				method.GetStatus(),
				method.GetResult(),
				method.GetStartTime(),
				method.GetEndTime(),
				//duration,
			)
			methodTable = append(methodTable, line)
		}

		columnLengths = calculateMaxLengthOfEachColumn(methodTable)

		for _, row := range methodTable {
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

	result = buff.String()
	log.Print(result)
	return result, err
}
