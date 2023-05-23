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
	FormatRuns(runs []galasaapi.Run, apiServerUrl string) (string, error)
	GetName() string

	// IsNeedingDetails - Does this formatter require all of the detailed fields to be filled-in,
	// so they can be displayed ? True if so, false otherwise.
	// The caller may need to make sure such things are gathered before calling, and some
	// formatters may not need all the detail.
	IsNeedingDetails() bool
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
