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

func SetUsers(loginId string, roleName string, apiClient *galasaapi.APIClient, console spi.Console) error {

	// We have the role name, but we need the role ID.
	roleWithThatName, err := getRoleFromRestApi(roleName, apiClient)

	if err == nil {

		// We have the user login id, but we need the user number
		var user *galasaapi.UserData
		user, err = getUserByLoginId(loginId, apiClient)
		if err == nil {

			userId := *user.Id
			roleId := *roleWithThatName.GetMetadata().Id

			// Send the update to the rest API
			err = sendUserUpdateToRestApi(userId, roleId, apiClient)
		}
	}
	return err
}

func getUserByLoginId(loginId string, apiClient *galasaapi.APIClient) (*galasaapi.UserData, error) {

	var user *galasaapi.UserData
	var err error

	var users []galasaapi.UserData
	users, err = getUserDataFromRestApi(loginId, apiClient)
	if len(users) < 1 {
		// Error: User not found, so cannot be updated.
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SERVER_USER_NOT_FOUND, loginId)
	} else {
		user = &(users[0])
	}
	return user, err
}

func sendUserUpdateToRestApi(
	userNumber string,
	roleId string,
	apiClient *galasaapi.APIClient,
) error {

	var context context.Context = nil
	var err error
	var restApiVersion string

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	var userUpdateData *galasaapi.UserUpdateData = galasaapi.NewUserUpdateData()
	userUpdateData.SetRole(roleId)

	apiCall := apiClient.UsersAPIApi.UpdateUser(context, userNumber).UserUpdateData(*userUpdateData).ClientApiVersion(restApiVersion)
	if err == nil {

		var resp *http.Response

		resp, err = apiCall.Execute()

		var statusCode int
		if resp != nil {
			defer resp.Body.Close()
			statusCode = resp.StatusCode
		}

		if err != nil {
			log.Println("sendUserUpdateToRestApi - Failed to update user record on from API server")
			err = galasaErrors.NewGalasaErrorWithHttpStatusCode(statusCode, galasaErrors.GALASA_ERROR_UPDATING_USER_RECORD, err.Error())
		} else {
			log.Println("sendUserUpdateToRestApi - User record updated ok.")
		}

	}

	return err
}

func getRoleFromRestApi(
	roleName string,
	apiClient *galasaapi.APIClient,
) (galasaapi.RBACRole, error) {

	var context context.Context = nil
	var err error
	var restApiVersion string
	var roleToReturn galasaapi.RBACRole

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	// Get the role which has that name off the server.
	apiCall := apiClient.RoleBasedAccessControlAPIApi.GetRBACRoles(context).Name(roleName).ClientApiVersion(restApiVersion)
	if err == nil {
		var roles []galasaapi.RBACRole
		var resp *http.Response

		roles, resp, err = apiCall.Execute()

		var statusCode int
		if resp != nil {
			defer resp.Body.Close()
			statusCode = resp.StatusCode
		}

		if err != nil {
			log.Println("getRoleFromRestApi - Failed to retrieve role " + roleName + " from from API server")
			err = galasaErrors.NewGalasaErrorWithHttpStatusCode(statusCode, galasaErrors.GALASA_ERROR_RETRIEVING_USER_LIST_FROM_API_SERVER, err.Error())
		} else {
			log.Printf("getRoleFromRestApi -  %v roles collected", len(roles))
			if len(roles) < 1 {
				log.Println("getRoleFromRestApi - Failed to retrieve role " + roleName + " from from API server")
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_ROLE_NAME_NOT_FOUND, roleName)
			} else {
				// The role we got back is good.
				roleToReturn = roles[0]
			}
		}

	}

	return roleToReturn, err
}
