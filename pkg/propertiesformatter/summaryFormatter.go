/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package propertiesformatter

import (
	"strconv"
	"strings"

	"github.com/galasa-dev/cli/pkg/galasaapi"
)

// -----------------------------------------------------
// Summary format.
const (
	SUMMARY_FORMATTER_NAME = "summary"
)

type PropertySummaryFormatter struct {
}

func NewPropertySummaryFormatter() PropertyFormatter {
	return new(PropertySummaryFormatter)
}

func (*PropertySummaryFormatter) GetName() string {
	return SUMMARY_FORMATTER_NAME
}

func (*PropertySummaryFormatter) FormatProperties(cpsProperties []galasaapi.CpsProperty) (string, error) {
	var result string = ""
	var err error = nil
	buff := strings.Builder{}
	totalProperties := len(cpsProperties)

	if totalProperties > 0 {
		var table [][]string

		var headers = []string{HEADER_PROPERTY_NAMESPACE, HEADER_PROPERTY_NAME, HEADER_PROPERTY_VALUE}

		table = append(table, headers)
		for _, property := range cpsProperties {
			var line []string
			namespace, name := GetNameAndNamespace(*property.Name)

			line = append(line, namespace)
			line = append(line, name, *property.Value)
			table = append(table, line)
		}

		columnLengths := calculateMaxLengthOfEachColumn(table)
		writeFormattedTableToStringBuilder(table, &buff, columnLengths)

		buff.WriteString("\n")

	}
	buff.WriteString("Total:" + strconv.Itoa(totalProperties) + "\n")

	result = buff.String()
	return result, err
}
