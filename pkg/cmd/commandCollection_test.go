/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommandCollectionReturnsNonNil(t *testing.T) {
	factory := NewMockFactory()
	commands, err := NewCommandCollection(factory)
	assert.Nil(t, err)
	assert.NotNil(t, commands)
}
