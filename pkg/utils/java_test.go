/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlankJavaHomeIsInvalid(t *testing.T) {
	javaHome := ""

	fileSystem := NewMockFileSystem()
	AddJavaRuntimeToMock(fileSystem, javaHome)

	err := ValidateJavaHome(fileSystem, "")
	if err == nil {
		assert.Fail(t, "Didn't fail the validation.")
	}
	assert.Contains(t, err.Error(), "GAL1050E", "Wrong error message")
}

func TestTrailingSlashInJavaHomeIsValid(t *testing.T) {
	javaHome := "/java"

	fileSystem := NewMockFileSystem()
	AddJavaRuntimeToMock(fileSystem, javaHome)

	err := ValidateJavaHome(fileSystem, javaHome+"/")
	if err != nil {
		assert.Fail(t, "Failed the validation. But should have passed ! %s", err.Error())
	}
}

func TestNonExistentJavaHomeFolderIsInValid(t *testing.T) {
	javaHome := "/java"

	fileSystem := NewMockFileSystem()
	// AddJavaRuntimeToMock(fileSystem, javaHome)

	err := ValidateJavaHome(fileSystem, javaHome)
	if err == nil {
		assert.Fail(t, "Didn't fail the validation.")
	}
	assert.Contains(t, err.Error(), "GAL1052E", "Wrong error message")
}

func TestNonExistentJavaBinaryHomeIsInValid(t *testing.T) {
	javaHome := "/java"

	fileSystem := NewMockFileSystem()
	fileSystem.MkdirAll("/java/bin")
	// AddJavaRuntimeToMock(fileSystem, javaHome)

	err := ValidateJavaHome(fileSystem, javaHome)
	if err == nil {
		assert.Fail(t, "Didn't fail the validation.")
	}
	assert.Contains(t, err.Error(), "GAL1054E", "Wrong error message")
}
