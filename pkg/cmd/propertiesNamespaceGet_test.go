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

func TestPropertiesNamespaceGetCommandInCommandCollectionHasName(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesNamespaceGetCommand := commands.GetCommand(COMMAND_NAME_PROPERTIES_NAMESPACE_GET)
	assert.Equal(t, COMMAND_NAME_PROPERTIES_NAMESPACE_GET, propertiesNamespaceGetCommand.GetName())
	assert.NotNil(t, propertiesNamespaceGetCommand.GetValues())
	assert.IsType(t, &PropertiesNamespaceGetCmdValues{}, propertiesNamespaceGetCommand.GetValues())
	assert.NotNil(t, propertiesNamespaceGetCommand.GetCobraCommand())
}
