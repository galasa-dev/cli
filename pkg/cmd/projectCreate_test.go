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

func TestCanCreateProjectFailsIfPackageNameInvalid(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false

	// When ...
	err := createProject(mockFileSystem, "very.INVALID_PACKAGE_NAME.very", isObrProjectRequired, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	assert.NotNil(t, err, "Should return an error if java package name is invalid")
	assert.Contains(t, err.Error(), "GAL1037E:", "Wrong error message reported.")
}

func TestCanCreateProjectGoldenPathNoOBR(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false

	// When ...
	err := createProject(mockFileSystem, "my.test.pkg", isObrProjectRequired, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	assert.Nil(t, err, "Golden path should not return an error")

	assertParentFolderAndContentsCreated(t, mockFileSystem, isObrProjectRequired)
	assertTestFolderAndContentsCreatedOk(t, mockFileSystem)
}

func assertParentFolderAndContentsCreated(t *testing.T, mockFileSystem utils.FileSystem, isObrProjectRequired bool) {
	parentFolderExists, err := mockFileSystem.DirExists("my.test.pkg")
	assert.Nil(t, err)
	assert.True(t, parentFolderExists, "Parent folder was not created.")

	parentPomXmlFileExists, err := mockFileSystem.Exists("my.test.pkg/pom.xml")
	assert.Nil(t, err)
	assert.True(t, parentPomXmlFileExists, "Parent folder pom.xml was not created.")

	text, err := mockFileSystem.ReadTextFile("my.test.pkg/pom.xml")
	assert.Nil(t, err)
	assert.Contains(t, text, "<groupId>my.test.pkg</groupId>", "parent pom.xml didn't substitute the group id")
	assert.Contains(t, text, "<artifactId>my.test.pkg</artifactId>", "parent pom.xml didn't substitute the artifact id")

	assert.Contains(t, text, "<module>my.test.pkg.test</module>", "parent pom.xml didn't have a test module included")

	if isObrProjectRequired {
		assert.Contains(t, text, "<module>my.test.pkg.obr</module>", "parent pom.xml didn't have an obr module included")
	} else {
		assert.NotContains(t, text, "<module>my.test.pkg.obr</module>", "parent pom.xml should not have an obr module included")
	}
}

func assertTestFolderAndContentsCreatedOk(t *testing.T, mockFileSystem utils.FileSystem) {

	testFolderExists, err := mockFileSystem.DirExists("my.test.pkg/my.test.pkg.test")
	assert.Nil(t, err)
	assert.True(t, testFolderExists, "Test folder was not created.")

	expectedPomFilePath := "my.test.pkg/my.test.pkg.test/pom.xml"
	testPomXmlFileExists, err := mockFileSystem.Exists(expectedPomFilePath)
	assert.Nil(t, err)
	assert.True(t, testPomXmlFileExists, "Test folder pom.xml was not created.")

	text, err := mockFileSystem.ReadTextFile(expectedPomFilePath)
	assert.Nil(t, err)
	assert.Contains(t, text, "<groupId>my.test.pkg</groupId>", "Test folder pom.xml didn't substitute the group id")
	assert.Contains(t, text, "<artifactId>my.test.pkg.test</artifactId>", "Test folder pom.xml didn't substitute the artifact id")

	testSrcFolderExists, err := mockFileSystem.DirExists("my.test.pkg/my.test.pkg.test/src/main/java/my/test/pkg/test")
	assert.Nil(t, err)
	assert.True(t, testSrcFolderExists, "Test src folder was not created.")

	expectedJavaFilePath := "my.test.pkg/my.test.pkg.test/src/main/java/my/test/pkg/test/SampleTest.java"
	testJavaFileExists, err := mockFileSystem.Exists(expectedJavaFilePath)
	assert.Nil(t, err)
	assert.True(t, testJavaFileExists, "Test java file was not created.")

	text, err = mockFileSystem.ReadTextFile(expectedJavaFilePath)
	assert.Nil(t, err)
	assert.Contains(t, text, "package my.test.pkg.test", "Test java file didn't substitute the java package")

}

func TestCreateProjectErrorsWhenMkAllDirsFails(t *testing.T) {

	// Given...
	mockFileSystem := utils.NewOverridableMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false

	// Over-ride the default MkdirAll function so that it fails...
	mockFileSystem.VirtualFunction_MkdirAll = func(targetFolderPath string) error {
		return errors.New("Simulated MkdirAll failure")
	}

	// When ...
	err := createProject(mockFileSystem, "my.test.pkg", isObrProjectRequired, forceOverwrite)

	// Then...
	assert.NotNil(t, err, "Sumulated error didn't bubble up to the top.")
	assert.Contains(t, err.Error(), "Simulated MkdirAll failure", "Failed to return reason for failure.")

}

func TestCreateProjectErrorsWhenWriteTextFileFails(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewOverridableMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false

	// Over-ride the default WriteTextFile function so that it fails...
	mockFileSystem.VirtualFunction_WriteTextFile = func(targetFilePath string, desiredContents string) error {
		return errors.New("Simulated WriteTextFile failure")
	}

	// When ...
	err := createProject(mockFileSystem, "my.test.pkg", isObrProjectRequired, forceOverwrite)

	// Then...
	assert.NotNil(t, err, "Sumulated error didn't bubble up to the top.")
	assert.Contains(t, err.Error(), "Simulated WriteTextFile failure", "Failed to return reason for failure.")
}

func TestCreateProjectPomFileAlreadyExistsNoForceOverwrite(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := false
	isObrProjectRequired := false
	testPackageName := "my.test.pkg"

	// Create a pom.xml file already...
	mockFileSystem.MkdirAll(testPackageName)
	mockFileSystem.WriteTextFile(testPackageName+"/pom.xml", "dummy test pom.xml")

	// When ...
	err := createProject(mockFileSystem, testPackageName, isObrProjectRequired, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1036E")
}

func TestCreateProjectPomFileAlreadyExistsWithForceOverwrite(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	isObrProjectRequired := false
	forceOverwrite := true
	testPackageName := "my.test.pkg"

	// Create a pom.xml file already...
	mockFileSystem.MkdirAll(testPackageName)
	mockFileSystem.WriteTextFile(testPackageName+"/pom.xml", "dummy test pom.xml")

	// When ...
	err := createProject(mockFileSystem, testPackageName, isObrProjectRequired, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, err.Error())
	}

	// Check that the pom with decent inputs has over-written the dummy test pom contents.
	parentPomXmlFileExists, err := mockFileSystem.Exists("my.test.pkg/pom.xml")
	assert.Nil(t, err)
	assert.True(t, parentPomXmlFileExists, "Parent folder pom.xml was not created.")

	text, err := mockFileSystem.ReadTextFile("my.test.pkg/pom.xml")
	assert.Nil(t, err)
	assert.Contains(t, text, "<groupId>my.test.pkg</groupId>", "parent pom.xml didn't substitute the group id")
	assert.Contains(t, text, "<artifactId>my.test.pkg</artifactId>", "parent pom.xml didn't substitute the artifact id")

}

func TestCanCreateProjectGoldenPathWithOBR(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := true

	// When ...
	err := createProject(mockFileSystem, "my.test.pkg", isObrProjectRequired, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assertParentFolderAndContentsCreated(t, mockFileSystem, isObrProjectRequired)
	assertTestFolderAndContentsCreatedOk(t, mockFileSystem)
	assertOBRFOlderAndContentsCreatedOK(t, mockFileSystem)
}

func assertOBRFOlderAndContentsCreatedOK(t *testing.T, mockFileSystem utils.FileSystem) {
	testFolderExists, err := mockFileSystem.DirExists("my.test.pkg/my.test.pkg.obr")
	assert.Nil(t, err)
	assert.True(t, testFolderExists, "Test folder was not created.")

	expectedPomFilePath := "my.test.pkg/my.test.pkg.obr/pom.xml"
	testPomXmlFileExists, err := mockFileSystem.Exists(expectedPomFilePath)
	assert.Nil(t, err)
	assert.True(t, testPomXmlFileExists, "Test folder pom.xml was not created.")

	text, err := mockFileSystem.ReadTextFile(expectedPomFilePath)
	assert.Nil(t, err)
	assert.Contains(t, text, "<groupId>my.test.pkg</groupId>", "Test folder pom.xml didn't substitute the group id")
	assert.Contains(t, text, "<artifactId>my.test.pkg.obr</artifactId>", "Test folder pom.xml didn't substitute the artifact id")
}
