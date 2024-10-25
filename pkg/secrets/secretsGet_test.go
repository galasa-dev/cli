/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package secrets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const (
    API_VERSION = "galasa-dev/v1alpha1"
    DUMMY_ENCODING = "myencoding"
    DUMMY_USERNAME = "dummy-username"
    DUMMY_PASSWORD = "dummy-password"
)

func createMockGalasaSecret(secretName string) galasaapi.GalasaSecret {
    secret := *galasaapi.NewGalasaSecret()

    secret.SetApiVersion(API_VERSION)
    secret.SetKind("GalasaSecret")

    secretMetadata := *galasaapi.NewGalasaSecretMetadata()
    secretMetadata.SetName(secretName)
    secretMetadata.SetEncoding(DUMMY_ENCODING)
    secretMetadata.SetType("UsernamePassword")

    secretData := *galasaapi.NewGalasaSecretData()
    secretData.SetUsername(DUMMY_USERNAME)
    secretData.SetPassword(DUMMY_PASSWORD)

    secret.SetMetadata(secretMetadata)
    secret.SetData(secretData)
    return secret
}

func generateExpectedSecretYaml(secretName string) string {
    return fmt.Sprintf(`apiVersion: %s
kind: GalasaSecret
metadata:
    name: %s
    encoding: %s
    type: UsernamePassword
data:
    username: %s
    password: %s
    token: null`, API_VERSION, secretName, DUMMY_ENCODING, DUMMY_USERNAME, DUMMY_PASSWORD)
}

func TestCanGetASecretByName(t *testing.T) {
    // Given...
    secretName := "SYSTEM1"
    outputFormat := "summary"

    // Create the mock secret to return
    secret := createMockGalasaSecret(secretName)
    secretBytes, _ := json.Marshal(secret)
    secretJson := string(secretBytes)

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusOK)
        writer.Write([]byte(secretJson))
    }

    interactions := []utils.HttpInteraction{
        getSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    expectedOutput := fmt.Sprintf(`name    type
%s UsernamePassword

Total:1
`, secretName)
    assert.Nil(t, err, "GetSecrets returned an unexpected error")
    assert.Equal(t, expectedOutput, console.ReadText())
}

func TestCanGetASecretByNameInYamlFormat(t *testing.T) {
    // Given...
    secretName := "SYSTEM1"
    outputFormat := "yaml"

    // Create the mock secret to return
    secret := createMockGalasaSecret(secretName)
    secretBytes, _ := json.Marshal(secret)
    secretJson := string(secretBytes)

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusOK)
        writer.Write([]byte(secretJson))
    }

    interactions := []utils.HttpInteraction{
        getSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    expectedOutput := generateExpectedSecretYaml(secretName) + "\n"
    assert.Nil(t, err, "GetSecrets returned an unexpected error")
    assert.Equal(t, expectedOutput, console.ReadText())
}

func TestCanGetAllSecretsOk(t *testing.T) {
    // Given...
    // Don't provide a secret name so that we can get all secrets
    secretName := ""
    outputFormat := "summary"

    // Create the mock secret to return
    secrets := make([]galasaapi.GalasaSecret, 0)
    secret1Name := "BOB"
    secret2Name := "BLAH"
    secret1 := createMockGalasaSecret(secret1Name)
    secret2 := createMockGalasaSecret(secret2Name)

    secrets = append(secrets, secret1, secret2)
    secretsBytes, _ := json.Marshal(secrets)
    secretsJson := string(secretsBytes)

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets", http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusOK)
        writer.Write([]byte(secretsJson))
    }

    interactions := []utils.HttpInteraction{
        getSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    expectedOutput := fmt.Sprintf(`name type
%s  UsernamePassword
%s UsernamePassword

Total:2
`, secret1Name, secret2Name)
    assert.Nil(t, err, "GetSecrets returned an unexpected error")
    assert.Equal(t, expectedOutput, console.ReadText())
}

func TestGetASecretWithUnknownFormatDisplaysError(t *testing.T) {
    // Given...
    secretName := "MYSECRET"
    outputFormat := "UNKNOWN FORMAT!"

    // The client-side validation should fail, so no HTTP interactions will be performed
    interactions := []utils.HttpInteraction{}

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "GetSecrets did not return an error as expected")
    consoleOutputText := console.ReadText()
    assert.Contains(t, consoleOutputText, "GAL1067E")
    assert.Contains(t, consoleOutputText, "Unsupported value 'UNKNOWN FORMAT!'")
    assert.Contains(t, consoleOutputText, "'summary', 'yaml'")
}

func TestGetASecretWithBlankNameDisplaysError(t *testing.T) {
    // Given...
    secretName := "    "
    outputFormat := "summary"

    // The client-side validation should fail, so no HTTP interactions will be performed
    interactions := []utils.HttpInteraction{}

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "GetSecrets did not return an error as expected")
    consoleOutputText := console.ReadText()
    assert.Contains(t, consoleOutputText, "GAL1172E")
    assert.Contains(t, consoleOutputText, " Invalid secret name provided")
}

