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
	createFolder(fileSystem, parentProjectFolder)

	createParentFolderPom(fileSystem, packageName, forceOverwrite)

	createTestProject(fileSystem, packageName, forceOverwrite)

	return nil
}

func createParentFolderPom(fileSystem utils.FileSystem, packageName string, forceOverwrite bool) {

	templateParameters := PomTemplateSubstitutionParameters{
		GroupId:    packageName,
		ArtifactId: packageName,
		Name:       packageName}

	targetFile := GeneratedFile{
		fileType:                 "pom",
		targetFilePath:           packageName + "/pom.xml",
		embeddedTemplateFilePath: "templates/projectCreate/parent-project/pom.xml",
		templateParameters:       templateParameters}

	createFile(fileSystem, targetFile, forceOverwrite)
}

func createFolder(fileSystem utils.FileSystem, targetFolderPath string) {
	err := fileSystem.MkdirAll(targetFolderPath)
	if err != nil {
		panic(err)
	}
}

func createTestProject(fileSystem utils.FileSystem, packageName string, forceOverwrite bool) {
	targetFolderPath := packageName + "/" + packageName + ".test"
	log.Printf("Creating tests project %s\n", targetFolderPath)

	// Create the base test folder
	createFolder(fileSystem, targetFolderPath)

	createTestFolderPom(fileSystem, targetFolderPath, packageName, forceOverwrite)

	createJavaSourceFolder(fileSystem, targetFolderPath, packageName, forceOverwrite)
}

func createJavaSourceFolder(fileSystem utils.FileSystem, testFolderPath string, packageName string, forceOverwrite bool) {

	// The folder is the package name but with slashes.
	// eg: my.package becomes my/package
	packageNameWithSlashes := strings.Replace(packageName, ".", "/", -1)
	targetSrcFolderPath := testFolderPath + "/src/main/java/" + packageNameWithSlashes + "/test"
	createFolder(fileSystem, targetSrcFolderPath)

	createJavaSourceFile(fileSystem, targetSrcFolderPath, packageName, forceOverwrite)
}

func createJavaSourceFile(fileSystem utils.FileSystem, targetSrcFolderPath string, packageName string, forceOverwrite bool) {
	templateParameters := JavaTestTemplateSubstitutionParameters{
		Package: packageName + ".test"}

	targetFile := GeneratedFile{
		fileType:                 "JavaSourceFile",
		targetFilePath:           targetSrcFolderPath + "/SampleTest.java",
		embeddedTemplateFilePath: "templates/projectCreate/parent-project/test-project/src/main/java/SampleTest.java",
		templateParameters:       templateParameters}

	createFile(fileSystem, targetFile, forceOverwrite)
}

func createTestFolderPom(fileSystem utils.FileSystem, targetTestFolderPath string, packageName string, forceOverwrite bool) {

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

	createFile(fileSystem, targetFile, forceOverwrite)
}

// createFile creates a file on the file system.
// If forceOverwrite is false, and there is already a file there, then an error will occur.
func createFile(
	fileSystem utils.FileSystem,
	generatedFile GeneratedFile,
	forceOverwrite bool) {

	log.Printf("Creating file of type %s at %s\n", generatedFile.fileType, generatedFile.targetFilePath)

	assertAllowedToWrite(fileSystem, generatedFile.targetFilePath, forceOverwrite)
	template := loadEmbeddedTemplate(generatedFile.embeddedTemplateFilePath)
	fileContents := substituteParametersIntoTemplate(template, generatedFile.templateParameters)

	// Write it out to the target file.
	err := fileSystem.WriteTextFile(generatedFile.targetFilePath, fileContents)
	if err != nil {
		panic(err)
	}
}

func substituteParametersIntoTemplate(template *template.Template, templateParameters interface{}) string {
	// Render the golang template into a string
	var buffer bytes.Buffer
	err := template.Execute(&buffer, templateParameters)
	if err != nil {
		panic(err)
	}
	fileContents := buffer.String()
	return fileContents
}

func loadEmbeddedTemplate(embeddedTemplateFilePath string) *template.Template {
	// Load-up the template file from the embedded file system.
	data, err := embeddedFileSystem.ReadFile(embeddedTemplateFilePath)
	if err != nil {
		panic(err)
	}

	// Parse the string data into a golang template
	template, err := template.New(embeddedTemplateFilePath).Parse(string(data))
	if err != nil {
		panic(err)
	}

	return template
}

func assertAllowedToWrite(fileSystem utils.FileSystem, targetFilePath string, forceOverwrite bool) {

	isAlreadyExists, err := fileSystem.Exists(targetFilePath)
	if err != nil {
		panic(err)
	}

	if isAlreadyExists && (!forceOverwrite) {
		panic("File '%s' already exists, so cannot be over-written")
	}
}
