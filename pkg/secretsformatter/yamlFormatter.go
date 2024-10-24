/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package secretsformatter

import (
	"strings"

	"github.com/galasa-dev/cli/pkg/galasaapi"
	"gopkg.in/yaml.v3"
)

const (
	YAML_FORMATTER_NAME = "yaml"
)

type SecretYamlFormatter struct {
}

func NewSecretYamlFormatter() SecretsFormatter {
	return new(SecretYamlFormatter)
}

func (*SecretYamlFormatter) GetName() string {
	return YAML_FORMATTER_NAME
}

func (*SecretYamlFormatter) FormatSecrets(secrets []galasaapi.GalasaSecret) (string, error) {
	var err error
	buff := strings.Builder{}

	for index, secret := range secrets {
		secretString := ""

		if index > 0 {
			secretString += "---\n"
		}

		var yamlRepresentationBytes []byte
		yamlRepresentationBytes, err = yaml.Marshal(secret)
		if err == nil {
			yamlStr := string(yamlRepresentationBytes)

			// The generated bean serialises in json as 'apiVersion' which is correct. In yaml it serialises as 'apiversion' (incorrect)
			// So this is a hack to correct that failure.
			// Note: This will corrupt any value string which also has 'apiversion' inside it !
			// TODO: The fix is to change the bean and add a 'yaml' annotation so it gets rendered correctly. Golang has yaml annotations, but does the generator support them ?
			yamlStr = strings.ReplaceAll(yamlStr, "apiversion", "apiVersion")
			secretString += yamlStr
		}

		buff.WriteString(secretString)
	}

	result := buff.String()
	return result, err
}
