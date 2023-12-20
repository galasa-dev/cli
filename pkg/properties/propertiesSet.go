/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"context"
	"log"
	"strings"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
)

var (
	APIVERSION = "galasa-dev/v1aplha1"
	KIND       = "GalasaProperty"
)

// SetProperty - performs all the logic to implement the `galasactl properties set` command,
// but in a unit-testable manner.
func SetProperty(
	namespace string,
	name string,
	value string,
	apiClient *galasaapi.APIClient,
) error {
	var err error

	err = validateInputsAreNotEmpty(namespace, name)
	if err == nil {
		galasaProperty := createGalasaProperty(namespace, name, value)

		log.Printf("SetProperty - Galasa Property created: ApiVersion:'%s', Kind:'%s', Namespace:'%s', Name:'%s', Value:'%s'", 
		galasaProperty.GetApiVersion(), galasaProperty.GetKind(), galasaProperty.Metadata.GetNamespace(), galasaProperty.Metadata.GetName(), galasaProperty.Data.GetValue())

		err = updateCpsProperty(namespace, name, galasaProperty, apiClient)

		// if updateProperty() returns an error containing "404 Not Found" due to receiving a
		// GAL5017E from the api, we know the property does not exist and
		// so we assume the user wants to create a new property
		if err != nil && strings.Contains(err.Error(), "404") {
			err = createCpsProperty(namespace, name, galasaProperty, apiClient)
		}
	}
	return err
}

func updateCpsProperty(namespace string,
	name string,
	property *galasaapi.GalasaProperty,
	apiClient *galasaapi.APIClient,
) error {
	var err error = nil
	var context context.Context = nil

	apicall := apiClient.ConfigurationPropertyStoreAPIApi.UpdateCpsProperty(context, namespace, name).GalasaProperty(*property)
	_, _, err = apicall.Execute()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PUT_PROPERTY_FAILED, name, err.Error())
	}

	return err
}

func createCpsProperty(namespace string,
	name string,
	property *galasaapi.GalasaProperty,
	apiClient *galasaapi.APIClient,
) error {
	var err error = nil
	var context context.Context = nil

	apicall := apiClient.ConfigurationPropertyStoreAPIApi.CreateCpsProperty(context, namespace).GalasaProperty(*property)
	_, _, err = apicall.Execute()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_POST_PROPERTY_FAILED, name, err.Error())
	}

	return err
}

func createGalasaProperty(namespace string, name string, value string) *galasaapi.GalasaProperty {
	var property = galasaapi.NewGalasaProperty()

	property.SetApiVersion(APIVERSION)
	property.SetKind(KIND)

	metadata := galasaapi.NewGalasaPropertyMetadata()
	metadata.SetNamespace(namespace)
	metadata.SetName(name)
	property.SetMetadata(*metadata)

	data := galasaapi.NewGalasaPropertyData()
	data.SetValue(value)
	property.SetData(*data)

	return property
}
