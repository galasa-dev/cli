/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteAppendsDataToFile(t *testing.T) {
	mockFileSystem := NewOverridableMockFileSystem()
	mockFile := NewOverridableMockFile(mockFileSystem, "/files/dummy.txt")
	
	mockFileSystem.Create(mockFile.path)

	desiredContents := []byte("dummy data")
	bytesWritten, err := mockFile.Write(desiredContents)

	assert.Nil(t, err)
	assert.Equal(t, len(desiredContents), bytesWritten)
	assert.Equal(t, mockFileSystem.data[mockFile.path].content, desiredContents)
}
