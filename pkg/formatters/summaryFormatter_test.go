/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSummaryFormatterNoDataReturnsTotalCountAllZeros(t *testing.T) {

	formatter := NewSummaryFormatter()
	// No data to format...
	formattableTest := make([]FormattableTest, 0)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

	assert.Nil(t, err)
	expectedFormattedOutput := "Total:0\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func createFormattableTestForSummary(queuedTimeUTC string, name string, testName string, status string, result string, requestor string) FormattableTest {
	//run1Id := "ar"
	//bundle := ""
	//testName := ""
	//requestor := ""
	// queued := ""
	// startTime := ""
	// endTime := ""
	formattableTest := FormattableTest{
		Name: name,
		//Bundle:        &bundle,
		//TestName:      &testName,
		TestName:      testName,
		Requestor:     requestor,
		Status:        status,
		Result:        result,
		QueuedTimeUTC: queuedTimeUTC,
		// StartTime:     &startTime,
		// EndTime:       &endTime,
	}
	return formattableTest
}

func TestSummaryFormatterLongResultStringReturnsExpectedFormat(t *testing.T) {
	formatter := NewSummaryFormatter()

	formattableTest := make([]FormattableTest, 0)
	formattableTest1 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName", "Finished", "MyLongResultString", "myUserId1")
	formattableTest = append(formattableTest, formattableTest1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

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

	formattableTest := make([]FormattableTest, 0)
	formattableTest1 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName", "Finished", "Short", "myUserId1")
	formattableTest = append(formattableTest, formattableTest1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

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

	formattableTest := make([]FormattableTest, 0)
	formattableTest1 := createFormattableTestForSummary("2023-05-04T10:45:29.545323Z", "LongRunName", "TestName", "LongStatus", "Short", "myUserId1")
	formattableTest2 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName", "short", "MyLongResultString", "myUserId1")
	formattableTest = append(formattableTest, formattableTest1, formattableTest2)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

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

	formattableTest := make([]FormattableTest, 0)
	formattableTest1 := createFormattableTestForSummary("2023-05-04T10:45:29.545323Z", "U123", "TestName", "Finished", "Passed", "myUserId1")
	formattableTest2 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName1", "Finished", "Failed", "myUserId2")
	formattableTest3 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U789", "MyTestName2", "Finished", "EnvFail", "myUserId1")
	formattableTest4 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L123", "MyTestName3", "UNKNOWN", "", "myUserId2")
	formattableTest5 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L456", "MyTestName4", "Building", "EnvFail", "myUserId1")
	formattableTest6 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L789", "MyTestName5", "Finished", "Passed With Defects", "myUserId2")
	formattableTest7 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C111", "MyTestName6", "Finished", "Failed", "myUserId1")
	formattableTest8 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C222", "MyTestName7", "Finished", "UNKNOWN", "myUserId2")
	formattableTest9 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C333", "MyTestName8", "Finished", "Ignored", "myUserId1")
	formattableTest = append(formattableTest, formattableTest1, formattableTest2, formattableTest3, formattableTest4, formattableTest5, formattableTest6, formattableTest7, formattableTest8, formattableTest9)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

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
