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

func TestPropertiesSetCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesSetCommand := commands.GetCommand(COMMAND_NAME_PROPERTIES_SET)
	assert.Equal(t, COMMAND_NAME_PROPERTIES_SET, propertiesSetCommand.Name())
	assert.NotNil(t, propertiesSetCommand.Values())
	assert.IsType(t, &PropertiesSetCmdValues{}, propertiesSetCommand.Values())
	assert.NotNil(t, propertiesSetCommand.CobraCommand())
}
