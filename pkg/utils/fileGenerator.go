/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"bytes"
	"embed"
	"log"
	"text/template"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
)

//---------------------------------------------------------------------------------------------------
// File generation types
//---------------------------------------------------------------------------------------------------

type GeneratedFileDef struct {
	FileType                 string
	TargetFilePath           string
	EmbeddedTemplateFilePath string
	TemplateParameters       interface{}
}

type FileGenerator struct {
	fileSystem         FileSystem
	embeddedFileSystem embed.FS
}

//-------------------------------------------------------------------------------------------------
// Public functions.
//-------------------------------------------------------------------------------------------------

func NewFileGenerator(fileSystem FileSystem, embeddedFileSystem embed.FS) *FileGenerator {
	fileGenerator := &FileGenerator{fileSystem: fileSystem, embeddedFileSystem: embeddedFileSystem}
	return fileGenerator
}

func (generator *FileGenerator) CreateFolder(targetFolderPath string) error {
	err := generator.fileSystem.MkdirAll(targetFolderPath)
	return err
}

// createFile creates a file on the file system.
// If forceOverwrite is false, and there is already a file there, then an error will occur.
func (generator *FileGenerator) CreateFile(
	generatedFile GeneratedFileDef,
	forceOverwrite bool) error {

	log.Printf("Creating file of type %s at %s\n", generatedFile.FileType, generatedFile.TargetFilePath)

	err := generator.checkAllowedToWrite(generatedFile.TargetFilePath, forceOverwrite)
	if err == nil {
		var template *template.Template
		template, err = generator.loadEmbeddedTemplate(generatedFile.EmbeddedTemplateFilePath)
		if err == nil {
			var fileContents string
			fileContents, err = generator.substituteParametersIntoTemplate(template, generatedFile.TemplateParameters)
			if err == nil {
				// Write it out to the target file.
				err = generator.fileSystem.WriteTextFile(generatedFile.TargetFilePath, fileContents)
			}
		}
	}
	if err == nil {
		log.Printf("Created file %s OK.", generatedFile.TargetFilePath)
	}
	return err
}

//-------------------------------------------------------------------------------------------------------------
// Internal logic
//-------------------------------------------------------------------------------------------------------------

// checkAllowedToWrite - Checks to see if we are allowed to write a file.
// The file may exist already. If it does, then we won't be able to over-write it unless
// the forceOverWrite flag is true.
func (generator *FileGenerator) checkAllowedToWrite(targetFilePath string, forceOverwrite bool) error {
	isAlreadyExists, err := generator.fileSystem.Exists(targetFilePath)
	if err == nil {
		if isAlreadyExists && (!forceOverwrite) {
			log.Printf("File %s exists, and we cannot over-write it as the --force flag is not set.", targetFilePath)
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CANNOT_OVERWRITE_FILE, targetFilePath)
		}
	}
	return err
}

func (generator *FileGenerator) substituteParametersIntoTemplate(template *template.Template, templateParameters interface{}) (string, error) {
	// Render the golang template into a string
	var buffer bytes.Buffer
	fileContents := ""
	err := template.Execute(&buffer, templateParameters)
	if err == nil {
		fileContents = buffer.String()
	}
	return fileContents, err
}

func (generator *FileGenerator) loadEmbeddedTemplate(embeddedTemplateFilePath string) (*template.Template, error) {
	// Load-up the template file from the embedded file system.
	data, err := generator.embeddedFileSystem.ReadFile(embeddedTemplateFilePath)
	var templ *template.Template = nil
	if err == nil {
		// Parse the string data into a golang template
		rawTemplate := template.New(embeddedTemplateFilePath)
		templ, err = rawTemplate.Parse(string(data))
	}
	return templ, err
}
