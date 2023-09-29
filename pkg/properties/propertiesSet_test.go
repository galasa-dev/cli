/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// MockServlet
func newUpdatePropertiesServletMock(t *testing.T) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockUpdatePropertiesServlet(t, w, r)
	}))

	return server
}

func mockUpdatePropertiesServlet(t *testing.T, w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.URL.Path, "/cps/") {
		t.Errorf("Expected to request '/cps/', got: %s", r.URL.Path)
	}
	if r.Header.Get("Accept") != "application/json" {
		t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
	}
	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var namespaceProperties string
	splitUrl := strings.Split(r.URL.Path, "/")
	namespace := splitUrl[2]

	//UPDATE -> cps/ns/properties/name
	statusCode, namespaceProperties = CheckNamespace(namespace)
	if len(splitUrl) == 5 {
		if namespace == "validNamespace" {
			propertyName := splitUrl[4]
			if propertyName == "invalidName" {
				statusCode = 404
				namespaceProperties = `{
					"error_code": 5018,
					"error_message": "GAL5018E: Error occured when trying to access property 'propertyName'. The property name provided is invalid."
					}`
			} else if propertyName == "validName" {
				statusCode = 200
				namespaceProperties = `{
						"name": "validNamespace.validName",
						"value": "updatedValue"
					}`
			} else if propertyName == "newName" {
				statusCode = 404
				namespaceProperties = `{
					"error_code": 5017,
					"error_message": "GAL5017E: Error occured when trying to access property 'newName'. The property name provided is invalid."
					}`
			}
		}
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(namespaceProperties))
}

// --------
// CREATING
func TestCreatePropertyWithValidNamespaceReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := "newName"
	value := "newValue"

	server := newUpdatePropertiesServletMock(t)
	apiServerUrl := server.URL
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := "Successfully created property '" + name + "' in namespace '" + namespace + "'"

	//When
	err := SetProperty(namespace, name, value, apiServerUrl, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestUpdatePropertyWithInvalidNamespaceAndInvalidPropertyNameReturnsError(t *testing.T) {
	//Given...
	namespace := "invalidNamespace"
	name := "newName"
	value := "newValue"

	server := newUpdatePropertiesServletMock(t)
	apiServerUrl := server.URL
	defer server.Close()

	console := utils.NewMockConsole()

	//When
	err := SetProperty(namespace, name, value, apiServerUrl, console)

	//Then
	assert.Error(t, err)
	assert.ErrorContains(t, err, "GAL1099E:")
}

// --------
// UPDATING
func TestUpdatePropertyWithValidNamespaceAndVaidNameValueReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := "validName"
	value := "updatedValue"

	server := newUpdatePropertiesServletMock(t)
	apiServerUrl := server.URL
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := "Successfully updated property '" + name + "' in namespace '" + namespace + "'"

	//When
	err := SetProperty(namespace, name, value, apiServerUrl, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestUpdatePropertyWithInvalidNamesapceAndValidNameReturnsError(t *testing.T) {
	//Given...
	namespace := "invalidNamespace"
	name := "validName"
	value := "updatedValue"

	server := newUpdatePropertiesServletMock(t)
	apiServerUrl := server.URL
	defer server.Close()

	console := utils.NewMockConsole()

	//When
	err := SetProperty(namespace, name, value, apiServerUrl, console)

	//Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1099E")
}
