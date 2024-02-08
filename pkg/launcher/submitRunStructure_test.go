/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package launcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func returnStructureType(structure interface{}) string {
	var result string

	_, isJavaStruct := structure.(JavaSubmitRunStructure)
	if isJavaStruct {
		result = "java structure"
	} else {
		result = "gherkin structure"
	}

	return result
}

func createSubmitRunStructure() SubmitRunStructure {

	overridesMap := make(map[string]interface{})
	overridesMap["map1"] = "string"

	runStruct := SubmitRunStructure{
		requestor:        "requestor",
		obrFromPortfolio: " obr",
		isTraceEnabled:   false,
		overrides:        overridesMap,
	}

	return runStruct
}

func TestCanDistinguishJavaStructWhenFuncOnlyAcceptBaseStruct(t *testing.T) {
	//Given...
	javaStruct := JavaSubmitRunStructure{}
	//When....
	structure := returnStructureType(javaStruct)
	//Then....
	assert.Equal(t, "java structure", structure)
}

func TestCanDistinguishCherkinStructWhenFuncOnlyAcceptBaseStruct(t *testing.T) {
	//Given...
	gherkinStruct := GherkinSubmitRunStructure{}
	//When....
	structure := returnStructureType(gherkinStruct)
	//Then....
	assert.Equal(t, "gherkin structure", structure)
}

func TestCanRetrieveJavaStructFields(t *testing.T) {
	//Given...
	runStruct := createSubmitRunStructure()

	javaStruct := JavaSubmitRunStructure{
		SubmitRunStructure: runStruct,
		groupName:          "groupName",
		stream:             "stream",
	}

	//When....
	structure := returnStructureType(javaStruct)

	//Then....
	assert.Equal(t, "java structure", structure)
	assert.Equal(t, "requestor", javaStruct.SubmitRunStructure.requestor)
	assert.Equal(t, false, javaStruct.SubmitRunStructure.isTraceEnabled)
	assert.Equal(t, "string", javaStruct.SubmitRunStructure.overrides["map1"])
	assert.Equal(t, "stream", javaStruct.stream)
	assert.Equal(t, "", javaStruct.className)
	assert.Equal(t, "", javaStruct.requestType)
}

func TestCanRetrieveGherkinStructFields(t *testing.T) {
	//Given...
	runStruct := createSubmitRunStructure()

	gherkinStruct := GherkinSubmitRunStructure{
		SubmitRunStructure: runStruct,
		gherkinFile:        "gherkin_file.feature",
	}

	//When....
	structure := returnStructureType(gherkinStruct)

	//Then....
	assert.Equal(t, "gherkin structure", structure)
	assert.Equal(t, "requestor", gherkinStruct.SubmitRunStructure.requestor)
	assert.Equal(t, "string", gherkinStruct.SubmitRunStructure.overrides["map1"])
	assert.Equal(t, "gherkin_file.feature", gherkinStruct.gherkinFile)

}
