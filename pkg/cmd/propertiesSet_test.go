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

func TestPropertiesSetCommandInCommandCollectionHasName(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesSetCommand := commands.GetCommand(COMMAND_NAME_PROPERTIES_SET)
	assert.Equal(t, COMMAND_NAME_PROPERTIES_SET, propertiesSetCommand.GetName())
	assert.NotNil(t, propertiesSetCommand.GetValues())
	assert.IsType(t, &PropertiesSetCmdValues{}, propertiesSetCommand.GetValues())
	assert.NotNil(t, propertiesSetCommand.GetCobraCommand())
}
