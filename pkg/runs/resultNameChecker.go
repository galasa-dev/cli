/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"sort"
	"strings"

	"github.com/galasa.dev/cli/pkg/errors"
)

const (
	RUN_RESULT_TOTAL               = "Total"
	RUN_RESULT_PASSED              = "Passed"
	RUN_RESULT_PASSED_WITH_DEFECTS = "Passed With Defects"
	RUN_RESULT_FAILED              = "Failed"
	RUN_RESULT_FAILED_WITH_DEFECTS = "Failed With Defects"
	RUN_RESULT_ENVFAIL             = "EnvFail"
	RUN_RESULT_UNKNOWN             = "UNKNOWN"
	RUN_RESULT_IGNORED             = "Ignored"
)

var RESULT_TYPES = []string{RUN_RESULT_PASSED, RUN_RESULT_PASSED_WITH_DEFECTS, RUN_RESULT_FAILED, RUN_RESULT_FAILED_WITH_DEFECTS, RUN_RESULT_ENVFAIL, RUN_RESULT_UNKNOWN, RUN_RESULT_IGNORED}

// ------------------------------------------------
func ValidateResultName(resultNameInput string) (string, error) {
	var err error = nil
	var resultOut string = ""
	resultInputs := strings.Split(resultNameInput, ",")
	for _, result := range resultInputs {
		result = strings.Trim(result, " ")
	}
	// compare the input result to the list of possibles - case insensitive
	for _, input := range resultInputs {
		for _, resultType := range RESULT_TYPES {
			if strings.EqualFold(input, resultType) {
				resultOut = resultType
				//resultOut = resultType+","
			}
		}
	}

	if resultOut == "" {
		err = errors.NewGalasaError(errors.GALASA_ERROR_INVALID_RESULT_ARGUMENT, resultNameInput, getResultNamesString())
	}
	return resultOut, err
}

func getResultNamesString() string {
	// extract names into a sorted slice
	names := make([]string, 0, len(RESULT_TYPES))
	names = append(names, RESULT_TYPES...)
	sort.Strings(names)

	// render list of sorted names into string
	resultNames := strings.Builder{}

	for count, resultName := range names {

		if count != 0 {
			resultNames.WriteString(", ")
		}
		resultNames.WriteString("'" + resultName + "'")
	}

	return resultNames.String()
}
