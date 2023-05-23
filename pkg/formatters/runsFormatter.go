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
