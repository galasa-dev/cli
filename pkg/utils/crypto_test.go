/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanRoundTripText(t *testing.T) {
	secret := "mysecretisreallylongandhardtoguess34567890$@"
	text := "this is my text"
	encrypted, err := Encrypt(secret, text)
	assert.Nil(t, err)
	if err == nil {
		var decrypted string
		decrypted, err = Decrypt(secret, encrypted)
		assert.Nil(t, err)
		assert.Equal(t, text, decrypted)
	}
}

func TestCryptoWithNoSecretFails(t *testing.T) {
	secret := "" // too small to use as an encryption key
	text := "this is my text"
	_, err := Encrypt(secret, text)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1149")
}
