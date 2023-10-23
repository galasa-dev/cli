/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package propertiesformatter

import "strings"

// -----------------------------------------------------
// Summary format.
 const (
	 RAW_FORMATTER_NAME = "raw"
 )
 
 type PropertyRawFormatter struct {
 }
 
 func NewPropertyRawFormatter() PropertyFormatter {
	 return new(PropertyRawFormatter)
 }
 
 func (*PropertyRawFormatter) GetName() string {
	 return RAW_FORMATTER_NAME
 }
 
 func (*PropertyRawFormatter) FormatProperties(cpsProperties []FormattableProperty) (string, error) {
	var result string = ""
	buff := strings.Builder{}
	var err error = nil
	
	for _, property := range cpsProperties {
		buff.WriteString(property.Namespace + "|" + 
		property.Name + "|" + 
		property.Value + "\n")
	}

	result = buff.String()
	return result, err
 }
 