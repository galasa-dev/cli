/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package secretsformatter

import (
    "fmt"
    "strings"

    "github.com/galasa-dev/cli/pkg/galasaapi"
)

// Displays secrets in the following format:
// name
// SYSTEM1
// MY_ZOS_SECRET
// ANOTHER-SECRET
// Total:3

// -----------------------------------------------------
// SecretsFormatter - implementations can take a collection of secrets
// and turn them into a string for display to the user.
const (
    HEADER_SECRET_NAME = "name"
    HEADER_SECRET_TYPE = "type"
    HEADER_SECRET_DESCRIPTION = "description"
    HEADER_LAST_UPDATED_TIME = "last-updated(UTC)"
	HEADER_LAST_UPDATED_BY = "last-updated-by"
)

type SecretsFormatter interface {
    FormatSecrets(secrets []galasaapi.GalasaSecret) (string, error)
    GetName() string
}

// -----------------------------------------------------
// Functions for tables
func calculateMaxLengthOfEachColumn(table [][]string) []int {
    columnLengths := make([]int, len(table[0]))
    for _, row := range table {
        for i, val := range row {
            if len(val) > columnLengths[i] {
                columnLengths[i] = len(val)
            }
        }
    }
    return columnLengths
}

func writeFormattedTableToStringBuilder(table [][]string, buff *strings.Builder, columnLengths []int) {
    for _, row := range table {
        for column, val := range row {

            // For every column except the last one, add spacing.
            if column < len(row)-1 {
                // %-*s : variable space-padding length, padding is on the right.
                buff.WriteString(fmt.Sprintf("%-*s", columnLengths[column], val))
                buff.WriteString(" ")
            } else {
                buff.WriteString(val)
            }
        }
        buff.WriteString("\n")
    }
}
