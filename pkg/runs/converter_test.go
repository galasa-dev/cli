/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/formatters"
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
	output := FormattableTestFromGalasaApi(runs, apiServerUrl)

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
	output := FormattableTestFromGalasaApi(runs, apiServerUrl)

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
	output := FormattableTestFromGalasaApi(runs, apiServerUrl)

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
	output := FormattableTestFromTestRun(finishedRunsMap, lostRunsMap)

	//Then
	assert.Equal(t, 0, len(output), "The input record is empty and so should be the output record")

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
	output := FormattableTestFromTestRun(finishedRunsMap, lostRunsMap)

	//Then
	assert.Equal(t, total, len(output), "The input record has a length of %v whilst the output has length of %v", total, len(output))
}

func TestFormattableTestsArePrintedInOrder(t *testing.T) {
	//Given
	var formattableTest []formatters.FormattableTest
	formattableTest1 := formatters.FormattableTest{
		RunId:    "id1",
		Name:     "formattableTest1",
		TestName: "testName1",
		Status:   "status1",
		Result:   "Failed",
		//StartTimeUTC  string
		//EndTimeUTC    string
		//QueuedTimeUTC string
		Requestor:    "Requestor1",
		Bundle:       "bundle1",
		ApiServerUrl: "127.0.0.1",
		Methods:      nil,
		Lost:         false,
	}
	formattableTest2 := formatters.FormattableTest{
		RunId:    "id2",
		Name:     "formattableTest2",
		TestName: "testName2",
		Status:   "status2",
		Result:   "Failed",
		//StartTimeUTC  string
		//EndTimeUTC    string
		//QueuedTimeUTC string
		Requestor:    "Requestor2",
		Bundle:       "bundle2",
		ApiServerUrl: "127.0.0.1",
		Methods:      nil,
		Lost:         false,
	}
	formattableTest3 := formatters.FormattableTest{
		RunId:    "id3",
		Name:     "formattableTest3",
		TestName: "testName3",
		Status:   "status3",
		Result:   "Passed",
		//StartTimeUTC  string
		//EndTimeUTC    string
		//QueuedTimeUTC string
		Requestor:    "Requestor3",
		Bundle:       "bundle3",
		ApiServerUrl: "137.0.0.1",
		Methods:      nil,
		Lost:         false,
	}
	formattableTest4 := formatters.FormattableTest{
		RunId:    "id4",
		Name:     "formattableTest4",
		TestName: "testName4",
		Status:   "status4",
		Result:   "Custard",
		//StartTimeUTC  string
		//EndTimeUTC    string
		//QueuedTimeUTC string
		Requestor:    "Requestor4",
		Bundle:       "bundle4",
		ApiServerUrl: "147.0.0.1",
		Methods:      nil,
		Lost:         false,
	}
	formattableTest5 := formatters.FormattableTest{
		RunId:    "id5",
		Name:     "formattableTest5",
		TestName: "testName5",
		Status:   "status5",
		Result:   "Doughnuts",
		//StartTimeUTC  string
		//EndTimeUTC    string
		//QueuedTimeUTC string
		Requestor:    "Requestor5",
		Bundle:       "bundle5",
		ApiServerUrl: "157.0.0.1",
		Methods:      nil,
		Lost:         false,
	}
	formattableTest6 := formatters.FormattableTest{
		RunId:    "id6",
		Name:     "formattableTest6",
		TestName: "testName6",
		Status:   "status6",
		Result:   "Custom",
		//StartTimeUTC  string
		//EndTimeUTC    string
		//QueuedTimeUTC string
		Requestor:    "Requestor6",
		Bundle:       "bundle6",
		ApiServerUrl: "167.0.0.1",
		Methods:      nil,
		Lost:         false,
	}
	formattableTest7 := formatters.FormattableTest{
		RunId:    "id7",
		Name:     "formattableTest7",
		TestName: "testName7",
		Status:   "status7",
		Result:   "Passed With Defects",
		//StartTimeUTC  string
		//EndTimeUTC    string
		//QueuedTimeUTC string
		Requestor:    "Requestor7",
		Bundle:       "bundle7",
		ApiServerUrl: "177.0.0.1",
		Methods:      nil,
		Lost:         false,
	}
	formattableTest = append(formattableTest, formattableTest1, formattableTest2, formattableTest3, formattableTest4, formattableTest5, formattableTest6, formattableTest7)
	//When
	orderedFormattableTest := orderFormattableTests(formattableTest)
	//Then
	assert.Equal(t, len(formattableTest), len(orderedFormattableTest))
	assert.Equal(t, "Passed", orderedFormattableTest[0].Result)
	assert.Equal(t, "Passed With Defects", orderedFormattableTest[1].Result)
	assert.Equal(t, "Failed", orderedFormattableTest[2].Result)
	assert.Equal(t, "Failed", orderedFormattableTest[3].Result)
	assert.Equal(t, "Custard", orderedFormattableTest[4].Result)
	assert.Equal(t, "Custom", orderedFormattableTest[5].Result)
	assert.Equal(t, "Doughnuts", orderedFormattableTest[6].Result)
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
	output := FormattableTestFromTestRun(finishedRunsMap, lostRunsMap)

	//Then
	assert.Equal(t, total, len(output), "The input record has a length of %v whilst the output has length of %v", total, len(output))
	//lostRuns are always appended last so checking if last appended test with len(output)-1
	assert.Equal(t, true, output[len(output)-1].Lost)
}

func TestRunsOfTestRunStructArePrintedInSortedOrder(t *testing.T) {
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

	finished4 := TestRun{
		Name:      "myTestRun4",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Failed",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finished5 := TestRun{
		Name:      "myTestRun5",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "Apples",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{{Method: "method1", Result: "passed"}, {Method: "method2", Result: "passed"}}}

	finishedRunsMap := make(map[string]*TestRun, 3)
	finishedRunsMap["myTestRun1"] = &finished1
	finishedRunsMap["myTestRun2"] = &finished2
	finishedRunsMap["myTestRun3"] = &finished3
	finishedRunsMap["myTestRun4"] = &finished4
	finishedRunsMap["myTestRun5"] = &finished5

	lostRunsMap := make(map[string]*TestRun)

	total := len(finishedRunsMap) + len(lostRunsMap)

	//When
	output := FormattableTestFromTestRun(finishedRunsMap, lostRunsMap)

	//Then
	assert.Equal(t, total, len(output), "The input record has a length of %v whilst the output has length of %v", total, len(output))
	assert.Equal(t, "Passed", output[0].Result)
	assert.Equal(t, "Failed", output[1].Result)
	assert.Equal(t, "Apples", output[2].Result)
	assert.Equal(t, "Custard", output[3].Result)
	assert.Equal(t, "Custom", output[4].Result)
}
