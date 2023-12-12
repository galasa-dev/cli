/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package propertiesformatter

import (
	"strings"

	"github.com/galasa-dev/cli/pkg/galasaapi"
	"gopkg.in/yaml.v3"
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

func (*PropertyYamlFormatter) FormatProperties(cpsProperties []galasaapi.GalasaProperty) (string, error) {
	var result string = ""
	var err error = nil
	buff := strings.Builder{}
	//totalProperties := len(cpsProperties)

	for index, property := range cpsProperties {
		propertyString := ""

		if index > 0 {
			propertyString += "---\n"
		}

		var yamlRepresentationBytes []byte
		yamlRepresentationBytes, err = yaml.Marshal(property)
		if err == nil {
			yamlStr := string(yamlRepresentationBytes)
			yamlStr = strings.ReplaceAll(yamlStr, "apiversion", "apiVersion")
			propertyString += yamlStr
		}

		buff.WriteString(propertyString)
	}

	result = buff.String()
	return result, err
}

func (*PropertyYamlFormatter) FormatNamespaces(namespaces []galasaapi.Namespace) (string, error) {
	var result string = ""
	var err error = nil
	buff := strings.Builder{}
	totalProperties := len(namespaces)

	if totalProperties > 0 {
		buff.WriteString("apiVersion: galasa-dev/v1alpha1\n")
	}
	for index, namespace := range namespaces {
		propertyString := ""

		if index > 0 {
			propertyString += "---\n"
		}

		var yamlRepresentationBytes []byte
		yamlRepresentationBytes, err = yaml.Marshal(namespace)
		if err == nil {
			yamlStr := string(yamlRepresentationBytes)
			yamlStr = strings.ReplaceAll(yamlStr, "propertiesurl", "propertiesUrl")
			propertyString += yamlStr
		}

		buff.WriteString(propertyString)
	}

	result = buff.String()
	return result, err
}
