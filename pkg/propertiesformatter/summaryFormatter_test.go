/*
* Copyright contributors to the Galasa project
*
* SPDX-License-Identifier: EPL-2.0
 */
package propertiesformatter

import (
	"testing"

	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/stretchr/testify/assert"
)

func TestSummaryFormatterNoDataReturnsTotalCountAllZeros(t *testing.T) {

	formatter := NewPropertySummaryFormatter()
	// No data to format...
	formattableProperty := make([]galasaapi.CpsProperty, 0)

	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperty)

	assert.Nil(t, err)
	expectedFormattedOutput := "Total:0\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterSingleDataReturnsCorrectly(t *testing.T) {
	// For..
	formatter := NewPropertySummaryFormatter()
	// No data to format...
	formattableProperties := make([]galasaapi.CpsProperty, 0)
	property1 := galasaapi.NewCpsProperty()
	property1.SetName("testNamespace.name1")
	property1.SetValue("value1")
	formattableProperties = append(formattableProperties, *property1)
 
	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperties)

	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput := 
`namespace     name  value
testNamespace name1 value1

Total:1
`
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterMultipleDataSeperatesWithNewLine(t *testing.T) {
	// For..
	formatter := NewPropertySummaryFormatter()
	// No data to format...
	formattableProperties := make([]galasaapi.CpsProperty, 0)
	property1 := galasaapi.NewCpsProperty()
	property1.SetName("testNamespace.name1")
	property1.SetValue("value1")
	formattableProperties = append(formattableProperties, *property1)
	property2 := galasaapi.NewCpsProperty()
	property2.SetName("testNamespace.name2")
	property2.SetValue("value2")
	formattableProperties = append(formattableProperties, *property2)
 
	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperties)

	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput := 
`namespace     name  value
testNamespace name1 value1
testNamespace name2 value2

Total:2
`
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}