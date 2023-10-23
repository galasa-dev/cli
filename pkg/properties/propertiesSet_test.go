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

	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// MockServlet
func newSetPropertiesServletMock(t *testing.T) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockSetPropertiesServlet(t, w, r)
	}))

	return server
}

func mockSetPropertiesServlet(t *testing.T, w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.URL.Path, "/cps/") {
		t.Errorf("Expected to request '/cps/', got: %s", r.URL.Path)
	}
	if r.Header.Get("Accept") != "application/json" {
		t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
	}
	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var response string
	splitUrl := strings.Split(r.URL.Path, "/")
	namespace := splitUrl[2]

	statusCode, response = CheckNamespace(namespace)
	if namespace == "validNamespace" {
		if len(splitUrl) == 5 {
			propertyName := splitUrl[4]
			//UPDATE -> cps/ns/properties/name
			statusCode, response = updateProperty(propertyName)
		} else if len(splitUrl) == 4 {
			statusCode, response = createProperty()
		}
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(response))
}

func createProperty() (int, string) {
	statusCode := 201
	response := "Successfully created property newName in validNamespace"
	return statusCode, response
}

func updateProperty(propertyName string) (int, string) {
	statusCode := 200
	response := "Successfully updated property validName in validNamespace"

	if propertyName == "invalidName" {
		statusCode = 404
		response = `{
			"error_code": 5018,
			"error_message": "GAL5018E: Error occured when trying to access property 'propertyName'. The property name provided is invalid."
			}`
	} else if propertyName == "newName" { //property should be created
		statusCode = 404
		response = `{
			"error_code": 5017,
			"error_message": "GAL5017E: Error occured when trying to access property 'newName'. The property name provided is invalid."
			}`
	}

	return statusCode, response
}

// --------
// CREATING
// bad value feturn error
func TestCreatePropertyWithValidNamespaceReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := "newName"
	value := "newValue"

	server := newSetPropertiesServletMock(t)
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

	server := newSetPropertiesServletMock(t)
	apiServerUrl := server.URL
	defer server.Close()

	console := utils.NewMockConsole()

	//When
	err := SetProperty(namespace, name, value, apiServerUrl, console)

	//Then
	assert.Error(t, err)
	assert.ErrorContains(t, err, "GAL1098E:")
}

// --------
// UPDATING
func TestUpdatePropertyWithValidNamespaceAndVaidNameValueReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := "validName"
	value := "updatedValue"

	server := newSetPropertiesServletMock(t)
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

	server := newSetPropertiesServletMock(t)
	apiServerUrl := server.URL
	defer server.Close()

	console := utils.NewMockConsole()

	//When
	err := SetProperty(namespace, name, value, apiServerUrl, console)

	//Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1098E")
}
