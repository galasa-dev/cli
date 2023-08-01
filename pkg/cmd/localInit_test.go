/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateGalasaHomeFolderWhenNotAlreadyInitialisedNonDevelopmentMode(t *testing.T) {
	mockFileSystem := files.NewMockFileSystem()
	mockEnv := utils.NewMockEnv()
	homeDir, _ := mockFileSystem.GetUserHomeDirPath()
	galasaDir := homeDir + "/.galasa/"
	m2Dir := homeDir + "/.m2/"
	isDevelopment := false
	cmdFlagGalasaHome := ""

	err := localEnvInit(mockFileSystem, mockEnv, cmdFlagGalasaHome, isDevelopment)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assertBootstrapPropertiesCreated(t, mockFileSystem, galasaDir)
	assertCPSPropertiesCreated(t, mockFileSystem, galasaDir)
	assertCredentialsPropertiesCreated(t, mockFileSystem, galasaDir)
	assertDSSPropertiesCreated(t, mockFileSystem, galasaDir)
	assertOverridesPropertiesCreated(t, mockFileSystem, galasaDir)
	assertSettingsXMLCreatedAndContentOk(t, mockFileSystem, m2Dir, isDevelopment)
}

func TestCanCreateGalasaHomeFolderWhenNotAlreadyInitialisedWithDevelopmentMode(t *testing.T) {
	mockFileSystem := files.NewMockFileSystem()
	mockEnv := utils.NewMockEnv()
	homeDir, _ := mockFileSystem.GetUserHomeDirPath()
	galasaDir := homeDir + "/.galasa/"
	m2Dir := homeDir + "/.m2/"
	isDevelopment := true
	cmdFlagGalasaHome := ""

	err := localEnvInit(mockFileSystem, mockEnv, cmdFlagGalasaHome, isDevelopment)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assertBootstrapPropertiesCreated(t, mockFileSystem, galasaDir)
	assertCPSPropertiesCreated(t, mockFileSystem, galasaDir)
	assertCredentialsPropertiesCreated(t, mockFileSystem, galasaDir)
	assertDSSPropertiesCreated(t, mockFileSystem, galasaDir)
	assertOverridesPropertiesCreated(t, mockFileSystem, galasaDir)
	assertSettingsXMLCreatedAndContentOk(t, mockFileSystem, m2Dir, isDevelopment)
}

func TestCanCreateGalasaHomeFolderWhenAlreadyInitialised(t *testing.T) {
	mockFileSystem := files.NewMockFileSystem()
	homeDir, _ := mockFileSystem.GetUserHomeDirPath()
	galasaDir := homeDir + "/.galasa/"
	m2Dir := homeDir + "/.m2/"
	cmdFlagGalasaHome := ""
	isDevelopment := false

	mockFileSystem.WriteTextFile(galasaDir+"bootstrap.properties", "")
	mockFileSystem.WriteTextFile(galasaDir+"dss.properties", "")
	mockFileSystem.WriteTextFile(galasaDir+"cps.properties", "")
	mockFileSystem.WriteTextFile(galasaDir+"credentials.properties", "")
	mockFileSystem.WriteTextFile(galasaDir+"overrides.properties", "")
	mockFileSystem.WriteTextFile(m2Dir+"settings.xml", "")

	mockEnv := utils.NewMockEnv()
	err := localEnvInit(mockFileSystem, mockEnv, cmdFlagGalasaHome, isDevelopment)
	if err != nil {
		assert.Fail(t, err.Error())
	}

}

func assertBootstrapPropertiesCreated(t *testing.T, mockFileSystem files.FileSystem, galasaDir string) {
	testBootstrapPropertiesExists, err := mockFileSystem.Exists(galasaDir + "bootstrap.properties")
	assert.Nil(t, err)
	assert.True(t, testBootstrapPropertiesExists, "Bootstrap properties was not created")
}

func assertCPSPropertiesCreated(t *testing.T, mockFileSystem files.FileSystem, galasaDir string) {
	testCPSPropertiesExists, err := mockFileSystem.Exists(galasaDir + "cps.properties")
	assert.Nil(t, err)
	assert.True(t, testCPSPropertiesExists, "CPS properties was not created")
}

func assertCredentialsPropertiesCreated(t *testing.T, mockFileSystem files.FileSystem, galasaDir string) {
	testCredentialsPropertiesExists, err := mockFileSystem.Exists(galasaDir + "credentials.properties")
	assert.Nil(t, err)
	assert.True(t, testCredentialsPropertiesExists, "Credentials properties was not created")
}

func assertDSSPropertiesCreated(t *testing.T, mockFileSystem files.FileSystem, galasaDir string) {
	testDSSPropertiesExists, err := mockFileSystem.Exists(galasaDir + "dss.properties")
	assert.Nil(t, err)
	assert.True(t, testDSSPropertiesExists, "DSS properties was not created")
}

func assertOverridesPropertiesCreated(t *testing.T, mockFileSystem files.FileSystem, galasaDir string) {
	testOverridesPropertiesExists, err := mockFileSystem.Exists(galasaDir + "overrides.properties")
	assert.Nil(t, err)
	assert.True(t, testOverridesPropertiesExists, "Overrides properties was not created")
}

func assertSettingsXMLCreatedAndContentOk(t *testing.T, mockFileSystem files.FileSystem, m2Dir string, isDevelopment bool) {
	testSettingsXMLExists, err := mockFileSystem.Exists(m2Dir + "settings.xml")
	assert.Nil(t, err)
	assert.True(t, testSettingsXMLExists, "Settings.xml was not created")
	settingsContent, err := mockFileSystem.ReadTextFile(m2Dir + "settings.xml")
	assert.Nil(t, err)
	assert.Contains(t, settingsContent, "<url>https://development.galasa.dev/main/maven-repo/obr</url>")

	if isDevelopment {
		assert.Contains(t, settingsContent, "<!-- Using the bleeding edge version of galasa. Comment out if not needed. -->")
	} else {
		assert.Contains(t, settingsContent, "<!-- To use the bleeding edge version of galasa, use the development obr")
	}
}
