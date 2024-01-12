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
		mockGalasaHome.GetNativeFolderPath()+"/galasactl.properties",
		fmt.Sprintf(
			"GALASA_CLIENT_ID=%s\n"+
				"GALASA_SECRET=%s\n"+
				"GALASA_ACCESS_TOKEN=%s", clientIdValue, secretValue, accessTokenValue))
	// When...
	authProperties, err := GetAuthProperties(mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.Nil(t, err, "Should not return an error if the galasactl.properties exists and all required properties are present")
	assert.Equal(t, clientIdValue, authProperties.GetClientId())
	assert.Equal(t, secretValue, authProperties.GetSecret())
	assert.Equal(t, accessTokenValue, authProperties.GetRefreshToken())
}

func TestGetAuthPropertiesWithNoClientIdValidPropertiesUnmarshalsAuthPropertiesFails(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	secretValue := "dummySecret"
	accessTokenValue := "abc"

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath()+"/galasactl.properties",
		fmt.Sprintf(
			// "GALASA_CLIENT_ID=%s\n"+
			"GALASA_SECRET=%s\n"+
				"GALASA_ACCESS_TOKEN=%s",
			// clientIdValue,
			secretValue, accessTokenValue))
	// When...
	_, err := GetAuthProperties(mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.NotNil(t, err, "Should return an error as the galasactl.properties exists but the client id is missing from the file.")
	assert.Contains(t, err.Error(), "GAL1122E")
	assert.Contains(t, err.Error(), "GALASA_CLIENT_ID")
}

func TestGetAuthPropertiesWithNoSecretIdValidPropertiesUnmarshalsAuthPropertiesFails(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	clientIdValue := "my-client-id"
	// secretValue := "dummySecret"
	accessTokenValue := "abc"

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath()+"/galasactl.properties",
		fmt.Sprintf(
			"GALASA_CLIENT_ID=%s\n"+
				// "GALASA_SECRET=%s\n"+
				"GALASA_ACCESS_TOKEN=%s",
			clientIdValue,
			// secretValue,
			accessTokenValue))
	// When...
	_, err := GetAuthProperties(mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.NotNil(t, err, "Should return an error as the galasactl.properties exists but the secret is missing from the file.")
	assert.Contains(t, err.Error(), "GAL1122E")
	assert.Contains(t, err.Error(), "GALASA_SECRET")
}

func TestGetAuthPropertiesWithNoRefreshTokenIdValidPropertiesUnmarshalsAuthPropertiesFails(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	clientIdValue := "my-client-id"
	secretValue := "dummySecret"
	// accessTokenValue := "abc"

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath()+"/galasactl.properties",
		fmt.Sprintf(
			"GALASA_CLIENT_ID=%s\n"+
				"GALASA_SECRET=%s\n",
			// "GALASA_ACCESS_TOKEN=%s",
			clientIdValue,
			secretValue,
		// accessTokenValue,
		))
	// When...
	_, err := GetAuthProperties(mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.NotNil(t, err, "Should return an error as the galasactl.properties exists but the secret is missing from the file.")
	assert.Contains(t, err.Error(), "GAL1122E")
	assert.Contains(t, err.Error(), "GALASA_ACCESS_TOKEN")
}

func TestGetAuthPropertiesWithEmptyGalasactlPropertiesReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	mockFileSystem.WriteTextFile(mockGalasaHome.GetNativeFolderPath()+"/galasactl.properties", "")

	// When...
	_, err := GetAuthProperties(mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.NotNil(t, err, "Should return an error if the galasactl.properties is empty")
	assert.ErrorContains(t, err, "GAL1122E")
}

func TestGetAuthPropertiesWithMissingPropertiesReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	clientIdValue := "dummyId"

	// Create a galasactl.properties file that is missing the GALASA_SECRET and GALASA_ACCESS_TOKEN properties
	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath()+"/galasactl.properties",
		fmt.Sprintf("GALASA_CLIENT_ID=%s", clientIdValue))

	// When...
	_, err := GetAuthProperties(mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.NotNil(t, err, "Should return an error if the galasactl.properties exists and some required properties are missing")
	assert.ErrorContains(t, err, "GAL1122E")
}

