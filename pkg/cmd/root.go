/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"fmt"
	"log"
	"os"
	"reflect"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:     "galasactl",
		Short:   "CLI for Galasa",
		Long:    `A tool for controlling Galasa resources using the command-line.`,
		Version: "unknowncliversion-unknowngithash",
	}

	logFileName string
)

func Execute() {

	// Catch execution if a panic happens.
	defer func() {
		err := recover()
		finalWord(err)
	}()

	// Execute the command
	err := RootCmd.Execute()
	finalWord(err)
}

func finalWord(obj interface{}) {
	text, exitCode := extractErrorDetails(obj)
	log.Println(text)
	if exitCode != 0 {
		fmt.Fprintln(os.Stderr, text)
	}

	log.Printf("Exit code is %v", exitCode)
	os.Exit(exitCode)
}

func extractErrorDetails(obj interface{}) (string, int) {
	exitCode := 0
	errorText := ""

	if obj == nil {
		errorText = "OK"
	} else {
		exitCode = 1

		// If it's a pointer to a galasa error.
		galasaErrorPtr, isGalasaError := obj.(*galasaErrors.GalasaError)
		if isGalasaError {
			errorType := (galasaErrorPtr).GetMessageType()
			if errorType.Ordinal == galasaErrors.GALASA_ERROR_TESTS_FAILED.Ordinal {
				// The failure was because some tests failed, rather than the tool or infrastructure failed.
				exitCode = 2
			}
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

	return errorText, exitCode
}

func IsInstanceOf(objectPtr interface{}, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&logFileName, "log", "l", "",
		"File to which log information will be sent. Any folder referred to must exist. "+
			"An existing file will be overwritten. "+
			"Specify \"-\" to log to stderr. "+
			"Defaults to not logging.")
	RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
