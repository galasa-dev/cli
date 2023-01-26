/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateProjectFailsIfPackageNameInvalid(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false
	featureNamesCommandSeparatedList := "test"

	// When ...
	err := createProject(mockFileSystem, "very.INVALID_PACKAGE_NAME.very",
		featureNamesCommandSeparatedList, isObrProjectRequired, forceOverwrite)

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
	featureNamesCommandSeparatedList := "test"

	// When ...
	err := createProject(mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList, isObrProjectRequired, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, "Golden path should not return an error. %s", err.Error())
	}

	assertParentFolderAndContentsCreated(t, mockFileSystem, isObrProjectRequired)
	assertTestFolderAndContentsCreatedOk(t, mockFileSystem, "test")
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

		// Make sure the OBR folder does not exist.
		expectedObrFolderPath := packageName + "/" + packageName + ".obr"
		var obrFolderExists bool
		obrFolderExists, _ = mockFileSystem.DirExists(expectedObrFolderPath)
		assert.False(t, obrFolderExists, "OBR folder exists, when it should not!")

		// OBR should not be mentioned in the parent pom.xml
	}
}

func assertTestFolderAndContentsCreatedOk(t *testing.T, mockFileSystem utils.FileSystem, featureName string) {

	testFolderExists, err := mockFileSystem.DirExists("my.test.pkg/my.test.pkg." + featureName)
	assert.Nil(t, err)
	assert.True(t, testFolderExists, "Test"+featureName+" folder was not created.")

	expectedPomFilePath := "my.test.pkg/my.test.pkg." + featureName + "/pom.xml"
	testPomXmlFileExists, err := mockFileSystem.Exists(expectedPomFilePath)
	assert.Nil(t, err)
	assert.True(t, testPomXmlFileExists, "Test folder pom.xml was not created for feature."+featureName)

	text, err := mockFileSystem.ReadTextFile(expectedPomFilePath)
	assert.Nil(t, err)
	assert.Contains(t, text, "<groupId>my.test.pkg</groupId>", "Test folder pom.xml didn't substitute the group id")
	assert.Contains(t, text, "<artifactId>my.test.pkg.test</artifactId>", "Test folder pom.xml didn't substitute the artifact id")

	testSrcFolderExists, err := mockFileSystem.DirExists("my.test.pkg/my.test.pkg.test/src/main/java/my/test/pkg/test")
	assert.Nil(t, err)
	assert.True(t, testSrcFolderExists, "Test src folder was not created.")

	// Examine the test java file generated.
	expectedJavaFilePath := "my.test.pkg/my.test.pkg." + featureName + "/src/main/java/my/test/pkg/" + featureName + "/Test" + utils.UppercaseFirstLetter(featureName) + ".java"
	assertJavaFileWasGenerated(t, mockFileSystem, expectedJavaFilePath, "my.test.pkg")

	// Examine the extended test java file generated.
	expectedJavaFilePath = "my.test.pkg/my.test.pkg." + featureName + "/src/main/java/my/test/pkg/" + featureName + "/Test" + utils.UppercaseFirstLetter(featureName) + "Extended.java"
	assertJavaFileWasGenerated(t, mockFileSystem, expectedJavaFilePath, "my.test.pkg")

	// Examine the resources file generated.
	expectedTextFilePath := "my.test.pkg/my.test.pkg." + featureName + "/src/main/resources/textfiles/sampleText.txt"
	isTestResourcesTextFileExists, err := mockFileSystem.Exists(expectedTextFilePath)
	assert.Nil(t, err)
	assert.True(t, isTestResourcesTextFileExists, "Test text resource file was not created.")
}

func assertJavaFileWasGenerated(t *testing.T, mockFileSystem utils.FileSystem, expectedJavaFilePath string, packageName string) {
	testJavaFileExists, err := mockFileSystem.Exists(expectedJavaFilePath)
	assert.Nil(t, err)
	assert.True(t, testJavaFileExists, "Test java file was not created.")

	var text string
	text, err = mockFileSystem.ReadTextFile(expectedJavaFilePath)
	assert.Nil(t, err)
	assert.Contains(t, text, "package "+packageName, "Test java file didn't substitute the java package")
}

