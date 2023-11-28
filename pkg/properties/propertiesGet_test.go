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

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// MockServlet
func NewPropertiesServletMock(t *testing.T) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MockPropertiesServlet(t, w, r)
	}))

	return server
}

func MockPropertiesServlet(t *testing.T, w http.ResponseWriter, r *http.Request) {
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

	//cps/ns/properties
	statusCode, namespaceProperties = CheckNamespace(namespace)
	if statusCode == 200 {
		if len(splitUrl) == 4 {
			query := r.URL.Query()
			//cps/ns/properties?prefix&suffix&infix
			if query.Has("prefix") || query.Has("suffix") || query.Has("infix") {
				queryValues := r.URL.Query()
				prefixParameter := queryValues.Get("prefix")
				suffixParameter := queryValues.Get("suffix")
				infixParameter := queryValues.Get("infix")

				namespaceProperties = checkQueryParameters(prefixParameter, suffixParameter, infixParameter)
			}
			//cps/ns/properties/propertyname
		} else if len(splitUrl) == 5 {
			propertyName := splitUrl[4]
			namespaceProperties, statusCode = CheckName(propertyName)
		}
	}
	w.WriteHeader(statusCode)
	w.Write([]byte(namespaceProperties))
}

func CheckName(name string) (string, int) {
	statusCode := 200
	namespaceProperties := "[]"
	switch name {
	case "property0":
		namespaceProperties = `[
		{
			"name": "validNamespace.property0",
			"value": "value0"
		}
	]`
	case "invalidName": //property name does not exist
		namespaceProperties = `[]`
	case "emptyValueName": //property name does not exist
		namespaceProperties = `[
			{
				"name": "validNamespace.emptyValueName",
				"value": ""
			}
		]`
	}
	return namespaceProperties, statusCode
}

func checkQueryParameters(prefixParameter string, suffixParameter string, infixParameter string) string {
	var namespaceProperties = ""
	//there are properties in the namespace that match a prefix and/or suffix
	if prefixParameter == "aPrefix" && suffixParameter == "aSuffix" {

		if infixParameter == "anInfix" { //for a single infix
			namespaceProperties = `[{"name": "validNamespace.aPrefix.anInfix.property.aSuffix","value": "prefixSuffixInfixVal"}]`
		} else { //no infix
			namespaceProperties = `[{"name": "validNamespace.aPrefix.property.aSuffix","value": "prefixSuffixVal"}]`
		}

	} else if suffixParameter == "aSuffix" {

		if infixParameter == "anInfix" {
			namespaceProperties = `[{"name":"validNamespace.property.anInfix.aSuffix", "value":"suffixInfixVal"}]`
		} else {
			namespaceProperties = `[{"name":"validNamespace.property.aSuffix", "value":"suffixVal"}]`
		}

	} else if prefixParameter == "aPrefix" {

		if infixParameter == "anInfix" {
			namespaceProperties = `[{"name":"validNamespace.aPrefix.property.anInfix", "value":"prefixInfixVal"}]`
		} else {
			namespaceProperties = `[{"name":"validNamespace.aPrefix.property", "value":"prefixVal"}]`
		}

	}
	//there are NO properties in the namespace that match the prefix and/or suffix
	if prefixParameter == "noMatchingPrefix" && suffixParameter == "noMatchingSuffix" {
		namespaceProperties = `[]`
	} else if suffixParameter == "noMatchingSuffix" {
		namespaceProperties = `[]`
	} else if prefixParameter == "noMatchingPrefix" {
		namespaceProperties = `[]`
	}

	//If only the infix parameter is supplied
	if prefixParameter == "" && suffixParameter == "" {
		if infixParameter == "anInfix" {
			namespaceProperties = `[{"name":"validNamespace.extra.anInfix.extra", "value":"infixVal"}]`
		} else if infixParameter == "noMatchingInfix" { //singular or multiple infixes that do not match
			namespaceProperties = `[]`
		}
	}

	return namespaceProperties
}

func CheckNamespace(namespace string) (int, string) {
	statusCode := 200
	namespaceProperties := "[]"

	switch namespace {
	case "validNamespace":
		namespaceProperties = `[
			{
				"name": "validNamespace.property0",
				"value": "value0"
			},
			{
				"name": "validNamespace.property1",
				"value": "value1"
			},
			{
				"name": "validNamespace.property2",
				"value": "value2"
			},
			{
				"name": "validNamespace.property3",
				"value": "value3"
			}
		]`
	case "invalidNamespace":
		statusCode = 404
		namespaceProperties = `{
		error_code: 5016,
		error_message: "GAL5016E: Error occured when trying to access namespace 'invalidNamespace'. The Namespace provided is invalid."
		}`
	}

	return statusCode, namespaceProperties
}