func TestGetAuthPropertiesWithMissingGalasactlPropertiesFileReturnsError(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	// When...
	_, err := GetAuthProperties(mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.NotNil(t, err, "Should return an error if the galasactl.properties does not exist")
	assert.ErrorContains(t, err, "GAL1043E")
}

func TestGetAuthPropertiesEnvVarsOverridesFileValues(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	clientIdValue := "client-id-from-file"
	secretValue := "secret-from-file"
	accessTokenValue := "token-from-file"

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath()+"/galasactl.properties",
		fmt.Sprintf(
			"GALASA_CLIENT_ID=%s\n"+
				"GALASA_SECRET=%s\n"+
				"GALASA_ACCESS_TOKEN=%s", clientIdValue, secretValue, accessTokenValue))

	clientIdValue = "client-id-from-env-var"
	secretValue = "secret-from-env-var"
	accessTokenValue = "token-from-env-var"

	mockEnvironment.SetEnv(CLIENT_ID_PROPERTY, clientIdValue)
	mockEnvironment.SetEnv(SECRET_PROPERTY, secretValue)
	mockEnvironment.SetEnv(ACCESS_TOKEN_PROPERTY, accessTokenValue)

	// When...
	authProperties, err := GetAuthProperties(mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.Nil(t, err, "Should not return an error if the galasactl.properties exists and all required properties are present")
	assert.Equal(t, clientIdValue, authProperties.GetClientId())
	assert.Equal(t, secretValue, authProperties.GetSecret())
	assert.Equal(t, accessTokenValue, authProperties.GetRefreshToken())
}

func TestGetAuthPropertiesClientIDEnvVarOverridesFileValue(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	clientIdValue := "client-id-from-file"
	secretValue := "secret-from-file"
	accessTokenValue := "token-from-file"

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath()+"/galasactl.properties",
		fmt.Sprintf(
			"GALASA_CLIENT_ID=%s\n"+
				"GALASA_SECRET=%s\n"+
				"GALASA_ACCESS_TOKEN=%s", clientIdValue, secretValue, accessTokenValue))

	clientIdValue = "client-id-from-env-var"
	// secretValue = "secret-from-env-var"
	// accessTokenValue = "token-from-env-var"

	mockEnvironment.SetEnv(CLIENT_ID_PROPERTY, clientIdValue)
	// mockEnvironment.SetEnv(SECRET_PROPERTY, secretValue)
	// mockEnvironment.SetEnv(ACCESS_TOKEN_PROPERTY, accessTokenValue)

	// When...
	authProperties, err := GetAuthProperties(mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.Nil(t, err, "Should not return an error if the galasactl.properties exists and all required properties are present")
	assert.Equal(t, clientIdValue, authProperties.GetClientId())
	assert.Equal(t, secretValue, authProperties.GetSecret())
	assert.Equal(t, accessTokenValue, authProperties.GetRefreshToken())
}

func TestGetAuthPropertiesSecretEnvVarOverridesFileValue(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	clientIdValue := "client-id-from-file"
	secretValue := "secret-from-file"
	accessTokenValue := "token-from-file"

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath()+"/galasactl.properties",
		fmt.Sprintf(
			"GALASA_CLIENT_ID=%s\n"+
				"GALASA_SECRET=%s\n"+
				"GALASA_ACCESS_TOKEN=%s", clientIdValue, secretValue, accessTokenValue))

	// clientIdValue = "client-id-from-env-var"
	secretValue = "secret-from-env-var"
	// accessTokenValue = "token-from-env-var"

	// mockEnvironment.SetEnv(CLIENT_ID_PROPERTY, clientIdValue)
	mockEnvironment.SetEnv(SECRET_PROPERTY, secretValue)
	// mockEnvironment.SetEnv(ACCESS_TOKEN_PROPERTY, accessTokenValue)

	// When...
	authProperties, err := GetAuthProperties(mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.Nil(t, err, "Should not return an error if the galasactl.properties exists and all required properties are present")
	assert.Equal(t, clientIdValue, authProperties.GetClientId())
	assert.Equal(t, secretValue, authProperties.GetSecret())
	assert.Equal(t, accessTokenValue, authProperties.GetRefreshToken())
}

func TestGetAuthPropertiesRefreshTokenEnvVarOverridesFileValue(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	mockEnvironment := utils.NewMockEnv()
	mockGalasaHome, _ := utils.NewGalasaHome(mockFileSystem, mockEnvironment, "")

	clientIdValue := "client-id-from-file"
	secretValue := "secret-from-file"
	accessTokenValue := "token-from-file"

	mockFileSystem.WriteTextFile(
		mockGalasaHome.GetNativeFolderPath()+"/galasactl.properties",
		fmt.Sprintf(
			"GALASA_CLIENT_ID=%s\n"+
				"GALASA_SECRET=%s\n"+
				"GALASA_ACCESS_TOKEN=%s", clientIdValue, secretValue, accessTokenValue))

	// clientIdValue = "client-id-from-env-var"
	// secretValue = "secret-from-env-var"
	accessTokenValue = "token-from-env-var"

	// mockEnvironment.SetEnv(CLIENT_ID_PROPERTY, clientIdValue)
	// mockEnvironment.SetEnv(SECRET_PROPERTY, secretValue)
	mockEnvironment.SetEnv(ACCESS_TOKEN_PROPERTY, accessTokenValue)

	// When...
	authProperties, err := GetAuthProperties(mockFileSystem, mockGalasaHome, mockEnvironment)

	// Then...
	assert.Nil(t, err, "Should not return an error if the galasactl.properties exists and all required properties are present")
	assert.Equal(t, clientIdValue, authProperties.GetClientId())
	assert.Equal(t, secretValue, authProperties.GetSecret())
	assert.Equal(t, accessTokenValue, authProperties.GetRefreshToken())
}
