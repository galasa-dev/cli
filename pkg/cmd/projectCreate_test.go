/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateProjectGoldenPath(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := true

	// When ...
	createProject(mockFileSystem, "my.test.package", forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.

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