// -------------------
// TESTS

func TestInvalidNamepsaceReturnsError(t *testing.T) {
	//Given...
	namespace := "invalidNamespace"
	name := ""
	prefix := ""
	suffix := ""
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Error(t, err)
	assert.ErrorContains(t, err, "GAL1096E")
}

func TestValidNamespaceReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := ""
	suffix := ""
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `namespace      name      value
validNamespace property0 value0
validNamespace property1 value1
validNamespace property2 value2
validNamespace property3 value3

Total:4
`
	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestEmptyNamespaceReturnsEmpty(t *testing.T) {
	//Given...
	namespace := "emptyNamespace"
	name := ""
	prefix := ""
	suffix := ""
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `Total:0
`
	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceAndPrefixReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := "aPrefix"
	suffix := ""
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `namespace      name             value
validNamespace aPrefix.property prefixVal

Total:1
`

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceAndSuffixReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := ""
	suffix := "aSuffix"
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `namespace      name             value
validNamespace property.aSuffix suffixVal

Total:1
`

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceWithMatchingPrefixAndSuffixReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := "aPrefix"
	suffix := "aSuffix"
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `namespace      name                     value
validNamespace aPrefix.property.aSuffix prefixSuffixVal

Total:1
`

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceWithNoMatchingPrefixAndSuffixReturnsEmpty(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := "noMatchingPrefix"
	suffix := "noMatchingSuffix"
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `Total:0
`

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceAndNoMatchingPrefixReturnsEmpty(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := "noMatchingPrefix"
	suffix := ""
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `Total:0
`

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceAndNoMatchingSuffixReturnsEmpty(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := ""
	suffix := "noMatchingSuffix"
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `Total:0
`

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceWithMatchingInfixReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := ""
	suffix := ""
	infix := "anInfix"
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `namespace      name                value
validNamespace extra.anInfix.extra infixVal

Total:1
`
	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceWithNoMatchingInfixReturnsEmpty(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := ""
	suffix := ""
	infix := "noMatchingInfix"
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `Total:0
`
	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceWithMatchingPrefixAndInfixReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := "aPrefix"
	suffix := ""
	infix := "anInfix"
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `namespace      name                     value
validNamespace aPrefix.property.anInfix prefixInfixVal

Total:1
`

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceWithMatchingSuffixAndInfixReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := ""
	suffix := "aSuffix"
	infix := "anInfix"
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `namespace      name                     value
validNamespace property.anInfix.aSuffix suffixInfixVal

Total:1
`

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceWithMatchingPrefixAndSuffixAndInfixReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := "aPrefix"
	suffix := "aSuffix"
	infix := "anInfix"
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `namespace      name                             value
validNamespace aPrefix.anInfix.property.aSuffix prefixSuffixInfixVal

Total:1
`

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceWithValidNameReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := "property0"
	prefix := ""
	suffix := ""
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `namespace      name      value
validNamespace property0 value0

Total:1
`
	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNameWithEmptyValueValidNameReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := "emptyValueName"
	prefix := ""
	suffix := ""
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `namespace      name           value
validNamespace emptyValueName 

Total:1
`

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestInvalidPropertyNameReturnsEmpty(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := "invalidName"
	prefix := ""
	suffix := ""
	infix := ""
	propertiesOutputFormat := "summary"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `Total:0
`

	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}

func TestValidNamespaceFormatRawReturnsOk(t *testing.T) {
	//Given...
	namespace := "validNamespace"
	name := ""
	prefix := ""
	suffix := ""
	infix := ""
	propertiesOutputFormat := "raw"

	server := NewPropertiesServletMock(t)
	apiClient := api.InitialiseAPI(server.URL)
	defer server.Close()

	console := utils.NewMockConsole()
	expectedOutput := `validNamespace|property0|value0
validNamespace|property1|value1
validNamespace|property2|value2
validNamespace|property3|value3
`
	//When
	err := GetProperties(namespace, name, prefix, suffix, infix, apiClient, propertiesOutputFormat, console)

	//Then
	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())
}
