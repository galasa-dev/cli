/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package secretsformatter

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

type SecretSummaryFormatter struct {
}

func NewSecretSummaryFormatter() SecretsFormatter {
    return new(SecretSummaryFormatter)
}

func (*SecretSummaryFormatter) GetName() string {
    return SUMMARY_FORMATTER_NAME
}

func (*SecretSummaryFormatter) FormatSecrets(secrets []galasaapi.GalasaSecret) (string, error) {
    var result string = ""
    var err error = nil
    buff := strings.Builder{}
    totalSecrets := len(secrets)

    if totalSecrets > 0 {
        var table [][]string

        var headers = []string{ HEADER_SECRET_NAME, HEADER_SECRET_TYPE }

        table = append(table, headers)
        for _, secret := range secrets {
            var line []string
            name := secret.Metadata.GetName()
            secretType := secret.Metadata.GetType()

            line = append(line, name, string(secretType))
            table = append(table, line)
        }

        columnLengths := calculateMaxLengthOfEachColumn(table)
        writeFormattedTableToStringBuilder(table, &buff, columnLengths)

        buff.WriteString("\n")

    }
    buff.WriteString("Total:" + strconv.Itoa(totalSecrets) + "\n")

    result = buff.String()
    return result, err
}
