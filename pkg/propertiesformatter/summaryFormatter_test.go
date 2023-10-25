/*
* Copyright contributors to the Galasa project
*
* SPDX-License-Identifier: EPL-2.0
 */
package propertiesformatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateFormattableTest (namespace string, name string, value string) FormattableProperty {
	FormattableProperty := FormattableProperty{
		Namespace : namespace,
		Name : name,
		Value : value,
	}
	return FormattableProperty
}

func TestSummaryFormatterNoDataReturnsTotalCountAllZeros(t *testing.T) {

	formatter := NewPropertySummaryFormatter()
	// No data to format...
	formattableProperty := make([]FormattableProperty, 0)

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
	formattableProperties := make([]FormattableProperty, 0)
	formattableProperty1 := CreateFormattableTest("namespace", "name1", "value1")
	formattableProperties = append(formattableProperties, formattableProperty1)
 
	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperties)

	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput := 
`namespace name  value
namespace name1 value1

Total:1
`
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSummaryFormatterMultipleDataSeperatesWithNewLine(t *testing.T) {
	// For..
	formatter := NewPropertySummaryFormatter()
	// No data to format...
	formattableProperties := make([]FormattableProperty, 0)
	formattableProperty1 := CreateFormattableTest("namespace", "name1", "value1")
	formattableProperties = append(formattableProperties, formattableProperty1)
	formattableProperty2 := CreateFormattableTest("namespace", "name2", "value2")
	formattableProperties = append(formattableProperties, formattableProperty2)
 
	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperties)

	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput :=  
`namespace name  value
namespace name1 value1
namespace name2 value2

Total:2
`
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}