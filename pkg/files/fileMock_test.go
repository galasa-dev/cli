/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteAppendsDataToFile(t *testing.T) {
	// Given...
	mockFileSystem := NewOverridableMockFileSystem()
	mockFile := NewOverridableMockFile(mockFileSystem, "/files/dummy.txt")
	mockFileSystem.Create(mockFile.path)

	// When...
	desiredContents := []byte("dummy data")
	bytesWritten, err := mockFile.Write(desiredContents)

	// Then...
	assert.Nil(t, err)
	assert.Equal(t, len(desiredContents), bytesWritten)
	assert.Equal(t, mockFileSystem.data[mockFile.path].content, desiredContents)
}
