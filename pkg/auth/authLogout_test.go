/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"errors"
	"testing"

	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestLogoutDeletesBearerTokenFile(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	bearerTokenFilePath := mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json"
	mockFileSystem.Create(bearerTokenFilePath)

	// When...
	err := Logout(mockFileSystem, mockGalasaHome)
	fileExists, _ := mockFileSystem.Exists(bearerTokenFilePath)

	// Then...
	assert.False(t, fileExists, "bearer token file should not exist")
	assert.Nil(t, err, "Should not return an error if the bearer token file has been successfully deleted")
}

func TestLogoutWithNoBearerTokenFileDoesNotThrowError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	bearerTokenFilePath := mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json"

	// When...
	err := Logout(mockFileSystem, mockGalasaHome)
	fileExists, _ := mockFileSystem.Exists(bearerTokenFilePath)

	// Then...
	assert.False(t, fileExists, "bearer token file should not exist")
	assert.Nil(t, err, "Should not return an error if the bearer token file does not already exist")
}

func TestLogoutWithFailingFileExistsReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewOverridableMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	mockFileSystem.VirtualFunction_Exists = func(targetFilePath string) (bool, error) {
		return false, errors.New("simulating a failed file exists check")
	}

	// When...
	err := Logout(mockFileSystem, mockGalasaHome)

	// Then...
	assert.NotNil(t, err, "Should return an error if the file exists check fails")
	assert.ErrorContains(t, err, "GAL1104E")
}