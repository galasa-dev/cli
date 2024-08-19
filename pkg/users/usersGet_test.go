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

// MockServlet
func NewUsersServletMock(t *testing.T) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MockUsersServlet(t, w, r)
	}))

	return server
}

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

	server := NewUsersServletMock(t)
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

	server := NewUsersServletMock(t)
	apiServerUrl := server.URL
	defer server.Close()

	apiClient := api.InitialiseAPI(apiServerUrl)

	//when
	err := GetUsers(loginId, apiClient, mockConsole)

	assert.ErrorContains(t, err, "GAL1156E")
}
