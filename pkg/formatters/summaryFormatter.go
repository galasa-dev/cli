/*
 * Copyright contributors to the Galasa project
 */
package formatters

import (
	"fmt"
	"strings"

	"github.com/galasa.dev/cli/pkg/galasaapi"
)

// -----------------------------------------------------
// Summary format.
const (
	SUMMARY_FORMATTER_NAME = "summary"
)

type SummaryFormatter struct {
}

func NewSummaryFormatter() RunsFormatter {
	return new(SummaryFormatter)
}

func (*SummaryFormatter) GetName() string {
	return SUMMARY_FORMATTER_NAME
}

func (*SummaryFormatter) FormatRuns(runs []galasaapi.Run, apiServerUrl string) (string, error) {
	var result string = ""
	var err error = nil

	if len(runs) < 1 {
		return result, err
	}

	var table [][]string

	var headers = []string{"name", "status", "result", "test-name"}

	table = append(table, headers)
	for _, run := range runs {
		var line []string
		line = append(line, run.TestStructure.GetRunName(), run.TestStructure.GetStatus(), run.TestStructure.GetResult(), run.TestStructure.GetTestName())
		table = append(table, line)
	}

	buff := strings.Builder{}

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
	result = buff.String()
	return result, err
}
