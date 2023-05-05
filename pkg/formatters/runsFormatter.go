/*
 * Copyright contributors to the Galasa project
 */
package formatters

import (
	"github.com/galasa.dev/cli/pkg/galasaapi"
)

// -----------------------------------------------------
// RunsFormatter - implementations can take a collection of run results
// and turn them into a string for display to the user.
type RunsFormatter interface {
	FormatRuns(runs []galasaapi.Run) (string, error)
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
