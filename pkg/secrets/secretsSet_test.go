/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package secrets

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func readSecretRequestBody(req *http.Request) galasaapi.SecretRequest {
    var secretRequest galasaapi.SecretRequest
    requestBodyBytes, _ := io.ReadAll(req.Body)
    defer req.Body.Close()

    _ = json.Unmarshal(requestBodyBytes, &secretRequest)
    return secretRequest
}

func TestCanCreateAUsernameSecret(t *testing.T) {
    // Given...
    secretName := "SYSTEM1"
    username := "my-username"
    password := ""
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := ""
    secretType := ""

    // Create the expected HTTP interactions with the API server
    createSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)

    // Validate the request body
    createSecretInteraction.ValidateRequestFunc = func(t *testing.T, req *http.Request) {
        secretRequest := readSecretRequestBody(req)
        assert.Equal(t, secretRequest.GetName(), secretName)
        assert.Empty(t, secretRequest.GetType())

        requestUsername := secretRequest.GetUsername()
        assert.Equal(t, requestUsername.GetValue(), username)
        assert.Empty(t, requestUsername.GetEncoding())

        requestPassword := secretRequest.GetPassword()
        assert.Empty(t, requestPassword.GetValue())
        assert.Empty(t, requestPassword.GetEncoding())

        requestToken := secretRequest.GetToken()
        assert.Empty(t, requestToken.GetValue())
        assert.Empty(t, requestToken.GetEncoding())
    }

    createSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusCreated)
    }

    interactions := []utils.HttpInteraction{
        createSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.Nil(t, err, "SetSecret returned an unexpected error")
    assert.Empty(t, console.ReadText(), "The console was written to on a successful creation, it should be empty")
}

func TestCanCreateAUsernamePasswordSecret(t *testing.T) {
    // Given...
    secretName := "SYSTEM1"
    username := "my-username"
    password := "my-password"
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := ""
    secretType := ""

    // Create the expected HTTP interactions with the API server
    createSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)

    // Validate the request body
    createSecretInteraction.ValidateRequestFunc = func(t *testing.T, req *http.Request) {
        secretRequest := readSecretRequestBody(req)
        assert.Equal(t, secretRequest.GetName(), secretName)
        assert.Empty(t, secretRequest.GetType())

        requestUsername := secretRequest.GetUsername()
        assert.Equal(t, requestUsername.GetValue(), username)
        assert.Empty(t, requestUsername.GetEncoding())

        requestPassword := secretRequest.GetPassword()
        assert.Equal(t, requestPassword.GetValue(), password)
        assert.Empty(t, requestPassword.GetEncoding())

        requestToken := secretRequest.GetToken()
        assert.Empty(t, requestToken.GetValue())
        assert.Empty(t, requestToken.GetEncoding())
    }

    createSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusCreated)
    }

    interactions := []utils.HttpInteraction{
        createSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.Nil(t, err, "SetSecret returned an unexpected error")
    assert.Empty(t, console.ReadText(), "The console was written to on a successful creation, it should be empty")
}

func TestCanCreateAUsernameTokenSecret(t *testing.T) {
    // Given...
    secretName := "SYSTEM1"
    username := "my-username"
    password := ""
    token := "my-token"
    base64Username := ""
    base64Password := ""
    base64Token := ""
    secretType := ""

    // Create the expected HTTP interactions with the API server
    createSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)

    // Validate the request body
    createSecretInteraction.ValidateRequestFunc = func(t *testing.T, req *http.Request) {
        secretRequest := readSecretRequestBody(req)
        assert.Equal(t, secretRequest.GetName(), secretName)
        assert.Empty(t, secretRequest.GetType())

        requestUsername := secretRequest.GetUsername()
        assert.Equal(t, requestUsername.GetValue(), username)
        assert.Empty(t, requestUsername.GetEncoding())

        requestPassword := secretRequest.GetPassword()
        assert.Empty(t, requestPassword.GetValue())
        assert.Empty(t, requestPassword.GetEncoding())

        requestToken := secretRequest.GetToken()
        assert.Equal(t, requestToken.GetValue(), token)
        assert.Empty(t, requestToken.GetEncoding())
    }

    createSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusCreated)
    }

    interactions := []utils.HttpInteraction{
        createSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.Nil(t, err, "SetSecret returned an unexpected error")
    assert.Empty(t, console.ReadText(), "The console was written to on a successful creation, it should be empty")
}

