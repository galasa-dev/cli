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

func TestPropertiesCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesCommand, err := commands.GetCommand(COMMAND_NAME_PROPERTIES)
	assert.NotNil(t, propertiesCommand)
	assert.Equal(t, COMMAND_NAME_PROPERTIES, propertiesCommand.Name())
	assert.NotNil(t, propertiesCommand.Values())
	assert.IsType(t, &PropertiesCmdValues{}, propertiesCommand.Values())
	assert.NotNil(t, propertiesCommand.CobraCommand())
	assert.Nil(t, err)
}
