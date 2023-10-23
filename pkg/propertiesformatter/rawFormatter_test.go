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

func createFormattableTestForRaw (namespace string, name string, value string) FormattableProperty {
	FormattableProperty := FormattableProperty{
		Namespace : namespace,
		Name : name,
		Value : value,
	}
	return FormattableProperty
}
 
func TestRawFormatterNoDataReturnsNothing(t *testing.T) {
	// For...
	formatter := NewPropertyRawFormatter()
	// No data to format...
	formattableProperty := make([]FormattableProperty, 0)
 
	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperty)
 
	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput := ""
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestRawFormatterSingleDataReturnsCorrectly(t *testing.T) {
	// For..
	formatter := NewPropertyRawFormatter()
	// No data to format...
	formattableProperties := make([]FormattableProperty, 0)
	formattableProperty1 := createFormattableTestForRaw("namespace", "name1", "value1")
	formattableProperties = append(formattableProperties, formattableProperty1)
 
	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperties)

	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput := "namespace|name1|value1\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestRawFormatterMultipleDataSeperatesWithNewLine(t *testing.T) {
	// For..
	formatter := NewPropertyRawFormatter()
	// No data to format...
	formattableProperties := make([]FormattableProperty, 0)
	formattableProperty1 := createFormattableTestForRaw("namespace", "name1", "value1")
	formattableProperties = append(formattableProperties, formattableProperty1)
	formattableProperty2 := createFormattableTestForRaw("namespace", "name2", "value2")
	formattableProperties = append(formattableProperties, formattableProperty2)
 
	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperties)

	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput :=  "namespace|name1|value1\n"+
								"namespace|name2|value2\n"
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}