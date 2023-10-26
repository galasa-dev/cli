/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package propertiesformatter

import (
	"strings"

	"github.com/galasa-dev/cli/pkg/galasaapi"
	"gopkg.in/yaml.v2"
)

// -----------------------------------------------------
// Yaml format.
const (
	YAML_FORMATTER_NAME = "yaml"
)

type PropertyYamlFormatter struct {
}

func NewPropertyYamlFormatter() PropertyFormatter {
	return new(PropertyYamlFormatter)
}

func (*PropertyYamlFormatter) GetName() string {
	return YAML_FORMATTER_NAME
}
 
 func (*PropertyYamlFormatter) FormatProperties(cpsProperties []galasaapi.CpsProperty) (string, error) {
	var result string = ""
	var err error = nil
	buff := strings.Builder{}
	totalProperties := len(cpsProperties)

	if totalProperties > 0 {
		buff.WriteString("apiVersion: galasa-dev/v1alpha1\n")
	}
	for index, property := range cpsProperties {
		propertyString := ""

		if index > 0 {
			propertyString += "---\n"
		}
		
		var yamlRepresentationBytes []byte
		yamlRepresentationBytes, err = yaml.Marshal(property)
		if err == nil {
			yamlStr := string(yamlRepresentationBytes)
			buff.WriteString(yamlStr)
		}

		buff.WriteString(propertyString)
	}


 
	result = buff.String()
	return result, err
}