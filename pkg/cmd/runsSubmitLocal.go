/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"embed"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/utils"
)

// executeSubmitLocal attempts to launch a JVM locally in which the Galasa framework is booted
// which will run the test(s) we wish to use. The tests execute locally, but may use values
// from the Galasa ecosystem, if the bootstrap points to it.
//
// JAVA_HOME must be set.
func executeSubmitLocal(
	fileSystem utils.FileSystem,
	embeddedFileSystem embed.FS,
	params RunsSubmitCmdParameters,
	javaHome string) error {

	var err error = nil

	err = validateJavaHome(fileSystem, javaHome)

	if err == nil {
		err = utils.InitialiseGalasaHomeFolder(fileSystem, embeddedFileSystem)
	}

	if err == nil {
		err = executeTestsInJvm(fileSystem, javaHome, params)
	}

	return err
}

func executeTestsInJvm(fileSystem utils.FileSystem, javaHome string, params RunsSubmitCmdParameters) error {
	var err error = nil

	
	// The Output method runs the command, waits for it to finish and collects its standard output. If there were no errors, dateOut will hold bytes with the date info.

    // dateOut, err := dateCmd.Output()
    // if err != nil {
    //     panic(err)
    // }
    // fmt.Println("> date")
    // fmt.Println(string(dateOut))
// Output and other methods of Command will return *exec.Error if there was a problem executing the command (e.g. wrong path), and *exec.ExitError if the command ran but exited with a non-zero return code.

    // _, err = exec.Command("date", "-x").Output()
    // if err != nil {
    //     switch e := err.(type) {
    //     case *exec.Error:
    //         fmt.Println("failed executing:", err)
    //     case *exec.ExitError:
    //         fmt.Println("command exit rc =", e.ExitCode())
    //     default:
    //         panic(err)
    //     }
    // }

	return err
}

// validateJavaHome validate that JAVA_HOME is set correctly.
// If $JAVA_HOME ends with a '/' (or '\' for windows, the trailing slash
// is removed before checks are made.
//
// Constraints:
// - It must be set.
// - JAVA_HOME/bin must be a folder which exists.
// - JAVA_HOME/bin/java must exist as a file.
func validateJavaHome(fileSystem utils.FileSystem, javaHome string) error {

	var err error = nil

	// Check that the javaHome string is well-formed.
	if javaHome == "" {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_JAVA_HOME_NOT_SET)
	}

	if err == nil {
		javaHome = sanitiseJavaHome(javaHome)
	}

	if err == nil {
		err = checkJavaHomeBinFolderExists(fileSystem, javaHome)
	}

	if err == nil {
		err = checkJavaHomeBinJavaProgram(fileSystem, javaHome)
	}

	return err
}

// checkJavaHomeBinJavaProgramExists check to make sure JAVA_HOME/bin/java is a program which
// - exists
func checkJavaHomeBinJavaProgram(fileSystem utils.FileSystem, javaHome string) error {
	var err error = nil

	// Check that the program $JAVA_HOME/bin/java exists
	javaProgramPath := javaHome + utils.FILE_SYSTEM_PATH_SEPARATOR + "bin" + utils.FILE_SYSTEM_PATH_SEPARATOR + "java"
	var isProgramThere bool
	isProgramThere, err = fileSystem.Exists(javaProgramPath)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_JAVA_PROGRAM_PRESENCE_FAIL, javaProgramPath, err.Error())
	} else {
		if !isProgramThere {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_JAVA_PROGRAM_MISSING, javaProgramPath)
		}
	}
	return err
}

// sanitiseJavaHome Massage the javaHome value to make it more valid.
// - strip off trailing path separators
func sanitiseJavaHome(initialJavaHome string) string {

	result := initialJavaHome

	// If the JAVA_HOME ends in a slash, strip it off.
	lastCharacter := initialJavaHome[len(initialJavaHome)-1:]
	if lastCharacter == utils.FILE_SYSTEM_PATH_SEPARATOR {
		// Last character is a file separator. Strip it off.
		result = initialJavaHome[0 : len(initialJavaHome)-1]
	}
	return result
}

// checkJavaHomeBinFolderExists Checks that the $JAVA_HOME/bin folder exists.
func checkJavaHomeBinFolderExists(fileSystem utils.FileSystem, javaHome string) error {
	var err error = nil
	javaBinFolder := javaHome + utils.FILE_SYSTEM_PATH_SEPARATOR + "bin"
	var isBinFolderThere bool
	isBinFolderThere, err = fileSystem.DirExists(javaBinFolder)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_JAVA_HOME_BIN_PRESENCE_FAIL, javaBinFolder, err.Error())
	} else {
		if !isBinFolderThere {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_JAVA_HOME_BIN_MISSING, javaBinFolder)
		}
	}
	return err
}
