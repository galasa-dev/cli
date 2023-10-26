/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"context"
	"strings"

	"github.com/galasa-dev/cli/pkg/api"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
)

// SetProperty - performs all the logic to implement the `galasactl properties set` command,
// but in a unit-testable manner.
func SetProperty(
	namespace string,
	name string,
	value string,
	apiServerUrl string,
	console utils.Console,
) error {
	var err error

	err = updateCpsProperty(namespace, name, value, apiServerUrl, console)

	// if updateProperty() returns an error containing "404 Not Found" due to receiving a
	// GAL5017E from the api, we know the property does not exist and
	// so we assume the user wants to create a new property
	if err != nil && strings.Contains(err.Error(), "404") {
		err = createCpsProperty(namespace, name, value, apiServerUrl, console)
	}

	if err != nil {
		console.WriteString(err.Error())
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

	apicall := restClient.ConfigurationPropertyStoreAPIApi.UpdateCpsProperty(context, namespace, name)
	apicall = apicall.Body(value)
	_, _, err = apicall.Execute()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PUT_PROPERTY_FAILED, name, err.Error())
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

	var cpsPropertyRequest = galasaapi.NewCreateCpsPropertyRequest()
	cpsPropertyRequest.SetName(name)
	cpsPropertyRequest.SetValue(value)

	apicall := restClient.ConfigurationPropertyStoreAPIApi.CreateCpsProperty(context, namespace).CreateCpsPropertyRequest(*cpsPropertyRequest)
	_, _, err = apicall.Execute()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_POST_PROPERTY_FAILED, name, value, err.Error())
	}

	return err
}
