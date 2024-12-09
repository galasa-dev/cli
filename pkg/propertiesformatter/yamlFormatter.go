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
	PropertyFormatter
}

func NewPropertyYamlFormatter() PropertyFormatter {
	return new(PropertyYamlFormatter)
}

func (*PropertyYamlFormatter) GetName() string {
	return YAML_FORMATTER_NAME
}

func (*PropertyYamlFormatter) FormatProperties(cpsProperties []galasaapi.GalasaProperty) (string, error) {
	var err error
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

			// The generated bean serialises in json as 'apiVersion' which is correct. In yaml it serialises as 'apiversion' (incorrect)
			// So this is a hack to correct that failure.
			// Note: This will corrupt any value string which also has 'apiversion' inside it !
			// TODO: The fix is to change the bean and add a 'yaml' annotation so it gets rendered correctly. Golang has yaml annotations, but does the generator support them ?
			yamlStr = strings.ReplaceAll(yamlStr, "apiversion", "apiVersion")
			propertyString += yamlStr
		}

		buff.WriteString(propertyString)
	}

	result := buff.String()
	return result, err
}
