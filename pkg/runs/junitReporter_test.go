/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"strings"
	"testing"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/stretchr/testify/assert"
)

func TestJunitReportWorks(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()

	finishedRuns := TestRun{
		Name:      "myTestRun",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "PASSED",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finishedRunsMap := make(map[string]*TestRun, 1)
	finishedRunsMap["myTestRun"] = &finishedRuns

	// When...
	err := ReportJunit(
		mockFileSystem,
		"myReportJunitFilename",
		"myGroup",
		finishedRunsMap, nil)

	// Then...
	if err != nil {
		assert.Fail(t, "ReportJunit failed when it should have passed. "+err.Error())
	}

	isExists, err := mockFileSystem.Exists("myReportJunitFilename")
	if err != nil {
		assert.Fail(t, "junit report does not exist in correct place. "+err.Error())
	}
	assert.True(t, isExists, "Junit report does not exist in the correct place.")

	// We expect a report like this:
	expectedReport := `<?xml version="1.0" encoding="UTF-8" ?>
	<testsuites id="myGroup" name="Galasa test run" tests="2" failures="2" time="0">
		<testsuite id="myTestRun" name="myStream/myBundle/com.myco.MyClass" tests="2" failures="2" time="0">
			<testcase id="method1" name="method1" time="0">
				<failure message="Failure messages are unavailable at this time" type="Unknown"></failure>
			</testcase>
			<testcase id="method2" name="method2" time="0">
				<failure message="Failure messages are unavailable at this time" type="Unknown"></failure>
			</testcase>
		</testsuite>
	</testsuites>`

	actualContents, err := mockFileSystem.ReadTextFile("myReportJunitFilename")
	if err != nil {
		assert.Fail(t, "Could not read the junit file. "+err.Error())
	}

	expected1 := strings.ReplaceAll(expectedReport, "\n", "")
	expected2 := strings.ReplaceAll(expected1, "\t", "")
	expected3 := strings.ReplaceAll(expected2, " ", "")

	actual1 := strings.ReplaceAll(actualContents, "\n", "")
	actual2 := strings.ReplaceAll(actual1, "\t", "")
	actual3 := strings.ReplaceAll(actual2, " ", "")

	assert.EqualValues(t, actual3, expected3)
}
