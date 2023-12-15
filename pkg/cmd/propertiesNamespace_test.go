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

func TestPropertiesNamespaceCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesNamespaceCommand, err := commands.GetCommand(COMMAND_NAME_PROPERTIES_NAMESPACE)
	assert.Nil(t, err)
	
	assert.Equal(t, COMMAND_NAME_PROPERTIES_NAMESPACE, propertiesNamespaceCommand.Name())
	assert.Nil(t, propertiesNamespaceCommand.Values())
	assert.NotNil(t, propertiesNamespaceCommand.CobraCommand())
}
