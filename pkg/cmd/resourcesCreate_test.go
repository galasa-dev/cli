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

func TestResourcesCreateCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	resourcesCreateCommand := commands.GetCommand(COMMAND_NAME_RESOURCES_CREATE)
	assert.NotNil(t, resourcesCreateCommand)
	assert.Equal(t, COMMAND_NAME_RESOURCES_CREATE, resourcesCreateCommand.Name())
	assert.NotNil(t, resourcesCreateCommand.Values())
	assert.IsType(t, &ResourcesCreateCmdValues{}, resourcesCreateCommand.Values())
	assert.NotNil(t, resourcesCreateCommand.CobraCommand())
}
