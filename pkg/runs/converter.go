/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"github.com/galasa.dev/cli/pkg/formatters"
	"github.com/galasa.dev/cli/pkg/galasaapi"
)

func NewFormattableTestFromGalasaApi(runs []galasaapi.Run, apiServerUrl string) []formatters.FormattableTest {
	var formattableTest []formatters.FormattableTest

	for _, run := range runs {
		//Get the data for each TestStructure in runs
		newFormattableTest := getTestStructureData(run, apiServerUrl)
		formattableTest = append(formattableTest, newFormattableTest)
	}
	return formattableTest
}

func getTestStructureData(run galasaapi.Run, apiServerUrl string) formatters.FormattableTest {
	newFormattableTest := formatters.NewFormattableTest()

	newFormattableTest.RunId = run.GetRunId()
	newFormattableTest.ApiServerUrl = apiServerUrl

	newFormattableTest.Name = run.TestStructure.GetRunName()
	newFormattableTest.TestName = run.TestStructure.GetTestName()
	newFormattableTest.Status = run.TestStructure.GetStatus()
	newFormattableTest.Result = run.TestStructure.GetResult()
	newFormattableTest.StartTimeUTC = run.TestStructure.GetStartTime()
	newFormattableTest.EndTimeUTC = run.TestStructure.GetEndTime()
	newFormattableTest.QueuedTimeUTC = run.TestStructure.GetQueued()
	newFormattableTest.Requestor = run.TestStructure.GetRequestor()
	newFormattableTest.Bundle = run.TestStructure.GetBundle()
	newFormattableTest.Methods = run.TestStructure.GetMethods()

	return newFormattableTest
}

func NewFormattableTestFromTestRun(finishedMap map[string]*TestRun, lostMap map[string]*TestRun) []formatters.FormattableTest {
	var formattableTest []formatters.FormattableTest
	for _, run := range finishedMap {
		isLost := false
		newFormattableTest := getTestRunData(*run, isLost)
		formattableTest = append(formattableTest, newFormattableTest)
	}
	for _, run := range lostMap {
		isLost := true
		newFormattableTest := getTestRunData(*run, isLost)
		formattableTest = append(formattableTest, newFormattableTest)
	}

	return formattableTest
}

func getTestRunData(run TestRun, isLost bool) formatters.FormattableTest {
	newFormattableTest := formatters.NewFormattableTest()

	newFormattableTest.RunId = ""
	newFormattableTest.ApiServerUrl = ""

	newFormattableTest.Name = run.Name
	newFormattableTest.TestName = run.Stream + "/" + run.Bundle + "/" + run.Class
	newFormattableTest.Status = run.Status
	newFormattableTest.Result = run.Result
	newFormattableTest.StartTimeUTC = ""
	newFormattableTest.EndTimeUTC = ""
	newFormattableTest.QueuedTimeUTC = ""
	newFormattableTest.Requestor = ""
	newFormattableTest.Bundle = run.Bundle
	newFormattableTest.Methods = nil
	newFormattableTest.Lost = isLost

	return newFormattableTest
}
