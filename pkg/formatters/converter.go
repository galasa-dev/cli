/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package formatters

import "github.com/galasa.dev/cli/pkg/galasaapi"

func NewFormattableTestFromGalasaApi(runs []galasaapi.Run, apiServerUrl string) []FormattableTest {
	var formattableTest []FormattableTest

	for _, run := range runs {
		//Get the data for each TestStructure in runs
		newFormattableTest := getTestStructureData(run, apiServerUrl)
		formattableTest = append(formattableTest, newFormattableTest)
	}
	return formattableTest
}

func getTestStructureData(run galasaapi.Run, apiServerUrl string) FormattableTest {
	newFormattableTest := NewFormattableTest()

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
