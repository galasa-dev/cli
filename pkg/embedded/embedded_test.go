/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package embedded

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockReadOnlyFileSystem struct {
	files map[string]string
}

func NewMockReadOnlyFileSystem() *MockReadOnlyFileSystem {
	fs := MockReadOnlyFileSystem{
		files: make(map[string]string, 0),
	}
	return &fs
}

// WriteFile - This function is not on the ReadOnlyFileSystem interface, but does allow unit tests
// to add data files to the mock file system, so the code under test can read it back.
func (fs *MockReadOnlyFileSystem) WriteFile(filePath string, content string) {
	fs.files[filePath] = content
}

func (fs *MockReadOnlyFileSystem) ReadFile(filePath string) ([]byte, error) {
	content := fs.files[filePath]
	return []byte(content), nil
}

func TestCanParseVersionsFromEmbeddedFS(t *testing.T) {
	propsFileName := "templates/version/build.properties"
	content := "galasactl.version=myVersion\n" +
		"galasa.boot.jar.version=0.1.2\n" +
		"galasa.framework.version=3.4.5\n"

	fs := NewMockReadOnlyFileSystem()
	fs.WriteFile(propsFileName, content)

	versions, err := readVersionsFromEmbeddedFile(fs, nil)

	assert.Nil(t, err)
	assert.NotNil(t, versions)
	assert.Equal(t, "0.1.2", versions.galasaBootJarVersion)
	assert.Equal(t, "3.4.5", versions.galasaFrameworkVersion)
	assert.Equal(t, "myVersion", versions.galasactlVersion)
}

func TestDoesntReReadVersionsFromEmbeddedFSWhenAlreadyKnowAnswers(t *testing.T) {
	propsFileName := "templates/version/build.properties"
	content := "galasactl.version=myVersion\n" +
		"galasa.boot.jar.version=0.1.2\n" +
		"galasa.framework.version=3.4.5\n"

	fs := NewMockReadOnlyFileSystem()
	fs.WriteFile(propsFileName, content)

	alreadyKnownVersions := &versions{
		galasaFrameworkVersion: "myFrameworkVersion",
		galasaBootJarVersion:   "myBootJarVersion",
		galasactlVersion:       "myGalasaCtlVersion",
	}

	versions, err := readVersionsFromEmbeddedFile(fs, alreadyKnownVersions)

	assert.Nil(t, err)
	assert.NotNil(t, versions)
	assert.Equal(t, "myBootJarVersion", versions.galasaBootJarVersion)
	assert.Equal(t, "myFrameworkVersion", versions.galasaFrameworkVersion)
	assert.Equal(t, "myGalasaCtlVersion", versions.galasactlVersion)
}
