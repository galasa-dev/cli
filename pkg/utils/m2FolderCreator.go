/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"embed"
	"log"
	"strings"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
)

const (
	MAVEN_REPO_URL_GALASA_BLEEDING_EDGE = "https://development.galasa.dev/main/maven-repo/obr"
	MAVEN_REPO_URL_MAVEN_CENTRAL        = "https://repo.maven.apache.org/maven2"
)

func InitialiseM2Folder(fileSystem FileSystem, embeddedFileSystem embed.FS, isDevelopment bool) error {

	var err error
	var userHomeDir string

	fileGenerator := NewFileGenerator(fileSystem, embeddedFileSystem)

	userHomeDir, err = fileSystem.GetUserHomeDirPath()
	if err == nil {
		m2Dir := userHomeDir + fileSystem.GetFilePathSeparator() + ".m2"
		err = fileGenerator.CreateFolder(m2Dir)

		if err == nil {
			err = createSettingsXMLFile(fileGenerator, fileSystem, m2Dir, isDevelopment)
		}
	}

	return err
}

func createSettingsXMLFile(
	fileGenerator *FileGenerator,
	fileSystem FileSystem,
	m2Dir string,
	isDevelopment bool,
) error {

	targetPath := m2Dir + fileSystem.GetFilePathSeparator() + "settings.xml"

	type M2SettingsXmlTemplateParams struct {
		IsDevelopment bool
	}

	templateParameters := M2SettingsXmlTemplateParams{
		IsDevelopment: isDevelopment,
	}

	xmlFile := GeneratedFileDef{
		FileType:                 "xml",
		TargetFilePath:           targetPath,
		EmbeddedTemplateFilePath: "templates/m2/settings.xml",
		TemplateParameters:       templateParameters,
	}

	settingsFileExists, err := fileSystem.Exists(targetPath)
	if err == nil {
		if settingsFileExists {
			err = warnIfFileDoesntContainGalasaOBRMavenRepository(fileSystem, targetPath)
		} else {
			err = fileGenerator.CreateFile(xmlFile, false, false)
		}
	}

	return err
}

func warnIfFileDoesntContainGalasaOBRMavenRepository(fileSystem FileSystem, filePath string) error {

	content, err := fileSystem.ReadTextFile(filePath)
	if err == nil {
		containsBleedingEdgeUrl := strings.Contains(content, MAVEN_REPO_URL_GALASA_BLEEDING_EDGE)
		containsMavenCentral := strings.Contains(content, MAVEN_REPO_URL_MAVEN_CENTRAL)

		log.Printf("Checking %s for obr maven repository references. containsBleedingEdgeUrl:%v containsMavenCentral:%v",
			filePath, containsBleedingEdgeUrl, containsMavenCentral)
		if !(containsBleedingEdgeUrl || containsMavenCentral) {
			// Neither of our magic urls are in the settings.xml , so the galasa obr probably can't be found.
			warningMessage := galasaErrors.NewGalasaError(galasaErrors.GALASA_WARNING_MAVEN_NO_GALASA_OBR_REPO,
				MAVEN_REPO_URL_MAVEN_CENTRAL, MAVEN_REPO_URL_GALASA_BLEEDING_EDGE).Error()
			fileSystem.OutputWarningMessage(warningMessage)
		}
	}
	return err
}
