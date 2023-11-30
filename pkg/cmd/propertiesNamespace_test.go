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

func TestPropertiesNamespaceCommandInCommandCollectionHasName(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesNamespaceCommand := commands.GetCommand(COMMAND_NAME_PROPERTIES_NAMESPACE)
	assert.Equal(t, COMMAND_NAME_PROPERTIES_NAMESPACE, propertiesNamespaceCommand.GetName())
	assert.Nil(t, propertiesNamespaceCommand.GetValues())
	assert.NotNil(t, propertiesNamespaceCommand.GetCobraCommand())
}
