/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

type ImageRenderer interface {
	RenderJsonBytesToImageBytes(jsonBinary []byte) ([]byte, error)
}

type ImageRendererImpl struct {
}

func NewImageRenderer() ImageRenderer {
	renderer := new(ImageRendererImpl)
	return renderer
}

func (renderer *ImageRendererImpl) RenderJsonBytesToImageBytes(jsonBinary []byte) ([]byte, error) {
	var err error
	return jsonBinary, err
}
