/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package launcher

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
)

// RemoteLauncher A launcher, which launches and monitors tests on a remote ecosystem via HTTP/HTTPS.
type RemoteLauncher struct {
	apiClient *galasaapi.APIClient
}

//----------------------------------------------------------------------------------
// Constructors
//----------------------------------------------------------------------------------

// NewRemoteLauncher create a remote launcher.
func NewRemoteLauncher(apiServerUrl string, apiClient *galasaapi.APIClient) *RemoteLauncher {
	log.Printf("NewRemoteLauncher(%s) entered.", apiServerUrl)

	launcher := new(RemoteLauncher)

	// An HTTP client which can communicate with the api server in an ecosystem.
	launcher.apiClient = apiClient

	return launcher
}

//----------------------------------------------------------------------------------
// Implementation of the launcher interface.
//----------------------------------------------------------------------------------

// GetRunsByGroup get all the testruns which are associated with a named group.
func (launcher *RemoteLauncher) GetRunsByGroup(groupName string) (*galasaapi.TestRuns, error) {
	log.Printf("GetRunsByGroup(%s) entered.", groupName)
	var (
		testRuns       *galasaapi.TestRuns
		err            error
		restApiVersion string
	)
	restApiVersion, err = embedded.GetGalasactlRestApiVersion()
	if err == nil {
		testRuns, _, err = launcher.apiClient.RunsAPIApi.GetRunsGroup(context.TODO(), groupName).ClientApiVersion(restApiVersion).Execute()
	}
	return testRuns, err
}

func (launcher *RemoteLauncher) SubmitTestRun(
	groupName string,
	className string,
	requestType string,
	requestor string,
	stream string,
	obrFromPortfolio string,
	isTraceEnabled bool,
	GherkinURL string,
	GherkinFeature string,
	overrides map[string]interface{},
) (*galasaapi.TestRuns, error) {

	classNames := make([]string, 1)
	classNames[0] = className

	testRunRequest := galasaapi.NewTestRunRequest()
	testRunRequest.SetClassNames(classNames)
	testRunRequest.SetRequestorType(requestType)
	testRunRequest.SetRequestor(requestor)
	testRunRequest.SetTestStream(stream)
	testRunRequest.SetTrace(isTraceEnabled)
	testRunRequest.SetOverrides(overrides)

	log.Printf("RemoteLauncher.SubmitTestRuns : using requestor %s\n", requestor)

	var resultGroup *galasaapi.TestRuns
	var err error
	var restApiVersion string

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {
		resultGroup, _, err = launcher.apiClient.RunsAPIApi.PostSubmitTestRuns(context.TODO(), groupName).TestRunRequest(*testRunRequest).ClientApiVersion(restApiVersion).Execute()
	}
	return resultGroup, err
}

func (launcher *RemoteLauncher) GetRunsById(runId string) (*galasaapi.Run, error) {
	var err error
	var rasRun *galasaapi.Run
	var restApiVersion string

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {
		rasRun, _, err = launcher.apiClient.ResultArchiveStoreAPIApi.GetRasRunById(context.TODO(), runId).ClientApiVersion(restApiVersion).Execute()
	}
	return rasRun, err
}

func (launcher *RemoteLauncher) GetStreams() ([]string, error) {

	var streams []string

	var restApiVersion string
	var err error

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {
		var properties []galasaapi.GalasaProperty
		properties, _, err = launcher.apiClient.ConfigurationPropertyStoreAPIApi.
			QueryCpsNamespaceProperties(context.TODO(), "framework").Prefix("test.stream").Suffix("repo").ClientApiVersion(restApiVersion).Execute()
		if err == nil {

			streams, err = getStreamNamesFromProperties(properties)
		}
	}
	return streams, err
}

// When passed an array of GalasaProperty objects, extract the stream names from them.
func getStreamNamesFromProperties(properties []galasaapi.GalasaProperty) ([]string, error) {
	var err error
	var streams []string = make([]string, 0)
	for _, property := range properties {
		propertyNamePtr := property.GetMetadata().Name

		propertyName := *propertyNamePtr
		// This is something like "test.stream.zosk8s.repo"

		streamName := propertyName[12 : len(propertyName)-5]
		streams = append(streams, streamName)
	}
	return streams, err
}

func (launcher *RemoteLauncher) GetTestCatalog(stream string) (TestCatalog, error) {

	var err error
	var testCatalog TestCatalog
	var cpsProperty []galasaapi.GalasaProperty
	var restApiVersion string

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {
		cpsProperty, _, err = launcher.apiClient.ConfigurationPropertyStoreAPIApi.QueryCpsNamespaceProperties(context.TODO(), "framework").Prefix("test.stream."+stream).Suffix("location").ClientApiVersion(restApiVersion).Execute()
		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PROPERTY_GET_FAILED, stream, err)
		} else if len(cpsProperty) < 1 {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CATALOG_NOT_FOUND, stream)
		}
		if err == nil {
			streamLocation :=cpsProperty[0].Data.Value
			catalogString := new(strings.Builder)
			var resp *http.Response
			resp, err = http.Get(*streamLocation)
			if err != nil {
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_GET_TEST_CATALOG_CONTENTS_FAILED, *streamLocation, stream, err)
			} else {
				defer resp.Body.Close()

				_, err = io.Copy(catalogString, resp.Body)
				if err != nil {
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CATALOG_COPY_FAILED, *streamLocation, stream, err)
				}
			}

			if err == nil {
				err = json.Unmarshal([]byte(catalogString.String()), &testCatalog)
				if err != nil {
					err = galasaErrors.NewGalasaError(
						galasaErrors.GALASA_ERROR_CATALOG_UNMARSHAL_FAILED, *streamLocation, stream, err)
				}
			}
		}
	}

	return testCatalog, err
}
