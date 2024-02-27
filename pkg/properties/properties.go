/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"regexp"
	"strings"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

const (
	PROPERTY_NAMESPACE_PATTERN = "^[a-z][a-z0-9]+$"
)

func validateInputsAreNotEmpty(namespace string, name string) error {
	var err error
	if len(strings.TrimSpace(name)) == 0 {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_MISSING_NAME_FLAG, name)
	} else {
		if len(strings.TrimSpace(namespace)) == 0 {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_MISSING_NAMESPACE_FLAG, namespace)
		}
	}
	return err
}

func validateNamespaceFormat(namespace string) error {
	var err error

	validNamespaceFormat, err := regexp.Compile(PROPERTY_NAMESPACE_PATTERN)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_COMPILE_NAMESPACE_REGEX, err.Error())
	} else {
		//check if the namespace format matches
		if !validNamespaceFormat.MatchString(namespace) {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_PROPERTY_NAMESPACE_FORMAT, namespace)
		}
	}

	return err
}
