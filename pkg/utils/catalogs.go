/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/galasaapi"
)

type TestCatalog map[string]interface{}

func FetchTestCatalog(apiClient *galasaapi.APIClient, stream string) (TestCatalog, error) {

	var testCatalog TestCatalog

	cpsProperty, _, err := apiClient.ConfigurationPropertyStoreAPIApi.GetCpsNamespaceCascadeProperty(nil, "framework", "test.stream."+stream, "location").Execute()
	if err != nil {
		err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PROPERTY_GET_FAILED, stream, err)
		panic(err)
	}

	if cpsProperty.Value == nil {
		err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CATALOG_NOT_FOUND, stream)
		panic(err)
	}

	catalogString := new(strings.Builder)

	resp, err := http.Get(*cpsProperty.Value)
	if err != nil {
		err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PROPERTY_GET_FAILED, *cpsProperty.Value, stream, err)
		panic(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(catalogString, resp.Body)
	if err != nil {
		err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CATALOG_COPY_FAILED, *cpsProperty.Value, stream, err)
		panic(err)
	}

	err = json.Unmarshal([]byte(catalogString.String()), &testCatalog)
	if err != nil {
		err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CATALOG_UNMARSHAL_FAILED, *cpsProperty.Value, stream, err)
		panic(err)
	}

	return testCatalog, nil
}
