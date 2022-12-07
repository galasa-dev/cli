/*
 * Copyright contributors to the Galasa project
 */
package errors

type ErrorMessage struct {
	Template string
	Ordinal  int
}

var (
	GALASA_ERROR_UNSUPPORTED_BOOTSTRAP_URL      = ErrorMessage{Template: "GAL1001E: unsupported bootstrap URL %s. Acceptable values start with 'http' or 'https'", Ordinal: 1001}
	GALASA_ERROR_BOOTSTRAP_URL_BAD_ENDING       = ErrorMessage{Template: "GAL1002E: bootstrap url does not end in '/bootstrap'. Bootstrap url is %s", Ordinal: 1002}
	GALASA_ERROR_BAD_BOOTSTRAP_CONTENT          = ErrorMessage{Template: "GAL1003E: bootstrap contents is badly formed. Bootstrap is at %s. Reason is: %s", Ordinal: 1003}
	GALASA_ERROR_FAILED_TO_GET_BOOTSTRAP        = ErrorMessage{Template: "GAL1004E: failed to load the bootstrap from %s. Reason is %s", Ordinal: 1004}
	GALASA_ERROR_THROTTLE_FILE_WRITE            = ErrorMessage{Template: "GAL1005E: failed to write to 'throttle' file %v. Reason is %s", Ordinal: 1005}
	GALASA_ERROR_SUBMIT_MIX_FLAGS_AND_PORTFOLIO = ErrorMessage{Template: "GAL1006E: The submit command does not support mixing of the test selection flags and a portfolio", Ordinal: 1006}
	GALASA_ERROR_SUBMIT_CREATE_REPORT_YAML      = ErrorMessage{Template: "GAL1007E: Failed to create report yaml file %s. Reason is %s", Ordinal: 1007}
	GALASA_ERROR_SUBMIT_RUNS_GROUP_CHECK        = ErrorMessage{Template: "GAL1008E: Failed to check if group %s exists already. Reason is %s", Ordinal: 1008}
	GALASA_ERROR_SUBMIT_MISSING_ACTION_FLAGS    = ErrorMessage{Template: "GAL1009E: The submit command requires either test selection flags or a portfolio", Ordinal: 1009}
	GALASA_ERROR_SUBMIT_INVALID_OVERRIDE        = ErrorMessage{Template: "GAL1010E: Invalid override '%v'", Ordinal: 1010}
	GALASA_ERROR_SUBMIT_REPORT_YAML_ENCODE      = ErrorMessage{Template: "GAL1011E: Failed to encode the yaml file %s. Reason is %s", Ordinal: 1011}
	GALASA_ERROR_SUBMIT_REPORT_JSON_MARSHAL     = ErrorMessage{Template: "GAL1012E: Failed to prepare test report for writing to json file %s. Reason is %s", Ordinal: 1012}
	GALASA_ERROR_SUBMIT_REPORT_JSON_WRITE_FAIL  = ErrorMessage{Template: "GAL1013E: Failed to write test report json file %s. Reason is %s", Ordinal: 1013}
	GALASA_ERROR_SUBMIT_REPORT_JUNIT_PREPARE    = ErrorMessage{Template: "GAL1014E: Failed to prepare test report for writing to the junit results file %s. Reason is %s", Ordinal: 1014}
	GALASA_ERROR_SUBMIT_REPORT_JUNIT_WRITE_FAIL = ErrorMessage{Template: "GAL1015E: Failed to  write test report junit results file %s. Reason is %s", Ordinal: 1015}
	GALASA_ERROR_EMPTY_PORTFOLIO                = ErrorMessage{Template: "GAL1016E: There are no tests in the test porfolio %s", Ordinal: 1016}
	GALASA_ERROR_TESTS_FAILED                   = ErrorMessage{Template: "GAL1017E: Not all runs passed %s", Ordinal: 1017}
	GALASA_ERROR_NO_TESTS_SELECTED              = ErrorMessage{Template: "GAL1018E: No tests were selected.", Ordinal: 1018}
	GALASA_ERROR_PREPARE_INVALID_OVERRIDE       = ErrorMessage{Template: "GAL1019E: Invalid override '%v'", Ordinal: 1019}
)
