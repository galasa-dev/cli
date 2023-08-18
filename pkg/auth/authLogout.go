package auth

import (
	"fmt"
	"log"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/utils"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
)

func Logout(fileSystem files.FileSystem, console utils.Console, env utils.Environment, galasaHome utils.GalasaHome) error {

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