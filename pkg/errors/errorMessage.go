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

func (err *GalasaError) GetMessageType() *MessageType {
	return err.msgType
}

// NewGalasaError creates a new GalasaError structure.
//
// params - are substituted into the message indicated by the message type.
func NewGalasaError(msgType *MessageType, params ...interface{}) *GalasaError {
	template := msgType.Template
	message := fmt.Sprintf(template, params...) // how to do this with variadic variable ?
	var galasaError *GalasaError = new(GalasaError)
	galasaError.msgType = msgType
	galasaError.message = message

	LogStackTrace()

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
	GALASA_ERROR_TESTS_FAILED                   = NewMessageType("GAL1017E: Not all runs passed. %v failed.", 1017)
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
	GALASA_ERROR_CANNOT_OVERWRITE_FILE          = NewMessageType("GAL1036E: File '%s' exists. Use the --force flag to overwrite it.", 1036)
	GALASA_ERROR_INVALID_CHAR_IN_PACKAGE_NAME   = NewMessageType("GAL1037E: Invalid Java package name '%s' should not contain the '%s' character.", 1037)
	GALASA_ERROR_INVALID_FIRST_CHAR_IN_PKG_NAME = NewMessageType("GAL1038E: Invalid Java package name '%s' should not start with the '%s' character.", 1038)
	GALASA_ERROR_INVALID_LAST_CHAR_IN_PKG_NAME  = NewMessageType("GAL1039E: Invalid Java package name '%s' should not end with the '%s' character.", 1039)
	GALASA_ERROR_PACKAGE_NAME_BLANK             = NewMessageType("GAL1040E: Invalid Java package name. Package name should not be blank.", 1040)
	GALASA_ERROR_FAILED_TO_CREATE_FOLDERS       = NewMessageType("GAL1041E: Failed to create folders '%s'. Reason is '%s'. Check that you have permissions to write to that folder, and that there is enough disk space available and try again.", 1041)
	GALASA_ERROR_FAILED_TO_WRITE_FILE           = NewMessageType("GAL1042E: Failed to open file '%s' for writing. Reason is '%s'. Check that you have permissions to write to that folder and file, and that there is enough disk space available and try again.", 1042)
	GALASA_ERROR_FAILED_TO_READ_FILE            = NewMessageType("GAL1043E: Failed to open file '%s' for reading. Reason is '%s'. Check that you have permissions to read the file and try again.", 1043)
	GALASA_ERROR_INVALID_PKG_RESERVED_WORD      = NewMessageType("GAL1044E: Invalid Java package name. Package name '%s' contains the reserved java keyword '%s'.", 1044)
	GALASA_ERROR_INVALID_FEATURE_NAME           = NewMessageType("GAL1045E: Invalid feature name. Feature name '%s' cannot be used as a java package name. '%s'", 1045)
	GALASA_ERROR_FAILED_TO_FIND_USER_HOME       = NewMessageType("GAL1046E: Failed to determine the home folder of this user. '%s'", 1046)
	GALASA_ERROR_CREATE_REPORT_YAML_EXISTS      = NewMessageType("GAL1047E: Cannot create the yaml report in file '%s' as that file already exists.", 1047)
	GALASA_ERROR_THROTTLE_FILE_READ             = NewMessageType("GAL1048E: Failed to read from 'throttle' file '%v'. Reason is '%s'", 1048)
	GALASA_ERROR_THROTTLE_FILE_INVALID          = NewMessageType("GAL1049E: Invalid value '%v' read from 'throttle' file '%v'. Reason is '%s'", 1049)
	GALASA_ERROR_JAVA_HOME_NOT_SET              = NewMessageType("GAL1050E: JAVA_HOME environment variable is not set. It must be for when --local flag is used.", 1050)
	GALASA_ERROR_JAVA_HOME_BIN_PRESENCE_FAIL    = NewMessageType("GAL1051E: Failed to determine if folder '%s' exists. Reason is '%s'", 1051)
	GALASA_ERROR_JAVA_HOME_BIN_MISSING          = NewMessageType("GAL1052E: Folder '%s' is missing. JAVA_HOME environment variable should refer to a folder which contains a 'bin' folder.", 1052)
	GALASA_ERROR_JAVA_PROGRAM_PRESENCE_FAIL     = NewMessageType("GAL1053E: Failed to determine if '%s' exists. Reason is '%s'", 1053)
	GALASA_ERROR_JAVA_PROGRAM_MISSING           = NewMessageType("GAL1054E: Program '%s' should exist. JAVA_HOME has been set incorrectly.", 1054)
	GALASA_ERROR_FAILED_TO_READ_BOOTSTRAP_FILE  = NewMessageType("GAL1055E: Failed to read bootstrap file '%s'. Reason is '%s'", 1055)
	GALASA_ERROR_RAS_FOLDER_NOT_DETECTED        = NewMessageType("GAL1056E: The RAS folder path could not be detected in trace output for runId '%s'", 1056)
	GALASA_ERROR_RUN_ID_NOT_DETECTED            = NewMessageType("GAL1057E: The run identifier could not be detected in trace output of the child process", 1057)
	GALASA_ERROR_FAILED_TO_LOAD_BOOTSTRAP_FILE  = NewMessageType("GAL1058E: Failed to load bootstrap file '%s'. Reason is '%s'", 1058)
	GALASA_ERROR_FAILED_TO_LOAD_OVERRIDES_FILE  = NewMessageType("GAL1059E: Failed to load overrides file '%s'. Reason is '%s'", 1059)
	GALASA_ERROR_INVALID_OBR_NOT_ENOUGH_PARTS   = NewMessageType("GAL1060E: Badly formed OBR parameter '%s'. Expected it to be of the form mvn:<GROUP_ID>/<ARTIFACT_ID>/<VERSION>/obr with 4 parts separated by slash characters.", 1060)
	GALASA_ERROR_INVALID_OBR_TOO_MANY_PARTS     = NewMessageType("GAL1061E: Badly formed OBR parameter '%s'. Expected it to be of the form mvn:<GROUP_ID>/<ARTIFACT_ID>/<VERSION>/obr with 4 parts separated by slash characters.", 1061)
	GALASA_ERROR_INVALID_OBR_NO_MVN_PREFIX      = NewMessageType("GAL1062E: Badly formed OBR parameter '%s'. Expected it to be of the form mvn:<GROUP_ID>/<ARTIFACT_ID>/<VERSION>/obr with a 'mvn:' prefix.", 1062)
	GALASA_ERROR_INVALID_OBR_NO_OBR_SUFFIX      = NewMessageType("GAL1063E: Badly formed OBR parameter '%s'. Expected it to be of the form mvn:<GROUP_ID>/<ARTIFACT_ID>/<VERSION>/obr with an '/obr' suffix.", 1063)
	GALASA_ERROR_INVALID_CLASS_INPUT_NO_SLASH   = NewMessageType("GAL1064E: Badly formed Class parameter '%s'. Expected it to be of the form <OSGiBundleId>/<FullyQualifiedJavaClass> with no .class suffix. No slash found.", 1064)
	GALASA_ERROR_INVALID_CLASS_TOO_MANY_SLASHES = NewMessageType("GAL1065E: Badly formed Class parameter '%s'. Expected it to be of the form <OSGiBundleId>/<FullyQualifiedJavaClass> with no .class suffix. Too many slashes found.", 1065)
	GALASA_ERROR_INVALID_CLASS_SUFFIX_FOUND     = NewMessageType("GAL1066E: Badly formed Class parameter '%s'. Expected it to be of the form <OSGiBundleId>/<FullyQualifiedJavaClass> with no .class suffix. Unwanted .class suffix detected.", 1066)

	// Warnings...
	GALASA_WARNING_MAVEN_NO_GALASA_OBR_REPO = NewMessageType("GAL2000W: Warning: Maven configuration file settings.xml should contain a reference to a Galasa repository so that the galasa OBR can be resolved. The official release repository is '%s', and 'bleeding edge' repository is '%s'", 2000)
)
