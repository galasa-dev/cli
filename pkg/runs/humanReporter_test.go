/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanCountFailuresNoLostOnePass(t *testing.T) {
	// Given...
	var lostRunsMap map[string]*TestRun = make(map[string]*TestRun, 0)

	finishedRuns := TestRun{
		Name:      "myTestRun",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Passed",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finishedRunsMap := make(map[string]*TestRun, 1)
	finishedRunsMap["myTestRun"] = &finishedRuns

	// When
	count := CountTotalFailedRuns(finishedRunsMap, lostRunsMap)

	assert.Equal(t, 0, count, "Failed to count failed test cases.")
}

func TestCanCountFailuresOneLostOneFailed(t *testing.T) {
	// Given...
	lostRun1 := TestRun{
		Name:      "myTestRun",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Failed",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finished1 := TestRun{
		Name:      "myTestRun",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Failed",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finishedRunsMap := make(map[string]*TestRun, 1)
	finishedRunsMap["myTestRun"] = &finished1

	lostRunsMap := make(map[string]*TestRun, 1)
	lostRunsMap["myTestRun"] = &lostRun1

	// When
	count := CountTotalFailedRuns(finishedRunsMap, lostRunsMap)

	assert.Equal(t, 2, count, "Failed to count failde test cases.")
}

func TestCanCallHumanReportNoErrors(t *testing.T) {
	// Given...
	lostRun1 := TestRun{
		Name:      "myTestRun",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Failed",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finished1 := TestRun{
		Name:      "myTestRun",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Failed",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finishedRunsMap := make(map[string]*TestRun, 1)
	finishedRunsMap["myTestRun"] = &finished1

	lostRunsMap := make(map[string]*TestRun, 1)
	lostRunsMap["myTestRun"] = &lostRun1

	// When
	reportText := FinalHumanReadableReportAsString(finishedRunsMap, lostRunsMap)

	assert.Contains(t, reportText, "Final report")
}

func TestCanCallHumanInterrimReportNoErrors(t *testing.T) {
	// Given...
	lostRun1 := TestRun{
		Name:      "myTestRun",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Failed",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finished1 := TestRun{
		Name:      "myTestRun",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Failed",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finishedRunsMap := make(map[string]*TestRun, 1)
	finishedRunsMap["myTestRun"] = &finished1

	lostRunsMap := make(map[string]*TestRun, 1)
	lostRunsMap["myTestRun"] = &lostRun1

	// When
	reportText := InterrimProgressReportAsString(make([]TestRun, 0), finishedRunsMap, finishedRunsMap, lostRunsMap, 5)

	assert.Contains(t, reportText, "Progress report")
}
