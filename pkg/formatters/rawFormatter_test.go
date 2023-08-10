/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package formatters

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/stretchr/testify/assert"
)

func createRunForRaw(runId string,
	runName string,
	status string,
	result string,
	bundle string,
	testName string,
	requestor string,
	queued string,
	startTime string,
	endTime string,
) galasaapi.Run {

	testStructure := galasaapi.TestStructure{
		RunName:   &runName,
		Bundle:    &bundle,
		TestName:  &testName,
		Requestor: &requestor,
		Status:    &status,
		Result:    &result,
		Queued:    &queued,
		StartTime: &startTime,
		EndTime:   &endTime,
	}
	run1 := galasaapi.Run{
		RunId:         &runId,
		TestStructure: &testStructure,
	}
	return run1
}

func TestRawFormatterNoDataReturnsNothing(t *testing.T) {

	formatter := NewRawFormatter()
	// No data to format...
	runs := make([]galasaapi.Run, 0)
	apiServerUrl := "https://127.0.0.1"

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerUrl)

	assert.Nil(t, err)
	expectedFormattedOutput := ""
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestRawFormatterReturnsExpectedFormat(t *testing.T) {
	formatter := NewRawFormatter()
	apiServerUrl := "https://127.0.0.1"

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForRaw("cbd-123", "U456", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z")
	runs = append(runs, run1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerUrl)

	assert.Nil(t, err)
	expectedFormattedOutput := "U456|Finished|Passed|2023-05-04T10:55:29.545323Z|2023-05-05T06:00:14.496953Z|2023-05-05T06:00:15.654565Z|1157|dev.galasa.Zos3270LocalJava11Ubuntu|galasa|dev.galasa|https://127.0.0.1/ras/runs/cbd-123/runlog\n"

	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestRawFormatterWithMultipleRunsSeparatesWithNewLine(t *testing.T) {
	formatter := NewRawFormatter()
	apiServerUrl := "https://127.0.0.1"

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForRaw("cbd-123", "U123", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z")
	run2 := createRunForRaw("cbd-456", "U456", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z")
	run3 := createRunForRaw("cbd-789", "U789", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z")
	runs = append(runs, run1, run2, run3)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerUrl)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"U123|Finished|Passed|2023-05-04T10:55:29.545323Z|2023-05-05T06:00:14.496953Z|2023-05-05T06:00:15.654565Z|1157|dev.galasa.Zos3270LocalJava11Ubuntu|galasa|dev.galasa|https://127.0.0.1/ras/runs/cbd-123/runlog\n" +
			"U456|Finished|Passed|2023-05-04T10:55:29.545323Z|2023-05-05T06:00:14.496953Z|2023-05-05T06:00:15.654565Z|1157|dev.galasa.Zos3270LocalJava11Ubuntu|galasa|dev.galasa|https://127.0.0.1/ras/runs/cbd-456/runlog\n" +
			"U789|Finished|Passed|2023-05-04T10:55:29.545323Z|2023-05-05T06:00:14.496953Z|2023-05-05T06:00:15.654565Z|1157|dev.galasa.Zos3270LocalJava11Ubuntu|galasa|dev.galasa|https://127.0.0.1/ras/runs/cbd-789/runlog\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestRawFormatterNoRunEndtimeReturnsBlankEndtimeFieldAndNoDuration(t *testing.T) {
	formatter := NewRawFormatter()
	apiServerUrl := "https://127.0.0.1"

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForRaw("cbd-123", "U456", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "")
	runs = append(runs, run1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerUrl)

	assert.Nil(t, err)
	expectedFormattedOutput := "U456|Finished|Passed|2023-05-04T10:55:29.545323Z|2023-05-05T06:00:14.496953Z|||dev.galasa.Zos3270LocalJava11Ubuntu|galasa|dev.galasa|https://127.0.0.1/ras/runs/cbd-123/runlog\n"

	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}
