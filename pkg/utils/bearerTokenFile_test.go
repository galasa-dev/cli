/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/galasa-dev/cli/pkg/files"
	"github.com/stretchr/testify/assert"
)

func TestWriteBearerTokenJsonFileWritesJwtJsonToFile(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := NewMockEnv()
	mockGalasaHome, _ := NewGalasaHome(mockFileSystem, mockEnvironment, "")

	token := "blah"

	// When...
	err := WriteBearerTokenJsonFile(mockFileSystem, mockGalasaHome, token)

	// Then...
	assert.Nil(t, err, "Should not return an error when writing a JWT to the bearer-token.json file in an existing galasa home directory")

	bearerTokenJson, _ := mockFileSystem.ReadTextFile(mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json")

	assert.Contains(t, bearerTokenJson, token)
}

func TestWriteBearerTokenJsonWithFailingWriteOperationReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewOverridableMockFileSystem()
	mockEnvironment := NewMockEnv()
	mockGalasaHome, _ := NewGalasaHome(mockFileSystem, mockEnvironment, "")

	mockFileSystem.VirtualFunction_WriteTextFile = func(path string, contents string) error {
		return errors.New("simulating a failed write operation")
	}

	token := "blah"

	// When...
	err := WriteBearerTokenJsonFile(mockFileSystem, mockGalasaHome, token)

	// Then...
	assert.NotNil(t, err, "Should return an error when writing the bearer-token.json file fails")
	assert.ErrorContains(t, err, "GAL1042E")
}

func TestGetBearerTokenFromTokenJsonFileReturnsBearerToken(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := NewMockEnv()
	mockGalasaHome, _ := NewGalasaHome(mockFileSystem, mockEnvironment, "")

	// This is a dummy JWT that expires 1 hour after the Unix epoch
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjM2MDB9._j3Fchdx5IIqgGrdEGWXHxdgVyoBEyoD2-IBvhlxF1s"
	mockCurrentTime := time.UnixMilli(0)
	mockTimeService := NewOverridableMockTimeService(mockCurrentTime)

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath()+"/bearer-token.json",
		fmt.Sprintf(`{"jwt":"%s"}`, expectedToken))

	// When...
	bearerToken, err := GetBearerTokenFromTokenJsonFile(mockFileSystem, mockGalasaHome, mockTimeService)

	// Then...
	assert.Nil(t, err, "Should not return an error when a valid bearer token exists and is valid in bearer-token.json")
	assert.Equal(t, expectedToken, bearerToken)
}

func TestGetBearerTokenFromTokenJsonFileWithMissingTokenFileReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := NewMockEnv()
	mockGalasaHome, _ := NewGalasaHome(mockFileSystem, mockEnvironment, "")
	mockTimeService := NewMockTimeService()

	// When...
	_, err := GetBearerTokenFromTokenJsonFile(mockFileSystem, mockGalasaHome, mockTimeService)

	// Then...
	assert.NotNil(t, err, "Should return an error when bearer token file does not exist")
	assert.ErrorContains(t, err, "GAL1107E")
}

func TestGetBearerTokenFromTokenJsonFileWithBadContentsReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := NewMockEnv()
	mockGalasaHome, _ := NewGalasaHome(mockFileSystem, mockEnvironment, "")
	mockTimeService := NewMockTimeService()

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath()+"/bearer-token.json",
		"notabearertoken")

	// When...
	_, err := GetBearerTokenFromTokenJsonFile(mockFileSystem, mockGalasaHome, mockTimeService)

	// Then...
	assert.NotNil(t, err, "Should return an error when the bearer token file exists but doesn't contain valid JSON")
	assert.ErrorContains(t, err, "GAL1107E")
}
