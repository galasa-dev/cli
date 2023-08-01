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

func (*SummaryFormatter) IsNeedingMethodDetails() bool {
	return false
}

func (*SummaryFormatter) FormatRuns(runs []galasaapi.Run, apiServerUrl string) (string, error) {
	var result string = ""
	var err error = nil
	buff := strings.Builder{}
	totalResults := len(runs)
	resultCountsMap := initialiseResultMap()
	if len(runs) > 0 {
		var table [][]string

		var headers = []string{HEADER_SUBMITTED_TIME, HEADER_RUNNAME, HEADER_STATUS, HEADER_RESULT, HEADER_TEST_NAME}

		table = append(table, headers)
		for _, run := range runs {
			var line []string
			submittedTime := run.TestStructure.GetQueued()
			submittedTimeReadable := formatTimeReadable(submittedTime)

			accumulateResults(resultCountsMap, run)

			line = append(line, submittedTimeReadable, run.TestStructure.GetRunName(), run.TestStructure.GetStatus(), run.TestStructure.GetResult(), run.TestStructure.GetTestName())
			table = append(table, line)
		}

		columnLengths := calculateMaxLengthOfEachColumn(table)
		writeFormattedTableToStringBuilder(table, &buff, columnLengths)

		buff.WriteString("\n")
	}
	totalReportString := generateResultTotalsReport(totalResults, resultCountsMap)
	buff.WriteString(totalReportString + "\n")

	result = buff.String()
	return result, err
}
