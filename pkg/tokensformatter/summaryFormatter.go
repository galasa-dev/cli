/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package tokensformatter

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

type TokenSummaryFormatter struct {
}

func NewTokenSummaryFormatter() TokenFormatter {
	return new(TokenSummaryFormatter)
}

func (*TokenSummaryFormatter) GetName() string {
	return SUMMARY_FORMATTER_NAME
}

func (*TokenSummaryFormatter) FormatTokens(authTokens []galasaapi.AuthToken) (string, error) {
	var result string = ""
	var err error = nil
	buff := strings.Builder{}
	totalTokens := len(authTokens)

	if totalTokens > 0 {
		var table [][]string

		var headers = []string{HEADER_TOKEN_ID, HEADER_TOKEN_DATE, HEADER_TOKEN_USER, HEADER_TOKEN_DESCRIPTION}

		table = append(table, headers)
		for _, token := range authTokens {
			var line []string
			id := *token.TokenId
			createdTime := *token.CreatedTime
			owner := *token.Owner.LoginId
			description := *token.Description

			line = append(line, id, createdTime, owner, description)
			table = append(table, line)
		}

		columnLengths := calculateMaxLengthOfEachColumn(table)
		writeFormattedTableToStringBuilder(table, &buff, columnLengths)

		buff.WriteString("\n")

	}
	buff.WriteString("Total:" + strconv.Itoa(totalTokens) + "\n")

	result = buff.String()
	return result, err
}
