/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"fmt"
	"testing"

	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetAuthPropertiesWithValidPropertiesUnmarshalsAuthProperties(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	clientIdValue := "dummyId"
	secretValue := "dummySecret"
	accessTokenValue := "abc"

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath() + "/galasactl.properties",
		fmt.Sprintf(
			"GALASA_CLIENT_ID=%s\n" +
			"GALASA_SECRET=%s\n" +
			"GALASA_ACCESS_TOKEN=%s", clientIdValue, secretValue, accessTokenValue))
	// When...
	authProperties, err := GetAuthProperties(mockFileSystem, mockGalasaHome)

	// Then...
	assert.Nil(t, err, "Should not return an error if the galasactl.properties exists and all required properties are present")
	assert.Equal(t, clientIdValue, authProperties.GetClientId())
	assert.Equal(t, secretValue, authProperties.GetSecret())
	assert.Equal(t, accessTokenValue, authProperties.GetRefreshToken())
}

func TestGetAuthPropertiesWithEmptyGalasactlPropertiesReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	mockFileSystem.WriteTextFile(mockGalasaHome.GetNativeFolderPath() + "/galasactl.properties", "")

	// When...
	_, err := GetAuthProperties(mockFileSystem, mockGalasaHome)

	// Then...
	assert.NotNil(t, err, "Should return an error if the galasactl.properties is empty")
	assert.ErrorContains(t, err, "GAL1105E")
}

func TestGetAuthPropertiesWithMissingPropertiesReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	clientIdValue := "dummyId"

	// Create a galasactl.properties file that is missing the GALASA_SECRET and GALASA_ACCESS_TOKEN properties
	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath() + "/galasactl.properties",
		fmt.Sprintf("GALASA_CLIENT_ID=%s", clientIdValue))

	// When...
	_, err := GetAuthProperties(mockFileSystem, mockGalasaHome)

	// Then...
	assert.NotNil(t, err, "Should return an error if the galasactl.properties exists and some required properties are missing")
	assert.ErrorContains(t, err, "GAL1105E")
}

func TestGetAuthPropertiesWithMissingGalasactlPropertiesFileReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	// When...
	_, err := GetAuthProperties(mockFileSystem, mockGalasaHome)

	// Then...
	assert.NotNil(t, err, "Should return an error if the galasactl.properties does not exist")
	assert.ErrorContains(t, err, "GAL1043E")
}
