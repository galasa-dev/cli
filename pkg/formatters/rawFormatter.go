/*
 * Copyright contributors to the Galasa project
 */
package formatters

import (
	"strings"

	"github.com/galasa.dev/cli/pkg/galasaapi"
)

// -----------------------------------------------------
// Summary format.
const (
	RAW_FORMATTER_NAME = "raw"
)

type RawFormatter struct {
}

func NewRawFormatter() RunsFormatter {
	return new(RawFormatter)
}

func (*RawFormatter) GetName() string {
	return RAW_FORMATTER_NAME
}

func (*RawFormatter) FormatRuns(runs []galasaapi.Run, apiServerUrl string) (string, error) {
	var result string = ""
	var err error = nil
	buff := strings.Builder{}

	for count, run := range runs {
		if count != 0 {
			buff.WriteString("\n")
		}

		var duration string = ""
		var startTimeString string = ""
		var endTimeString string = ""
		startTimeStringRaw := run.TestStructure.GetStartTime()
		endTimeStringRaw := run.TestStructure.GetEndTime()

		if len(startTimeStringRaw) > 0 {
			startTimeString = formatTime(startTimeStringRaw)
			if len(endTimeStringRaw) > 0 {
				endTimeString = formatTime(endTimeStringRaw)
			}
		}

		duration = calculateDurationMilliseconds(startTimeString, endTimeString)

		runLog := apiServerUrl + "/ras/run/" + run.GetRunId() + "/runlog"

		buff.WriteString(run.TestStructure.GetRunName() + "|" +
			run.TestStructure.GetStatus() + "|" +
			run.TestStructure.GetResult() + "|" +
			formatTime(run.TestStructure.GetQueued()) + "|" +
			startTimeString + "|" +
			endTimeString + "|" +
			duration + "|" +
			run.TestStructure.GetTestName() + "|" +
			run.TestStructure.GetRequestor() + "|" +
			run.TestStructure.GetBundle() + "|" +
			runLog,
		)
	}
	result = buff.String()
	return result, err
}
