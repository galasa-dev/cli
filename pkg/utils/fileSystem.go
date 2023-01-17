/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"errors"
	"os"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
)

// FileSystem is a thin interface layer above the os package which can be mocked out
type FileSystem interface {
	// MkdirAll creates all folders in the file system if they don't already exist.
	MkdirAll(targetFolderPath string) error

	ReadTextFile(filePath string) (string, error)

	WriteTextFile(targetFilePath string, desiredContents string) error

	Exists(path string) (bool, error)

	DirExists(path string) (bool, error)
}

//------------------------------------------------------------------------------------
// The implementation of the real os-delegating variant of the FileSystem interface
//------------------------------------------------------------------------------------

type OSFileSystem struct {
}

// NewOSFileSystem creates an implementation of the thin file system layer which delegates
// to the real os package calls.
func NewOSFileSystem() FileSystem {
	return OSFileSystem{}
}

//------------------------------------------------------------------------------------
// Interface methods...
//------------------------------------------------------------------------------------

func (OSFileSystem) MkdirAll(targetFolderPath string) error {
	err := os.MkdirAll(targetFolderPath, 0755)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_CREATE_FOLDERS, targetFolderPath, err.Error())
	}
	return err
}

func (OSFileSystem) WriteTextFile(targetFilePath string, desiredContents string) error {
	bytes := []byte(desiredContents)
	err := os.WriteFile(targetFilePath, bytes, 0644)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_WRITE_FILE, targetFilePath, err.Error())
	}
	return err
}

func (OSFileSystem) ReadTextFile(filePath string) (string, error) {
	text := ""
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_READ_FILE, filePath, err.Error())
	} else {
		text = string(bytes)
	}
	return text, err
}

func (OSFileSystem) Exists(path string) (bool, error) {
	isExists := true
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does not exist
			isExists = false
			err = nil
		}
	}
	return isExists, err
}

func (OSFileSystem) DirExists(path string) (bool, error) {
	isDirExists := true
	metadata, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does not exist
			isDirExists = false
			err = nil
		}
	} else {
		isDirExists = metadata.IsDir()
	}
	return isDirExists, err
}
