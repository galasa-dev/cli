/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"fmt"
	"log"
	"os"
)

/*
 * CaptureLog(logFileName) decides whether to re-direct the log information to the
 * specified file, or if the file name is "-" or empty, the log information won't be
 * re-directed, but will appear on stderr.
 */
func CaptureLog(logFileName string) *os.File {
	var logFile = os.Stderr

	// Send the log to a file
	if logFileName != "" && logFileName != "-" {
		// The user has set the logFileName using the --log xxxx syntax
		f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(f)
		} else {
			panic(fmt.Sprintf("error opening log file: %v", err))
		}
	}

	return logFile
}
