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
