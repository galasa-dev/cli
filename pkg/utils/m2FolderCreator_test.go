/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/embedded"
	"github.com/stretchr/testify/assert"
)

// To validate the string as a valid java package name before we start to use it.
func TestCanCreateM2FolderAndSettingsXML(t *testing.T) {
	// Given...
	mockFileSystem := NewMockFileSystem()
	embeddedFileSystem := embedded.GetEmbeddedFileSystem()

	// When ...
	err := InitialiseM2Folder(mockFileSystem, embeddedFileSystem)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, err.Error())
	}

	// Check the home folder has been created.
	homeDir, _ := mockFileSystem.GetUserHomeDir()
	m2Dir := homeDir + "/.m2"
	assertFolderExists(t, mockFileSystem, m2Dir, "Didn't create "+m2Dir+" folder in home directory.")

	isExists, _ := mockFileSystem.Exists(m2Dir + "/settings.xml")
	assert.True(t, isExists, "Failed to create file "+m2Dir+"/settings.xml")

	settingsContent, _ := mockFileSystem.ReadTextFile(m2Dir + "/settings.xml")
	assert.Contains(t, settingsContent, "<url>https://development.galasa.dev/main/maven-repo/obr</url>", "Test settings.xml didn't have the correct obr url")

}
