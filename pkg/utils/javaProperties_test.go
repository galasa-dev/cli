/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanCreateAPropsFileAndReadItBack(t *testing.T) {
	var props map[string]interface{} = make(map[string]interface{})
	props["a"] = "b"
	props["d"] = "e"

	fs := NewMockFileSystem()

	WritePropertiesFile(fs, "myPropsFile.properties", props)

	propsGotBack, err := ReadPropertiesFile(fs, "myPropsFile.properties")

	assert.Nil(t, err)
	assert.Contains(t, propsGotBack, "a")
	assert.Contains(t, propsGotBack, "d")
	assert.Equal(t, propsGotBack["a"], "b")
	assert.Equal(t, propsGotBack["d"], "e")
}

