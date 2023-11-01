/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package propertiesformatter

import (
	"strings"

	"github.com/galasa-dev/cli/pkg/galasaapi"
)

// -----------------------------------------------------
// Raw format.
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

func (*PropertyRawFormatter) FormatProperties(cpsProperties []galasaapi.CpsProperty) (string, error) {
	var result string = ""
	buff := strings.Builder{}
	var err error = nil

	for _, property := range cpsProperties {
		namespace, name := GetNameAndNamespace(*property.Name)
		buff.WriteString(namespace + "|" +
			name + "|" +
			*property.Value + "\n")
	}

	result = buff.String()
	return result, err
}

func (*PropertyRawFormatter) FormatNamespaces(namespaces []galasaapi.Namespace) (string, error) {
	var result string = ""
	buff := strings.Builder{}
	var err error = nil

	for _, namespace := range namespaces {
		buff.WriteString(*namespace.Name + "|" +
			*namespace.Type + "\n")
	}

	result = buff.String()
	return result, err
}
 