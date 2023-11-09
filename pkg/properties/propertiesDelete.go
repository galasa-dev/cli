/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"context"
	"net/http"

	"github.com/galasa-dev/cli/pkg/api"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

// DeleteProperty - performs all the logic to implement the `galasactl properties delete` command,
// but in a unit-testable manner.
func DeleteProperty(
	namespace string,
	name string,
	apiServerUrl string,
) error {
	var err error
	err = validateInputsAreNotEmpty(namespace, name)
	if err == nil {
		err = deleteCpsProperty(namespace, name, apiServerUrl)
	}
	return err
}

func deleteCpsProperty(namespace string,
	name string,
	apiServerUrl string,
) error {
	var err error = nil
	var resp *http.Response
	var context context.Context = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	apicall := restClient.ConfigurationPropertyStoreAPIApi.DeleteCpsProperty(context, namespace, name)
	_, resp, err = apicall.Execute()

	defer resp.Body.Close()

	if (resp != nil) && (resp.StatusCode != 200) {
		var apiError galasaErrors.GalasaAPIError
		err = apiError.UnmarshalApiError(resp)
		if err == nil {
			//Ensure that the conversion of the error doesn't raise another exception
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_DELETE_PROPERTY_FAILED, name, apiError.Message)
		}
	}
	return err
}
