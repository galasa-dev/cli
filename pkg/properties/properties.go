/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"strings"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
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
