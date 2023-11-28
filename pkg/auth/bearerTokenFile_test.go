/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestWriteBearerTokenJsonFileWritesJwtJsonToFile(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	jwtJsonToWrite := `{"jwt":"blah"}`

	// When...
	err := WriteBearerTokenJsonFile(mockFileSystem, mockGalasaHome, jwtJsonToWrite)

	// Then...
	assert.Nil(t, err, "Should not return an error when writing a JWT to the bearer-token.json file in an existing galasa home directory")

	bearerTokenJson, _ := mockFileSystem.ReadTextFile(mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json")
	assert.Equal(t, jwtJsonToWrite, bearerTokenJson)
}

func TestWriteBearerTokenJsonWithFailingWriteOperationReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewOverridableMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	mockFileSystem.VirtualFunction_WriteTextFile = func(path string, contents string) error {
		return errors.New("simulating a failed write operation")
	}

	jwtJsonToWrite := `{"jwt":"blah"}`

	// When...
	err := WriteBearerTokenJsonFile(mockFileSystem, mockGalasaHome, jwtJsonToWrite)

	// Then...
	assert.NotNil(t, err, "Should return an error when writing the bearer-token.json file fails")
	assert.ErrorContains(t, err, "GAL1042E")
}

func TestGetBearerTokenFromTokenJsonFileReturnsBearerToken(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	// This is a dummy JWT that expires 1 hour after the Unix epoch
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjM2MDB9._j3Fchdx5IIqgGrdEGWXHxdgVyoBEyoD2-IBvhlxF1s"
	mockCurrentTime := time.UnixMilli(0)
	mockTimeService := utils.NewOverridableMockTimeService(mockCurrentTime)

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json",
		fmt.Sprintf(`{"jwt":"%s"}`, expectedToken))

	// When...
	bearerToken, err := GetBearerTokenFromTokenJsonFile(mockFileSystem, mockGalasaHome, mockTimeService)

	// Then...
	assert.Nil(t, err, "Should not return an error when a valid bearer token exists and is valid in bearer-token.json")
	assert.Equal(t, expectedToken, bearerToken)
}

func TestGetBearerTokenFromTokenJsonFileWithExpiredTokenReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	// This is a dummy JWT that expires 1 second after the Unix epoch
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjF9.2H0EJnt58ApysedXcvNUAy6FhgBIbDmPfq9d79qF4yQ"
	mockTime := time.UnixMilli(0)
	mockTimeService := utils.NewOverridableMockTimeService(mockTime)

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json",
		fmt.Sprintf(`{"jwt":"%s"}`, expectedToken))

	// When...
	_, err := GetBearerTokenFromTokenJsonFile(mockFileSystem, mockGalasaHome, mockTimeService)

	// Then...
	assert.NotNil(t, err, "Should return an error when a bearer token has expired")
	assert.ErrorContains(t, err, "GAL1108E")
}

func TestGetBearerTokenFromTokenJsonFileWithMissingTokenFileReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")
	mockTimeService := utils.NewMockTimeService()

	// When...
	_, err := GetBearerTokenFromTokenJsonFile(mockFileSystem, mockGalasaHome, mockTimeService)

	// Then...
	assert.NotNil(t, err, "Should return an error when bearer token file does not exist")
	assert.ErrorContains(t, err, "GAL1107E")
}

func TestGetBearerTokenFromTokenJsonFileWithBadContentsReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")
	mockTimeService := utils.NewMockTimeService()

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath() + "/bearer-token.json",
		"notabearertoken")

	// When...
	_, err := GetBearerTokenFromTokenJsonFile(mockFileSystem, mockGalasaHome, mockTimeService)

	// Then...
	assert.NotNil(t, err, "Should return an error when the bearer token file exists but doesn't contain valid JSON")
	assert.ErrorContains(t, err, "GAL1107E")
}