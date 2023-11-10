/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"fmt"
	"log"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/utils"
)

func Logout(fileSystem files.FileSystem, galasaHome utils.GalasaHome) error {

    var err error = nil
    bearerTokenFilePath := fmt.Sprintf("%s/%s", galasaHome.GetNativeFolderPath(), "bearer-token.json")
    if _, err = fileSystem.Exists(bearerTokenFilePath); err == nil {
        log.Printf("Deleting bearer token file '%s'", bearerTokenFilePath)
        fileSystem.DeleteFile(bearerTokenFilePath)
        log.Printf("Deleted bearer token file '%s' OK", bearerTokenFilePath)
    }

    if err != nil {
        err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_DELETE_BEARER_TOKEN_FILE, bearerTokenFilePath)
    }

    return err
}