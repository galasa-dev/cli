/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package embedded

import (
	"embed"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
)

type ReadOnlyFileSystem interface {
	ReadFile(filePath string) ([]byte, error)
}

type EmbeddedFileSystem struct {
	embeddedFileSystem embed.FS
}

func NewReadOnlyFileSystem() ReadOnlyFileSystem {
	result := EmbeddedFileSystem{
		embeddedFileSystem: embeddedFileSystem,
	}
	return &result
}

//------------------------------------------------------------------------------------
// Interface methods...
//------------------------------------------------------------------------------------

// The only thing which this class actually supports.
func (fs *EmbeddedFileSystem) ReadFile(filePath string) ([]byte, error) {

	bytes, err := fs.embeddedFileSystem.ReadFile(filePath)
	if err != nil {
		galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_EMBEDDED_FS_READ_FAILED, err.Error())
	}
	return bytes, err
}
