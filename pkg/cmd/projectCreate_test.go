/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateProjectFailsIfPackageNameInvalid(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false
	featureNamesCommandSeparatedList := "test"
	maven := true
	gradle := false
	isDevelopment := false

	// When ...
	err := createProject(mockFileSystem, "very.INVALID_PACKAGE_NAME.very",
		featureNamesCommandSeparatedList, isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	// Should have created a folder for the parent package.
	assert.NotNil(t, err, "Should return an error if java package name is invalid")
	assert.Contains(t, err.Error(), "GAL1037E:", "Wrong error message reported.")
}

func TestCanCreateProjectGoldenPathNoOBR(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false
	featureNamesCommandSeparatedList := "test"
	maven := true
	gradle := false
	isDevelopment := false

	// When ...
	err := createProject(
		mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, "Golden path should not return an error. %s", err.Error())
	}

	assertParentFolderAndContentsCreated(t, mockFileSystem, isObrProjectRequired, maven, gradle)
	assertTestFolderAndContentsCreatedOk(t, mockFileSystem, "test", maven, gradle)
}

func assertParentFolderAndContentsCreated(t *testing.T, mockFileSystem files.FileSystem, isObrProjectRequired bool, isMaven bool, isGradle bool) {
	parentFolderExists, err := mockFileSystem.DirExists("my.test.pkg")
	assert.Nil(t, err)
	assert.True(t, parentFolderExists, "Parent folder was not created.")

	if isMaven {
		parentPomXmlFileExists, err := mockFileSystem.Exists("my.test.pkg/pom.xml")
		assert.Nil(t, err)
		assert.True(t, parentPomXmlFileExists, "Parent folder pom.xml was not created.")

		gitIgnoreFileExists, err := mockFileSystem.Exists("my.test.pkg/.gitignore")
		assert.Nil(t, err)
		assert.True(t, gitIgnoreFileExists, "Parent folder .gitignore was not created.")

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

	if isGradle {
		parentSettingsGradleFileExists, err := mockFileSystem.Exists("my.test.pkg/settings.gradle")
		assert.Nil(t, err)
		assert.True(t, parentSettingsGradleFileExists, "Parent folder settings.gradle was not created.")

		text, err := mockFileSystem.ReadTextFile("my.test.pkg/settings.gradle")
		assert.Nil(t, err)
		assert.Contains(t, text, "include 'my.test.pkg.test'", "parent settings.gradle didn't have a test module included")

		if isObrProjectRequired {
			assert.Contains(t, text, "include 'my.test.pkg.obr'", "parent settings.gradle didn't have an obr module included")
		} else {
			assert.NotContains(t, text, "include 'my.test.pkg.obr'", "parent settings.gradle should not have an obr module included")

			// Make sure the OBR folder does not exist.
			expectedObrFolderPath := packageName + "/" + packageName + ".obr"
			var obrFolderExists bool
			obrFolderExists, _ = mockFileSystem.DirExists(expectedObrFolderPath)
			assert.False(t, obrFolderExists, "OBR folder exists, when it should not!")
		}
	}
}

func assertTestFolderAndContentsCreatedOk(t *testing.T, mockFileSystem files.FileSystem, featureName string, isMaven bool, isGradle bool) {

	testFolderExists, err := mockFileSystem.DirExists("my.test.pkg/my.test.pkg." + featureName)
	assert.Nil(t, err)
	assert.True(t, testFolderExists, "Test"+featureName+" folder was not created.")

	if isMaven {
		expectedPomFilePath := "my.test.pkg/my.test.pkg." + featureName + "/pom.xml"
		testPomXmlFileExists, err := mockFileSystem.Exists(expectedPomFilePath)
		assert.Nil(t, err)
		assert.True(t, testPomXmlFileExists, "Test folder pom.xml was not created for feature."+featureName)

		text, err := mockFileSystem.ReadTextFile(expectedPomFilePath)
		assert.Nil(t, err)
		assert.Contains(t, text, "<groupId>my.test.pkg</groupId>", "Test folder pom.xml didn't substitute the group id")
		assert.Contains(t, text, "<artifactId>my.test.pkg.test</artifactId>", "Test folder pom.xml didn't substitute the artifact id")
	}

	if isGradle {
		expectedBuildGradleFilePath := "my.test.pkg/my.test.pkg." + featureName + "/build.gradle"
		testBuildGradleFileExists, err := mockFileSystem.Exists(expectedBuildGradleFilePath)
		assert.Nil(t, err)
		assert.True(t, testBuildGradleFileExists, "Test folder build.gradle was not created for feature."+featureName)

		expectedBndFilePath := "my.test.pkg/my.test.pkg." + featureName + "/bnd.bnd"
		testBndFileExists, err := mockFileSystem.Exists(expectedBndFilePath)
		assert.Nil(t, err)
		assert.True(t, testBndFileExists, "Test folder bnd.bnd was not created for feature."+featureName)

		buildGradleText, err := mockFileSystem.ReadTextFile(expectedBuildGradleFilePath)
		assert.Nil(t, err)
		assert.Contains(t, buildGradleText, "group = 'my.test.pkg'", "Test folder build.gradle didn't substitute the group id")

		bndFileText, err := mockFileSystem.ReadTextFile(expectedBndFilePath)
		assert.Nil(t, err)
		assert.Contains(t, bndFileText, "Bundle-Name: my.test.pkg", "Test folder bnd.bnd didn't substitute the bundle name")
	}

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

func assertJavaFileWasGenerated(t *testing.T, mockFileSystem files.FileSystem, expectedJavaFilePath string, packageName string) {
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
	mockFileSystem := files.NewOverridableMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false
	featureNamesCommandSeparatedList := "test"
	maven := true
	gradle := false
	isDevelopment := false

	// Over-ride the default MkdirAll function so that it fails...
	mockFileSystem.VirtualFunction_MkdirAll = func(targetFolderPath string) error {
		return errors.New("Simulated MkdirAll failure")
	}

	// When ...
	err := createProject(
		mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	assert.NotNil(t, err, "Sumulated error didn't bubble up to the top.")
	assert.Contains(t, err.Error(), "Simulated MkdirAll failure", "Failed to return reason for failure.")

}

func TestCreateProjectErrorsWhenWriteTextFileFails(t *testing.T) {
	// Given...
	mockFileSystem := files.NewOverridableMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false
	featureNamesCommandSeparatedList := "test"
	maven := true
	gradle := false
	isDevelopment := false

	// Over-ride the default WriteTextFile function so that it fails...
	mockFileSystem.VirtualFunction_WriteTextFile = func(targetFilePath string, desiredContents string) error {
		return errors.New("Simulated WriteTextFile failure")
	}

	// When ...
	err := createProject(
		mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	assert.NotNil(t, err, "Sumulated error didn't bubble up to the top.")
	assert.Contains(t, err.Error(), "Simulated WriteTextFile failure", "Failed to return reason for failure.")
}

func TestCreateProjectPomFileAlreadyExistsNoForceOverwrite(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := false
	isObrProjectRequired := false
	testPackageName := "my.test.pkg"
	featureNamesCommandSeparatedList := "test"
	maven := true
	gradle := false
	isDevelopment := false

	// Create a pom.xml file already...
	mockFileSystem.MkdirAll(testPackageName)
	mockFileSystem.WriteTextFile(testPackageName+"/pom.xml", "dummy test pom.xml")

	// When ...
	err := createProject(
		mockFileSystem, testPackageName, featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	// Should have created a folder for the parent package.
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1036E")
}

func TestCreateProjectPomFileAlreadyExistsWithForceOverwrite(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	isObrProjectRequired := false
	forceOverwrite := true
	testPackageName := "my.test.pkg"
	featureNamesCommandSeparatedList := "test"
	maven := true
	gradle := false
	isDevelopment := false

	// Create a pom.xml file already...
	mockFileSystem.MkdirAll(testPackageName)
	mockFileSystem.WriteTextFile(testPackageName+"/pom.xml", "dummy test pom.xml")

	// When ...
	err := createProject(
		mockFileSystem, testPackageName, featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

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
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := true
	featureNamesCommandSeparatedList := "test"
	maven := true
	gradle := false
	isDevelopment := false

	// When ...
	err := createProject(
		mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assertParentFolderAndContentsCreated(t, mockFileSystem, isObrProjectRequired, maven, gradle)
	assertTestFolderAndContentsCreatedOk(t, mockFileSystem, "test", maven, gradle)
	assertOBRFOlderAndContentsCreatedOK(t, mockFileSystem, maven, gradle)
}

func assertOBRFOlderAndContentsCreatedOK(t *testing.T, mockFileSystem files.FileSystem, isMaven bool, isGradle bool) {
	testFolderExists, err := mockFileSystem.DirExists("my.test.pkg/my.test.pkg.obr")
	assert.Nil(t, err)
	assert.True(t, testFolderExists, "Test folder was not created.")

	if isMaven {
		expectedPomFilePath := "my.test.pkg/my.test.pkg.obr/pom.xml"
		testPomXmlFileExists, err := mockFileSystem.Exists(expectedPomFilePath)
		assert.Nil(t, err)
		assert.True(t, testPomXmlFileExists, "Test folder pom.xml was not created.")

		text, err := mockFileSystem.ReadTextFile(expectedPomFilePath)
		assert.Nil(t, err)
		assert.Contains(t, text, "<groupId>my.test.pkg</groupId>", "Test folder pom.xml didn't substitute the group id")
		assert.Contains(t, text, "<artifactId>my.test.pkg.obr</artifactId>", "Test folder pom.xml didn't substitute the artifact id")
	}

	if isGradle {
		expectedBuildGradleFilePath := "my.test.pkg/my.test.pkg.obr/build.gradle"
		testBuildGradleFileExists, err := mockFileSystem.Exists(expectedBuildGradleFilePath)
		assert.Nil(t, err)
		assert.True(t, testBuildGradleFileExists, "Test folder build.gradle was not created.")

		text, err := mockFileSystem.ReadTextFile(expectedBuildGradleFilePath)
		assert.Nil(t, err)
		assert.Contains(t, text, "group = 'my.test.pkg'", "Test folder build.gradle didn't substitute the group id")
	}
}

func TestCreateProjectWithTwoFeaturesWorks(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := false
	isObrProjectRequired := false
	testPackageName := "my.test.pkg"
	featureNamesCommandSeparatedList := "account,payee"
	maven := true
	gradle := false
	isDevelopment := false

	// When ...
	err := createProject(
		mockFileSystem, testPackageName,
		featureNamesCommandSeparatedList, isObrProjectRequired,
		forceOverwrite, maven, gradle, isDevelopment)

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
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := false
	isObrProjectRequired := false
	testPackageName := "my.test.pkg"
	featureNamesCommandSeparatedList := "Account"
	maven := true
	gradle := false
	isDevelopment := false

	// When ...
	err := createProject(mockFileSystem, testPackageName,
		featureNamesCommandSeparatedList, isObrProjectRequired,
		forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	// Should have created a folder for the parent package.
	assert.NotNil(t, err, "err should have been set!")
	assert.Contains(t, err.Error(), "GAL1045E")

}

func TestCanCreateGradleProjectWithNoOBR(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false
	featureNamesCommandSeparatedList := "test"
	maven := false
	gradle := true
	isDevelopment := false

	// When ...
	err := createProject(mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, "Creating a Maven project should not return an error. %s", err.Error())
	}

	assertParentFolderAndContentsCreated(t, mockFileSystem, isObrProjectRequired, maven, gradle)
	assertTestFolderAndContentsCreatedOk(t, mockFileSystem, "test", maven, gradle)
}

func TestCanCreateGradleProjectWithOBR(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := true
	featureNamesCommandSeparatedList := "test"
	maven := false
	gradle := true
	isDevelopment := false

	// When ...
	err := createProject(
		mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assertParentFolderAndContentsCreated(t, mockFileSystem, isObrProjectRequired, maven, gradle)
	assertTestFolderAndContentsCreatedOk(t, mockFileSystem, "test", maven, gradle)
	assertOBRFOlderAndContentsCreatedOK(t, mockFileSystem, maven, gradle)
}

func TestCanCreateMavenAndGradleProject(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false
	featureNamesCommandSeparatedList := "test"
	maven := true
	gradle := true
	isDevelopment := false

	// When ...
	err := createProject(
		mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assertParentFolderAndContentsCreated(t, mockFileSystem, isObrProjectRequired, maven, gradle)
	assertTestFolderAndContentsCreatedOk(t, mockFileSystem, "test", maven, gradle)
}

func TestCreateProjectDefaultsToMavenProject(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := false
	featureNamesCommandSeparatedList := "test"
	maven := false
	gradle := false
	expectMaven := true
	isDevelopment := false

	// When ...
	err := createProject(
		mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, err.Error())
	}

	assertParentFolderAndContentsCreated(t, mockFileSystem, isObrProjectRequired, expectMaven, gradle)
	assertTestFolderAndContentsCreatedOk(t, mockFileSystem, "test", expectMaven, gradle)
}

func TestCanCreateGradleProjectNonDevelopmentModeGeneratesCommentedOutMavenRepoReference(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := true
	featureNamesCommandSeparatedList := "test"
	maven := false
	gradle := true
	isDevelopment := false

	// When ...
	err := createProject(
		mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, err.Error())
	}

	settingsGradleText, err := mockFileSystem.ReadTextFile("my.test.pkg/settings.gradle")
	assert.Nil(t, err)
	assert.Contains(t, settingsGradleText, "//    url 'https://development.galasa.dev/main/maven-repo/obr'", "parent settings.gradle didn't have a commented-out bleeding edge repo ref.")

	buildGradleText, err := mockFileSystem.ReadTextFile("my.test.pkg/my.test.pkg.test/build.gradle")
	assert.Nil(t, err)
	assert.Contains(t, buildGradleText, "//    url 'https://development.galasa.dev/main/maven-repo/obr'", "child build.gradle didn't have a commented-out bleeding edge repo ref.")

}

func TestCanCreateGradleProjectDevelopmentModeGeneratesMavenRepoReference(t *testing.T) {
	// Given...
	mockFileSystem := files.NewMockFileSystem()
	forceOverwrite := true
	isObrProjectRequired := true
	featureNamesCommandSeparatedList := "test"
	maven := false
	gradle := true
	isDevelopment := true

	// When ...
	err := createProject(
		mockFileSystem, "my.test.pkg", featureNamesCommandSeparatedList,
		isObrProjectRequired, forceOverwrite, maven, gradle, isDevelopment)

	// Then...
	// Should have created a folder for the parent package.
	if err != nil {
		assert.Fail(t, err.Error())
	}

	settingsGradleText, err := mockFileSystem.ReadTextFile("my.test.pkg/settings.gradle")
	assert.Nil(t, err)
	assert.Contains(t, settingsGradleText, "           url 'https://development.galasa.dev/main/maven-repo/obr'", "parent settings.gradle didn't have an uncommented bleeding edge repo ref.")

	buildGradleText, err := mockFileSystem.ReadTextFile("my.test.pkg/my.test.pkg.test/build.gradle")
	assert.Nil(t, err)
	assert.Contains(t, buildGradleText, "       url 'https://development.galasa.dev/main/maven-repo/obr'", "child build.gradle didn't have an uncommented bleeding edge repo ref.")
}