func TestGetNonExistantSecretDisplaysError(t *testing.T) {
    // Given...
    nonExistantSecret := "secretDoesNotExist123"
    outputFormat := "summary"

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets/" + nonExistantSecret, http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusNotFound)
        writer.Write([]byte(`{ "error_message": "No such secret exists" }`))
    }


    interactions := []utils.HttpInteraction{ getSecretInteraction }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        nonExistantSecret,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SecretsGet did not return an error but it should have")
    consoleOutputText := console.ReadText()
    assert.Contains(t, consoleOutputText, nonExistantSecret)
    assert.Contains(t, consoleOutputText, "GAL1177E")
    assert.Contains(t, consoleOutputText, "Error details from the server are: 'No such secret exists'")
}

func TestSecretsGetFailsWithNoExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := "MYSECRET"
    outputFormat := "summary"

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusInternalServerError)
    }

    interactions := []utils.HttpInteraction{
        getSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SecretsGet did not return an error but it should have")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText , secretName)
    assert.Contains(t, consoleText , "GAL1174E")
}

func TestSecretsGetFailsWithNonJsonContentTypeExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := "MYSECRET"
    outputFormat := "summary"

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Header().Set("Content-Type", "application/notJsonOnPurpose")
        writer.Write([]byte("something not json but non-zero-length."))
    }

    interactions := []utils.HttpInteraction{
        getSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SecretsGet did not return an error but it should have")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText, secretName)
    assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
    assert.Contains(t, consoleText, "GAL1178E")
    assert.Contains(t, consoleText, "Error details from the server are not in the json format")
}

func TestSecretsGetFailsWithBadlyFormedJsonContentExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := "MYSECRET"
    outputFormat := "summary"

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(`{ "this": "isBadJson because it doesnt end in a close braces" `))
    }

    interactions := []utils.HttpInteraction{
        getSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SecretsGet did not return an error but it should have")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText, secretName)
    assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
    assert.Contains(t, consoleText, "GAL1176E")
    assert.Contains(t, consoleText, "Error details from the server are not in a valid json format")
    assert.Contains(t, consoleText, "Cause: 'unexpected end of JSON input'")
}

func TestSecretsGetFailsWithFailureToReadResponseBodyGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := "MYSECRET"
    outputFormat := "summary"

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(`{}`))
    }

    interactions := []utils.HttpInteraction{
        getSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReaderAsMock(true)

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SecretsGet returned an unexpected error")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText, secretName)
    assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
    assert.Contains(t, consoleText, "GAL1175E")
    assert.Contains(t, consoleText, "Error details from the server could not be read")
}

func TestGetAllSecretsFailsWithNoExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := ""
    outputFormat := "summary"

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets", http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusInternalServerError)
    }

    interactions := []utils.HttpInteraction{
        getSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SecretsGet did not return an error but it should have")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText , "GAL1180E")
}

func TestGetAllSecretsFailsWithNonJsonContentTypeExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := ""
    outputFormat := "summary"

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets", http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Header().Set("Content-Type", "application/notJsonOnPurpose")
        writer.Write([]byte("something not json but non-zero-length."))
    }

    interactions := []utils.HttpInteraction{
        getSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SecretsGet did not return an error but it should have")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
    assert.Contains(t, consoleText, "GAL1184E")
    assert.Contains(t, consoleText, "Error details from the server are not in the json format")
}

func TestGetAllSecretsFailsWithBadlyFormedJsonContentExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := ""
    outputFormat := "summary"

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets", http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(`{ "this": "isBadJson because it doesnt end in a close braces" `))
    }

    interactions := []utils.HttpInteraction{
        getSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SecretsGet did not return an error but it should have")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
    assert.Contains(t, consoleText, "GAL1182E")
    assert.Contains(t, consoleText, "Error details from the server are not in a valid json format")
    assert.Contains(t, consoleText, "Cause: 'unexpected end of JSON input'")
}

func TestGetAllSecretsFailsWithFailureToReadResponseBodyGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := ""
    outputFormat := "summary"

    // Create the expected HTTP interactions with the API server
    getSecretInteraction := utils.NewHttpInteraction("/secrets", http.MethodGet)
    getSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(`{}`))
    }

    interactions := []utils.HttpInteraction{
        getSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReaderAsMock(true)

    // When...
    err := GetSecrets(
        secretName,
        outputFormat,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SecretsGet returned an unexpected error")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
    assert.Contains(t, consoleText, "GAL1181E")
    assert.Contains(t, consoleText, "Error details from the server could not be read")
}
