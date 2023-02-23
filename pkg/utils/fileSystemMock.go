/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"bytes"
	"log"
	"os"
)

// ------------------------------------------------------------------------------------
// The implementation of the file system interface built on an in-memory map.
// ------------------------------------------------------------------------------------
type Node struct {
	content []byte
	isDir   bool
}

type MockFileSystem struct {
	// Where the in-memory data is kept.
	data map[string]*Node

	// Collects warnings messages
	warningMessageBuffer *bytes.Buffer

	// The mock struct contains methods which can be over-ridden on a per-test basis.
	// The New
	VirtualFunction_MkdirAll             func(targetFolderPath string) error
	VirtualFunction_WriteTextFile        func(targetFilePath string, desiredContents string) error
	VirtualFunction_ReadTextFile         func(filePath string) (string, error)
	VirtualFunction_Exists               func(path string) (bool, error)
	VirtualFunction_DirExists            func(path string) (bool, error)
	VirtualFunction_GetUserHomeDir       func() (string, error)
	VirtualFunction_WriteBinaryFile      func(targetFilePath string, desiredContents []byte) error
	VirtualFunction_OutputWarningMessage func(string) error
}

// NewMockFileSystem creates an implementation of the thin file system layer which delegates
// to a memory map. This uses the default behaviour for all the virtual functions in our
// MockFileSystem
func NewMockFileSystem() FileSystem {
	mockFileSystem := NewOverridableMockFileSystem()
	return mockFileSystem
}

// NewOverridableMockFileSystem creates an implementation of the thin file system layer which delegates
// to delegates to a memory map, but because the MockFileSystem is returned (rather than a FileSystem)
// it means the caller can set up different virtual functions, to change the behaviour.
func NewOverridableMockFileSystem() *MockFileSystem {

	// Allocate the structure
	mockFileSystem := MockFileSystem{
		data: make(map[string]*Node)}

	mockFileSystem.warningMessageBuffer = &bytes.Buffer{}

	// Set up functions inside the structure to call the basic/default mock versions...
	// These can later be over-ridden on a test-by-test basis.
	mockFileSystem.VirtualFunction_MkdirAll = func(targetFolderPath string) error {
		return mockFSMkdirAll(mockFileSystem, targetFolderPath)
	}
	mockFileSystem.VirtualFunction_WriteTextFile = func(targetFilePath string, desiredContents string) error {
		return mockFSWriteTextFile(mockFileSystem, targetFilePath, desiredContents)
	}
	mockFileSystem.VirtualFunction_ReadTextFile = func(filePath string) (string, error) {
		return mockFSReadTextFile(mockFileSystem, filePath)
	}
	mockFileSystem.VirtualFunction_Exists = func(path string) (bool, error) {
		return mockFSExists(mockFileSystem, path)
	}
	mockFileSystem.VirtualFunction_DirExists = func(path string) (bool, error) {
		return mockFSDirExists(mockFileSystem, path)
	}
	mockFileSystem.VirtualFunction_GetUserHomeDir = func() (string, error) {
		return mockFSGetUserHomeDir()
	}
	mockFileSystem.VirtualFunction_WriteBinaryFile = func(path string, content []byte) error {
		return mockFSWriteBinaryFile(mockFileSystem, path, content)
	}
	mockFileSystem.VirtualFunction_OutputWarningMessage = func(message string) error {
		return mockFSOutputWarningMessage(mockFileSystem, message)
	}

	return &mockFileSystem
}

//------------------------------------------------------------------------------------
// Interface methods...
//------------------------------------------------------------------------------------

func (fs *MockFileSystem) MkdirAll(targetFolderPath string) error {
	// Call the virtual function.
	return fs.VirtualFunction_MkdirAll(targetFolderPath)
}

func (fs *MockFileSystem) WriteBinaryFile(targetFilePath string, desiredContents []byte) error {
	return fs.VirtualFunction_WriteBinaryFile(targetFilePath, desiredContents)
}

// WriteTextFile writes a string to a text file
func (fs *MockFileSystem) WriteTextFile(targetFilePath string, desiredContents string) error {
	// Call the virtual function.
	return fs.VirtualFunction_WriteTextFile(targetFilePath, desiredContents)
}

func (fs *MockFileSystem) ReadTextFile(filePath string) (string, error) {
	// Call the virtual function.
	return fs.VirtualFunction_ReadTextFile(filePath)
}

func (fs *MockFileSystem) Exists(path string) (bool, error) {
	// Call the virtual function.
	return fs.VirtualFunction_Exists(path)
}

func (fs *MockFileSystem) DirExists(path string) (bool, error) {
	// Call the virtual function.
	return fs.VirtualFunction_DirExists(path)
}

func (fs *MockFileSystem) GetUserHomeDir() (string, error) {
	return fs.VirtualFunction_GetUserHomeDir()
}

func (fs MockFileSystem) OutputWarningMessage(message string) error {
	return fs.VirtualFunction_OutputWarningMessage(message)
}

//------------------------------------------------------------------------------------
// Default implementations of the methods...
//------------------------------------------------------------------------------------

func mockFSMkdirAll(fs MockFileSystem, targetFolderPath string) error {
	nodeToAdd := Node{content: []byte(""), isDir: true}
	fs.data[targetFolderPath] = &nodeToAdd
	return nil
}

func mockFSWriteBinaryFile(fs MockFileSystem, targetFilePath string, desiredContents []byte) error {
	nodeToAdd := Node{content: desiredContents, isDir: false}
	fs.data[targetFilePath] = &nodeToAdd
	return nil
}

func mockFSWriteTextFile(fs MockFileSystem, targetFilePath string, desiredContents string) error {
	nodeToAdd := Node{content: []byte(desiredContents), isDir: false}
	fs.data[targetFilePath] = &nodeToAdd
	return nil
}

func mockFSReadTextFile(fs MockFileSystem, filePath string) (string, error) {
	text := ""
	var err error = nil
	node := fs.data[filePath]
	if node == nil {
		err = os.ErrNotExist
	} else {
		text = string(node.content)
	}
	return text, err
}

func mockFSExists(fs MockFileSystem, path string) (bool, error) {
	isExists := true
	var err error = nil
	node := fs.data[path]
	if node == nil {
		isExists = false
	}
	return isExists, err
}

func mockFSDirExists(fs MockFileSystem, path string) (bool, error) {
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

func mockFSGetUserHomeDir() (string, error) {
	return "/User/Home/testuser", nil
}

func mockFSOutputWarningMessage(fs MockFileSystem, message string) error {
	log.Printf("Mock warning message collected: %s", message)
	fs.warningMessageBuffer.WriteString(message)
	return nil
}

//------------------------------------------------------------------------------------
// Extra methods on the mock to allow unit tests to get data out of the mock object.
//------------------------------------------------------------------------------------

func (fs MockFileSystem) GetAllWarningMessages() string {
	messages := fs.warningMessageBuffer.String()
	log.Printf("Mock reading back previously collected warnings messages: %s", messages)
	return messages
}
