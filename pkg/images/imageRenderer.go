/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

type ImageRenderer interface {
	RenderJsonBytesToImageFiles(jsonBinary []byte, writer ImageFileWriter) error
}

type ImageRendererImpl struct {
}

func NewImageRenderer() ImageRenderer {
	renderer := new(ImageRendererImpl)
	return renderer
}

func (renderer *ImageRendererImpl) RenderJsonBytesToImageFiles(jsonBinary []byte, writer ImageFileWriter) error {
	var err error

	// Call the image file writer here...

	return err
}
