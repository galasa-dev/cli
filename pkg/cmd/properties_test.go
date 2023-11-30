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

func TestPropertiesCommandInCommandCollectionHasName(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesCommand := commands.GetCommand(COMMAND_NAME_PROPERTIES)
	assert.NotNil(t, propertiesCommand)
	assert.Equal(t, COMMAND_NAME_PROPERTIES, propertiesCommand.GetName())
	assert.NotNil(t, propertiesCommand.GetValues())
	assert.IsType(t, &PropertiesCmdValues{}, propertiesCommand.GetValues())
	assert.NotNil(t, propertiesCommand.GetCobraCommand())
}
