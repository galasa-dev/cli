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

func TestSummaryFormatterNoDataReturnsTotalCountAllZeros(t *testing.T) {

	formatter := NewSummaryFormatter()
	// No data to format...
	runs := make([]galasaapi.Run, 0)
	apiServerURL := ""

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerURL)

	assert.Nil(t, err)
	expectedFormattedOutput := "Total:0\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func createRunForSummary(queued string, runName string, testName string, status string, result string) galasaapi.Run {
	//run1Id := "ar"
	//bundle := ""
	//testName := ""
	//requestor := ""
	// queued := ""
	// startTime := ""
	// endTime := ""
	testStructure := galasaapi.TestStructure{
		RunName: &runName,
		//Bundle:        &bundle,
		//TestName:      &testName,
		TestName: &testName,
		//Requestor:     &requestor,
		Status: &status,
		Result: &result,
		Queued: &queued,
		// StartTime:     &startTime,
		// EndTime:       &endTime,
	}
	run1 := galasaapi.Run{
		//RunId:         &run1Id,
		TestStructure: &testStructure,
	}
	return run1
}

func TestSummaryFormatterLongResultStringReturnsExpectedFormat(t *testing.T) {
	formatter := NewSummaryFormatter()

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName", "Finished", "MyLongResultString")
	runs = append(runs, run1)
	apiServerURL := ""

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerURL)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"submitted-time      name status   result             test-name\n" +
			"2023-05-04 10:55:29 U456 Finished MyLongResultString MyTestName\n" +
			"\n" +
			"Total:1\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterShortResultStringReturnsExpectedFormat(t *testing.T) {
	formatter := NewSummaryFormatter()

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName", "Finished", "Short")
	runs = append(runs, run1)
	apiServerURL := ""

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerURL)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"submitted-time      name status   result test-name\n" +
			"2023-05-04 10:55:29 U456 Finished Short  MyTestName\n" +
			"\n" +
			"Total:1\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterShortAndLongStatusReturnsExpectedFormat(t *testing.T) {
	formatter := NewSummaryFormatter()

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForSummary("2023-05-04T10:45:29.545323Z", "LongRunName", "TestName", "LongStatus", "Short")
	run2 := createRunForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName", "short", "MyLongResultString")
	runs = append(runs, run1, run2)
	apiServerURL := ""

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerURL)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"submitted-time      name        status     result             test-name\n" +
			"2023-05-04 10:45:29 LongRunName LongStatus Short              TestName\n" +
			"2023-05-04 10:55:29 U456        short      MyLongResultString MyTestName\n" +
			"\n" +
			"Total:2\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterMultipleRunsDifferentResultsProducesExpectedTotalsCount(t *testing.T) {
	formatter := NewSummaryFormatter()

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForSummary("2023-05-04T10:45:29.545323Z", "U123", "TestName", "Finished", "Passed")
	run2 := createRunForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName1", "Finished", "Failed")
	run3 := createRunForSummary("2023-05-04T10:55:29.545323Z", "U789", "MyTestName2", "Finished", "EnvFail")
	run4 := createRunForSummary("2023-05-04T10:55:29.545323Z", "L123", "MyTestName3", "UNKNOWN", "")
	run5 := createRunForSummary("2023-05-04T10:55:29.545323Z", "L456", "MyTestName4", "Building", "EnvFail")
	run6 := createRunForSummary("2023-05-04T10:55:29.545323Z", "L789", "MyTestName5", "Finished", "Passed With Defects")
	run7 := createRunForSummary("2023-05-04T10:55:29.545323Z", "C111", "MyTestName6", "Finished", "Failed")
	run8 := createRunForSummary("2023-05-04T10:55:29.545323Z", "C222", "MyTestName7", "Finished", "UNKNOWN")
	run9 := createRunForSummary("2023-05-04T10:55:29.545323Z", "C333", "MyTestName8", "Finished", "Ignored")
	runs = append(runs, run1, run2, run3, run4, run5, run6, run7, run8, run9)
	apiServerURL := ""

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerURL)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"submitted-time      name status   result              test-name\n" +
			"2023-05-04 10:45:29 U123 Finished Passed              TestName\n" +
			"2023-05-04 10:55:29 U456 Finished Failed              MyTestName1\n" +
			"2023-05-04 10:55:29 U789 Finished EnvFail             MyTestName2\n" +
			"2023-05-04 10:55:29 L123 UNKNOWN                      MyTestName3\n" +
			"2023-05-04 10:55:29 L456 Building EnvFail             MyTestName4\n" +
			"2023-05-04 10:55:29 L789 Finished Passed With Defects MyTestName5\n" +
			"2023-05-04 10:55:29 C111 Finished Failed              MyTestName6\n" +
			"2023-05-04 10:55:29 C222 Finished UNKNOWN             MyTestName7\n" +
			"2023-05-04 10:55:29 C333 Finished Ignored             MyTestName8\n" +
			"\n" +
			"Total:9 Passed:1 PassedWithDefects:1 Failed:2 EnvFail:2 UNKNOWN:1 Active:1 Ignored:1\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}
