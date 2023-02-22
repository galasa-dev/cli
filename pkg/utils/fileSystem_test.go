/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTildaExpansionWhenFilenameBlankReturnsBlank(t *testing.T) {
	pathIn := ""
	fs := NewMockFileSystem()
	pathGotBack, err := TildaExpansion(fs, pathIn)
	assert.Nil(t, err)
	assert.Empty(t, pathGotBack)
}

func TestTildaExpansionWhenFilenameNormalBlankReturnsBlank(t *testing.T) {
	pathIn := "normal"
	fs := NewMockFileSystem()
	pathGotBack, err := TildaExpansion(fs, pathIn)
	assert.Nil(t, err)
	assert.Equal(t, pathGotBack, "normal")
}
