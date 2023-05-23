/*
 * Copyright contributors to the Galasa project
 */
package formatters

import (
	"strconv"
	"time"

	"github.com/galasa.dev/cli/pkg/galasaapi"
)

// -----------------------------------------------------
// RunsFormatter - implementations can take a collection of run results
// and turn them into a string for display to the user.
const (
	DATE_FORMAT = "2006-01-02 15:04:05"
)

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

func formatTime(rawTime string) string {
	formattedTimeString := rawTime[0:10] + " " + rawTime[11:19]
	return formattedTimeString
}

func calculateDurationMilliseconds(startTimeString string, endTimeString string) string {
	var duration string = ""

	startTime, err := time.Parse(DATE_FORMAT, startTimeString)
	if err == nil {
		endTime, err := time.Parse(DATE_FORMAT, endTimeString)
		if err == nil {
			duration = strconv.FormatInt(endTime.Sub(startTime).Milliseconds(), 10)
		}
	}
	return duration
}
