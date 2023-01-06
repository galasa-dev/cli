/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"os"
)

// ------------------------------------------------------------------------------------
// The implementation of the file system interface built on an in-memory map.
// ------------------------------------------------------------------------------------
type Node struct {
	content string
	isDir   bool
}

type MockFileSystem struct {
	data map[string]*Node
}

// NewOSFileSystem creates an implementation of the thin file system layer which delegates
// to the real os package calls.
func NewMockFileSystem() FileSystem {
	return MockFileSystem{data: make(map[string]*Node)}
}

//------------------------------------------------------------------------------------
// Interface methods...
//------------------------------------------------------------------------------------

func (fs MockFileSystem) MkdirAll(targetFolderPath string) error {
	nodeToAdd := Node{content: "", isDir: true}
	fs.data[targetFolderPath] = &nodeToAdd
	return nil
}

func (fs MockFileSystem) WriteTextFile(targetFilePath string, desiredContents string) error {
	nodeToAdd := Node{content: desiredContents, isDir: false}
	fs.data[targetFilePath] = &nodeToAdd
	return nil
}

func (fs MockFileSystem) ReadTextFile(filePath string) (string, error) {
	text := ""
	var err error = nil
	node := fs.data[filePath]
	if node == nil {
		err = os.ErrNotExist
	} else {
		text = node.content
	}
	return text, err
}

func (fs MockFileSystem) Exists(path string) (bool, error) {
	isExists := true
	var err error = nil
	node := fs.data[path]
	if node == nil {
		isExists = false
	}
	return isExists, err
}

func (fs MockFileSystem) DirExists(path string) (bool, error) {
	isDirExists := true
	var err error = nil
	node := fs.data[path]
	if node == nil {
		isDirExists = false
	} else {
		isDirExists = node.isDir
	}
	return isDirExists, err
}
