/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"bytes"
	"embed"
	"log"
	"strings"
	"text/template"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"

	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

// PomTemplateSubstitutionParameters holds all the substitution parameters a pom.xml file
// template uses.
type PomTemplateSubstitutionParameters struct {
	GroupId          string
	ArtifactId       string
	Name             string
	ParentArtifactId string
	GalasaVersion    string
}

// JavaTestTemplateSubstitutionParameters holds all the substitution parameters a java test file
// template uses
type JavaTestTemplateSubstitutionParameters struct {
	Package string
}

type GeneratedFile struct {
	fileType                 string
	targetFilePath           string
	embeddedTemplateFilePath string
	templateParameters       interface{}
}

// Embed all the template files into the go executable, so there are no extra files
// we need to ship/install/locate on the target machine.
// We can access the "embedded" file system as if they are normal files.
//
//go:embed templates/*
var embeddedFileSystem embed.FS

var (
	projectCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Creates a new Galasa project",
		Long:  "Creates a new Galasa test project with optional OBR project and build process files",
		Args:  cobra.NoArgs,
		Run:   executeCreateProject,
	}

	packageName   string
	force         bool
	galasaVersion = RootCmd.Version
)

func init() {
	cmd := projectCreateCmd
	parentCommand := projectCmd

	cmd.Flags().StringVar(&packageName, "package", "", "Java package name for tests we create. Forms part of the project name, maven/gradle group/artifact ID, and OSGi bundle name. For example: com.myco.myproduct.myapp")
	cmd.MarkFlagRequired("package")

	cmd.Flags().BoolVar(&force, "force", false, "Force-overwrite files which already exist.")

	parentCommand.AddCommand(cmd)
}

func executeCreateProject(cmd *cobra.Command, args []string) {

	utils.CaptureLog(logFileName)

	log.Println("Galasa CLI - Create project")

	// Operations on the file system will all be relative to the current folder.
	fileSystem := utils.NewOSFileSystem()

	createProject(fileSystem, packageName, force)
}

// createProject will create the following artifacts in the specified file system:
//
// .							 All files are relative to the current directory.
// └── packageName               The parent package.
//
//	├── pom.xml               The parent pom.
//	├── packageName.test      The tests project.
//	│   └── pom.xml
//	└── packageName.obr       The OBR project.
//	    └── pom.xml
func createProject(fileSystem utils.FileSystem, packageName string, forceOverwrite bool) error {
	log.Printf("Creating project using packageName:%s\n", packageName)

	// Create the parent folder
	parentProjectFolder := packageName
	err := createFolder(fileSystem, parentProjectFolder)
	if err == nil {
		err = createParentFolderPom(fileSystem, packageName, forceOverwrite)
		if err == nil {
			err = createTestProject(fileSystem, packageName, forceOverwrite)
		}
	}

	return err
}

func createParentFolderPom(fileSystem utils.FileSystem, packageName string, forceOverwrite bool) error {

	templateParameters := PomTemplateSubstitutionParameters{
		GroupId:    packageName,
		ArtifactId: packageName,
		Name:       packageName}

	targetFile := GeneratedFile{
		fileType:                 "pom",
		targetFilePath:           packageName + "/pom.xml",
		embeddedTemplateFilePath: "templates/projectCreate/parent-project/pom.xml",
		templateParameters:       templateParameters}

	err := createFile(fileSystem, targetFile, forceOverwrite)
	return err
}

func createFolder(fileSystem utils.FileSystem, targetFolderPath string) error {
	err := fileSystem.MkdirAll(targetFolderPath)
	return err
}

func createTestProject(fileSystem utils.FileSystem, packageName string, forceOverwrite bool) error {
	targetFolderPath := packageName + "/" + packageName + ".test"
	log.Printf("Creating tests project %s\n", targetFolderPath)

	// Create the base test folder
	err := createFolder(fileSystem, targetFolderPath)
	if err == nil {
		err = createTestFolderPom(fileSystem, targetFolderPath, packageName, forceOverwrite)
		if err == nil {
			err = createJavaSourceFolder(fileSystem, targetFolderPath, packageName, forceOverwrite)
		}
	}
	return err
}

