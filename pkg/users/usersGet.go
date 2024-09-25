/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package users

import (
	"context"
	"fmt"

	"github.com/galasa-dev/cli/pkg/embedded"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

func GetUsers(loginId string, apiClient *galasaapi.APIClient, console spi.Console) error {
	var err error
	var userData galasaapi.UserData

	loginId, err = validateLoginIdFlag(loginId)
	if err == nil {
		userData, err = getUserDataFromRestApi(loginId, apiClient)

		if err == nil {
			extractedUserId := userData.GetLoginId()
			console.WriteString(fmt.Sprintf("id: %s\n", extractedUserId))
		}

	}

	return err
}

func getUserDataFromRestApi(
	loginId string,
	apiClient *galasaapi.APIClient,
) (galasaapi.UserData, error) {

	var err error
	var context context.Context = nil

	var restApiVersion string

	var userProperties = make([]galasaapi.UserData, 0)

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {

		apiCall := apiClient.UsersAPIApi.GetUserByLoginId(context).LoginId(loginId).ClientApiVersion(restApiVersion)
		userProperties, _, err = apiCall.Execute()

	}

	//Since we have the loginId filter, we can return the first index.
	//It will always be the currently logged in user
	return userProperties[0], err
}
