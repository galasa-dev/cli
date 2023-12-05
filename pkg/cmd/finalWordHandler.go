/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"
	"os"
	"reflect"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

// A final word handler can set the exit code of the entire process.
// Or it could be mocked-out to just collect it and checked in tests.
type FinalWordHandler interface {
	FinalWord(interface{})
}

// The real implementation of the interface.
type RealFinalWordHandler struct {
}

func NewRealFinalWordHandler() FinalWordHandler {
	return new(RealFinalWordHandler)
}

func (*RealFinalWordHandler) FinalWord(obj interface{}) {
	text, exitCode, isStackTraceWanted := extractErrorDetails(obj)
	if rootCmdValues.isCapturingLogs {
		log.Println(text)
	}

	if isStackTraceWanted && rootCmdValues.isCapturingLogs {
		galasaErrors.LogStackTrace()
	}

	if rootCmdValues.isCapturingLogs {
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
