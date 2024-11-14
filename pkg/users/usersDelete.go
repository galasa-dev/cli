/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package users

import (
	"context"
	"log"
	"net/http"

	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

func DeleteUser(loginId string, apiClient *galasaapi.APIClient, console spi.Console) error {

	userData, err := getUserDataFromRestApi(loginId, apiClient)

	if err == nil {
		err = deletUserFromRestApi(userData, apiClient)
	}

	return err

}

func deletUserFromRestApi(
	users []galasaapi.UserData,
	apiClient *galasaapi.APIClient,
) error {

	var context context.Context = nil
	var resp *http.Response

	restApiVersion, err := embedded.GetGalasactlRestApiVersion()

	if len(users) > 0 {
		userNumber := users[0].GetId()
		apiCall := apiClient.UsersAPIApi.DeleteUserByNumber(context, userNumber).ClientApiVersion(restApiVersion)
		resp, err = apiCall.Execute()

		if err != nil {
			log.Println("deleteUserFromRestApi - Failed to delete user from API server")
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_DELETE_USER, err.Error())
		} else {
			defer resp.Body.Close()
			log.Printf("deleteUserFromRestApi")
		}
	} else {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SERVER_DELETE_USER_NOT_FOUND)
	}

	return err
}
