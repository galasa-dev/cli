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
// RunsFormatter - implementations can take a collection of run results
// and turn them into a string for display to the user.
type RunsFormatter interface {
	FormatRuns(runs []galasaapi.Run) (string, error)
}

// -----------------------------------------------------
// Summary format.
type SummaryFormatter struct {
}

func NewSummaryFormatter() RunsFormatter {
	return new(SummaryFormatter)
}

func (*SummaryFormatter) FormatRuns(runs []galasaapi.Run) (string, error) {
	var err error = nil
	var table [][]string

	var headers = []string{"RunName", "Status", "Result", "ShortTestName"}

	table = append(table, headers)
	for _, run := range runs {
		var line []string
		line = append(line, run.TestStructure.GetRunName(), run.TestStructure.GetStatus(), run.TestStructure.GetResult(), run.TestStructure.GetTestShortName())
		table = append(table, line)
	}

	buff := strings.Builder{}

	columnLengths := calculateMaxLengthOfEachColumn(table)

	buff.WriteString("\n")
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
	result := buff.String()

	return result, err
}

func calculateMaxLengthOfEachColumn(table [][]string) []int {
	columnLengths := make([]int, len(table[0]))
	for _, row := range table {
		for i, val := range row {
			if len(val) > columnLengths[i] {
				columnLengths[i] = len(val)
			}
		}
	}
	return columnLengths
}
