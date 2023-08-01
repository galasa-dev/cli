/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package launcher

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/galasaapi"
)

// RemoteLauncher A launcher, which launches and monitors tests on a remote ecosystem via HTTP/HTTPS.
type RemoteLauncher struct {
	apiClient *galasaapi.APIClient
}

//----------------------------------------------------------------------------------
// Constructors
//----------------------------------------------------------------------------------

// NewRemoteLauncher create a remote launcher.
func NewRemoteLauncher(apiServerUrl string) *RemoteLauncher {
	log.Printf("NewRemoteLauncher(%s) entered.", apiServerUrl)

	launcher := new(RemoteLauncher)

	// An HTTP client which can communicate with the api server in an ecosystem.
	launcher.apiClient = api.InitialiseAPI(apiServerUrl)

	return launcher
}

//----------------------------------------------------------------------------------
// Implementation of the launcher interface.
//----------------------------------------------------------------------------------

// GetRunsByGroup get all the testruns which are associated with a named group.
func (launcher *RemoteLauncher) GetRunsByGroup(groupName string) (*galasaapi.TestRuns, error) {
	log.Printf("GetRunsByGroup(%s) entered.", groupName)
	var (
		testRuns *galasaapi.TestRuns
		err      error
	)
	testRuns, _, err = launcher.apiClient.RunsAPIApi.GetRunsGroup(nil, groupName).Execute()
	return testRuns, err
}

func (launcher *RemoteLauncher) SubmitTestRuns(
	groupName string,
	classNames []string,
	requestType string,
	requestor string,
	stream string,
	isTraceEnabled bool,
	overrides map[string]interface{},
) (*galasaapi.TestRuns, error) {

	testRunRequest := galasaapi.NewTestRunRequest()
	testRunRequest.SetClassNames(classNames)
	testRunRequest.SetRequestorType(requestType)
	testRunRequest.SetRequestor(requestor)
	testRunRequest.SetTestStream(stream)
	testRunRequest.SetTrace(isTraceEnabled)
	testRunRequest.SetOverrides(overrides)

	var resultGroup *galasaapi.TestRuns
	var err error

	resultGroup, _, err = launcher.apiClient.RunsAPIApi.PostSubmitTestRuns(nil, groupName).TestRunRequest(*testRunRequest).Execute()

	return resultGroup, err
}

func (launcher *RemoteLauncher) GetRunsById(runId string) (*galasaapi.Run, error) {
	rasRun, _, err := launcher.apiClient.ResultArchiveStoreAPIApi.GetRasRunById(nil, runId).Execute()
	return rasRun, err
}

func (launcher *RemoteLauncher) GetStreams() ([]string, error) {

	var streams []string

	cpsProperty, _, err := launcher.apiClient.ConfigurationPropertyStoreAPIApi.
		GetCpsNamespaceCascadeProperty(nil, "framework", "test", "streams").Execute()
	if err == nil {
		if cpsProperty.Value == nil {
			streams = make([]string, 0)
		} else {
			streams = strings.Split(*cpsProperty.Value, ",")
		}
	}
	return streams, err
}

func (launcher *RemoteLauncher) GetTestCatalog(stream string) (TestCatalog, error) {

	var err error = nil
	var testCatalog TestCatalog
	var cpsProperty *galasaapi.CpsProperty

	cpsProperty, _, err = launcher.apiClient.ConfigurationPropertyStoreAPIApi.GetCpsNamespaceCascadeProperty(
		nil, "framework", "test.stream."+stream, "location").Execute()
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PROPERTY_GET_FAILED, stream, err)
	} else if cpsProperty.Value == nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CATALOG_NOT_FOUND, stream)
	}

	if err == nil {
		catalogString := new(strings.Builder)
		var resp *http.Response
		resp, err = http.Get(*cpsProperty.Value)
		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PROPERTY_GET_FAILED, *cpsProperty.Value, stream, err)
		} else {
			defer resp.Body.Close()

			_, err = io.Copy(catalogString, resp.Body)
			if err != nil {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CATALOG_COPY_FAILED, *cpsProperty.Value, stream, err)
			}
		}

		if err == nil {
			err = json.Unmarshal([]byte(catalogString.String()), &testCatalog)
			if err != nil {
				err = galasaErrors.NewGalasaError(
					galasaErrors.GALASA_ERROR_CATALOG_UNMARSHAL_FAILED, *cpsProperty.Value, stream, err)
			}
		}
	}

	return testCatalog, nil
}
