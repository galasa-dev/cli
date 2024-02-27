/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"sort"

	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/runsformatter"
)

func orderFormattableTests(formattableTest []runsformatter.FormattableTest) []runsformatter.FormattableTest {
	var orderedFormattableTest []runsformatter.FormattableTest

	//get slice of all result labels in ordered form
	orderedResultLabels := getAvailableResultLabelsinOrder(formattableTest)

	//formattableTest runs grouped by results
	//map["passed"] = [run1, run2, ...]
	runsGroupedByResultsMap := make(map[string][]runsformatter.FormattableTest)
	for _, run := range formattableTest {
		runsGroupedByResultsMap[run.Result] = append(runsGroupedByResultsMap[run.Result], run)
	}

	//append tests in order
	for _, result := range orderedResultLabels {
		orderedFormattableTest = append(orderedFormattableTest, runsGroupedByResultsMap[result]...)
	}
	return orderedFormattableTest
}

func getAvailableResultLabelsinOrder(formattableTest []runsformatter.FormattableTest) []string {
	var orderedResultLabels []string
	orderedResultLabels = append(orderedResultLabels, RESULT_PASSED)
	orderedResultLabels = append(orderedResultLabels, RESULT_PASSED_WITH_DEFECTS)
	orderedResultLabels = append(orderedResultLabels, RESULT_FAILED)
	orderedResultLabels = append(orderedResultLabels, RESULT_FAILED_WITH_DEFECTS)
	orderedResultLabels = append(orderedResultLabels, RESULT_ENVFAIL)

	//Build a list of standard labels to prevent duplication
	var standardResultLabels = make(map[string]struct{})
	for _, key := range orderedResultLabels {
		//'struct{}{}' allocates no storage. In Go 1.19 we can use just '{}' instead
		standardResultLabels[key] = struct{}{}
	}

	//Gathering custom labels
	var customResultLabels []string
	for _, run := range formattableTest {
		_, isStandardLabel := standardResultLabels[run.Result]
		if !isStandardLabel {
			customResultLabels = append(customResultLabels, run.Result)
		}
	}

	sort.Strings(customResultLabels)
	orderedResultLabels = append(orderedResultLabels, customResultLabels...)

	return orderedResultLabels
}

func FormattableTestFromGalasaApi(runs []galasaapi.Run, apiServerUrl string) []runsformatter.FormattableTest {
	var formattableTest []runsformatter.FormattableTest

	for _, run := range runs {
		//Get the data for each TestStructure in runs
		newFormattableTest := getTestStructureData(run, apiServerUrl)
		formattableTest = append(formattableTest, newFormattableTest)
	}

	orderedFormattableTest := orderFormattableTests(formattableTest)

	return orderedFormattableTest
}

func getTestStructureData(run galasaapi.Run, apiServerUrl string) runsformatter.FormattableTest {
	newFormattableTest := runsformatter.NewFormattableTest()

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

func FormattableTestFromTestRun(finishedMap map[string]*TestRun, lostMap map[string]*TestRun) []runsformatter.FormattableTest {
	var formattableTest []runsformatter.FormattableTest
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

	orderedFormattableTest := orderFormattableTests(formattableTest)

	return orderedFormattableTest
}

func getTestRunData(run TestRun, isLost bool) runsformatter.FormattableTest {
	newFormattableTest := runsformatter.NewFormattableTest()

	newFormattableTest.RunId = ""
	newFormattableTest.ApiServerUrl = ""

	newFormattableTest.Name = run.Name
	if run.GherkinUrl != "" {
		newFormattableTest.TestName = run.GherkinFeature
	} else {
		newFormattableTest.TestName = run.Stream + "/" + run.Bundle + "/" + run.Class
	}
	newFormattableTest.Status = run.Status
	newFormattableTest.Result = run.Result
	newFormattableTest.StartTimeUTC = ""
	newFormattableTest.EndTimeUTC = ""
	newFormattableTest.QueuedTimeUTC = run.QueuedTimeUTC
	newFormattableTest.Requestor = run.Requestor
	newFormattableTest.Bundle = run.Bundle
	newFormattableTest.Methods = nil
	newFormattableTest.Lost = isLost

	return newFormattableTest
}
