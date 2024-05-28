/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runsformatter

import "strings"

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

func (*SummaryFormatter) FormatRuns(testResultsData []FormattableTest) (string, error) {
	var result string = ""
	var err error
	buff := strings.Builder{}
	totalResults := len(testResultsData)
	resultCountsMap := initialiseResultMap()

	if totalResults > 0 {
		var table [][]string

		var headers = []string{HEADER_SUBMITTED_TIME, HEADER_RUNNAME, HEADER_REQUESTOR, HEADER_STATUS, HEADER_RESULT, HEADER_TEST_NAME}

		table = append(table, headers)
		for _, run := range testResultsData {
			if run.Lost {
				resultCountsMap[RUN_RESULT_LOST] += 1
			} else {
				var line []string
				submittedTime := run.QueuedTimeUTC
				submittedTimeReadable := formatTimeReadable(submittedTime)

				accumulateResults(resultCountsMap, run)

				line = append(line, submittedTimeReadable, run.Name, run.Requestor, run.Status, run.Result, run.TestName)
				table = append(table, line)
			}
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
