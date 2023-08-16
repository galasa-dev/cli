/*
* Copyright contributors to the Galasa project
 */
package auth

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/utils"
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
			assert.Contains(t, requestBodyStr, "secret")
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
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome)

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
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome)

	// Then...
	assert.NotNil(t, err, "Should return an error if the galasactl.properties file does not contain valid YAML")
	assert.ErrorContains(t, err, "GAL1089E")
}

func TestLoginCreatesBearerTokenFileContainingJWT(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	galasactlPropertiesFilePath := mockGalasaHome.GetNativeFolderPath() + "/galasactl.properties"

	mockClientId := "dummyId"
	mockSecret := "shhhh"
	mockRefreshToken := "abcdefg"
	mockFileSystem.WriteTextFile(galasactlPropertiesFilePath, fmt.Sprintf(
		"auth.client.id=%s\n"+
		"auth.secret=%s\n"+
		"auth.access.token=%s", mockClientId, mockSecret, mockRefreshToken))

	mockResponse := `{"jwt":"blah"}`
	server := NewAuthServletMock(t, 200, mockResponse)
	defer server.Close()

	apiServerUrl := server.URL

	// When...
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome)

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
	mockSecret := "shhhh"
	mockRefreshToken := "abcdefg"
	mockFileSystem.WriteTextFile(galasactlPropertiesFilePath, fmt.Sprintf(
		"auth.client.id=%s\n"+
		"auth.secret=%s\n"+
		"auth.access.token=%s", mockClientId, mockSecret, mockRefreshToken))

	mockFileSystem.VirtualFunction_WriteTextFile = func(path string, contents string) error {
		return errors.New("simulating a failed write operation")
	}

	mockResponse := `{"jwt":"blah"}`
	server := NewAuthServletMock(t, 200, mockResponse)
	defer server.Close()

	apiServerUrl := server.URL

	// When...
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome)

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
	mockSecret := "shhhh"
	mockRefreshToken := "abcdefg"


	mockFileSystem.WriteTextFile(galasactlPropertiesFilePath, fmt.Sprintf(
		"auth.client.id=%s\n"+
		"auth.secret=%s\n"+
		"auth.access.token=%s", mockClientId, mockSecret, mockRefreshToken))

	mockResponse := `{"error":"something went wrong!"}`
	server := NewAuthServletMock(t, 500, mockResponse)
	defer server.Close()

	apiServerUrl := server.URL

	// When...
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome)

	// Then...
	assert.NotNil(t, err, "Should return an error if the API request returns an error")
	assert.ErrorContains(t, err, "GAL1090E")
}

func TestLoginWithMissingAuthPropertyReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	galasactlPropertiesFilePath := mockGalasaHome.GetNativeFolderPath() + "/galasactl.properties"

	mockClientId := "dummyId"
	mockSecret := "shhhh"
	mockFileSystem.WriteTextFile(galasactlPropertiesFilePath, fmt.Sprintf(
		"auth.client.id=%s\n"+
		"auth.secret=%s\n", mockClientId, mockSecret))

	mockResponse := `{"jwt":"blah"}`
	server := NewAuthServletMock(t, 200, mockResponse)
	defer server.Close()

	apiServerUrl := server.URL

	// When...
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome)

	// Then...
	assert.NotNil(t, err, "Should return an error if the auth.access.token property is missing")
	assert.ErrorContains(t, err, "GAL1089E")
}