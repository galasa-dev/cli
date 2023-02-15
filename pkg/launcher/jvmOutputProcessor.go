/*
 * Copyright contributors to the Galasa project
 */
package launcher

import (
	"bytes"
	"log"
	"regexp"
	"strings"
)

type JVMOutputProcessor struct {

	// The data which the processor has been passed so far.
	bytesCollected *bytes.Buffer

	// The runid which has been detected.
	detectedRunId string

	// The location of the RAS folder for this test.
	// "" if it isn't known yet.
	detectedRasFolderPath string

	// The channel that goroutines can wait on until the processor
	// has detected something of interest, or the buffer is closed.
	// The channel can post "ALERT" when something is detected.
	publishResultChannel chan string
}

func NewJVMOutputProcessor() *JVMOutputProcessor {
	processor := new(JVMOutputProcessor)
	processor.detectedRunId = ""
	processor.bytesCollected = bytes.NewBuffer([]byte{})
	processor.publishResultChannel = make(chan string, 10)
	processor.detectedRasFolderPath = ""
	return processor
}

var (
	runIdRegex         *regexp.Regexp = regexp.MustCompile(`Allocated Run Name (?P<runid>\S*) to this run`)
	runIdIndex         int            = runIdRegex.SubexpIndex("runid")
	rasFolderPathRegex *regexp.Regexp = regexp.MustCompile(`Result Archive Stores are \[(?P<ras_location>.*)]`)
	rasLocationIndex   int            = rasFolderPathRegex.SubexpIndex("ras_location")
)

const (
	SHUTDOWN_FRAMEWORK_EYE_CATCHER = `d.g.f.Framework - Framework shutdown`
)

func (processor *JVMOutputProcessor) Write(bytesToWrite []byte) (int, error) {

	bytesWrittenCount, err := processor.bytesCollected.Write(bytesToWrite)

	if err == nil {
		// See if we can gather the runId from the trace output.
		// We would expect it to appear in a string like this:
		// "d.g.f.FrameworkInitialisation - Allocated Run Name U525 to this run"
		stringToSearch := string(bytesToWrite)

		isAlertable := false

		runId := detectRunId(stringToSearch)
		if runId != "" {
			processor.detectedRunId = runId
			isAlertable = true
		}

		rasFolderPath := detectRasFolderPath(stringToSearch)
		if rasFolderPath != "" {
			processor.detectedRasFolderPath = rasFolderPath
			isAlertable = true
		}

		isShutdownDetected := detectShutdown(stringToSearch)
		if isShutdownDetected {
			isAlertable = true
		}

		if isAlertable {
			// Now alert anyone who may be listening on the go channel.
			processor.publishResultChannel <- "ALERT"
		}
	}

	return bytesWrittenCount, err
}

// We expect each test to trace the following:
// "Result Archive Stores are [file:///Users/mcobbett/.galasa/ras]"
// So we should pick up this location and use it to find the json
// file containing the status of this test.
func detectRasFolderPath(stringToSearch string) string {
	var rasFolderPath string = ""

	rasFolderPathMatches := rasFolderPathRegex.FindStringSubmatch(stringToSearch)
	if rasFolderPathMatches == nil {
		// No matches in this string.
	} else {
		rasFolderPath = rasFolderPathMatches[rasLocationIndex]
		log.Printf("JVM Output processor discovered that the RAS folder path for the test is %s\n", rasFolderPath)
	}
	return rasFolderPath
}

func detectRunId(stringToSearch string) string {
	var runId string = ""

	runIdMatches := runIdRegex.FindStringSubmatch(stringToSearch)
	if runIdMatches == nil {
		// No matches in this string.
	} else {
		runId = runIdMatches[runIdIndex]
		log.Printf("JVM Output processor discovered that the RunId for the test is %s\n", runId)
	}
	return runId
}

func detectShutdown(stringToSearch string) bool {
	return strings.Contains(stringToSearch, SHUTDOWN_FRAMEWORK_EYE_CATCHER)
}