func TestCanCreateATokenSecret(t *testing.T) {
    // Given...
    secretName := "SYSTEM1"
    username := ""
    password := ""
    token := "my-token"
    base64Username := ""
    base64Password := ""
    base64Token := ""
    secretType := ""

    // Create the expected HTTP interactions with the API server
    createSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)

    // Validate the request body
    createSecretInteraction.ValidateRequestFunc = func(t *testing.T, req *http.Request) {
        secretRequest := readSecretRequestBody(req)
        assert.Equal(t, secretRequest.GetName(), secretName)
        assert.Empty(t, secretRequest.GetType())

        requestUsername := secretRequest.GetUsername()
        assert.Empty(t, requestUsername.GetValue())
        assert.Empty(t, requestUsername.GetEncoding())

        requestPassword := secretRequest.GetPassword()
        assert.Empty(t, requestPassword.GetValue())
        assert.Empty(t, requestPassword.GetEncoding())

        requestToken := secretRequest.GetToken()
        assert.Equal(t, requestToken.GetValue(), token)
        assert.Empty(t, requestToken.GetEncoding())
    }

    createSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusCreated)
    }

    interactions := []utils.HttpInteraction{
        createSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.Nil(t, err, "SetSecret returned an unexpected error")
    assert.Empty(t, console.ReadText(), "The console was written to on a successful creation, it should be empty")
}

func TestCanUpdateASecret(t *testing.T) {
    // Given...
    secretName := "SYSTEM1"
    username := ""
    password := "my-new-password"
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := ""
    secretType := ""

    // Create the expected HTTP interactions with the API server
    updateSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)

    // Validate the request body
    updateSecretInteraction.ValidateRequestFunc = func(t *testing.T, req *http.Request) {
        secretRequest := readSecretRequestBody(req)
        assert.Equal(t, secretRequest.GetName(), secretName)
        assert.Empty(t, secretRequest.GetType())

        requestUsername := secretRequest.GetUsername()
        assert.Empty(t, requestUsername.GetValue())
        assert.Empty(t, requestUsername.GetEncoding())

        requestPassword := secretRequest.GetPassword()
        assert.Equal(t, requestPassword.GetValue(), password)
        assert.Empty(t, requestPassword.GetEncoding())

        requestToken := secretRequest.GetToken()
        assert.Empty(t, requestToken.GetValue())
        assert.Empty(t, requestToken.GetEncoding())
    }

    updateSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusNoContent)
    }

    interactions := []utils.HttpInteraction{
        updateSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.Nil(t, err, "SetSecret returned an unexpected error")
    assert.Empty(t, console.ReadText(), "The console was written to on a successful creation, it should be empty")
}

func TestCanUpdateAUsernamePasswordSecretInBase64Format(t *testing.T) {
    // Given...
    secretName := "SYSTEM1"
    username := ""
    password := ""
    token := ""
    base64Username := "my-base64-username"
    base64Password := "my-base64-password"
    base64Token := ""
    secretType := ""

    // Create the expected HTTP interactions with the API server
    updateSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)

    // Validate the request body
    updateSecretInteraction.ValidateRequestFunc = func(t *testing.T, req *http.Request) {
        secretRequest := readSecretRequestBody(req)
        assert.Equal(t, secretRequest.GetName(), secretName)
        assert.Empty(t, secretRequest.GetType())

        requestUsername := secretRequest.GetUsername()
        assert.Equal(t, requestUsername.GetValue(), base64Username)
        assert.Equal(t, requestUsername.GetEncoding(), BASE64_ENCODING)

        requestPassword := secretRequest.GetPassword()
        assert.Equal(t, requestPassword.GetValue(), base64Password)
        assert.Equal(t, requestPassword.GetEncoding(), BASE64_ENCODING)

        requestToken := secretRequest.GetToken()
        assert.Empty(t, requestToken.GetValue())
        assert.Empty(t, requestToken.GetEncoding())
    }

    updateSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusNoContent)
    }

    interactions := []utils.HttpInteraction{
        updateSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.Nil(t, err, "SetSecret returned an unexpected error")
    assert.Empty(t, console.ReadText(), "The console was written to on a successful creation, it should be empty")
}

