/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogRedirectorFailsWhenLogFileIsAFolder(t *testing.T) {

	var err error = nil

	fileSystem := NewOverridableMockFileSystem()

	// Create a fake folder '.'
	logFileName := "."
	fileSystem.MkdirAll(logFileName)

	err = CaptureLog(fileSystem, logFileName)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1069")
}
