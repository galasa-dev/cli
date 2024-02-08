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
	_, ok := structure.(JavaSubmitRunStructure)
	if ok {
        result = "java structure"
	} else {
        result = "gherkin structure"
    }
    return result
}
func TestCanDistinguishJavaStructWhenGivenBaseStruct(t *testing.T) {
	//Given...
	javaStruct := JavaSubmitRunStructure{}
	//When....
	structure := returnStructureType(javaStruct)
	//Then....
	assert.Equal(t, "java structure", structure)
}

func TestCanDistinguishCherkinStructWhenGivenBaseStruct(t *testing.T) {
	//Given...
	gherkinStruct := GherkinSubmitRunStructure{}
	//When....
	structure := returnStructureType(gherkinStruct)
	//Then....
	assert.Equal(t, "gherkin structure", structure)
}

func TestCanRetrieveJavaStructFieldsWhenGivenBaseStruct(t *testing.T){
	//Given...
	overridesMap := make(map[string]interface{})

	runStruct := SubmitRunStructure{
		requestor : "requestor",
		obrFromPortfolio : " obr",
		isTraceEnabled   : false,
		overrides : overridesMap,
	}

	javaStruct := JavaSubmitRunStructure{
		SubmitRunStructure: runStruct,
		groupName : "groupName",
		stream : "stream",
	}

	//When....

	//Then....
	
}





