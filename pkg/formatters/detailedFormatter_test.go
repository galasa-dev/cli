/*
 * Copyright contributors to the Galasa project
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

func createRunForDetailed(runName string,
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
		//RunId:         &run1Id,
		TestStructure: &testStructure,
	}
	return run1
}

func TestDetailedFormatterNoDataReturnsHeadersOnly(t *testing.T) {

	formatter := NewDetailedFormatter()
	// No data to format...
	runs := make([]galasaapi.Run, 0)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs)

	assert.Nil(t, err)
	expectedFormattedOutput := ""
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestDetailedFormatterReturnsExpectedFormat(t *testing.T) {
	formatter := NewDetailedFormatter()

	methods := make([]galasaapi.TestMethod, 0)
	method1 := createMethod("testCoreIvtTest", "test", "finished", "passed", "2023-05-05T06:03:38.872894Z", "2023-05-05T06:03:47.222758Z")
	methods = append(methods, method1)

	runs := make([]galasaapi.Run, 0)
	run1 := createRunForDetailed("U456", "Finished", "Passed", "dev.galasa", "dev.galasa.Zos3270LocalJava11Ubuntu", "galasa", "2023-05-05T06:00:14.496953Z", "2023-05-05T06:04:37.654565Z", "2023-05-04T10:55:29.545323Z", methods)
	runs = append(runs, run1)

	// When...
	actualFormattedOutput, err := formatter.FormatRuns(runs)

	assert.Nil(t, err)
	expectedFormattedOutput :=
		"name        :  U456\n" +
			"status      :  Finished\n" +
			"result      :  Passed\n" +
			"queued-time :  2023-05-05T06:00:14.496953Z\n" +
			"start-time  :  2023-05-05T06:04:37.654565Z\n" +
			"end-time    :  2023-05-04T10:55:29.545323Z\n" +
			"duration(ms):  \n" +
			"test-name   :  dev.galasa.Zos3270LocalJava11Ubuntu\n" +
			"requestor   :  galasa\n" +
			"bundle      :  dev.galasa\n" +
			"run-log     :  \n" +
			"\n" +
			"method          type status   result start-time                  end-time                    duration(ms)\n" +
			"testCoreIvtTest test finished passed 2023-05-05T06:03:38.872894Z 2023-05-05T06:03:47.222758Z \n"

	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}
