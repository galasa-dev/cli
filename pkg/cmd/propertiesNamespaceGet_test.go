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

func TestPropertiesNamespaceGetCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesNamespaceGetCommand, err := commands.GetCommand(COMMAND_NAME_PROPERTIES_NAMESPACE_GET)
	assert.Equal(t, COMMAND_NAME_PROPERTIES_NAMESPACE_GET, propertiesNamespaceGetCommand.Name())
	assert.NotNil(t, propertiesNamespaceGetCommand.Values())
	assert.IsType(t, &PropertiesNamespaceGetCmdValues{}, propertiesNamespaceGetCommand.Values())
	assert.NotNil(t, propertiesNamespaceGetCommand.CobraCommand())
	assert.Nil(t, err)
}
