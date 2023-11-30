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

func TestPropertiesGetCommandInCommandCollectionHasName(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesGetCommand := commands.GetCommand(COMMAND_NAME_PROPERTIES_GET)
	assert.Equal(t, COMMAND_NAME_PROPERTIES_GET, propertiesGetCommand.GetName())
	assert.NotNil(t, propertiesGetCommand.GetValues())
	assert.IsType(t, &PropertiesGetCmdValues{}, propertiesGetCommand.GetValues())
	assert.NotNil(t, propertiesGetCommand.GetCobraCommand())
}
