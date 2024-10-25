/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package secretsformatter

import (
    "testing"

    "github.com/galasa-dev/cli/pkg/galasaapi"
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

func TestSecretSummaryFormatterNoDataReturnsTotalCountAllZeros(t *testing.T) {
    // Given...
    formatter := NewSecretSummaryFormatter()
    secrets := make([]galasaapi.GalasaSecret, 0)

    // When...
    actualFormattedOutput, err := formatter.FormatSecrets(secrets)

    // Then...
    assert.Nil(t, err)
    expectedFormattedOutput := "Total:0\n"
    assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSecretSummaryFormatterSingleDataReturnsCorrectly(t *testing.T) {
    // Given...
    formatter := NewSecretSummaryFormatter()
    secret1 := createMockGalasaSecret("MYSECRET")
    secrets := []galasaapi.GalasaSecret{ secret1 }

    // When...
    actualFormattedOutput, err := formatter.FormatSecrets(secrets)

    // Then...
    assert.Nil(t, err)
    expectedFormattedOutput :=
        `name     type
MYSECRET UsernamePassword

Total:1
`
    assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSecretSummaryFormatterMultipleDataSeperatesWithNewLine(t *testing.T) {
    // Given..
    formatter := NewSecretSummaryFormatter()
    secrets := make([]galasaapi.GalasaSecret, 0)
    secret1 := createMockGalasaSecret("SECRET1")
    secret2 := createMockGalasaSecret("SECRET_2")
    secret3 := createMockGalasaSecret("SECRET-3")
    secrets = append(secrets, secret1, secret2, secret3)

    // When...
    actualFormattedOutput, err := formatter.FormatSecrets(secrets)

    // Then...
    assert.Nil(t, err)
    expectedFormattedOutput := `name     type
SECRET1  UsernamePassword
SECRET_2 UsernamePassword
SECRET-3 UsernamePassword

Total:3
`
    assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}
