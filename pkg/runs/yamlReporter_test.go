/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"log"
	"strings"
	"testing"

	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestYamlReportWorks(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()

	finishedRuns := TestRun{
		Name:      "myTestRun",
		Bundle:    "myBundle",
		Class:     "com.myco.MyClass",
		Stream:    "myStream",
		Status:    "myStatus",
		Result:    "PASSED",
		Overrides: make(map[string]string, 1),
		Tests:     []TestMethod{TestMethod{Method: "method1", Result: "passed"}, TestMethod{Method: "method2", Result: "passed"}}}

	finishedRunsMap := make(map[string]*TestRun, 1)
	finishedRunsMap["myTestRun"] = &finishedRuns

	// When...
	err := ReportYaml(
		mockFileSystem,
		"myReportYamlFilename",
		finishedRunsMap, nil)

	// Then...
	if err != nil {
		assert.Fail(t, "ReportYaml failed when it should have passed. "+err.Error())
	}

	isExists, err := mockFileSystem.Exists("myReportYamlFilename")
	if err != nil {
		assert.Fail(t, "yaml report does not exist in correct place. "+err.Error())
	}
	assert.True(t, isExists, "Yaml report does not exist in the correct place.")

	// We expect a report like this:
	expectedReport := `
	tests:
	- name: myTestRun
	  bundle: myBundle
	  class: com.myco.MyClass
	  stream: myStream
	  status: myStatus
	  result: PASSED
	  overrides: {}
	  tests:
	  - name: method1
		result: passed
	  - name: method2
		result: passed`

	actualContents, err := mockFileSystem.ReadTextFile("myReportYamlFilename")
	if err != nil {
		assert.Fail(t, "Could not read the yaml file. "+err.Error())
	}

	log.Print(actualContents)

	expected1 := strings.ReplaceAll(expectedReport, "\n", "")
	expected2 := strings.ReplaceAll(expected1, "\t", "")
	expected3 := strings.ReplaceAll(expected2, " ", "")

	actual1 := strings.ReplaceAll(actualContents, "\n", "")
	actual2 := strings.ReplaceAll(actual1, "\t", "")
	actual3 := strings.ReplaceAll(actual2, " ", "")

	assert.EqualValues(t, expected3, actual3)
}