func TestCanUpdateATokenSecretInBase64Format(t *testing.T) {
    // Given...
    secretName := "SYSTEM1"
    username := ""
    password := ""
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := "my-base64-token"
    secretType := ""

    // Create the expected HTTP interactions with the API server
    updateSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)

    // Validate the request body
    updateSecretInteraction.ValidateRequestFunc = func(t *testing.T, req *http.Request) {
        secretRequest := readSecretRequestBody(req)
        assert.Equal(t, secretRequest.GetName(), secretName)
        assert.Empty(t, secretRequest.GetType())

        requestUsername := secretRequest.GetUsername()
        assert.Empty(t, requestUsername.GetValue())
        assert.Empty(t, requestUsername.GetEncoding())

        requestPassword := secretRequest.GetPassword()
        assert.Empty(t, requestPassword.GetValue())
        assert.Empty(t, requestPassword.GetEncoding())

        requestToken := secretRequest.GetToken()
        assert.Equal(t, requestToken.GetValue(), base64Token)
        assert.Equal(t, requestToken.GetEncoding(), BASE64_ENCODING)
    }

    updateSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusNoContent)
    }

    interactions := []utils.HttpInteraction{
        updateSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.Nil(t, err, "SetSecret returned an unexpected error")
    assert.Empty(t, console.ReadText(), "The console was written to on a successful creation, it should be empty")
}

func TestCanUpdateASecretsTypeOk(t *testing.T) {
    // Given...
    secretName := "SYSTEM1"
    username := ""
    password := ""
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := "my-base64-token"
    secretType := "token"

    // Create the expected HTTP interactions with the API server
    updateSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)

    // Validate the request body
    updateSecretInteraction.ValidateRequestFunc = func(t *testing.T, req *http.Request) {
        secretRequest := readSecretRequestBody(req)
        assert.Equal(t, secretRequest.GetName(), secretName)
        assert.Equal(t, secretRequest.GetType(), galasaapi.TOKEN)

        requestUsername := secretRequest.GetUsername()
        assert.Empty(t, requestUsername.GetValue())
        assert.Empty(t, requestUsername.GetEncoding())

        requestPassword := secretRequest.GetPassword()
        assert.Empty(t, requestPassword.GetValue())
        assert.Empty(t, requestPassword.GetEncoding())

        requestToken := secretRequest.GetToken()
        assert.Equal(t, requestToken.GetValue(), base64Token)
        assert.Equal(t, requestToken.GetEncoding(), BASE64_ENCODING)
    }

    updateSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusNoContent)
    }

    interactions := []utils.HttpInteraction{
        updateSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.Nil(t, err, "SetSecret returned an unexpected error")
    assert.Empty(t, console.ReadText(), "The console was written to on a successful creation, it should be empty")
}

func TestUpdateSecretWithNoNameThrowsError(t *testing.T) {
    // Given...
    secretName := ""
    username := ""
    password := ""
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := "my-base64-token"
    secretType := ""

    // Validation should fail, so no HTTP interactions should take place
    interactions := []utils.HttpInteraction{}

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SetSecret did not return an error as expected")
    consoleOutputText := console.ReadText()
    assert.Contains(t, consoleOutputText, "GAL1172E")
    assert.Contains(t, consoleOutputText, "Invalid secret name provided")
}

func TestUpdateSecretWithUnknownTypeThrowsError(t *testing.T) {
    // Given...
    secretName := "MYSECRET"
    username := ""
    password := ""
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := "my-base64-token"
    secretType := "UNKNOWN"

    // Validation should fail, so no HTTP interactions should take place
    interactions := []utils.HttpInteraction{}

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SetSecret did not return an error as expected")
    consoleOutputText := console.ReadText()
    assert.Contains(t, consoleOutputText, "GAL1186E")
    assert.Contains(t, consoleOutputText, "Invalid secret type provided")
}

func TestUpdateSecretWithInvalidFlagCombinationThrowsError(t *testing.T) {
    // Given...
    // Provide a unencoded credentials and base64-encoded ones
    secretName := "MYSECRET"
    username := "my-username"
    password := "my-password"
    token := "my-token"
    base64Username := "my-base64-username"
    base64Password := "my-base64-password"
    base64Token := "my-base64-token"
    secretType := ""

    // Validation should fail, so no HTTP interactions should take place
    interactions := []utils.HttpInteraction{}

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SetSecret did not return an error as expected")
    consoleOutputText := console.ReadText()
    assert.Contains(t, consoleOutputText, "GAL1193E")
    assert.Contains(t, consoleOutputText, "Invalid flag combination provided")
}

func TestSetSecretFailsWithNoExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := "MYSECRET"
    username := ""
    password := ""
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := "my-base64-token"
    secretType := ""

    // Create the expected HTTP interactions with the API server
    updateSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)
    updateSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusInternalServerError)
    }

    interactions := []utils.HttpInteraction{
        updateSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SetSecret did not return an error but it should have")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText , secretName)
    assert.Contains(t, consoleText , "GAL1187E")
    assert.Contains(t, consoleText , "Unexpected http status code 500 received from the server")
}

func TestSetSecretFailsWithNonJsonContentTypeExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := "MYSECRET"
    username := ""
    password := ""
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := "my-base64-token"
    secretType := ""

    // Create the expected HTTP interactions with the API server
    updateSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)
    updateSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Header().Set("Content-Type", "application/notJsonOnPurpose")
        writer.Write([]byte("something not json but non-zero-length."))
    }

    interactions := []utils.HttpInteraction{
        updateSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SetSecret did not return an error but it should have")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText, secretName)
    assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
    assert.Contains(t, consoleText, "GAL1191E")
    assert.Contains(t, consoleText, "Error details from the server are not in the json format")
}

func TestSetSecretFailsWithBadlyFormedJsonContentExplanationErrorPayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := "MYSECRET"
    username := ""
    password := ""
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := "my-base64-token"
    secretType := ""

    // Create the expected HTTP interactions with the API server
    updateSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)
    updateSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(`{ "this": "isBadJson because it doesnt end in a close braces" `))
    }

    interactions := []utils.HttpInteraction{
        updateSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SetSecret did not return an error but it should have")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText, secretName)
    assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
    assert.Contains(t, consoleText, "GAL1189E")
    assert.Contains(t, consoleText, "Error details from the server are not in a valid json format")
    assert.Contains(t, consoleText, "Cause: 'unexpected end of JSON input'")
}

func TestSetSecretFailsWithValidErrorResponsePayloadGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := "MYSECRET"
    username := ""
    password := ""
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := "my-base64-token"
    secretType := ""
    apiErrorCode := 5000
    apiErrorMessage := "this is an error from the API server"

    // Create the expected HTTP interactions with the API server
    updateSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)
    updateSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusInternalServerError)

        apiError := errors.GalasaAPIError{
            Code: apiErrorCode,
            Message: apiErrorMessage,
        }
        apiErrorBytes, _ := json.Marshal(apiError)
        writer.Write(apiErrorBytes)
    }

    interactions := []utils.HttpInteraction{
        updateSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReader()

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SetSecret did not return an error but it should have")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText, secretName)
    assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
    assert.Contains(t, consoleText, "GAL1190E")
    assert.Contains(t, consoleText, apiErrorMessage)
}

func TestSecretsSetFailsWithFailureToReadResponseBodyGivesCorrectMessage(t *testing.T) {
    // Given...
    secretName := "MYSECRET"
    username := ""
    password := ""
    token := ""
    base64Username := ""
    base64Password := ""
    base64Token := "my-base64-token"
    secretType := ""

    // Create the expected HTTP interactions with the API server
    updateSecretInteraction := utils.NewHttpInteraction("/secrets/" + secretName, http.MethodPut)
    updateSecretInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(http.StatusInternalServerError)
        writer.Write([]byte(`{}`))
    }

    interactions := []utils.HttpInteraction{
        updateSecretInteraction,
    }

    server := utils.NewMockHttpServer(t, interactions)
    defer server.Server.Close()

    console := utils.NewMockConsole()
    apiServerUrl := server.Server.URL
    apiClient := api.InitialiseAPI(apiServerUrl)
    mockByteReader := utils.NewMockByteReaderAsMock(true)

    // When...
    err := SetSecret(
        secretName,
        username,
        password,
        token,
        base64Username,
        base64Password,
        base64Token,
        secretType,
        console,
        apiClient,
        mockByteReader)

    // Then...
    assert.NotNil(t, err, "SetSecret did not return an error but it should have")
    consoleText := console.ReadText()
    assert.Contains(t, consoleText, secretName)
    assert.Contains(t, consoleText, strconv.Itoa(http.StatusInternalServerError))
    assert.Contains(t, consoleText, "GAL1188E")
    assert.Contains(t, consoleText, "Error details from the server could not be read")
}
