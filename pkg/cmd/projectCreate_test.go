/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"errors"
	"testing"

	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateProjectGoldenPath(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := true

	// When ...
	err := createProject(mockFileSystem, "my.test.package", forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	assert.Nil(t, err, "Golden path should not return an error")

	parentFolderExists, err := mockFileSystem.DirExists("my.test.package")
	assert.Nil(t, err)
	assert.True(t, parentFolderExists, "Parent folder was not created.")

	parentPomXmlFileExists, err := mockFileSystem.Exists("my.test.package/pom.xml")
	assert.Nil(t, err)
	assert.True(t, parentPomXmlFileExists, "Parent folder pom.xml was not created.")

	text, err := mockFileSystem.ReadTextFile("my.test.package/pom.xml")
	assert.Nil(t, err)
	assert.Contains(t, text, "<groupId>my.test.package</groupId>", "parent pom.xml didn't substitute the group id")
	assert.Contains(t, text, "<artifactId>my.test.package</artifactId>", "parent pom.xml didn't substitute the artifact id")

	testFolderExists, err := mockFileSystem.DirExists("my.test.package/my.test.package.test")
	assert.Nil(t, err)
	assert.True(t, testFolderExists, "Test folder was not created.")

	expectedPomFilePath := "my.test.package/my.test.package.test/pom.xml"
	testPomXmlFileExists, err := mockFileSystem.Exists(expectedPomFilePath)
	assert.Nil(t, err)
	assert.True(t, testPomXmlFileExists, "Test folder pom.xml was not created.")

	text, err = mockFileSystem.ReadTextFile(expectedPomFilePath)
	assert.Nil(t, err)
	assert.Contains(t, text, "<groupId>my.test.package</groupId>", "Test folder pom.xml didn't substitute the group id")
	assert.Contains(t, text, "<artifactId>my.test.package.test</artifactId>", "Test folder pom.xml didn't substitute the artifact id")

	testSrcFolderExists, err := mockFileSystem.DirExists("my.test.package/my.test.package.test/src/main/java/my/test/package/test")
	assert.Nil(t, err)
	assert.True(t, testSrcFolderExists, "Test src folder was not created.")

	expectedJavaFilePath := "my.test.package/my.test.package.test/src/main/java/my/test/package/test/SampleTest.java"
	testJavaFileExists, err := mockFileSystem.Exists(expectedJavaFilePath)
	assert.Nil(t, err)
	assert.True(t, testJavaFileExists, "Test java file was not created.")

	text, err = mockFileSystem.ReadTextFile(expectedJavaFilePath)
	assert.Nil(t, err)
	assert.Contains(t, text, "package my.test.package.test", "Test java file didn't substitute the java package")

}

func TestCreateProjectErrorsWhenMkAllDirsFails(t *testing.T) {

	// Given...
	mockFileSystem := utils.NewOverridableMockFileSystem()
	forceOverwrite := true

	// Over-ride the default MkdirAll function so that it fails...
	mockFileSystem.VirtualFunction_MkdirAll = func(targetFolderPath string) error {
		return errors.New("Simulated MkdirAll failure")
	}

	// When ...
	err := createProject(mockFileSystem, "my.test.package", forceOverwrite)

	// Then...
	assert.NotNil(t, err, "Sumulated error didn't bubble up to the top.")
	assert.Contains(t, err.Error(), "Simulated MkdirAll failure", "Failed to return reason for failure.")

}

func TestCreateProjectErrorsWhenWriteTextFileFails(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewOverridableMockFileSystem()
	forceOverwrite := true

	// Over-ride the default WriteTextFile function so that it fails...
	mockFileSystem.VirtualFunction_WriteTextFile = func(targetFilePath string, desiredContents string) error {
		return errors.New("Simulated WriteTextFile failure")
	}

	// When ...
	err := createProject(mockFileSystem, "my.test.package", forceOverwrite)

	// Then...
	assert.NotNil(t, err, "Sumulated error didn't bubble up to the top.")
	assert.Contains(t, err.Error(), "Simulated WriteTextFile failure", "Failed to return reason for failure.")
}

func TestCreateProjectPomFileAlreadyExistsNoForceOverwrite(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := false
	testPackageName := "my.test.package"

	// Create a pom.xml file already...
	mockFileSystem.MkdirAll(testPackageName)
	mockFileSystem.WriteTextFile(testPackageName+"/pom.xml", "dummy test pom.xml")

	// When ...
	err := createProject(mockFileSystem, testPackageName, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1036E")
}

func TestCreateProjectPomFileAlreadyExistsWithForceOverwrite(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := true
	testPackageName := "my.test.package"

	// Create a pom.xml file already...
	mockFileSystem.MkdirAll(testPackageName)
	mockFileSystem.WriteTextFile(testPackageName+"/pom.xml", "dummy test pom.xml")

	// When ...
	err := createProject(mockFileSystem, testPackageName, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	assert.Nil(t, err)

	parentPomXmlFileExists, err := mockFileSystem.Exists("my.test.package/pom.xml")
	assert.Nil(t, err)
	assert.True(t, parentPomXmlFileExists, "Parent folder pom.xml was not created.")

	text, err := mockFileSystem.ReadTextFile("my.test.package/pom.xml")
	assert.Nil(t, err)
	assert.Contains(t, text, "<groupId>my.test.package</groupId>", "parent pom.xml didn't substitute the group id")
	assert.Contains(t, text, "<artifactId>my.test.package</artifactId>", "parent pom.xml didn't substitute the artifact id")

}
