/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package secretsformatter

import (
	"fmt"
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

func createMockGalasaSecretWithDescription(secretName string, description string) galasaapi.GalasaSecret {
    secret := *galasaapi.NewGalasaSecret()

    secret.SetApiVersion(API_VERSION)
    secret.SetKind("GalasaSecret")

    secretMetadata := *galasaapi.NewGalasaSecretMetadata()
    secretMetadata.SetName(secretName)
    secretMetadata.SetEncoding(DUMMY_ENCODING)
    secretMetadata.SetType("UsernamePassword")

	if description != "" {
		secretMetadata.SetDescription(description)
	}

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
	description := "secret for system1"
	secretName := "MYSECRET"
    secret1 := createMockGalasaSecretWithDescription(secretName, description)
    secrets := []galasaapi.GalasaSecret{ secret1 }

    // When...
    actualFormattedOutput, err := formatter.FormatSecrets(secrets)

    // Then...
    assert.Nil(t, err)
    expectedFormattedOutput := fmt.Sprintf(
`name     type             description
%s UsernamePassword %s

Total:1
`, secretName, description)
    assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}

func TestSecretSummaryFormatterMultipleDataSeperatesWithNewLine(t *testing.T) {
    // Given..
    formatter := NewSecretSummaryFormatter()
    secrets := make([]galasaapi.GalasaSecret, 0)

	secret1Name := "SECRET1"
	secret1Description := "my first secret"
	secret2Name := "SECRET2"
	secret2Description := "my second secret"
	secret3Name := "SECRET3"
	secret3Description := "my third secret"

    secret1 := createMockGalasaSecretWithDescription(secret1Name, secret1Description)
    secret2 := createMockGalasaSecretWithDescription(secret2Name, secret2Description)
    secret3 := createMockGalasaSecretWithDescription(secret3Name, secret3Description)
    secrets = append(secrets, secret1, secret2, secret3)

    // When...
    actualFormattedOutput, err := formatter.FormatSecrets(secrets)

    // Then...
    assert.Nil(t, err)
    expectedFormattedOutput := fmt.Sprintf(
`name    type             description
%s UsernamePassword %s
%s UsernamePassword %s
%s UsernamePassword %s

Total:3
`, secret1Name, secret1Description, secret2Name, secret2Description, secret3Name, secret3Description)
    assert.Equal(t, expectedFormattedOutput, actualFormattedOutput)
}
