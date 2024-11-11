/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package usersformatter

import (
	"strconv"
	"strings"

	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
)

// -----------------------------------------------------
// Summary format.
const (
	SUMMARY_FORMATTER_NAME = "summary"
)

const CLIENT_WEB_UI = "web-ui"
const CLIENT_REST_API = "rest-api"

type UserSummaryFormatter struct {
}

func NewUserSummaryFormatter() UserFormatter {
	return new(UserSummaryFormatter)
}

func (*UserSummaryFormatter) GetName() string {
	return SUMMARY_FORMATTER_NAME
}

func (*UserSummaryFormatter) FormatUsers(users []galasaapi.UserData) (string, error) {
	var result string = ""
	var err error = nil
	buff := strings.Builder{}
	totalUsers := len(users)

	if totalUsers > 0 {
		var table [][]string

		var headers = []string{HEADER_USER_LOGIN_ID, HEADER_WEBUI_LAST_LOGIN, HEADER_RESTAPI_LAST_LOGIN}

		table = append(table, headers)
		for _, user := range users {

			var line []string

			clients := user.GetClients()

			loginId := user.GetLoginId()
			var webLastLogin, restLastLogin string

			for _, client := range clients {
				switch client.GetClientName() {
				case CLIENT_WEB_UI:
					webLastLogin = utils.FormatTimeToNearestDateTimeMins(client.GetLastLogin().String())
				case CLIENT_REST_API:
					restLastLogin = utils.FormatTimeToNearestDateTimeMins(client.GetLastLogin().String())
				}
			}

			line = append(line, loginId, webLastLogin, restLastLogin)
			table = append(table, line)
		}

		columnLengths := calculateMaxLengthOfEachColumn(table)
		writeFormattedTableToStringBuilder(table, &buff, columnLengths)

		buff.WriteString("\n")

	}
	buff.WriteString("Total:" + strconv.Itoa(totalUsers) + "\n")

	result = buff.String()
	return result, err
}
