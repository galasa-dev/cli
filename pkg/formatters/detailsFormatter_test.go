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

func createMethod(methodName string,
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

func createRunForDetails(runId string,
	runName string,
	status string,
	result string,
	bundle string,
	testName string,
	requestor string,
	queued string,
	startTime string,
	endTime string,
	methods []galasaapi.TestMethod) galasaapi.Run {

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
		RunId:         &runId,
		TestStructure: &testStructure,
	}
	return run1
}

func TestDetailsFormatterNoDataReturnsTotalsCountAllZeros(t *testing.T) {

	formatter := NewDetailsFormatter()
	// No data to format...
	runs := make([]galasaapi.Run, 0)
	apiServerUrl := "https://127.0.0.1"

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerUrl)

	assert.Nil(t, err)
	expectedFormattedOutput := "Total:0\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestDetailsFormatterReturnsExpectedFormat(t *testing.T) {
	formatter := NewDetailsFormatter()
	apiServerUrl := "https://127.0.0.1"

	methods := make([]galasaapi.TestMethod, 0)
	method1 := createMethod("testCoreIvtTest", "test", "finished", "passed",
		"2023-05-05T06:03:38.872894Z", "2023-05-05T06:03:39.222758Z")
	methods = append(methods, method1)

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForDetails("cbd-123", "U456", "Finished", "Passed", "dev.galasa",
		"dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z",
		"2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	runs = append(runs, run1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerUrl)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"name           : U456\n" +
			"status         : Finished\n" +
			"result         : Passed\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-123/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"Total:1 Passed:1\n"

	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestDetailsFormatterWithMultipleRunsReturnsSeparatedWithDashes(t *testing.T) {
	formatter := NewDetailsFormatter()
	apiServerUrl := "https://127.0.0.1"

	methods := make([]galasaapi.TestMethod, 0)
	method1 := createMethod("testCoreIvtTest", "test", "finished", "passed", "2023-05-05T06:03:38.872894Z", "2023-05-05T06:03:39.222758Z")
	methods = append(methods, method1)

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForDetails("cbd-123", "U123", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	run2 := createRunForDetails("cbd-456", "U456", "Finished", "Failed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	run3 := createRunForDetails("cbd-789", "U789", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	runs = append(runs, run1, run2, run3)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerUrl)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"name           : U123\n" +
			"status         : Finished\n" +
			"result         : Passed\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-123/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"---" +
			"\n\n" +
			"name           : U456\n" +
			"status         : Finished\n" +
			"result         : Failed\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-456/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"---" +
			"\n\n" +
			"name           : U789\n" +
			"status         : Finished\n" +
			"result         : Passed\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-789/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"Total:3 Passed:2 Failed:1\n"

	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestDetailsNoRunEndtimeReturnsBlankEndtimeFieldAndNoDuration(t *testing.T) {
	formatter := NewDetailsFormatter()
	apiServerUrl := "https://127.0.0.1"

	methods := make([]galasaapi.TestMethod, 0)
	method1 := createMethod("testCoreIvtTest", "test", "finished", "passed", "2023-05-05T06:03:38.872894Z", "2023-05-05T06:03:39.222758Z")
	methods = append(methods, method1)

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForDetails("cbd-123", "U456", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "", methods)
	runs = append(runs, run1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerUrl)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"name           : U456\n" +
			"status         : Finished\n" +
			"result         : Passed\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : \n" +
			"duration(ms)   : \n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-123/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"Total:1 Passed:1\n"

	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestMethodTableRendersOkIfNoEndtime(t *testing.T) {
	formatter := NewDetailsFormatter()
	apiServerUrl := "https://127.0.0.1"

	methods := make([]galasaapi.TestMethod, 0)
	method1 := createMethod("testCoreIvtTest", "test", "finished", "passed", "2023-05-05T06:03:38.872894Z", "")
	methods = append(methods, method1)

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForDetails("cbd-123", "U456", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	runs = append(runs, run1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerUrl)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"name           : U456\n" +
			"status         : Finished\n" +
			"result         : Passed\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-123/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38          \n" +
			"\n" +
			"Total:1 Passed:1\n"

	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestDetailsFormatterMultipleRunsDifferentResultsProducesExpectedTotalsCount(t *testing.T) {
	formatter := NewDetailsFormatter()
	apiServerUrl := "https://127.0.0.1"

	methods := make([]galasaapi.TestMethod, 0)
	method1 := createMethod("testCoreIvtTest", "test", "finished", "passed", "2023-05-05T06:03:38.872894Z", "2023-05-05T06:03:39.222758Z")
	methods = append(methods, method1)

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForDetails("cbd-123", "U123", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	run2 := createRunForDetails("cbd-456", "U456", "Finished", "Failed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	run3 := createRunForDetails("cbd-789", "U789", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	run4 := createRunForDetails("cbd-12345", "C123", "Finished", "Passed With Defects", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	run5 := createRunForDetails("cbd-67890", "C456", "UNKNOWN", "EnvFail", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	run6 := createRunForDetails("cbd-98765", "C789", "Finished", "Failed With Defects", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	run7 := createRunForDetails("cbd-543210", "L111", "Finished", "Failed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:00:15.654565Z", methods)
	run8 := createRunForDetails("cbd-222", "L222", "Building", "", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "", methods)
	run9 := createRunForDetails("cbd-333", "L333", "Generating", "", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-04T10:55:29.545323Z", "2023-05-05T06:00:14.496953Z", "", methods)
	runs = append(runs, run1, run2, run3, run4, run5, run6, run7, run8, run9)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs, apiServerUrl)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"name           : U123\n" +
			"status         : Finished\n" +
			"result         : Passed\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-123/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"---" +
			"\n\n" +
			"name           : U456\n" +
			"status         : Finished\n" +
			"result         : Failed\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-456/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"---" +
			"\n\n" +
			"name           : U789\n" +
			"status         : Finished\n" +
			"result         : Passed\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-789/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"---" +
			"\n\n" +
			"name           : C123\n" +
			"status         : Finished\n" +
			"result         : Passed With Defects\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-12345/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"---" +
			"\n\n" +
			"name           : C456\n" +
			"status         : UNKNOWN\n" +
			"result         : EnvFail\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-67890/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"---" +
			"\n\n" +
			"name           : C789\n" +
			"status         : Finished\n" +
			"result         : Failed With Defects\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-98765/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"---" +
			"\n\n" +
			"name           : L111\n" +
			"status         : Finished\n" +
			"result         : Failed\n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : 2023-05-05 06:00:15\n" +
			"duration(ms)   : 1157\n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-543210/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"---" +
			"\n\n" +
			"name           : L222\n" +
			"status         : Building\n" +
			"result         : \n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : \n" +
			"duration(ms)   : \n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-222/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"---" +
			"\n\n" +
			"name           : L333\n" +
			"status         : Generating\n" +
			"result         : \n" +
			"submitted-time : 2023-05-04 10:55:29\n" +
			"start-time     : 2023-05-05 06:00:14\n" +
			"end-time       : \n" +
			"duration(ms)   : \n" +
			"test-name      : dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor      : galasa\n" +
			"bundle         : dev.galasa\n" +
			"run-log        : https://127.0.0.1/ras/runs/cbd-333/runlog\n" +
			"\n" +
			"method          type status   result start-time          end-time            duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05 06:03:38 2023-05-05 06:03:39 349\n" +
			"\n" +
			"Total:9 Passed:2 PassedWithDefects:1 Failed:2 FailedWithDefects:1 EnvFail:1 Active:2\n"

	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}
