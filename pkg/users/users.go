/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package users

import (
	"strings"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

func validateLoginIdFlag(loginId string) (string, error) {

	var err error

	loginId = strings.TrimSpace(loginId)

	if loginId == "" {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_MISSING_USER_LOGIN_ID_FLAG)
	}

	if err == nil {

		if loginId != "me" {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_LOGIN_ID_NOT_SUPPORTED, loginId)
		}
	}

	return loginId, err

}
