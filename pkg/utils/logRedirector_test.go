/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/stretchr/testify/assert"
)

func TestLogRedirectorFailsWhenLogFileIsAFolder(t *testing.T) {

	var err error = nil

	fileSystem := files.NewOverridableMockFileSystem()

	// Create a fake folder '.'
	logFileName := "."
	fileSystem.MkdirAll(logFileName)

	err = CaptureLog(fileSystem, logFileName)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1069")
}