func TestCreateProjectErrorsWhenMkAllDirsFails(t *testing.T) {

	// Given...
	mockFileSystem := utils.NewOverridableMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false
	featureNamesCommandSeparatedList := "test"

	// Over-ride the default MkdirAll function so that it fails...
	mockFileSystem.VirtualFunction_MkdirAll = func(targetFolderPath string) error {
		return errors.New("Simulated MkdirAll failure")
	}

	// When ...
	err := createProject(mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList, isObrProjectRequired, forceOverwrite)

	// Then...
	assert.NotNil(t, err, "Sumulated error didn't bubble up to the top.")
	assert.Contains(t, err.Error(), "Simulated MkdirAll failure", "Failed to return reason for failure.")

}

func TestCreateProjectErrorsWhenWriteTextFileFails(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewOverridableMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false
	featureNamesCommandSeparatedList := "test"

	// Over-ride the default WriteTextFile function so that it fails...
	mockFileSystem.VirtualFunction_WriteTextFile = func(targetFilePath string, desiredContents string) error {
		return errors.New("Simulated WriteTextFile failure")
	}

	// When ...
	err := createProject(mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList, isObrProjectRequired, forceOverwrite)

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
	featureNamesCommandSeparatedList := "test"

	// Create a pom.xml file already...
	mockFileSystem.MkdirAll(testPackageName)
	mockFileSystem.WriteTextFile(testPackageName+"/pom.xml", "dummy test pom.xml")

	// When ...
	err := createProject(mockFileSystem, testPackageName, featureNamesCommandSeparatedList, isObrProjectRequired, forceOverwrite)

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
	featureNamesCommandSeparatedList := "test"

	// Create a pom.xml file already...
	mockFileSystem.MkdirAll(testPackageName)
	mockFileSystem.WriteTextFile(testPackageName+"/pom.xml", "dummy test pom.xml")

	// When ...
	err := createProject(mockFileSystem, testPackageName, featureNamesCommandSeparatedList, isObrProjectRequired, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, err.Error())
	}

	// Check that the pom with decent inputs has overwritten the dummy test pom contents.
	parentPomXmlFileExists, err := mockFileSystem.Exists("my.test.pkg/pom.xml")
	assert.Nil(t, err)
	assert.True(t, parentPomXmlFileExists, "Parent folder pom.xml was not created.")

	text, err := mockFileSystem.ReadTextFile("my.test.pkg/pom.xml")
	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(text, "<?xml"), "pom template is expanding as HTML by accident!")
	assert.Contains(t, text, "<groupId>my.test.pkg</groupId>", "parent pom.xml didn't substitute the group id")
	assert.Contains(t, text, "<artifactId>my.test.pkg</artifactId>", "parent pom.xml didn't substitute the artifact id")

}

func TestCanCreateProjectGoldenPathWithOBR(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := true
	featureNamesCommandSeparatedList := "test"

	// When ...
	err := createProject(mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList, isObrProjectRequired, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assertParentFolderAndContentsCreated(t, mockFileSystem, isObrProjectRequired)
	assertTestFolderAndContentsCreatedOk(t, mockFileSystem, "test")
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

func TestCreateProjectWithTwoFeaturesWorks(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := false
	isObrProjectRequired := false
	testPackageName := "my.test.pkg"
	featureNamesCommandSeparatedList := "account,payee"

	// When ...
	err := createProject(mockFileSystem, testPackageName,
		featureNamesCommandSeparatedList, isObrProjectRequired, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, "err should not have been set. "+err.Error())
	}

	isAccountModuleExists, _ := mockFileSystem.DirExists(testPackageName + "/" + testPackageName + ".account")
	assert.True(t, isAccountModuleExists, "account feature module does not exist.")

	isPayeeModuleExists, _ := mockFileSystem.DirExists(testPackageName + "/" + testPackageName + ".payee")
	assert.True(t, isPayeeModuleExists, "payee feature module does not exist.")
}

func TestCreateProjectWithInvalidFeaturesFails(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()
	forceOverwrite := false
	isObrProjectRequired := false
	testPackageName := "my.test.pkg"
	featureNamesCommandSeparatedList := "Account"

	// When ...
	err := createProject(mockFileSystem, testPackageName,
		featureNamesCommandSeparatedList, isObrProjectRequired, forceOverwrite)

	// Then...
	// Should have created a folder for the parent package.
	assert.NotNil(t, err, "err should have been set!")
	assert.Contains(t, err.Error(), "GAL1045E")

}
