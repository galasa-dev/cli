/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
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

func (*RawFormatter) IsNeedingMethodDetails() bool {
	return false
}

func (*RawFormatter) FormatRuns(runs []galasaapi.Run, apiServerUrl string) (string, error) {
	var result string = ""
	var err error = nil
	buff := strings.Builder{}

	for _, run := range runs {
		startTimeStringRaw := run.TestStructure.GetStartTime()
		endTimeStringRaw := run.TestStructure.GetEndTime()

		duration := getDuration(startTimeStringRaw, endTimeStringRaw)

		runLog := apiServerUrl + RAS_RUNS_URL + run.GetRunId() + "/runlog"

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
