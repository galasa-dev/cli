/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package errors

import (
	"fmt"
	"strconv"
	"strings"
)

// The 'type' of a message, used inside Galasa errors
type MessageType struct {
	Template           string
	Ordinal            int
	IsStackTraceWanted bool
}

func NewMessageType(template string, ordinal int, isStackTraceWanted bool) *MessageType {
	messageType := new(MessageType)
	messageType.Ordinal = ordinal
	messageType.Template = template

	// As a sanity check... Make sure that the ordinal declared by the message template
	// is also rendered within the message, and that the numbers are the same...
	if !strings.Contains(template, strconv.Itoa(ordinal)) {
		panic(
			fmt.Sprintf("Programming error: Template does not contain a number the same as the ordinal number. "+
				"Ordinal is %d. Template is %s", ordinal, template))
	}

	messageType.IsStackTraceWanted = isStackTraceWanted

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

	if galasaError.msgType.IsStackTraceWanted {
		LogStackTrace()
	}

	return galasaError
}

// Render a galasa error into a string, so the GalasaError structure can be used
// as a normal error.
func (err *GalasaError) Error() string {
	return err.message
}

const (
	STACK_TRACE_WANTED     = true
	STACK_TRACE_NOT_WANTED = false
)

var (
	// A map of all the messages. Indexed by ordinal number.
	GALASA_ALL_MESSAGES = make(map[int]*MessageType)

	GALASA_ERROR_UNSUPPORTED_BOOTSTRAP_URL                = NewMessageType("GAL1001E: Unsupported bootstrap URL %s. Acceptable values start with 'http' or 'https'", 1001, STACK_TRACE_WANTED)
	GALASA_ERROR_BOOTSTRAP_URL_BAD_ENDING                 = NewMessageType("GAL1002E: Bootstrap url does not end in '/bootstrap'. Bootstrap url is %s", 1002, STACK_TRACE_WANTED)
	GALASA_ERROR_BAD_BOOTSTRAP_CONTENT                    = NewMessageType("GAL1003E: Bootstrap contents is badly formed. Bootstrap is at %s. Reason is: %s", 1003, STACK_TRACE_WANTED)
	GALASA_ERROR_FAILED_TO_GET_BOOTSTRAP                  = NewMessageType("GAL1004E: Failed to load the bootstrap from %s. Reason is %s. If the URL is not resolving, try adding the hostname to your /etc/hosts file. This might especially be needed if communicating over a VPN connection.", 1004, STACK_TRACE_WANTED)
	GALASA_ERROR_THROTTLE_FILE_WRITE                      = NewMessageType("GAL1005E: Failed to write to 'throttle' file %v. Reason is %s", 1005, STACK_TRACE_WANTED)
	GALASA_ERROR_SUBMIT_MIX_FLAGS_AND_PORTFOLIO           = NewMessageType("GAL1006E: The submit command does not support mixing of the test selection flags and a portfolio", 1006, STACK_TRACE_WANTED)
	GALASA_ERROR_SUBMIT_CREATE_REPORT_YAML                = NewMessageType("GAL1007E: Failed to create report yaml file %s. Reason is %s", 1007, STACK_TRACE_WANTED)
	GALASA_ERROR_SUBMIT_RUNS_GROUP_CHECK                  = NewMessageType("GAL1008E: Failed to check if group %s exists already. Reason is %s", 1008, STACK_TRACE_WANTED)
	GALASA_ERROR_SUBMIT_MISSING_ACTION_FLAGS              = NewMessageType("GAL1009E: The submit command requires either test selection flags (eg: --stream, --class, --bundle, --package, --tag, --regex, --test) or --portfolio flag to be specified. Use the --help flag for more details.", 1009, STACK_TRACE_WANTED)
	GALASA_ERROR_SUBMIT_INVALID_OVERRIDE                  = NewMessageType("GAL1010E: Invalid override '%v'", 1010, STACK_TRACE_WANTED)
	GALASA_ERROR_SUBMIT_REPORT_YAML_ENCODE                = NewMessageType("GAL1011E: Failed to encode the yaml file %s. Reason is %s", 1011, STACK_TRACE_WANTED)
	GALASA_ERROR_SUBMIT_REPORT_JSON_MARSHAL               = NewMessageType("GAL1012E: Failed to prepare test report for writing to json file %s. Reason is %s", 1012, STACK_TRACE_WANTED)
	GALASA_ERROR_SUBMIT_REPORT_JSON_WRITE_FAIL            = NewMessageType("GAL1013E: Failed to write test report json file %s. Reason is %s", 1013, STACK_TRACE_WANTED)
	GALASA_ERROR_SUBMIT_REPORT_JUNIT_PREPARE              = NewMessageType("GAL1014E: Failed to prepare test report for writing to the junit results file %s. Reason is %s", 1014, STACK_TRACE_WANTED)
	GALASA_ERROR_SUBMIT_REPORT_JUNIT_WRITE_FAIL           = NewMessageType("GAL1015E: Failed to  write test report junit results file %s. Reason is %s", 1015, STACK_TRACE_WANTED)
	GALASA_ERROR_EMPTY_PORTFOLIO                          = NewMessageType("GAL1016E: There are no tests in the test porfolio %s", 1016, STACK_TRACE_WANTED)
	GALASA_ERROR_TESTS_FAILED                             = NewMessageType("GAL1017E: Not all runs passed. %v failed.", 1017, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_NO_TESTS_SELECTED                        = NewMessageType("GAL1018E: No tests were selected.", 1018, STACK_TRACE_WANTED)
	GALASA_ERROR_PREPARE_INVALID_OVERRIDE                 = NewMessageType("GAL1019E: Invalid override '%v'", 1019, STACK_TRACE_WANTED)
	GALASA_ERROR_OPEN_LOG_FILE_FAILED                     = NewMessageType("GAL1020E: Failed to open log file '%s' for writing. Reason is %s", 1020, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_OPEN_PORTFOLIO_FILE_FAILED               = NewMessageType("GAL1021E: Failed to open portfolio file '%s' for reading. Reason is %s", 1021, STACK_TRACE_WANTED)
	GALASA_ERROR_PORTFOLIO_BAD_FORMAT                     = NewMessageType("GAL1022E: Failed to read portfolio file '%s' because the content is in the wrong format. Reason is %s", 1022, STACK_TRACE_WANTED)
	GALASA_ERROR_PORTFOLIO_BAD_FORMAT_VERSION             = NewMessageType("GAL1023E: Failed to read portfolio file '%s' because the content are not using format '%s'.", 1023, STACK_TRACE_WANTED)
	GALASA_ERROR_PORTFOLIO_BAD_RESOURCE_KIND              = NewMessageType("GAL1024E: Failed to read portfolio file '%s' because the content are not a resource of type '%s'.", 1024, STACK_TRACE_WANTED)
	GALASA_ERROR_CATALOG_NOT_FOUND                        = NewMessageType("GAL1025E: Unable to locate test stream '%s' catalog location", 1025, STACK_TRACE_WANTED)
	GALASA_ERROR_PROPERTY_GET_FAILED                      = NewMessageType("GAL1026E: Failed to find location of tests in stream '%s'. Reason is %s", 1026, STACK_TRACE_WANTED)
	GALASA_ERROR_CATALOG_COPY_FAILED                      = NewMessageType("GAL1027E: Failed to copy test catalog from REST reply for property '%s', stream '%s'. Reason is %s", 1027, STACK_TRACE_WANTED)
	GALASA_ERROR_CATALOG_UNMARSHAL_FAILED                 = NewMessageType("GAL1028E: Failed to unmarshal test catalog from REST reply for property '%s', stream '%s'. Reason is %s", 1028, STACK_TRACE_WANTED)
	GALASA_ERROR_NO_STREAMS_CONFIGURED                    = NewMessageType("GAL1029E: Stream '%s' is not found in the ecosystem. There are no streams set up. Ask your Galasa system administrator to add a new stream with the desired name.", 1029, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_STREAM                           = NewMessageType("GAL1030E: Stream '%s' is not found in the ecosystem. Valid streams are:%s. Try again using a valid stream, or ask your Galasa system administrator to add a new stream with the desired name.", 1030, STACK_TRACE_WANTED)
	GALASA_ERROR_STREAM_FLAG_REQUIRED                     = NewMessageType("GAL1031E: Invalid flags. --bundle, --package, --test, --tag, and --class flags can only be specified if --stream is provided.", 1031, STACK_TRACE_WANTED)
	GALASA_ERROR_SELECTION_REGEX_ERROR                    = NewMessageType("GAL1032E: Invalid select regex '%v'. Reason is %v", 1032, STACK_TRACE_WANTED)
	GALASA_ERROR_SELECTION_REGEX_QUOTED_ERROR             = NewMessageType("GAL1033E: Invalid select quoted regex '%v'. Reason is %v", 1033, STACK_TRACE_WANTED)
	GALASA_ERROR_CLASS_FORMAT                             = NewMessageType("GAL1034E: Class '%v' is not format 'bundle/class'", 1034, STACK_TRACE_WANTED)
	GALASA_ERROR_CLASS_NAME_BLANK                         = NewMessageType("GAL1035E: Class '%v' is not format. Name is blank", 1035, STACK_TRACE_WANTED)
	GALASA_ERROR_CANNOT_OVERWRITE_FILE                    = NewMessageType("GAL1036E: File '%s' exists. Use the --force flag to overwrite it.", 1036, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_CHAR_IN_PACKAGE_NAME             = NewMessageType("GAL1037E: Invalid Java package name '%s' should not contain the '%s' character.", 1037, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_FIRST_CHAR_IN_PKG_NAME           = NewMessageType("GAL1038E: Invalid Java package name '%s' should not start with the '%s' character.", 1038, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_LAST_CHAR_IN_PKG_NAME            = NewMessageType("GAL1039E: Invalid Java package name '%s' should not end with the '%s' character.", 1039, STACK_TRACE_WANTED)
	GALASA_ERROR_PACKAGE_NAME_BLANK                       = NewMessageType("GAL1040E: Invalid Java package name. Package name should not be blank.", 1040, STACK_TRACE_WANTED)
	GALASA_ERROR_FAILED_TO_CREATE_FOLDERS                 = NewMessageType("GAL1041E: Failed to create folders '%s'. Reason is '%s'. Check that you have permissions to write to that folder, and that there is enough disk space available and try again.", 1041, STACK_TRACE_WANTED)
	GALASA_ERROR_FAILED_TO_WRITE_FILE                     = NewMessageType("GAL1042E: Failed to open file '%s' for writing. Reason is '%s'. Check that you have permissions to write to that folder and file, and that there is enough disk space available and try again.", 1042, STACK_TRACE_WANTED)
	GALASA_ERROR_FAILED_TO_READ_FILE                      = NewMessageType("GAL1043E: Failed to open file '%s' for reading. Reason is '%s'. Check that you have permissions to read the file and try again.", 1043, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_PKG_RESERVED_WORD                = NewMessageType("GAL1044E: Invalid Java package name. Package name '%s' contains the reserved java keyword '%s'.", 1044, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_FEATURE_NAME                     = NewMessageType("GAL1045E: Invalid feature name. Feature name '%s' cannot be used as a java package name. '%s'", 1045, STACK_TRACE_WANTED)
	GALASA_ERROR_FAILED_TO_FIND_USER_HOME                 = NewMessageType("GAL1046E: Failed to determine the home folder of this user. '%s'", 1046, STACK_TRACE_WANTED)
	GALASA_ERROR_CREATE_REPORT_YAML_EXISTS                = NewMessageType("GAL1047E: Cannot create the yaml report in file '%s' as that file already exists.", 1047, STACK_TRACE_WANTED)
	GALASA_ERROR_THROTTLE_FILE_READ                       = NewMessageType("GAL1048E: Failed to read from 'throttle' file '%v'. Reason is '%s'", 1048, STACK_TRACE_WANTED)
	GALASA_ERROR_THROTTLE_FILE_INVALID                    = NewMessageType("GAL1049E: Invalid value '%v' read from 'throttle' file '%v'. Reason is '%s'", 1049, STACK_TRACE_WANTED)
	GALASA_ERROR_JAVA_HOME_NOT_SET                        = NewMessageType("GAL1050E: JAVA_HOME environment variable is not set. It must be for when --local flag is used.", 1050, STACK_TRACE_WANTED)
	GALASA_ERROR_JAVA_HOME_BIN_PRESENCE_FAIL              = NewMessageType("GAL1051E: Failed to determine if folder '%s' exists. Reason is '%s'", 1051, STACK_TRACE_WANTED)
	GALASA_ERROR_JAVA_HOME_BIN_MISSING                    = NewMessageType("GAL1052E: Folder '%s' is missing. JAVA_HOME environment variable should refer to a folder which contains a 'bin' folder.", 1052, STACK_TRACE_WANTED)
	GALASA_ERROR_JAVA_PROGRAM_PRESENCE_FAIL               = NewMessageType("GAL1053E: Failed to determine if '%s' exists. Reason is '%s'", 1053, STACK_TRACE_WANTED)
	GALASA_ERROR_JAVA_PROGRAM_MISSING                     = NewMessageType("GAL1054E: Program '%s' should exist. JAVA_HOME has been set incorrectly.", 1054, STACK_TRACE_WANTED)
	GALASA_ERROR_FAILED_TO_READ_BOOTSTRAP_FILE            = NewMessageType("GAL1055E: Failed to read bootstrap file '%s'. Reason is '%s'", 1055, STACK_TRACE_WANTED)
	GALASA_ERROR_RAS_FOLDER_NOT_DETECTED                  = NewMessageType("GAL1056E: The RAS folder path could not be detected in trace output for runId '%s'", 1056, STACK_TRACE_WANTED)
	GALASA_ERROR_RUN_ID_NOT_DETECTED                      = NewMessageType("GAL1057E: The run identifier could not be detected in trace output of the child process", 1057, STACK_TRACE_WANTED)
	GALASA_ERROR_FAILED_TO_LOAD_BOOTSTRAP_FILE            = NewMessageType("GAL1058E: Failed to load bootstrap file '%s'. Reason is '%s'", 1058, STACK_TRACE_WANTED)
	GALASA_ERROR_FAILED_TO_LOAD_OVERRIDES_FILE            = NewMessageType("GAL1059E: Failed to load overrides file '%s'. Reason is '%s'", 1059, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_OBR_NOT_ENOUGH_PARTS             = NewMessageType("GAL1060E: Badly formed OBR parameter '%s'. Expected it to be of the form mvn:<GROUP_ID>/<ARTIFACT_ID>/<VERSION>/obr with 4 parts separated by slash characters.", 1060, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_OBR_TOO_MANY_PARTS               = NewMessageType("GAL1061E: Badly formed OBR parameter '%s'. Expected it to be of the form mvn:<GROUP_ID>/<ARTIFACT_ID>/<VERSION>/obr with 4 parts separated by slash characters.", 1061, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_OBR_NO_MVN_PREFIX                = NewMessageType("GAL1062E: Badly formed OBR parameter '%s'. Expected it to be of the form mvn:<GROUP_ID>/<ARTIFACT_ID>/<VERSION>/obr with a 'mvn:' prefix.", 1062, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_OBR_NO_OBR_SUFFIX                = NewMessageType("GAL1063E: Badly formed OBR parameter '%s'. Expected it to be of the form mvn:<GROUP_ID>/<ARTIFACT_ID>/<VERSION>/obr with an '/obr' suffix.", 1063, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_CLASS_INPUT_NO_SLASH             = NewMessageType("GAL1064E: Badly formed Class parameter '%s'. Expected it to be of the form <OSGiBundleId>/<FullyQualifiedJavaClass> with no .class suffix. No slash found.", 1064, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_CLASS_TOO_MANY_SLASHES           = NewMessageType("GAL1065E: Badly formed Class parameter '%s'. Expected it to be of the form <OSGiBundleId>/<FullyQualifiedJavaClass> with no .class suffix. Too many slashes found.", 1065, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_CLASS_SUFFIX_FOUND               = NewMessageType("GAL1066E: Badly formed Class parameter '%s'. Expected it to be of the form <OSGiBundleId>/<FullyQualifiedJavaClass> with no .class suffix. Unwanted .class suffix detected.", 1066, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_OUTPUT_FORMAT                    = NewMessageType("GAL1067E: Unsupported value '%s' for parameter --format. Supported values are: %s", 1067, STACK_TRACE_WANTED)
	GALASA_ERROR_QUERY_RUNS_FAILED                        = NewMessageType("GAL1068E: Could not query run results. Reason: '%s'", 1068, STACK_TRACE_WANTED)
	GALASA_ERROR_LOG_FILE_IS_A_FOLDER                     = NewMessageType("GAL1069E: Could not open log file for writing. '%s' is a directory, the --log parameter should refer to a file path (existing or not), or '-' (the console)", 1069, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_BOOTSTRAP_BAD_DEBUG_MODE_VALUE           = NewMessageType("GAL1070E: Invalid value '%s' detected for optional property '%s' in bootstrap properties. Valid values are 'listen' or 'attach'. Only used when --debug flag is set. Defaults to 'listen'. Can be overridden with the --debugMode flag.", 1070, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_ARG_BAD_DEBUG_MODE_VALUE                 = NewMessageType("GAL1071E: Invalid value '%s' detected for optional --debugMode argument. Valid values are 'listen' or 'attach'. Only used when --debug flag is set. Defaults to 'listen'. Default can be set with an optional property '%s' in bootstrap properties.", 1071, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_BOOTSTRAP_BAD_DEBUG_PORT_VALUE           = NewMessageType("GAL1072E: Invalid value '%s' detected for optional property '%s' in bootstrap properties. Valid values must be a non-zero positive integer, for a port number. Only used when --debug flag is set. Defaults to '%s'. Can be overridden with the --debugPort flag.", 1072, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_RETRIEVING_ARTIFACTS_FAILED              = NewMessageType("GAL1073E: Could not get run artifacts. Reason: '%s'", 1073, STACK_TRACE_WANTED)
	GALASA_ERROR_DOWNLOADING_ARTIFACT_FAILED              = NewMessageType("GAL1074E: Could not download artifact '%s'. Reason: '%s'", 1074, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_RUN_NAME                         = NewMessageType("GAL1075E: run name '%s' is invalid. Expected it to be in format starting with letters, and ending in a number with no non-alphanumeric characters.", 1075, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_INVALID_FROM_OR_TO_PARAMETER             = NewMessageType("GAL1076E: Badly formed from or to value '%s' specified in the age parameter. The value could not be converted into an integer value.", 1076, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_FROM_AGE_SMALLER_THAN_TO_AGE             = NewMessageType("GAL1077E: Invalid value '%s' detected for age parameter. The 'from' value must be greater than the 'to' value.", 1077, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_INVALID_AGE_PARAMETER                    = NewMessageType("GAL1078E: Badly formed '--age' parameter value '%s' specified. Age of the test runs should be specified in the format '{FROM}{TIME-UNIT}:{TO}{TIME-UNIT}' or '{FROM}{TIME-UNIT}', where 'FROM' is a positive, non-zero integer, 'TO' is a non-negative integer, and 'TIME-UNIT' can be %s. 'FROM' must be greater than 'TO'. 'TO' defaults to 0 if not specified.", 1078, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_NO_RUNNAME_OR_AGE_SPECIFIED              = NewMessageType("GAL1079E: The --age or the --name parameter must be used to identify which test run(s) you want see.", 1079, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_INVALID_FROM_AGE_SPECIFIED               = NewMessageType("GAL1080E: Invalid 'from' value '%s' in the '--age' parameter. Age of the test runs should be specified in the format '{FROM}{TIME-UNIT}:{TO}{TIME-UNIT}' or '{FROM}{TIME-UNIT}', where 'FROM' is a positive, non-zero integer, 'TO' is a non-negative integer, and 'TIME-UNIT' can be %s. 'FROM' must be greater than 'TO'. 'TO' defaults to 0 if not specified.", 1080, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_NEGATIVE_AGE_SPECIFIED                   = NewMessageType("GAL1081E: Unable use a negative value '%s' in the '--age' parameter. Age of the test runs should be specified in the format '{FROM}{TIME-UNIT}:{TO}{TIME-UNIT}' or '{FROM}{TIME-UNIT}', where 'FROM' is a positive, non-zero integer, 'TO' is a non-negative integer, and 'TIME-UNIT' can be %s. 'FROM' must be greater than 'TO'. 'TO' defaults to 0 if not specified.", 1081, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_BAD_TIME_UNIT_AGE_SPECIFIED              = NewMessageType("GAL1082E: Invalid time unit specified '%s' in the '--age' parameter. Age of the test runs should be specified in the format '{FROM}{TIME-UNIT}:{TO}{TIME-UNIT}' or '{FROM}{TIME-UNIT}', where 'FROM' is a positive, non-zero integer, 'TO' is a non-negative integer, and 'TIME-UNIT' can be %s. 'FROM' must be greater than 'TO'. 'TO' defaults to 0 if not specified.", 1082, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_NO_ARTIFACTS_TO_DOWNLOAD                 = NewMessageType("GAL1083E: No artifacts to download for run: '%s'", 1083, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_HTTP_RESPONSE_CLOSE_FAILED               = NewMessageType("GAL1084E: Communications failure while closing an HTTP response. '%s'", 1084, STACK_TRACE_WANTED)
	GALASA_ERROR_EMBEDDED_FS_READ_FAILED                  = NewMessageType("GAL1085E: Programming logic error. ReadTextFile operation on embedded file system failed. Reason is %s", 1085, STACK_TRACE_WANTED)
	GALASA_ERROR_QUERY_RESULTNAMES_FAILED                 = NewMessageType("GAL1086E: Communications problem between the command-line tool and the target Galasa Ecosystem. The tool could not retrieve the list of valid result names. Reason: '%s'", 1086, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_RESULT_ARGUMENT                  = NewMessageType("GAL1087E: Invalid '--result' parameter value: %s. The possible result values currently in the Ecosystem Result Archive Store (RAS) are: %s", 1087, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_ACTIVE_AND_RESULT_ARE_MUTUALLY_EXCLUSIVE = NewMessageType("GAL1088E: --active and --result must not be used at the same time, they are mutually exclusive.", 1088, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_MAVEN_AND_OR_GRADLE_FLAG_MUST_BE_SET     = NewMessageType("GAL1089E: Need to use --maven and/or --gradle parameter", 1089, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_RETRIEVING_USERNAME_FAILED               = NewMessageType("GAL1090E: Could not get username of current requestor. Reason is '%s'", 1090, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_BAD_BOOTSTRAP_FILE_URL                   = NewMessageType("GAL1091E: '%s' is not a properly formed file URL", 1091, STACK_TRACE_WANTED)
	GALASA_ERROR_TEST_NOT_IN_RUN_GROUP_LOST               = NewMessageType("GAL1092E: Galasa Ecosystem error: A test was submitted for launch. The galasa runtime is not reporting test progress. "+
		"The test is lost and may execute but test progress cannot be monitored from this tool. (bundle: %s, class: %s).", 1092, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_FAILED_TO_SUBMIT_TEST               = NewMessageType("GAL1093E: Failed to submit test (bundle: %s, class: %s). Reason is: %s", 1093, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_NO_OBR_SPECIFIED_ON_INPUTS          = NewMessageType("GAL1094E: User error: Cannot run test %s on a local JVM because no OBR information is available. Supply an OBR using the --obr parameter, or (if using a portfolio) ensure the portfolio contains an OBR for this test.", 1094, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_INVALID_NAMESPACE                   = NewMessageType("GAL1095E: namespace '%s' is invalid. Expected it to be in a format starting with letters, and ending in a number with no non-alphanumeric characters.", 1095, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_QUERY_NAMESPACE_FAILED              = NewMessageType("GAL1096E: Could not query namespace results.", 1096, STACK_TRACE_WANTED)
	GALASA_ERROR_INVALID_PROPERTIES_FLAG_COMBINATION = NewMessageType("GAL1097E: Property --name cannot be set alongside --prefix, --suffix, or infix.", 1097, STACK_TRACE_WANTED)
	GALASA_ERROR_PUT_PROPERTY_FAILED                 = NewMessageType("GAL1098E: Property '%v' could not be updated. Reason: '%s'", 1098, STACK_TRACE_NOT_WANTED)
	GALASA_ERROR_DELETE_PROPERTY_FAILED              = NewMessageType("GAL1099E: Property '%v' could not be deleted. Reason: '%s'", 1099, STACK_TRACE_WANTED)
	GALASA_ERROR_POST_PROPERTY_FAILED                = NewMessageType("GAL1100E: Property '%v' could not be created. Reason: '%s'", 1100, STACK_TRACE_WANTED)

	// Warnings...
	GALASA_WARNING_MAVEN_NO_GALASA_OBR_REPO = NewMessageType("GAL2000W: Warning: Maven configuration file settings.xml should contain a reference to a Galasa repository so that the galasa OBR can be resolved. The official release repository is '%s', and 'pre-release' repository is '%s'", 2000, STACK_TRACE_WANTED)
	// Information messages...
	GALASA_INFO_FOLDER_DOWNLOADED_TO = NewMessageType("GAL2501I: Downloaded %d artifacts to folder '%s'\n", 2501, STACK_TRACE_NOT_WANTED)
)
