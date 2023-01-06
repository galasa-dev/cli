/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateProjectGoldenPath(t *testing.T) {
	// Given...
	mockFileSystem := afero.NewMemMapFs()
	forceOverwrite := true

	// When ...
	createProject(mockFileSystem, "my.test.package", forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	utility := afero.Afero{Fs: mockFileSystem}

	parentFolderExists, err := utility.DirExists("my.test.package")
	assert.Nil(t, err)
	assert.True(t, parentFolderExists, "Parent folder was not created.")

	parentPomXmlFileExists, err := utility.Exists("my.test.package/pom.xml")
	assert.Nil(t, err)
	assert.True(t, parentPomXmlFileExists, "Parent folder pom.xml was not created.")

	bytes, err := utility.ReadFile("my.test.package/pom.xml")
	assert.Nil(t, err)
	assert.Contains(t, string(bytes), "<groupId>my.test.package</groupId>", "parent pom.xml didn't substitute the group id")
	assert.Contains(t, string(bytes), "<artifactId>my.test.package</artifactId>", "parent pom.xml didn't substitute the artifact id")

	testFolderExists, err := utility.DirExists("my.test.package/my.test.package.test")
	assert.Nil(t, err)
	assert.True(t, testFolderExists, "Test folder was not created.")

	expectedPomFilePath := "my.test.package/my.test.package.test/pom.xml"
	testPomXmlFileExists, err := utility.Exists(expectedPomFilePath)
	assert.Nil(t, err)
	assert.True(t, testPomXmlFileExists, "Test folder pom.xml was not created.")

	bytes, err = utility.ReadFile(expectedPomFilePath)
	assert.Nil(t, err)
	assert.Contains(t, string(bytes), "<groupId>my.test.package</groupId>", "Test folder pom.xml didn't substitute the group id")
	assert.Contains(t, string(bytes), "<artifactId>my.test.package.test</artifactId>", "Test folder pom.xml didn't substitute the artifact id")

	testSrcFolderExists, err := utility.DirExists("my.test.package/my.test.package.test/src/main/java/my/test/package/test")
	assert.Nil(t, err)
	assert.True(t, testSrcFolderExists, "Test src folder was not created.")

	expectedJavaFilePath := "my.test.package/my.test.package.test/src/main/java/my/test/package/test/SampleTest.java"
	testJavaFileExists, err := utility.Exists(expectedJavaFilePath)
	assert.Nil(t, err)
	assert.True(t, testJavaFileExists, "Test java file was not created.")

	bytes, err = utility.ReadFile(expectedJavaFilePath)
	assert.Nil(t, err)
	assert.Contains(t, string(bytes), "package my.test.package.test", "Test java file didn't substitute the java package")

}
