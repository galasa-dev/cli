/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/stretchr/testify/assert"
)

func CreateMethod(methodName string,
	methodType string,
	status string,
	result string,
	startTime string,
	endTime string) galasaapi.TestMethod {

	method := galasaapi.TestMethod{
		MethodName: &methodName,
		Type:       &methodType,
		Status:     &status,
		Result:     &result,
		StartTime:  &startTime,
		EndTime:    &endTime,
	}
	return method
}

func createRunForConverter(queued string, runName string,
	testName string,
	status string,
	result string,
	methods []galasaapi.TestMethod) galasaapi.Run {
	run1Id := "ar"
	bundle := ""
	// testName = ""
	requestor := ""
	// queued := ""
	startTime := ""
	endTime := ""

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
		Methods:   methods,
	}
	run1 := galasaapi.Run{
		RunId:         &run1Id,
		TestStructure: &testStructure,
	}
	return run1
}

func TestGalasaapiRunHasNoRecordsReturnsNoRecords(t *testing.T) {
	//Given
	runs := make([]galasaapi.Run, 0)
	apiServerUrl := ""

	//When
	output := NewFormattableTestFromGalasaApi(runs, apiServerUrl)

	//Then
	assert.Equal(t, len(output), 0, "The input record is empty and so should be the output record")

}

func TestGalasaapiRunHasRecordsReturnsSameAmountOfRecordsWithNoMethods(t *testing.T) {
	//Given
	methods := make([]galasaapi.TestMethod, 0)

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForConverter("2023-05-04T10:45:29.545323Z", "LongRunName", "TestName", "LongStatus", "Short", methods)
	run2 := createRunForConverter("2023-05-04T10:55:29.545323Z", "U456", "MyTestName", "short", "MyLongResultString", methods)
	runs = append(runs, run1, run2)
	apiServerUrl := ""

	//When
	output := NewFormattableTestFromGalasaApi(runs, apiServerUrl)

	//Then
	assert.Equal(t, len(runs), len(output), "The input record has a length of %v whilst the output has length of %v", len(runs), len(output))
}

func TestGalasaapiRunHasRecordsReturnsSameAmountOfRecordsWithMethods(t *testing.T) {
	//Given
	methods := make([]galasaapi.TestMethod, 0)
	method1 := CreateMethod("testCoreIvtTest", "test", "finished", "passed",
		"2023-05-05T06:03:38.872894Z", "2023-05-05T06:03:39.222758Z")
	method2 := CreateMethod("testCoreIvtTest2", "test2", "finished2", "passed2",
		"2023-05-05T06:03:38.872894Z", "2023-05-05T06:03:39.222758Z")
	methods = append(methods, method1, method2)

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForConverter("2023-05-04T10:45:29.545323Z", "LongRunName", "TestName", "LongStatus", "Short", methods)
	run2 := createRunForConverter("2023-05-04T10:55:29.545323Z", "U456", "MyTestName", "short", "MyLongResultString", methods)
	runs = append(runs, run1, run2)
	apiServerUrl := ""

	//When
	output := NewFormattableTestFromGalasaApi(runs, apiServerUrl)

	//Then
	assert.Equal(t, len(runs), len(output), "The input record has a length of %v whilst the output has length of %v", len(runs), len(output))
	assert.Equal(t, len(methods), len(output[0].Methods))
	//check status of first method of first run
	assert.Equal(t, "finished", output[0].Methods[0].GetStatus())
	//check result of second method of first run
	assert.Equal(t, "passed2", output[0].Methods[1].GetResult())
}

func TestTestRunHasNoRecordsReturnsNoRecords(t *testing.T) {
	//Given
	var finishedRunsMap map[string]*TestRun = make(map[string]*TestRun, 0)
	var lostRunsMap map[string]*TestRun = make(map[string]*TestRun, 0)

	//When
	output := NewFormattableTestFromTestRun(finishedRunsMap, lostRunsMap)

	//Then
	assert.Equal(t, len(output), 0, "The input record is empty and so should be the output record")

}

func TestRunsTestRunHasRecordsReturnsSameAmountofRecords(t *testing.T) {
	// Given...
	finished1 := TestRun{
		Name:      "myTestRun1",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Passed",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finished2 := TestRun{
		Name:      "myTestRun2",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Custom",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finished3 := TestRun{
		Name:      "myTestRun3",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Custard",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finishedRunsMap := make(map[string]*TestRun, 3)
	finishedRunsMap["myTestRun1"] = &finished1
	finishedRunsMap["myTestRun2"] = &finished2
	finishedRunsMap["myTestRun3"] = &finished3

	lostRunsMap := make(map[string]*TestRun, 0)

	total := len(finishedRunsMap) + len(lostRunsMap)
	//When
	output := NewFormattableTestFromTestRun(finishedRunsMap, lostRunsMap)

	//Then
	assert.Equal(t, total, len(output), "The input record has a length of %v whilst the output has length of %v", total, len(output))
}

func TestRunsTestRunHasRecordsReturnsTrueForLostRecord(t *testing.T) {
	//Given
	// Given...
	lostRun1 := TestRun{
		Name:      "myLostTestRun",
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
	lostRunsMap["myLostTestRun"] = &lostRun1

	total := len(finishedRunsMap) + len(lostRunsMap)
	//When
	output := NewFormattableTestFromTestRun(finishedRunsMap, lostRunsMap)

	//Then
	assert.Equal(t, total, len(output), "The input record has a length of %v whilst the output has length of %v", total, len(output))
	//lostRuns are always appended last so checking if last appended test with len(output)-1
	assert.Equal(t, output[len(output)-1].Lost, true)
}
