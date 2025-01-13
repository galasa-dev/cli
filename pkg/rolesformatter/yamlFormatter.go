/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package rolesformatter

import (
	"strings"

	"github.com/galasa-dev/cli/pkg/galasaapi"
	"gopkg.in/yaml.v3"
)

const (
	YAML_FORMATTER_NAME = "yaml"
)

type RolesYamlFormatter struct {
}

func NewRolesYamlFormatter() RolesFormatter {
	return new(RolesYamlFormatter)
}

func (*RolesYamlFormatter) GetName() string {
	return YAML_FORMATTER_NAME
}

func (*RolesYamlFormatter) FormatRoles(roles []galasaapi.RBACRole) (string, error) {
	var err error
	buff := strings.Builder{}

	for index, role := range roles {
		content := ""

		if index > 0 {
			content += "---\n"
		}

		var yamlRepresentationBytes []byte
		yamlRepresentationBytes, err = yaml.Marshal(role)
		if err == nil {
			yamlStr := string(yamlRepresentationBytes)
			content += yamlStr
		}

		buff.WriteString(content)
	}

	result := buff.String()
	return result, err
}
