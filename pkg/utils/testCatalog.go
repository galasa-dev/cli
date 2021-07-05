//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/galasa.dev/cli/pkg/galasaapi"
)

type TestCatalog map[string]interface{}

func FetchTestCatalog(apiClient *galasaapi.APIClient, stream string) (TestCatalog, error) {

	var testCatalog TestCatalog

	cpsProperty, _, err := apiClient.ConfigurationPropertyStoreAPIApi.GetCpsNamespaceCascadeProperty(nil, "framework", "test.stream." + stream, "location").Execute()
	if (err != nil) {
		panic(err)
	}

	if cpsProperty.Value == nil {
		return testCatalog, errors.New("Unable to locate test stream \"" + stream + "\" catalog location")
	}

	catalogString := new(strings.Builder)

	resp, err := http.Get(*cpsProperty.Value)
	if (err != nil) {
		panic(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(catalogString, resp.Body)
	if (err != nil) {
		panic(err)
	}

	err = json.Unmarshal([]byte(catalogString.String()), &testCatalog)
	if (err != nil) {
		panic(err)
	}

	return testCatalog, nil
}
