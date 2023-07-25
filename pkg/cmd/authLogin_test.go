/*
* Copyright contributors to the Galasa project
 */
package cmd

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"strings"
)

func NewAuthServletMock(t *testing.T) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/auth") {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"jwt\":\"blah\"}"))
		}

	}))

	return server
}

func TestLoginCreatesBearerTokenFile(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	bearerTokenFilePath := mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json"
	mockFileSystem.Create(bearerTokenFilePath)

	server := NewAuthServletMock(t)
	defer server.Close()

	apiServerUrl := server.URL

	// When ...
	err := Login(apiServerUrl, mockFileSystem, mockGalasaHome)
	fileExists, _ := mockFileSystem.Exists(bearerTokenFilePath)

	// Then...
	// Should have created a folder for the parent package.
	assert.True(t, fileExists, "Bearer token file should exist")
	assert.Nil(t, err, "Should not return an error if the bearer token file has been successfully created")
}
