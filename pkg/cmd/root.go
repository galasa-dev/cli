/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:     "galasactl",
		Short:   "CLI for Galasa",
		Long:    `A tool for controlling Galasa resources using the command-line.`,
		Version: "unknowncliversion-unknowngithash",
	}

	// The file to which logs are being directed, if any. "" if not.
	logFileName string

	// We don't trace anything until this flag is true.
	// This means that any errors which occur in the cobra framework are not
	// followed by stack traces all the time.
	isCapturingLogs bool = false

	// The path to GALASA_HOME. Over-rides the environment variable.
	CmdParamGalasaHomePath string
)

func Execute() {

	// Catch execution if a panic happens.
	defer func() {
		err := recover()

		// Display the error and exit.
		finalWord(err)
	}()

	// Execute the command
	err := RootCmd.Execute()
	finalWord(err)
}

func finalWord(obj interface{}) {
	text, exitCode, isStackTraceWanted := extractErrorDetails(obj)
	if isCapturingLogs {
		log.Println(text)
	}

	if exitCode != 0 {
		fmt.Fprintln(os.Stderr, text)
	}

	if isStackTraceWanted && isCapturingLogs {
		galasaErrors.LogStackTrace()
	}

	if isCapturingLogs {
		log.Printf("Exit code is %v", exitCode)
	}

	os.Exit(exitCode)
}

func extractErrorDetails(obj interface{}) (string, int, bool) {
	exitCode := 0
	errorText := ""
	var isStackTraceWanted bool = false

	if obj == nil {
		errorText = "OK"
	} else {
		exitCode = 1
		isStackTraceWanted = true

		// If it's a pointer to a galasa error.
		galasaErrorPtr, isGalasaError := obj.(*galasaErrors.GalasaError)
		if isGalasaError {
			errorType := (galasaErrorPtr).GetMessageType()
			if errorType.Ordinal == galasaErrors.GALASA_ERROR_TESTS_FAILED.Ordinal {
				// The failure was because some tests failed, rather than the tool or infrastructure failed.
				exitCode = 2
			}
			// Don't log a stack trace for Galasa errors. We know where they come from.
			isStackTraceWanted = false
		}

		err, isErrorType := obj.(error)
		if isErrorType {
			errorText = err.Error()
		} else {

			stringValue, isString := obj.(string)
			if isString {
				errorText = stringValue
			} else {
				errorText = "unknown error."
			}
		}
	}

	return errorText, exitCode, isStackTraceWanted
}

func IsInstanceOf(objectPtr interface{}, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}

func init() {
	galasaCtlVersion, err := embedded.GetGalasaCtlVersion()
	if err != nil {
		// If that failed, something is very wrong...
		// like we can't read the file from the embedded file system.
		// Give up out now with a bad exit code
		finalWord(err)
	}

	RootCmd.Version = galasaCtlVersion

	RootCmd.PersistentFlags().StringVarP(&logFileName, "log", "l", "",
		"File to which log information will be sent. Any folder referred to must exist. "+
			"An existing file will be overwritten. "+
			"Specify \"-\" to log to stderr. "+
			"Defaults to not logging.")

	RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	SetHelpFlagForAllCommands(RootCmd, func(cobra *cobra.Command) {
		alias := cobra.NameAndAliases()
		//if the command has an alias,
		//the format would be cobra.Name, cobra.Aliases
		//otherwise it is just cobra.Name
		nameAndAliases := strings.Split(alias, ", ")
		if len(nameAndAliases) > 1 {
			alias = nameAndAliases[1]
		}

		cobra.Flags().BoolP("help", "h", false, "Displays the options for the "+alias+" command.")
	})

	RootCmd.PersistentFlags().StringVarP(&CmdParamGalasaHomePath, "galasahome", "", "",
		"Path to a folder where Galasa will read and write files and configuration settings. "+
			"The default is '${HOME}/.galasa'. "+
			"This overrides the GALASA_HOME environment variable which may be set instead.",
	)
}

func SetHelpFlagForAllCommands(command *cobra.Command, setHelpFlag func(*cobra.Command)) {
	setHelpFlag(command)

	//for all the commands eg properties get, set etc
	for _, cobraCommand := range command.Commands() {
		SetHelpFlagForAllCommands(cobraCommand, setHelpFlag)
	}
}