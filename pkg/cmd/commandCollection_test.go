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

func TestCommandCollectionGetCommandInvalidNameReturnsError(t *testing.T){
	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)
	command, err := commands.GetCommand("bogus command name")
	assert.Nil(t, command)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "GAL1120E: 'bogus command name' not found in command collection")
}

func TestCommandCollectionGetCommandValidCmdNameReturnsOk(t *testing.T){
	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)
	command, err := commands.GetCommand("galasactl")
	assert.NotNil(t, command)
	assert.Nil(t, err)
}