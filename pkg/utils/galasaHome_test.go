/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/stretchr/testify/assert"
)

func TestDefaultHomePathTakenFromFileSystem(t *testing.T) {
	// Given
	fs := files.NewMockFileSystem()
	env := NewMockEnv()

	// When
	galasaHome, err := NewGalasaHome(fs, env, "")

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, galasaHome)
	assert.Equal(t, galasaHome.GetNativeFolderPath(), "/User/Home/testuser/.galasa")
}

func TestHomePathTakenFromEnvVarIfSet(t *testing.T) {
	// Given
	fs := files.NewMockFileSystem()
	env := NewMockEnv()
	env.SetEnv("GALASA_HOME", "AnyWhereIwantItToBe")

	// When
	galasaHome, err := NewGalasaHome(fs, env, "")

	// Then
	assert.Nil(t, err)
	assert.NotNil(t, galasaHome)
	assert.Equal(t, galasaHome.GetNativeFolderPath(), "AnyWhereIwantItToBe")
}
