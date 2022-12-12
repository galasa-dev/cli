/*
 * Copyright contributors to the Galasa project
 */
package errors

import (
	"fmt"
)

// The 'type' of a message, used inside Galasa errors
type MessageType struct {
	Template string
	Ordinal  int
}

func NewMessageType(template string, ordinal int) *MessageType {
	messageType := new(MessageType)
	messageType.Ordinal = ordinal
	messageType.Template = template

	// Add the error message to the list of all messages.
	// (So we can render documentation for each message.)
	GALASA_ALL_MESSAGES[ordinal] = messageType

	return messageType
}

// A Galasa errror instance.
// It can be treated as a normal error, but also holds a message-type
// So we can programmatically tell the difference between various errors as required.
type GalasaError struct {
	msgType *MessageType
	message string
}

// NewGalasaError creates a new GalasaError structure.
//
// params - are substituted into the message indicated by the message type.
func NewGalasaError(msgType *MessageType, params ...interface{}) *GalasaError {
	template := msgType.Template
	message := fmt.Sprintf(template, params...) // how to do this with variadic variable ?
	var galasaError *GalasaError
	galasaError = new(GalasaError)
	galasaError.msgType = msgType
	galasaError.message = message

	return galasaError
}

// Render a galasa error into a string, so the GalasaError structure can be used
// as a normal error.
func (err *GalasaError) Error() string {
	return err.message
}

