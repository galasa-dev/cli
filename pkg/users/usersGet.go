/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package users

import (
	"context"
	"encoding/json"

	"github.com/galasa-dev/cli/pkg/embedded"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

type UserData struct {
	LoginId string `json:"login_id"`
}

func GetUsers(loginId string, apiClient *galasaapi.APIClient, console spi.Console) error {
	var err error
	var outputText []byte

	err = validateLoginIdFlag(loginId)
	if err == nil {
		outputText, err = getUserDataFromRestApi(loginId, apiClient)

		if err == nil {
			var userData []UserData
			err = json.Unmarshal(outputText, &userData)

			if err == nil {
				console.WriteString("id: " + userData[0].LoginId)
			}

		}

	}

	return err
}

func getUserDataFromRestApi(
	loginId string,
	apiClient *galasaapi.APIClient,
) ([]byte, error) {

	var err error
	var context context.Context = nil

	var restApiVersion string

	var userProperties = make([]galasaapi.UserData, 0)

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {

		apiCall := apiClient.UsersAPIApi.GetUserByLoginId(context).LoginId(loginId).ClientApiVersion(restApiVersion)
		userProperties, _, err = apiCall.Execute()

		if err == nil {
			userPropertiesJSON, err := json.Marshal(userProperties)
			if err == nil {
				return userPropertiesJSON, err
			}
		}

	}

	return nil, err
}
