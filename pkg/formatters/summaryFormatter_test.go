/*
 * Copyright contributors to the Galasa project
 */
package formatters

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/stretchr/testify/assert"
)

func TestSummaryFormatterNoDataReturnsHeadersOnly(t *testing.T) {

	formatter := NewSummaryFormatter()
	// No data to format...
	runs := make([]galasaapi.Run, 0)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs)

	assert.Nil(t, err)
	expectedFormattedOutput := ""
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func createRunForSummary(runName string, testShortName string, status string, result string) galasaapi.Run {
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
		TestShortName: &testShortName,
		//Requestor:     &requestor,
		Status: &status,
		Result: &result,
		// Queued:        &queued,
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
	run1 := createRunForSummary("U456", "MyTestName", "Finished", "MyLongResultString")
	runs = append(runs, run1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"name status   result             test-name\n" +
			"U456 Finished MyLongResultString MyTestName\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterShortResultStringReturnsExpectedFormat(t *testing.T) {
	formatter := NewSummaryFormatter()

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForSummary("U456", "MyTestName", "Finished", "Short")
	runs = append(runs, run1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"name status   result test-name\n" +
			"U456 Finished Short  MyTestName\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterShortAndLongStatusReturnsExpectedFormat(t *testing.T) {
	formatter := NewSummaryFormatter()

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForSummary("LongRunName", "TestName", "LongStatus", "Short")
	run2 := createRunForSummary("U456", "MyTestName", "short", "MyLongResultString")
	runs = append(runs, run1, run2)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"name        status     result             test-name\n" +
			"LongRunName LongStatus Short              TestName\n" +
			"U456        short      MyLongResultString MyTestName\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}
