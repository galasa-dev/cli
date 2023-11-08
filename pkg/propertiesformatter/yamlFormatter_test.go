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

// PROPERTIES
func TestPropertiesYamlFormatterNoDataReturnsBlankString(t *testing.T) {

	formatter := NewPropertyYamlFormatter()
	// No data to format...
	formattableProperty := make([]galasaapi.CpsProperty, 0)

	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperty)

	assert.Nil(t, err)
	expectedFormattedOutput := ""
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestPropertiesYamlFormatterSingleDataReturnsCorrectly(t *testing.T) {
	// For..
	formatter := NewPropertyYamlFormatter()
	// No data to format...
	formattableProperties := make([]galasaapi.CpsProperty, 0)
	property1 := galasaapi.NewCpsProperty()
	property1.SetName("namespace.name1")
	property1.SetValue("value1")
	formattableProperties = append(formattableProperties, *property1)

	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperties)

	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput := `apiVersion: galasa-dev/v1alpha1
name: namespace.name1
value: value1
`
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestPropertiesYamlFormatterMultipleDataSeperatesWithNewLine(t *testing.T) {
	// For..
	formatter := NewPropertyYamlFormatter()
	// No data to format...
	formattableProperties := make([]galasaapi.CpsProperty, 0)
	property1 := galasaapi.NewCpsProperty()
	property1.SetName("namespace.name1")
	property1.SetValue("value1")
	formattableProperties = append(formattableProperties, *property1)
	property2 := galasaapi.NewCpsProperty()
	property2.SetName("namespace.name2")
	property2.SetValue("value2")
	formattableProperties = append(formattableProperties, *property2)

	// When...
	actualFormattedOutput, err := formatter.FormatProperties(formattableProperties)

	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput := `apiVersion: galasa-dev/v1alpha1
name: namespace.name1
value: value1
---
name: namespace.name2
value: value2
`
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

// NAMESPACES
func TestNamespacesYamlFormatterNoDataReturnsBlankString(t *testing.T) {

	formatter := NewPropertyYamlFormatter()
	// No data to format...
	namespace := make([]galasaapi.Namespace, 0)

	// When...
	actualFormattedOutput, err := formatter.FormatNamespaces(namespace)

	assert.Nil(t, err)
	expectedFormattedOutput := ""
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestNamespacesYamlFormatterSingleDataReturnsCorrectly(t *testing.T) {
	// For..
	formatter := NewPropertyYamlFormatter()

	namespaces := make([]galasaapi.Namespace, 0)
	namespace1 := CreateNamespace("framework", "normal", "")
	namespaces = append(namespaces, *namespace1)

	// When...
	actualFormattedOutput, err := formatter.FormatNamespaces(namespaces)

	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput := `apiVersion: galasa-dev/v1alpha1
name: framework
propertiesUrl: null
type: normal
`
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestNamespacesYamlFormatterMultipleDataSeperatesWithNewLine(t *testing.T) {
	// For..
	formatter := NewPropertyYamlFormatter()

	namespaces := make([]galasaapi.Namespace, 0)
	namespace1 := CreateNamespace("framework", "normal", "cps/namespaces/normal")
	namespace2 := CreateNamespace("secure", "secure", "")
	namespaces = append(namespaces, *namespace1, *namespace2)

	// When...
	actualFormattedOutput, err := formatter.FormatNamespaces(namespaces)

	// Then...
	assert.Nil(t, err)
	expectedFormattedOutput := `apiVersion: galasa-dev/v1alpha1
name: framework
propertiesUrl: cps/namespaces/normal
type: normal
---
name: secure
propertiesUrl: null
type: secure
`
	assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}
