/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
*/
package propertiesformatter

import (
	"strings"
	"github.com/galasa-dev/cli/pkg/utils"
)

// -----------------------------------------------------
// Summary format.
const (
	YAML_FORMATTER_NAME = "summary"
)

type PropertyYamlFormatter struct {
}

func NewPropertyYamlFormatter() PropertyFormatter {
	return new(PropertyYamlFormatter)
}

func (*PropertyYamlFormatter) GetName() string {
	return YAML_FORMATTER_NAME
}
 
 func (*PropertyYamlFormatter) FormatProperties(cpsProperties []FormattableProperty) (string, error) {
	var result string = ""
	var err error = nil
	buff := strings.Builder{}
	totalProperties := len(cpsProperties)
	counter := 0

	if totalProperties > 0 {
		buff.WriteString("apiVersion: galasa-dev/v1alpha1\n")
	}
	for _, property := range cpsProperties {
		propertyString := ""

		propertyString += utils.ConvertToYaml("GalasaProperty", property.Namespace, property.Name, property.Value)
		if counter != totalProperties -1 {
			propertyString += "---\n"
		}
		counter ++
		buff.WriteString(propertyString)
	}


 
	result = buff.String()
	return result, err
}