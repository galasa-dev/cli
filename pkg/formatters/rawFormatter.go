/*
 * Copyright contributors to the Galasa project
 */
package formatters

import (
	"strings"
	"time"

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

func (*RawFormatter) IsNeedingDetails() bool {
	return false
}

func (*RawFormatter) FormatRuns(runs []galasaapi.Run, apiServerUrl string) (string, error) {
	var result string = ""
	var err error = nil
	buff := strings.Builder{}

	for _, run := range runs {

		var duration string = ""

		var startTimeStringForDuration time.Time
		var endTimeStringForDuration time.Time

		startTimeStringRaw := run.TestStructure.GetStartTime()
		endTimeStringRaw := run.TestStructure.GetEndTime()

		if len(startTimeStringRaw) > 0 {
			startTimeStringForDuration = formatTimeForDurationCalculation(startTimeStringRaw)
			if len(endTimeStringRaw) > 0 {
				endTimeStringForDuration = formatTimeForDurationCalculation(endTimeStringRaw)
				duration = calculateDurationMilliseconds(startTimeStringForDuration, endTimeStringForDuration)
			}
		}

		runLog := apiServerUrl + "/ras/run/" + run.GetRunId() + "/runlog"

		buff.WriteString(run.TestStructure.GetRunName() + "|" +
			run.TestStructure.GetStatus() + "|" +
			run.TestStructure.GetResult() + "|" +
			run.TestStructure.GetQueued() + "|" +
			startTimeStringRaw + "|" +
			endTimeStringRaw + "|" +
			duration + "|" +
			run.TestStructure.GetTestName() + "|" +
			run.TestStructure.GetRequestor() + "|" +
			run.TestStructure.GetBundle() + "|" +
			runLog,
		)

		buff.WriteString("\n")

	}
	result = buff.String()
	return result, err
}
