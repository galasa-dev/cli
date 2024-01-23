/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

import (
	"testing"

	"github.com/galasa-dev/cli/pkg/files"
	"github.com/stretchr/testify/assert"
)

func TestCanWriteAFile(t *testing.T) {
	fs := files.NewMockFileSystem()
	fs.MkdirAll("/my/images/folder")
	writer := NewImageFileWriter(fs, "/my/images/folder")
	buff := []byte("hello world")

	err := writer.WriteImageFile("File1.png", buff)

	if assert.Nil(t, err, "should have been able to write an image file with no error!") {
		var isExists bool
		isExists, err = fs.Exists("/my/images/folder/File1.png")
		if assert.Nil(t, err, "The image file should have been created! Error trying to see if it's there or not.") {
			assert.True(t, isExists, "The image file should have been created. It's not there.")
		}
	}

	assert.Equal(t, writer.GetImageFilesWrittenCount(), 1)
}
