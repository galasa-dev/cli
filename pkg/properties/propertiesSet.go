/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

// SetProperty - performs all the logic to implement the `galasactl properties update` command,
// but in a unit-testable manner.
func SetProperty(
	namespace string,
	name string,
	value string,
	apiServerUrl string,
	console utils.Console,
) error {
	var err error

	if err == nil {

		err = updateCpsProperty(namespace, name, value, apiServerUrl, console)

		if err != nil && strings.Contains(err.Error(), "404") {
			//if the property does not exist we expect a 404 Not Found error (5017 api error), so we attempt to create it
			err = createCpsProperty(namespace, name, value, apiServerUrl, console)
		}

		if err == nil {
			console.WriteString("Successfully set the value of '" + name + "' in namespace '" + namespace + "'")
		} else {
			//for any other error
			console.WriteString(err.Error())
		}
	}
	return err
}

func updateCpsProperty(namespace string,
	name string,
	value string,
	apiServerUrl string,
	console utils.Console,
) error {
	var err error = nil

	var context context.Context = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	var httpResponse *http.Response

	apicall := restClient.ConfigurationPropertyStoreAPIApi.UpdateCpsProperty(context, namespace, name)
	apicall = apicall.Body(value)
	_, httpResponse, err = apicall.Execute()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PUT_PROPERTY_FAILED, name, value, err.Error())
	} else {
		if httpResponse.StatusCode != http.StatusOK {
			httpError := "\nhttp response status code: " + strconv.Itoa(httpResponse.StatusCode)
			errString := err.Error() + httpError
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PUT_PROPERTY_STATUS_CODE_NOT_OK, errString)
		}
	}

	return err
}

func createCpsProperty(namespace string,
	name string,
	value string,
	apiServerUrl string,
	console utils.Console,
) error {
	var err error = nil

	var context context.Context = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	var httpResponse *http.Response
	var cpsPropertyRequest = galasaapi.NewCreateCpsPropertyRequest()
	cpsPropertyRequest.SetName(name)
	cpsPropertyRequest.SetValue(value)

	apicall := restClient.ConfigurationPropertyStoreAPIApi.CreateCpsPropertyRequest(cpsPropertyRequest)
	apicall = restClient.ConfigurationPropertyStoreAPIApi.CreateCpsProperty(context, namespace)
	_, httpResponse, err = apicall.Execute()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_POST_PROPERTY_FAILED, name, value, err.Error())
	} else {
		if httpResponse.StatusCode != http.StatusOK {
			httpError := "\nhttp response status code: " + strconv.Itoa(httpResponse.StatusCode)
			errString := err.Error() + httpError
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_POST_PROPERTY_STATUS_CODE_NOT_OK, errString)
		}
	}

	return err
}
