/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package tokensformatter

import (
	"fmt"
	"strings"

	"github.com/galasa-dev/cli/pkg/galasaapi"
)

//Print in the following fashion:
// tokenid                   created(YYYY-MM-DD)  user     description
// 098234980123-1283182389   2023-12-03           mcobbett So I can access ecosystem1 from my laptop.
// 8218971d287s1-dhj32er2323 2024-03-03           mcobbett Automated build of example repo can change CPS properties
// 87a6sd87ahq2-2y8hqwdjj273 2023-08-04           savvas   CLI access from vscode
// Total:3

// -----------------------------------------------------
// TokensFormatter - implementations can take a collection of auth tokens results
// and turn them into a string for display to the user.
const (
	HEADER_TOKEN_ID            = "tokenid"
	HEADER_TOKEN_CREATION_TIME = "created(YYYY-MM-DD)"
	HEADER_TOKEN_USER          = "user"
	HEADER_TOKEN_DESCRIPTION   = "description"
)

type TokenFormatter interface {
	FormatTokens(tokenResults []galasaapi.AuthToken) (string, error)
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

// -----------------------------------------------------
// Functions for time formats and duration
func formatTimeToNearestDate(rawTime string) string {
	var formattedTimeString string
	if len(rawTime) < 19 {
		formattedTimeString = ""
	} else {
		formattedTimeString = rawTime[0:10]
	}
	return formattedTimeString
}