var (
	// A map of all the messages. Indexed by ordinal number.
	GALASA_ALL_MESSAGES = make(map[int]*MessageType)

	GALASA_ERROR_UNSUPPORTED_BOOTSTRAP_URL      = NewMessageType("GAL1001E: Unsupported bootstrap URL %s. Acceptable values start with 'http' or 'https'", 1001)
	GALASA_ERROR_BOOTSTRAP_URL_BAD_ENDING       = NewMessageType("GAL1002E: Bootstrap url does not end in '/bootstrap'. Bootstrap url is %s", 1002)
	GALASA_ERROR_BAD_BOOTSTRAP_CONTENT          = NewMessageType("GAL1003E: Bootstrap contents is badly formed. Bootstrap is at %s. Reason is: %s", 1003)
	GALASA_ERROR_FAILED_TO_GET_BOOTSTRAP        = NewMessageType("GAL1004E: Failed to load the bootstrap from %s. Reason is %s. If the URL is not resolving, try adding the hostname to your /etc/hosts file. This might especially be needed if communicating over a VPN connection.", 1004)
	GALASA_ERROR_THROTTLE_FILE_WRITE            = NewMessageType("GAL1005E: Failed to write to 'throttle' file %v. Reason is %s", 1005)
	GALASA_ERROR_SUBMIT_MIX_FLAGS_AND_PORTFOLIO = NewMessageType("GAL1006E: The submit command does not support mixing of the test selection flags and a portfolio", 1006)
	GALASA_ERROR_SUBMIT_CREATE_REPORT_YAML      = NewMessageType("GAL1007E: Failed to create report yaml file %s. Reason is %s", 1007)
	GALASA_ERROR_SUBMIT_RUNS_GROUP_CHECK        = NewMessageType("GAL1008E: Failed to check if group %s exists already. Reason is %s", 1008)
	GALASA_ERROR_SUBMIT_MISSING_ACTION_FLAGS    = NewMessageType("GAL1009E: The submit command requires either test selection flags or a portfolio", 1009)
	GALASA_ERROR_SUBMIT_INVALID_OVERRIDE        = NewMessageType("GAL1010E: Invalid override '%v'", 1010)
	GALASA_ERROR_SUBMIT_REPORT_YAML_ENCODE      = NewMessageType("GAL1011E: Failed to encode the yaml file %s. Reason is %s", 1011)
	GALASA_ERROR_SUBMIT_REPORT_JSON_MARSHAL     = NewMessageType("GAL1012E: Failed to prepare test report for writing to json file %s. Reason is %s", 1012)
	GALASA_ERROR_SUBMIT_REPORT_JSON_WRITE_FAIL  = NewMessageType("GAL1013E: Failed to write test report json file %s. Reason is %s", 1013)
	GALASA_ERROR_SUBMIT_REPORT_JUNIT_PREPARE    = NewMessageType("GAL1014E: Failed to prepare test report for writing to the junit results file %s. Reason is %s", 1014)
	GALASA_ERROR_SUBMIT_REPORT_JUNIT_WRITE_FAIL = NewMessageType("GAL1015E: Failed to  write test report junit results file %s. Reason is %s", 1015)
	GALASA_ERROR_EMPTY_PORTFOLIO                = NewMessageType("GAL1016E: There are no tests in the test porfolio %s", 1016)
	GALASA_ERROR_TESTS_FAILED                   = NewMessageType("GAL1017E: Not all runs passed %s", 1017)
	GALASA_ERROR_NO_TESTS_SELECTED              = NewMessageType("GAL1018E: No tests were selected.", 1018)
	GALASA_ERROR_PREPARE_INVALID_OVERRIDE       = NewMessageType("GAL1019E: Invalid override '%v'", 1019)
	GALASA_ERROR_OPEN_LOG_FILE_FAILED           = NewMessageType("GAL1020E: Failed to open log file '%s' for writing. Reason is %s", 1020)
	GALASA_ERROR_OPEN_PORTFOLIO_FILE_FAILED     = NewMessageType("GAL1021E: Failed to open portfolio file '%s' for reading. Reason is %s", 1021)
	GALASA_ERROR_PORTFOLIO_BAD_FORMAT           = NewMessageType("GAL1022E: Failed to read portfolio file '%s' because the content is in the wrong format. Reason is %s", 1022)
	GALASA_ERROR_PORTFOLIO_BAD_FORMAT_VERSION   = NewMessageType("GAL1023E: Failed to read portfolio file '%s' because the content are not using format '%s'.", 1023)
	GALASA_ERROR_PORTFOLIO_BAD_RESOURCE_KIND    = NewMessageType("GAL1024E: Failed to read portfolio file '%s' because the content are not a resource of type '%s'.", 1024)
	GALASA_ERROR_CATALOG_NOT_FOUND              = NewMessageType("GAL1025E: Unable to locate test stream '%s' catalog location", 1025)
	GALASA_ERROR_PROPERTY_GET_FAILED            = NewMessageType("GAL1026E: Failed to find location of tests in stream '%s'. Reason is %s", 1026)
	GALASA_ERROR_CATALOG_COPY_FAILED            = NewMessageType("GAL1027E: Failed to copy test catalog from REST reply for property '%s', stream '%s'. Reason is %s", 1027)
	GALASA_ERROR_CATALOG_UNMARSHAL_FAILED       = NewMessageType("GAL1028E: Failed to unmarshal test catalog from REST reply for property '%s', stream '%s'. Reason is %s", 1028)
	GALASA_ERROR_NO_STREAMS_CONFIGURED          = NewMessageType("GAL1029E: Stream '%s' is not found in the ecosystem. There are no streams set up. Ask your Galasa system administrator to add a new stream with the desired name.", 1029)
	GALASA_ERROR_INVALID_STREAM                 = NewMessageType("GAL1030E: Stream '%s' is not found in the ecosystem. Valid streams are:%s. Try again using a valid stream, or ask your Galasa system administrator to add a new stream with the desired name.", 1030)
	GALASA_ERROR_STREAM_FLAG_REQUIRED           = NewMessageType("GAL1031E: Invalid flags. --bundle, --package, --test and --tag flags can only be specified if --stream is provided.", 1031)
	GALASA_ERROR_SELECTION_REGEX_ERROR          = NewMessageType("GAL1032E: Invalid select regex '%v'. Reason is %v", 1032)
	GALASA_ERROR_SELECTION_REGEX_QUOTED_ERROR   = NewMessageType("GAL1033E: Invalid select quoted regex '%v'. Reason is %v", 1033)
	GALASA_ERROR_CLASS_FORMAT                   = NewMessageType("GAL1034E: Class '%v' is not format 'bundle/class'", 1034)
	GALASA_ERROR_CLASS_NAME_BLANK               = NewMessageType("GAL1035E: Class '%v' is not format. Name is blank", 1035)
)
