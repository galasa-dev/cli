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

func TestCanCalculateTargetPathsOk(t *testing.T) {
	fs := files.NewMockFileSystem()
	renderer := NewImageRenderer()
	expander := NewImageExpander(fs, renderer).(*ImageExpanderImpl)
	folderPath, err := expander.calculateTargetImagePaths("a/b/terminals/c/e.gz")

	assert.Nil(t, err, "could not get paths when we should have been able to.")
	if err == nil {
		assert.Equal(t, "a/b/images/c", folderPath)
	}
}

func TestCanExpandAGzFileToAnImageFile(t *testing.T) {

	realFs := files.NewOSFileSystem()

	var gzContents []byte
	var err error

	gzContents, err = realFs.ReadBinaryFile("./testdata/term1-00001.gz")

	assert.Nil(t, err, "could not load the real gz file data")
	if err == nil {

		// Load the real gz contents into the mock file system...
		// so any image files we create don't infect the real file system.
		fs := files.NewMockFileSystem()
		err = fs.WriteBinaryFile("/U423/zos3270/terminals/term1/term1-0001.gz", gzContents)
		assert.Nil(t, err, "could not write real gz contents into the mock file system")
		if err == nil {

			renderer := NewImageRenderer()
			expander := NewImageExpander(fs, renderer)

			// When...
			err = expander.ExpandImages("/U423")
			assert.Nil(t, err, "could not expand images")
			if err == nil {

				// Then...
				var isExists bool
				isExists, err = fs.DirExists("/U423/zos3270/images/term1")
				assert.Nil(t, err, "could not find out if file exists or not")
				if err == nil {

					assert.True(t, isExists, "Image folder %s was not created.", "/U423/zos3270/images/term1")
					if isExists {

						// Read the rendered file contents.
						var renderedContents []byte
						renderedContents, err = fs.ReadBinaryFile("/U423/zos3270/images/term1/term1-0001.png")
						assert.Nil(t, err, "could not read rendered file")
						if err == nil {

							// Read the file contents which we think the image should be once rendered.
							var expectedContents []byte
							expectedContents, err = realFs.ReadBinaryFile("./testdata/term1-00001.png")
							assert.Nil(t, err, "could not read image file we will compare against")
							if err == nil {

								// Compare the files.
								if assert.Equal(t, expectedContents, renderedContents, "Rendered file does not match the expected") {

									assert.Equal(t, 1, expander.GetExpandedImageFileCount(), "wrong number of expanded files counted.")
								}
							}
						}
					}
				}
			}
		}
	}
}
