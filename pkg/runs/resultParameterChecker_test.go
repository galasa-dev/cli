/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------
// Functions

func NewResultNamesServletMock(t *testing.T, status int) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ras/resultnames" {
			t.Errorf("Expected to request '/ras/resultnames', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		w.Write([]byte(`
			{
				"resultnames":["UNKNOWN","Passed","Failed","EnvFail"]
			}
		`))

	}))

	return server
}

func TestValidResultNamePassesValidation(t *testing.T) {
	// Given...
	server := NewResultNamesServletMock(t, http.StatusOK)
	defer server.Close()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()
	mockAuthenticator := utils.NewMockAuthenticatorWithAPIClient(apiClient)
	commsRetrier, _ := api.NewCommsRetrierWithAPIClient(1, 0, mockTimeService, mockAuthenticator)

	// When...
	result, err := ValidateResultParameter("Passed", commsRetrier)

	// Then...
	assert.Nil(t, err)
	assert.Equal(t, "Passed", result)
}

func TestValidResultNameLowercasePassesValidation(t *testing.T) {
	// Given...
	server := NewResultNamesServletMock(t, http.StatusOK)
	defer server.Close()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()
	mockAuthenticator := utils.NewMockAuthenticatorWithAPIClient(apiClient)
	commsRetrier, _ := api.NewCommsRetrierWithAPIClient(1, 0, mockTimeService, mockAuthenticator)

	// When...
	result, err := ValidateResultParameter("envfail", commsRetrier)

	// Then...
	assert.Nil(t, err)
	assert.Equal(t, "EnvFail", result)
}

func TestValidResultNameUppercasePassesValidation(t *testing.T) {
	// Given...
	server := NewResultNamesServletMock(t, http.StatusOK)
	defer server.Close()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()
	mockAuthenticator := utils.NewMockAuthenticatorWithAPIClient(apiClient)
	commsRetrier, _ := api.NewCommsRetrierWithAPIClient(1, 0, mockTimeService, mockAuthenticator)

	// When...
	result, err := ValidateResultParameter("FAILED", commsRetrier)

	// Then...
	assert.Nil(t, err)
	assert.Equal(t, "Failed", result)
}

func TestMultipleValidResultNamesMixedCasePassesValidation(t *testing.T) {
	// Given...
	server := NewResultNamesServletMock(t, http.StatusOK)
	defer server.Close()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()
	mockAuthenticator := utils.NewMockAuthenticatorWithAPIClient(apiClient)
	commsRetrier, _ := api.NewCommsRetrierWithAPIClient(1, 0, mockTimeService, mockAuthenticator)

	// When...
	result, err := ValidateResultParameter("unKnown,PASSed,Failed", commsRetrier)

	// Then...
	assert.Nil(t, err)
	assert.Equal(t, "UNKNOWN,Passed,Failed", result)
}

func TestGarbageResultNameFailsWithError(t *testing.T) {
	// Given...
	server := NewResultNamesServletMock(t, http.StatusOK)
	defer server.Close()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()
	mockAuthenticator := utils.NewMockAuthenticatorWithAPIClient(apiClient)
	commsRetrier, _ := api.NewCommsRetrierWithAPIClient(1, 0, mockTimeService, mockAuthenticator)

	// When...
	result, err := ValidateResultParameter("garbage", commsRetrier)

	// Then...
	assert.NotNil(t, err, "Should not have validated OK.")
	assert.ErrorContains(t, err, "GAL1087E")
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "'garbage'")
	assert.Contains(t, err.Error(), "'Passed'")
	assert.Contains(t, err.Error(), "'Failed'")
	assert.Contains(t, err.Error(), "'EnvFail'")
	assert.Contains(t, err.Error(), "'UNKNOWN'")
}

func TestValidResultFollowedByGarbageResultNameFailsWithError(t *testing.T) {
	// Given...
	server := NewResultNamesServletMock(t, http.StatusOK)
	defer server.Close()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()
	mockAuthenticator := utils.NewMockAuthenticatorWithAPIClient(apiClient)
	commsRetrier, _ := api.NewCommsRetrierWithAPIClient(1, 0, mockTimeService, mockAuthenticator)

	// When...
	result, err := ValidateResultParameter("passed,garbage", commsRetrier)

	// Then...
	assert.NotNil(t, err, "Should not have validated OK.")
	assert.ErrorContains(t, err, "GAL1087E")
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "'garbage'")
	assert.Contains(t, err.Error(), "'Passed'")
	assert.Contains(t, err.Error(), "'Failed'")
	assert.Contains(t, err.Error(), "'EnvFail'")
	assert.Contains(t, err.Error(), "'UNKNOWN'")
}

func TestMultipleGarbageResultNameFailsWithError(t *testing.T) {
	// Given...
	server := NewResultNamesServletMock(t, http.StatusOK)
	defer server.Close()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()
	mockAuthenticator := utils.NewMockAuthenticatorWithAPIClient(apiClient)
	commsRetrier, _ := api.NewCommsRetrierWithAPIClient(1, 0, mockTimeService, mockAuthenticator)

	// When...
	result, err := ValidateResultParameter("fail,garbage", commsRetrier)

	// Then...
	assert.NotNil(t, err, "Should not have validated OK.")
	assert.ErrorContains(t, err, "GAL1087E")
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "'fail'")
	assert.Contains(t, err.Error(), "'garbage'")
	assert.Contains(t, err.Error(), "'Passed'")
	assert.Contains(t, err.Error(), "'Failed'")
	assert.Contains(t, err.Error(), "'EnvFail'")
	assert.Contains(t, err.Error(), "'UNKNOWN'")
}

func TestMixOfValidAndGarbageResultNameFailsWithError(t *testing.T) {
	// Given...
	server := NewResultNamesServletMock(t, http.StatusOK)
	defer server.Close()
	apiServerUrl := server.URL
	apiClient := api.InitialiseAPI(apiServerUrl)
	mockTimeService := utils.NewMockTimeService()
	mockAuthenticator := utils.NewMockAuthenticatorWithAPIClient(apiClient)
	commsRetrier, _ := api.NewCommsRetrierWithAPIClient(1, 0, mockTimeService, mockAuthenticator)

	// When...
	result, err := ValidateResultParameter("passed,garbage,envfail,rubbish", commsRetrier)

	// Then...
	assert.NotNil(t, err, "Should not have validated OK.")
	assert.ErrorContains(t, err, "GAL1087E")
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "'garbage'")
	assert.Contains(t, err.Error(), "'rubbish'")
	assert.Contains(t, err.Error(), "'Passed'")
	assert.Contains(t, err.Error(), "'Failed'")
	assert.Contains(t, err.Error(), "'EnvFail'")
	assert.Contains(t, err.Error(), "'UNKNOWN'")
}
