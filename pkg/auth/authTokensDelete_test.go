/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func newDeleteTokensServletMock(t *testing.T) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockDeleteTokensServlet(t, w, r)
	}))

	return server
}

func mockDeleteTokensServlet(t *testing.T, w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.URL.Path, "/auth/tokens") {
		t.Errorf("Expected to request '/auth/tokens', got: %s", r.URL.Path)
	}
	if r.Header.Get("Accept") != "application/json" {
		t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
	}
	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var jsonRespBody string
	splitUrl := strings.Split(r.URL.Path, "/")
	tokenId := splitUrl[3]

	//auth/tokens/{tokenId}
	if tokenId == "unauthorizedToken" {
		statusCode = 401
		jsonRespBody = `{
			"error_code": 5401,
			"error_message": "GAL5401E: Unauthorized. Please ensure you have provided a valid 'Authorization' header with a valid bearer token and try again."
			}`
	} else if tokenId == "serverErrToken" {
		statusCode = 500
		jsonRespBody = `{
				"error_code": 5000,
				"error_message": "GAL5000E: Error occured when trying to access the endpoint. Report the problem to your Galasa Ecosystem owner."
			}`
	} else if tokenId == "notFoundToken" {
		statusCode = 404
		jsonRespBody = ""
	} else if tokenId == "validToken" {
		statusCode = 200
		jsonRespBody = ""
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(jsonRespBody))
}

func TestValidTokenIdFormatReturnsOk(t *testing.T) {
	//Given
	validTokenIdFormat := "siuIOHUDH98Y73GioudusIUH"

	//When
	err := validateTokenId(validTokenIdFormat)

	//then
	assert.Nil(t, err)
}
func TestInvalidTokenIdWithSpecialCharactersReturnsError(t *testing.T) {
	//Given
	invalidTokenIdFormat := "siuIOHUDH98Y:73GioudusIUH"

	//When
	err := validateTokenId(invalidTokenIdFormat)

	//then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The provided token id, '"+invalidTokenIdFormat+"', provided does not match formatting requirements. The token id must be an alphanumeric string only.")
	assert.Contains(t, err.Error(), "GAL1156E")
}

func TestDeleteValidTokenReturnsOk(t *testing.T) {
	//Given...
	tokenId := "validToken"

	server := newDeleteTokensServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()

	//When
	err := DeleteToken(tokenId, apiClient, console)

	//Then
	assert.Nil(t, err)
}

func TestDeleteValidNotFoundTokenReturnsOk(t *testing.T) {
	//Given...
	tokenId := "notFoundToken"

	server := newDeleteTokensServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()

	//When
	err := DeleteToken(tokenId, apiClient, console)

	//Then
	assert.Nil(t, err)
}

func TestDeleteUnauthorizedTokenReturnsError(t *testing.T) {
	//Given...
	tokenId := "unauthorizedToken"

	server := newDeleteTokensServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()

	//When
	err := DeleteToken(tokenId, apiClient, console)

	//Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Unauthorized. Please ensure you have provided a valid 'Authorization' header with a valid bearer token and try again.")
	assert.Contains(t, err.Error(), "GAL5401E")
}

func TestDeleteTokenServerErrorTokenReturnedReturnsError(t *testing.T) {
	//Given...
	tokenId := "serverErrToken"

	server := newDeleteTokensServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()

	//When
	err := DeleteToken(tokenId, apiClient, console)

	//Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Error occured when trying to access the endpoint. Report the problem to your Galasa Ecosystem owner.")
	assert.Contains(t, err.Error(), "GAL5000E")
}
