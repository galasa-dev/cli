/*
 * Copyright contributors to the Galasa project
 */
package auth

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestLogoutDeletesBearerTokenFile(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockConsole := utils.NewMockConsole()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	bearerTokenFilePath := mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json"
	mockFileSystem.Create(bearerTokenFilePath)

	// When...
	err := Logout(mockFileSystem, mockConsole, mockEnvironment, mockGalasaHome)
	fileExists, _ := mockFileSystem.Exists(bearerTokenFilePath)

	// Then...
	assert.False(t, fileExists, "bearer token file should not exist")
	assert.Nil(t, err, "Should not return an error if the bearer token file has been successfully deleted")
}

func TestLogoutWithNoBearerTokenFileDoesNotThrowError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockConsole := utils.NewMockConsole()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	bearerTokenFilePath := mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json"

	// When...
	err := Logout(mockFileSystem, mockConsole, mockEnvironment, mockGalasaHome)
	fileExists, _ := mockFileSystem.Exists(bearerTokenFilePath)

	// Then...
	assert.False(t, fileExists, "bearer token file should not exist")
	assert.Nil(t, err, "Should not return an error if the bearer token file does not already exist")
}