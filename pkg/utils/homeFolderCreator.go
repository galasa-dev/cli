/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"embed"

	"github.com/galasa.dev/cli/pkg/embedded"
)

//  galasaErrors "github.com/galasa.dev/cli/pkg/errors"

func InitialiseGalasaHomeFolder(fileSystem FileSystem, embeddedFileSystem embed.FS) error {

	var err error
	var userHomeDir string

	fileGenerator := NewFileGenerator(fileSystem, embeddedFileSystem)

	userHomeDir, err = fileSystem.GetUserHomeDir()
	if err == nil {
		galasaHomeDir := userHomeDir + FILE_SYSTEM_PATH_SEPARATOR + ".galasa"
		err = fileGenerator.CreateFolder(galasaHomeDir)

		if err == nil {
			err = createLibDirAndContent(fileGenerator, galasaHomeDir+FILE_SYSTEM_PATH_SEPARATOR+"lib")
		}

		if err == nil {
			err = createBootstrapPropertiesFile(fileGenerator, galasaHomeDir)
		}

		if err == nil {
			err = createOverridesPropertiesFile(fileGenerator, galasaHomeDir)
		}

		if err == nil {
			err = createCPSPropertiesFile(fileGenerator, galasaHomeDir)
		}

		if err == nil {
			err = createDSSPropertiesFile(fileGenerator, galasaHomeDir)
		}

		if err == nil {
			err = createCredentialsPropertiesFile(fileGenerator, galasaHomeDir)
		}

	}

	return err
}

func createBootstrapPropertiesFile(fileGenerator *FileGenerator, galasaHomeDir string) error {

	targetPath := galasaHomeDir + FILE_SYSTEM_PATH_SEPARATOR + "bootstrap.properties"

	propertyFile := GeneratedFileDef{
		FileType:                 "properties",
		TargetFilePath:           targetPath,
		EmbeddedTemplateFilePath: "templates/galasahome/bootstrap.properties",
		TemplateParameters:       nil,
	}

	err := fileGenerator.CreateFile(propertyFile, false, false)

	return err
}

func createCPSPropertiesFile(fileGenerator *FileGenerator, galasaHomeDir string) error {

	targetPath := galasaHomeDir + FILE_SYSTEM_PATH_SEPARATOR + "cps.properties"

	propertyFile := GeneratedFileDef{
		FileType:                 "properties",
		TargetFilePath:           targetPath,
		EmbeddedTemplateFilePath: "templates/galasahome/cps.properties",
		TemplateParameters:       nil,
	}

	err := fileGenerator.CreateFile(propertyFile, false, false)

	return err
}

func createDSSPropertiesFile(fileGenerator *FileGenerator, galasaHomeDir string) error {

	targetPath := galasaHomeDir + FILE_SYSTEM_PATH_SEPARATOR + "dss.properties"

	propertyFile := GeneratedFileDef{
		FileType:                 "properties",
		TargetFilePath:           targetPath,
		EmbeddedTemplateFilePath: "templates/galasahome/dss.properties",
		TemplateParameters:       nil,
	}

	err := fileGenerator.CreateFile(propertyFile, false, false)

	return err
}

func createCredentialsPropertiesFile(fileGenerator *FileGenerator, galasaHomeDir string) error {

	targetPath := galasaHomeDir + FILE_SYSTEM_PATH_SEPARATOR + "credentials.properties"

	propertyFile := GeneratedFileDef{
		FileType:                 "properties",
		TargetFilePath:           targetPath,
		EmbeddedTemplateFilePath: "templates/galasahome/credentials.properties",
		TemplateParameters:       nil,
	}

	err := fileGenerator.CreateFile(propertyFile, false, false)

	return err
}

func createOverridesPropertiesFile(fileGenerator *FileGenerator, galasaHomeDir string) error {

	targetPath := galasaHomeDir + FILE_SYSTEM_PATH_SEPARATOR + "overrides.properties"

	propertyFile := GeneratedFileDef{
		FileType:                 "properties",
		TargetFilePath:           targetPath,
		EmbeddedTemplateFilePath: "templates/galasahome/overrides.properties",
		TemplateParameters:       nil,
	}

	err := fileGenerator.CreateFile(
		propertyFile,
		false, // Don't force overwrite
		false) // Don't error if the file already exists.

	return err
}

func createLibDirAndContent(fileGenerator *FileGenerator, galasaLibDir string) error {

	err := fileGenerator.CreateFolder(galasaLibDir)

	if err == nil {
		galasaVersion := embedded.GetGalasaVersion()
		galasaVersionLibDir := galasaLibDir + FILE_SYSTEM_PATH_SEPARATOR + galasaVersion
		err = fileGenerator.CreateFolder(galasaVersionLibDir)

		if err == nil {
			bootJarVersion := embedded.GetBootJarVersion()

			installedBootJar := GeneratedFileDef{
				FileType: "jar",
				TargetFilePath: galasaVersionLibDir + FILE_SYSTEM_PATH_SEPARATOR +
					"galasa-boot-" + bootJarVersion + ".jar",
				EmbeddedTemplateFilePath: "templates/galasahome/lib/galasa-boot-" + bootJarVersion + ".jar",
				TemplateParameters:       nil,
			}

			err = fileGenerator.CreateFile(
				installedBootJar,
				false, // don't force overwrite
				false) // don't error if it already exists.
		}
	}

	return err
}

func GetGalasaBootJarPath(fileSystem FileSystem) (string, error) {
	var galasaBootJarPath string = ""
	var userHomeDir, err = fileSystem.GetUserHomeDir()
	if err == nil {
		galasaVersion := embedded.GetGalasaVersion()
		bootJarVersion := embedded.GetBootJarVersion()

		galasaBootJarPath = userHomeDir +
			FILE_SYSTEM_PATH_SEPARATOR + ".galasa" +
			FILE_SYSTEM_PATH_SEPARATOR + "lib" +
			FILE_SYSTEM_PATH_SEPARATOR + galasaVersion +
			FILE_SYSTEM_PATH_SEPARATOR + "galasa-boot-" + bootJarVersion + ".jar"
	}
	return galasaBootJarPath, err
}
