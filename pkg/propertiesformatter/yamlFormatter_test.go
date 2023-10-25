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

func TestYamlFormatterNoDataReturnsBlankString(t *testing.T) {
 
	formatter := NewPropertyYamlFormatter()
	// No data to format...
	formattableProperty := make([]FormattableProperty, 0)
 
	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperty)
 
	assert.Nil(t, err)
	expectedFormattedOutput := ""
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestYamlFormatterSingleDataReturnsCorrectly(t *testing.T) {
	// For..
	formatter := NewPropertyYamlFormatter()
	// No data to format...
	formattableProperties := make([]FormattableProperty, 0)
	formattableProperty1 := CreateFormattableTest("namespace", "name1", "value1")
	formattableProperties = append(formattableProperties, formattableProperty1)
 
	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperties)

	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput := `apiVersion: galasa-dev/v1alpha1
Kind: GalasaProperty
metadata:
	cpsnamespace: namespace
	name: name1
data:
	value: value1
`
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestYamlFormatterMultipleDataSeperatesWithNewLine(t *testing.T) {
	// For..
	formatter := NewPropertyYamlFormatter()
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
	expectedFormattedOutput := `apiVersion: galasa-dev/v1alpha1
Kind: GalasaProperty
metadata:
	cpsnamespace: namespace
	name: name1
data:
	value: value1
---
Kind: GalasaProperty
metadata:
	cpsnamespace: namespace
	name: name2
data:
	value: value2
`
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}