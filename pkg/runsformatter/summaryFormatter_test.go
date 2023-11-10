/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runsformatter

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

func createFormattableTestForSummary(queuedTimeUTC string, name string, testName string, status string, result string, requestor string, isLost bool) FormattableTest {
	// run1Id := "ar"
	// bundle := ""
	// testName := ""
	// requestor := ""
	// queued := ""
	// startTime := ""
	// endTime := ""
	// isLost := false
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
		Lost: isLost,
	}
	return formattableTest
}

func TestSummaryFormatterLongResultStringReturnsExpectedFormat(t *testing.T) {
	formatter := NewSummaryFormatter()

	formattableTest := make([]FormattableTest, 0)
	formattableTest1 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName", "Finished", "MyLongResultString", "myUserId1", false)
	formattableTest = append(formattableTest, formattableTest1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"submitted-time(UTC) name requestor status   result             test-name\n" +
			"2023-05-04 10:55:29 U456 myUserId1 Finished MyLongResultString MyTestName\n" +
			"\n" +
			"Total:1\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterShortResultStringReturnsExpectedFormat(t *testing.T) {
	formatter := NewSummaryFormatter()

	formattableTest := make([]FormattableTest, 0)
	formattableTest1 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName", "Finished", "Short", "myUserId1", false)
	formattableTest = append(formattableTest, formattableTest1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"submitted-time(UTC) name requestor status   result test-name\n" +
			"2023-05-04 10:55:29 U456 myUserId1 Finished Short  MyTestName\n" +
			"\n" +
			"Total:1\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterShortAndLongStatusReturnsExpectedFormat(t *testing.T) {
	formatter := NewSummaryFormatter()

	formattableTest := make([]FormattableTest, 0)
	formattableTest1 := createFormattableTestForSummary("2023-05-04T10:45:29.545323Z", "LongRunName", "TestName", "LongStatus", "Short", "myUserId1", false)
	formattableTest2 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName", "short", "MyLongResultString", "myUserId1", false)
	formattableTest = append(formattableTest, formattableTest1, formattableTest2)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"submitted-time(UTC) name        requestor status     result             test-name\n" +
			"2023-05-04 10:45:29 LongRunName myUserId1 LongStatus Short              TestName\n" +
			"2023-05-04 10:55:29 U456        myUserId1 short      MyLongResultString MyTestName\n" +
			"\n" +
			"Total:2\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterWithMultipleRunsPrintsOnlyFinishedRuns(t *testing.T) {
	formatter := NewSummaryFormatter()

	formattableTest := make([]FormattableTest, 0)
	formattableTest1 := createFormattableTestForSummary("2023-05-04T10:45:29.545323Z", "U123", "TestName", "Finished", "Passed", "myUserId1", false)
	formattableTest2 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName1", "Finished", "Failed", "myUserId2", false)
	formattableTest3 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U789", "MyTestName2", "Finished", "EnvFail", "myUserId1", false)
	formattableTest4 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L123", "MyTestName3", "UNKNOWN", "", "myUserId2", false)
	formattableTest5 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L456", "MyTestName4", "Building", "EnvFail", "myUserId1", false)
	formattableTest6 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L789", "MyTestName5", "Finished", "Passed With Defects", "myUserId2", false)
	formattableTest7 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C111", "MyTestName6", "Finished", "Failed", "myUserId1", false)
	formattableTest8 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C222", "MyTestName7", "Finished", "UNKNOWN", "myUserId2", false)
	formattableTest9 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C333", "MyTestName8", "Finished", "Ignored", "myUserId1", false)
	formattableTest = append(formattableTest, formattableTest1, formattableTest2, formattableTest3, formattableTest4, formattableTest5, formattableTest6, formattableTest7, formattableTest8, formattableTest9)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"submitted-time(UTC) name requestor status   result              test-name\n" +
			"2023-05-04 10:45:29 U123 myUserId1 Finished Passed              TestName\n" +
			"2023-05-04 10:55:29 U456 myUserId2 Finished Failed              MyTestName1\n" +
			"2023-05-04 10:55:29 U789 myUserId1 Finished EnvFail             MyTestName2\n" +
			"2023-05-04 10:55:29 L123 myUserId2 UNKNOWN                      MyTestName3\n" +
			"2023-05-04 10:55:29 L456 myUserId1 Building EnvFail             MyTestName4\n" +
			"2023-05-04 10:55:29 L789 myUserId2 Finished Passed With Defects MyTestName5\n" +
			"2023-05-04 10:55:29 C111 myUserId1 Finished Failed              MyTestName6\n" +
			"2023-05-04 10:55:29 C222 myUserId2 Finished UNKNOWN             MyTestName7\n" +
			"2023-05-04 10:55:29 C333 myUserId1 Finished Ignored             MyTestName8\n" +
			"\n" +
			"Total:9 Passed:1 PassedWithDefects:1 Failed:2 EnvFail:2 UNKNOWN:1 Active:1 Ignored:1\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterMultipleRunsWithLostRunsDoesNotDisplayLostRunsAndCountsThem(t *testing.T) {
	formatter := NewSummaryFormatter()

	formattableTest := make([]FormattableTest, 0)
	formattableTest1 := createFormattableTestForSummary("2023-05-04T10:45:29.545323Z", "U123", "TestName", "Finished", "Passed", "myUserId1", false)
	formattableTest2 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName1", "Finished", "Failed", "myUserId2", true)
	formattableTest3 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U789", "MyTestName2", "Finished", "EnvFail", "myUserId1", true)
	formattableTest4 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L123", "MyTestName3", "UNKNOWN", "", "myUserId2", false)
	formattableTest5 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L456", "MyTestName4", "Building", "EnvFail", "myUserId1", false)
	formattableTest6 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L789", "MyTestName5", "Finished", "Passed With Defects", "myUserId2", false)
	formattableTest7 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C111", "MyTestName6", "Finished", "Failed", "myUserId1", false)
	formattableTest8 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C222", "MyTestName7", "Finished", "UNKNOWN", "myUserId2", false)
	formattableTest9 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C333", "MyTestName8", "Finished", "Ignored", "myUserId1", true)
	//formattableTest10 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L321", "MyTestName9", "UNKNOWN", "", "myUserId2", true)
	formattableTest = append(formattableTest, formattableTest1, formattableTest2, formattableTest3, formattableTest4, formattableTest5, formattableTest6, formattableTest7, formattableTest8, formattableTest9)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"submitted-time(UTC) name requestor status   result              test-name\n" +
			"2023-05-04 10:45:29 U123 myUserId1 Finished Passed              TestName\n" +
			"2023-05-04 10:55:29 L123 myUserId2 UNKNOWN                      MyTestName3\n" +
			"2023-05-04 10:55:29 L456 myUserId1 Building EnvFail             MyTestName4\n" +
			"2023-05-04 10:55:29 L789 myUserId2 Finished Passed With Defects MyTestName5\n" +
			"2023-05-04 10:55:29 C111 myUserId1 Finished Failed              MyTestName6\n" +
			"2023-05-04 10:55:29 C222 myUserId2 Finished UNKNOWN             MyTestName7\n" +
			"\n" +
			"Total:9 Passed:1 PassedWithDefects:1 Failed:1 Lost:3 EnvFail:1 UNKNOWN:1 Active:1\n"

	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterMultipleRunsWithUnknownStatusOfLostRunsDoesNotDisplayLostRuns(t *testing.T) {
	formatter := NewSummaryFormatter()

	formattableTest := make([]FormattableTest, 0)
	formattableTest1 := createFormattableTestForSummary("2023-05-04T10:45:29.545323Z", "U123", "TestName", "Finished", "Passed", "myUserId1", false)
	formattableTest2 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U456", "MyTestName1", "Finished", "Failed", "myUserId2", true)
	formattableTest3 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "U789", "MyTestName2", "Finished", "EnvFail", "myUserId1", true)
	formattableTest4 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L123", "MyTestName3", "UNKNOWN", "", "myUserId2", false)
	formattableTest5 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L456", "MyTestName4", "Building", "EnvFail", "myUserId1", false)
	formattableTest6 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L789", "MyTestName5", "Finished", "Passed With Defects", "myUserId2", false)
	formattableTest7 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C111", "MyTestName6", "Finished", "Failed", "myUserId1", false)
	formattableTest8 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C222", "MyTestName7", "Finished", "UNKNOWN", "myUserId2", false)
	formattableTest9 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "C333", "MyTestName8", "Finished", "Ignored", "myUserId1", true)
	formattableTest10 := createFormattableTestForSummary("2023-05-04T10:55:29.545323Z", "L321", "MyTestName9", "UNKNOWN", "", "myUserId2", true)
	formattableTest = append(formattableTest, formattableTest1, formattableTest2, formattableTest3, formattableTest4, formattableTest5, formattableTest6, formattableTest7, formattableTest8, formattableTest9, formattableTest10)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"submitted-time(UTC) name requestor status   result              test-name\n" +
			"2023-05-04 10:45:29 U123 myUserId1 Finished Passed              TestName\n" +
			"2023-05-04 10:55:29 L123 myUserId2 UNKNOWN                      MyTestName3\n" +
			"2023-05-04 10:55:29 L456 myUserId1 Building EnvFail             MyTestName4\n" +
			"2023-05-04 10:55:29 L789 myUserId2 Finished Passed With Defects MyTestName5\n" +
			"2023-05-04 10:55:29 C111 myUserId1 Finished Failed              MyTestName6\n" +
			"2023-05-04 10:55:29 C222 myUserId2 Finished UNKNOWN             MyTestName7\n" +
			"\n" +
			"Total:10 Passed:1 PassedWithDefects:1 Failed:1 Lost:4 EnvFail:1 UNKNOWN:1 Active:1\n"

	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterHasTestWithoutTimeStamps(t *testing.T) {
	formatter := NewSummaryFormatter()

	formattableTest := make([]FormattableTest, 0)
	formattableTest1 := createFormattableTestForSummary("", "U123", "TestName", "Finished", "Passed", "myUserId1", false)
	formattableTest = append(formattableTest, formattableTest1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(formattableTest)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"submitted-time(UTC) name requestor status   result test-name\n" +
			"                    U123 myUserId1 Finished Passed TestName\n" +
			"\n" +
			"Total:1 Passed:1\n"

	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}
