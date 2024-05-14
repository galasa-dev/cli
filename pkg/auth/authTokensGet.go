/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
)

// GetTokens - performs all the logic to implement the `galasactl auth tokens get` command
func GetTokens(apiServerUrl string, fileSystem files.FileSystem, galasaHome utils.GalasaHome, env utils.Environment) error {

	var err error = nil
	var authToken galasaapi.AuthToken
	var authTokens galasaapi.AuthTokens
	authTokens, err = GetAuthTokensFromRestApi(fileSystem, galasaHome, env)
	if err == nil {
		var jwt string
		jwt, err = GetJwtFromRestApi(apiServerUrl, authProperties)
		if err == nil {
			err = utils.WriteBearerTokenJsonFile(fileSystem, galasaHome, jwt)
		}
	}
	return err
}

func GetAuthTokensFromRestApi(fileSystem files.FileSystem, galasaHome utils.GalasaHome, env utils.Environment) (galasaapi.AuthTokens, error){

}