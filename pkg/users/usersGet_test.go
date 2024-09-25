/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package users

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func MockUsersServlet(t *testing.T, w http.ResponseWriter, r *http.Request) {

	assert.NotEmpty(t, r.Header.Get("ClientApiVersion"))

	if !strings.Contains(r.URL.Path, "/users") {
		t.Errorf("Expected to request '/users', got: %s", r.URL.Path)
	}
	if r.Header.Get("Accept") != "application/json" {
		t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
	}
	w.Header().Set("Content-Type", "application/json")
}

func TestNullOrEmptyLoginIdReturnsError(t *testing.T) {

	//given
	mockConsole := utils.NewMockConsole()
	loginId := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MockUsersServlet(t, w, r)
	}))
	apiServerUrl := server.URL
	defer server.Close()

	apiClient := api.InitialiseAPI(apiServerUrl)

	//when
	err := GetUsers(loginId, apiClient, mockConsole)

	assert.ErrorContains(t, err, "GAL1155E")

}

func TestNotMeInputLoginIdReturnsError(t *testing.T) {

	//given
	mockConsole := utils.NewMockConsole()
	loginId := "notMe"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MockUsersServlet(t, w, r)
	}))
	apiServerUrl := server.URL
	defer server.Close()

	apiClient := api.InitialiseAPI(apiServerUrl)

	//when
	err := GetUsers(loginId, apiClient, mockConsole)

	assert.ErrorContains(t, err, "GAL1156E")
}

func TestMeInputLoginIdPrintsDetailsOnConsole(t *testing.T) {

	//given
	mockConsole := utils.NewMockConsole()
	loginId := "me"

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		jsonToReturn := `[{ "login_id": "myUserId" }]`
		writer.Write([]byte(jsonToReturn))

	}))
	apiServerUrl := server.URL
	defer server.Close()

	apiClient := api.InitialiseAPI(apiServerUrl)

	//when
	err := GetUsers(loginId, apiClient, mockConsole)

	assert.Nil(t, err)
	text := mockConsole.ReadText()
	assert.Equal(t, "id: myUserId\n", text)
}