func createJavaSourceFolder(fileSystem utils.FileSystem, testFolderPath string, packageName string, forceOverwrite bool) error {

	// The folder is the package name but with slashes.
	// eg: my.package becomes my/package
	packageNameWithSlashes := strings.Replace(packageName, ".", "/", -1)
	targetSrcFolderPath := testFolderPath + "/src/main/java/" + packageNameWithSlashes + "/test"
	err := createFolder(fileSystem, targetSrcFolderPath)
	if err == nil {
		err = createJavaSourceFile(fileSystem, targetSrcFolderPath, packageName, forceOverwrite)
	}
	return err
}

func createJavaSourceFile(fileSystem utils.FileSystem, targetSrcFolderPath string, packageName string, forceOverwrite bool) error {
	templateParameters := JavaTestTemplateSubstitutionParameters{
		Package: packageName + ".test"}

	targetFile := GeneratedFile{
		fileType:                 "JavaSourceFile",
		targetFilePath:           targetSrcFolderPath + "/SampleTest.java",
		embeddedTemplateFilePath: "templates/projectCreate/parent-project/test-project/src/main/java/SampleTest.java",
		templateParameters:       templateParameters}

	err := createFile(fileSystem, targetFile, forceOverwrite)
	return err
}

func createTestFolderPom(fileSystem utils.FileSystem, targetTestFolderPath string, packageName string, forceOverwrite bool) error {

	pomTemplateParameters := PomTemplateSubstitutionParameters{
		GroupId:          packageName,
		ParentArtifactId: packageName,
		ArtifactId:       packageName + ".test",
		Name:             packageName + ".test",
		GalasaVersion:    galasaVersion}

	targetFile := GeneratedFile{
		fileType:                 "pom",
		targetFilePath:           targetTestFolderPath + "/pom.xml",
		embeddedTemplateFilePath: "templates/projectCreate/parent-project/test-project/pom.xml",
		templateParameters:       pomTemplateParameters}

	err := createFile(fileSystem, targetFile, forceOverwrite)
	return err
}

// createFile creates a file on the file system.
// If forceOverwrite is false, and there is already a file there, then an error will occur.
func createFile(
	fileSystem utils.FileSystem,
	generatedFile GeneratedFile,
	forceOverwrite bool) error {

	log.Printf("Creating file of type %s at %s\n", generatedFile.fileType, generatedFile.targetFilePath)

	err := assertAllowedToWrite(fileSystem, generatedFile.targetFilePath, forceOverwrite)
	if err == nil {
		var template *template.Template
		template, err = loadEmbeddedTemplate(generatedFile.embeddedTemplateFilePath)
		if err == nil {
			var fileContents string
			fileContents, err = substituteParametersIntoTemplate(template, generatedFile.templateParameters)
			if err == nil {
				// Write it out to the target file.
				err = fileSystem.WriteTextFile(generatedFile.targetFilePath, fileContents)
			}
		}
	}
	return err
}

func substituteParametersIntoTemplate(template *template.Template, templateParameters interface{}) (string, error) {
	// Render the golang template into a string
	var buffer bytes.Buffer
	fileContents := ""
	err := template.Execute(&buffer, templateParameters)
	if err == nil {
		fileContents = buffer.String()
	}
	return fileContents, err
}

func loadEmbeddedTemplate(embeddedTemplateFilePath string) (*template.Template, error) {
	// Load-up the template file from the embedded file system.
	data, err := embeddedFileSystem.ReadFile(embeddedTemplateFilePath)
	var templ *template.Template = nil
	if err == nil {
		// Parse the string data into a golang template
		rawTemplate := template.New(embeddedTemplateFilePath)
		templ, err = rawTemplate.Parse(string(data))
	}
	return templ, err
}

func assertAllowedToWrite(fileSystem utils.FileSystem, targetFilePath string, forceOverwrite bool) error {
	isAlreadyExists, err := fileSystem.Exists(targetFilePath)
	if err == nil {
		if isAlreadyExists && (!forceOverwrite) {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CANNOT_OVERWRITE_FILE, targetFilePath)
		}
	}
	return err
}
