/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"strings"
)

func NewAuthServletMock(t *testing.T, status int, mockResponse string) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		if strings.Contains(request.URL.Path, "/auth") {
			requestBody, err := io.ReadAll(request.Body)
			assert.Nil(t, err, "Error reading request body")

			requestBodyStr := string(requestBody)
			assert.Contains(t, requestBodyStr, "client_id")
			assert.Contains(t, requestBodyStr, "refresh_token")

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(status)
			writer.Write([]byte(mockResponse))
		}
	}))

	return server
}

func TestLoginWithNoGalasactlPropertiesFileReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	mockResponse := `{"jwt":"blah"}`
	server := NewAuthServletMock(t, 200, mockResponse)
	defer server.Close()

	apiServerUrl := server.URL

	// When...
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.NotNil(t, err, "Should return an error if the galasactl.properties file does not exist")
	assert.ErrorContains(t, err, "GAL1043E")
}

func TestLoginWithBadGalasactlPropertiesFileReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	galasactlPropertiesFilePath := mockGalasaHome.GetNativeFolderPath() + "/galasactl.properties"
	mockFileSystem.WriteTextFile(galasactlPropertiesFilePath, "here are some bad galasactl.properties contents!")

	mockResponse := `{"jwt":"blah"}`
	server := NewAuthServletMock(t, 200, mockResponse)
	defer server.Close()

	apiServerUrl := server.URL

	// When...
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.NotNil(t, err, "Should return an error if the galasactl.properties file does not contain valid YAML")
	assert.ErrorContains(t, err, "GAL1122E")
}

func TestLoginCreatesBearerTokenFileContainingJWT(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	galasactlPropertiesFilePath := mockGalasaHome.GetNativeFolderPath() + "/galasactl.properties"

	mockClientId := "dummyId"
	mockRefreshToken := "abcdefg"
	tokenPropertyValue := mockRefreshToken + TOKEN_SEPARATOR + mockClientId
	mockFileSystem.WriteTextFile(galasactlPropertiesFilePath, fmt.Sprintf("GALASA_TOKEN=%s", tokenPropertyValue))

	mockResponse := `{"jwt":"blah"}`
	server := NewAuthServletMock(t, 200, mockResponse)
	defer server.Close()

	apiServerUrl := server.URL

	// When...
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome, mockEnvironment)

	bearerTokenFilePath := mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json"
	bearerTokenFileExists, _ := mockFileSystem.Exists(bearerTokenFilePath)
	bearerTokenFileContents, _ := mockFileSystem.ReadTextFile(bearerTokenFilePath)

	// Then...
	assert.Nil(t, err, "Should not return an error if the bearer token file has been successfully created")
	assert.True(t, bearerTokenFileExists, "Bearer token file should exist")
	assert.Equal(t, mockResponse, bearerTokenFileContents)
}

func TestLoginWithFailedFileWriteReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewOverridableMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	galasactlPropertiesFilePath := mockGalasaHome.GetNativeFolderPath() + "/galasactl.properties"

	mockClientId := "dummyId"
	mockRefreshToken := "abcdefg"
	tokenPropertyValue := mockRefreshToken + TOKEN_SEPARATOR + mockClientId
	mockFileSystem.WriteTextFile(galasactlPropertiesFilePath, fmt.Sprintf("GALASA_TOKEN=%s", tokenPropertyValue))

	mockFileSystem.VirtualFunction_WriteTextFile = func(path string, contents string) error {
		return errors.New("simulating a failed write operation")
	}

	mockResponse := `{"jwt":"blah"}`
	server := NewAuthServletMock(t, 200, mockResponse)
	defer server.Close()

	apiServerUrl := server.URL

	// When...
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.NotNil(t, err, "Should return an error if writing the bearer token file fails")
	assert.ErrorContains(t, err, "GAL1042E")
}

func TestLoginWithFailedTokenRequestReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	galasactlPropertiesFilePath := mockGalasaHome.GetNativeFolderPath() + "/galasactl.properties"

	mockClientId := "dummyId"
	mockRefreshToken := "abcdefg"
	tokenPropertyValue := mockRefreshToken + TOKEN_SEPARATOR + mockClientId
	mockFileSystem.WriteTextFile(galasactlPropertiesFilePath, fmt.Sprintf("GALASA_TOKEN=%s", tokenPropertyValue))

	mockResponse := `{"error":"something went wrong!"}`
	server := NewAuthServletMock(t, 500, mockResponse)
	defer server.Close()

	apiServerUrl := server.URL

	// When...
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.NotNil(t, err, "Should return an error if the API request returns an error")
	assert.ErrorContains(t, err, "GAL1106E")
}

func TestLoginWithMissingAuthPropertyReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	galasactlPropertiesFilePath := mockGalasaHome.GetNativeFolderPath() + "/galasactl.properties"

	mockFileSystem.WriteTextFile(galasactlPropertiesFilePath, "unknown.value=blah")

	mockResponse := `{"jwt":"blah"}`
	server := NewAuthServletMock(t, 200, mockResponse)
	defer server.Close()

	apiServerUrl := server.URL

	// When...
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.NotNil(t, err, "Should return an error if the GALASA_ACCESS_TOKEN property is missing")
	assert.ErrorContains(t, err, "GAL1122E")
}

func TestGetAuthenticatedAPIClientWithBearerTokenFileReturnsClient(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")
	apiServerUrl := "http://dummy-url"

	mockCurrentTime := time.UnixMilli(0)
	mockTimeService := utils.NewOverridableMockTimeService(mockCurrentTime)

	// This is a dummy JWT that expires 1 hour after the Unix epoch
	mockJwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjM2MDB9._j3Fchdx5IIqgGrdEGWXHxdgVyoBEyoD2-IBvhlxF1s"

	bearerTokenFilePath := mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json"
	mockFileSystem.WriteTextFile(bearerTokenFilePath, fmt.Sprintf(`{"jwt":"%s"}`, mockJwt))

	// When...
	apiClient, err := GetAuthenticatedAPIClient(apiServerUrl, mockFileSystem, mockGalasaHome, mockTimeService, mockEnvironment)

	// Then...
	assert.Nil(t, err, "No error should have been thrown")
	assert.NotNil(t, apiClient, "API client should not be nil")
}

func TestGetAuthenticatedAPIClientWithMissingBearerTokenFileAttemptsLogin(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	mockCurrentTime := time.UnixMilli(0)
	mockTimeService := utils.NewOverridableMockTimeService(mockCurrentTime)

	galasactlPropertiesFilePath := mockGalasaHome.GetNativeFolderPath() + "/galasactl.properties"

	accessTokenValue := "abc"
	clientIdValue := "dummyId"
	tokenPropertyValue := accessTokenValue + TOKEN_SEPARATOR + clientIdValue

	mockFileSystem.WriteTextFile(galasactlPropertiesFilePath, fmt.Sprintf("GALASA_TOKEN=%s", tokenPropertyValue))

	// This is a dummy JWT that expires 1 hour after the Unix epoch
	mockJwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjM2MDB9._j3Fchdx5IIqgGrdEGWXHxdgVyoBEyoD2-IBvhlxF1s"
	mockResponse := fmt.Sprintf(`{"jwt":"%s"}`, mockJwt)

	server := NewAuthServletMock(t, 200, mockResponse)
	defer server.Close()

	apiServerUrl := server.URL

	// When...
	apiClient, err := GetAuthenticatedAPIClient(apiServerUrl, mockFileSystem, mockGalasaHome, mockTimeService, mockEnvironment)

	// Then...
	assert.Nil(t, err, "No error should have been thrown")
	assert.NotNil(t, apiClient, "API client should not be nil if the login was successful")
}

