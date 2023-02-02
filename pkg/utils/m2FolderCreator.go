/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"embed"
)

//  galasaErrors "github.com/galasa.dev/cli/pkg/errors"

func InitialiseM2Folder(fileSystem FileSystem, embeddedFileSystem embed.FS) error {

	var err error
	var userHomeDir string

	fileGenerator := NewFileGenerator(fileSystem, embeddedFileSystem)

	userHomeDir, err = fileSystem.GetUserHomeDir()
	if err == nil {
		m2Dir := userHomeDir + FILE_SYSTEM_PATH_SEPARATOR + ".m2"

		if err == nil {
			err = fileGenerator.CreateFolder(m2Dir)
		}
		if err == nil {
			err = createSettingsXMLFile(fileGenerator, m2Dir)
		}

	}

	return err
}

func createSettingsXMLFile(fileGenerator *FileGenerator, m2Dir string) error {

	targetPath := m2Dir + FILE_SYSTEM_PATH_SEPARATOR + "settings.xml"

	xmlFile := GeneratedFileDef{
		FileType:                 "xml",
		TargetFilePath:           targetPath,
		EmbeddedTemplateFilePath: "templates/m2/settings.xml",
		TemplateParameters:       nil,
	}

	err := fileGenerator.CreateFile(xmlFile, false, false)

	return err
}
