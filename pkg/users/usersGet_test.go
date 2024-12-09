/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func NewUsersServletMock(t *testing.T, state string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockUsersServlet(t, w, r, state)
	}))
	return server
}

func mockUsersServlet(t *testing.T, writer http.ResponseWriter, request *http.Request, state string) {
	writer.Header().Set("Content-Type", "application/json")
	var statusCode int
	var body string
	statusCode = 200
	if state == "populated" {
		body = `
[
    {
		"url": "http://localhost:8080/users/d2055afbc0ae6e513fa9b23c1a000d9f",
		"login-id": "test-user",
		"id": "d2055afbc0ae6e513fa9b23c1a000d9f",
		"clients": [
			{
				"last-login": "2024-10-28T14:54:49.546029Z",
				"client-name": "web-ui"
			}
		]
    },    
	{
		"url": "http://localhost:8080/users/d2055afbc0ae6e513fa9b23c1a000d9f",
		"login-id": "test-user2",
		"id": "d2055afbc0ae6e513fa9b23c1a000d9f",
		"clients": [
			{
				"last-login": "2024-10-28T14:54:49.546029Z",
				"client-name": "web-ui"
			}
		]
    },
	{
		"url": "http://localhost:8080/users/d2055afbc0ae6e513fa9b23c1a000d9f",
		"login-id": "test-user3",
		"id": "d2055afbc0ae6e513fa9b23c1a000d9f",
		"clients": [
			{
				"last-login": "2024-10-28T14:54:49.546029Z",
				"client-name": "web-ui"
			},
			{
				"last-login": "2024-10-28T15:32:49.546029Z",
				"client-name": "rest-api"
			}
		]
    }    
]
`
	} else if state == "empty" {
		body = `{
    []
}`
	} else if state == "missingLoginIdFlag" {
		statusCode = 400
		body = `{"error_code": 1155,"error_message": "GAL1155E: The id provided by the --login-id field cannot be an empty string."}`
	} else if state == "invalidLoginIdFlag" {
		statusCode = 400
		body = `{"error_code": 1157,"error_message": "GAL1157E: '%s' is not supported as a valid value. Valid value should not contain spaces. A value of 'admin' is valid but 'galasa admin' is not."}`
	} else if state == "populatedByLoginId" {
		body = `
			[
				{
					"url": "http://localhost:8080/users/d2055afbc0ae6e513fa9b23c1a000d9f",
					"login-id": "test-user",
					"id": "d2055afbc0ae6e513fa9b23c1a000d9f",
					"clients": [
						{
							"last-login": "2024-10-28T14:54:49.546029Z",
							"client-name": "web-ui"
						}
					]
    			}
			]
		`
	} else {
		statusCode = 500
		body = `{"error_code": 5000,"error_message": "GAL5000E: Error occured when trying to access the endpoint. Report the problem to your Galasa Ecosystem owner."}`
	}
	writer.WriteHeader(statusCode)
	writer.Write([]byte(body))
}

func TestMultipleNamespacesPathReturnsOk(t *testing.T) {
	//Given...
	serverState := "populated"
	server := NewUsersServletMock(t, serverState)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `login-id   web-last-login(UTC) rest-api-last-login(UTC)
test-user  2024-10-28 14:54    
test-user2 2024-10-28 14:54    
test-user3 2024-10-28 14:54    2024-10-28 15:32

Total:3
`

	//When
	err := GetUsers("", apiClient, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestMissingLoginIdFlagReturnsBadRequest(t *testing.T) {
	//Given...
	serverState := "missingLoginId"
	server := NewUsersServletMock(t, serverState)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `GAL1155E: The loginId provided by the --login-id field cannot be an empty string.`

	//When
	err := GetUsers("     ", apiClient, console)

	//Then
	assert.NotNil(t, err)
	assert.Equal(t, expectedOutput, err.Error())
}

func TestGetTokensByLoginIdReturnsOK(t *testing.T) {
	//Given...
	serverState := "populatedByLoginId"
	server := NewUsersServletMock(t, serverState)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `login-id  web-last-login(UTC) rest-api-last-login(UTC)
test-user 2024-10-28 14:54    

Total:1
`

	//When
	err := GetUsers("test-user", apiClient, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}
