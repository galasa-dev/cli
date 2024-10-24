/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package secrets

import (
    "strings"

    galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

func validateSecretName(secretName string) error {
    var err error
    secretName = strings.TrimSpace(secretName)

    if secretName == "" {
        err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_MISSING_SECRET_NAME)
    }
    return err
}
