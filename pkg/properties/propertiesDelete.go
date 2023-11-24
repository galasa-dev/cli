/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"context"
	"net/http"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
)

// DeleteProperty - performs all the logic to implement the `galasactl properties delete` command,
// but in a unit-testable manner.
func DeleteProperty(
	namespace string,
	name string,
	apiClient *galasaapi.APIClient,
) error {
	var err error
	err = validateInputsAreNotEmpty(namespace, name)
	if err == nil {
		err = deleteCpsProperty(namespace, name, apiClient)
	}
	return err
}

func deleteCpsProperty(namespace string,
	name string,
	apiClient *galasaapi.APIClient,
) error {
	var err error = nil
	var resp *http.Response
	var context context.Context = nil

	apicall := apiClient.ConfigurationPropertyStoreAPIApi.DeleteCpsProperty(context, namespace, name)
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
