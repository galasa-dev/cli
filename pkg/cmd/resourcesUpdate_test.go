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

func TestResourcesUpdateCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	resourcesUpdateCommand, err := commands.GetCommand(COMMAND_NAME_RESOURCES_UPDATE)
	assert.Nil(t, err)
	
	assert.NotNil(t, resourcesUpdateCommand)
	assert.Equal(t, COMMAND_NAME_RESOURCES_UPDATE, resourcesUpdateCommand.Name())
	assert.NotNil(t, resourcesUpdateCommand.Values())
	assert.IsType(t, &ResourcesUpdateCmdValues{}, resourcesUpdateCommand.Values())
	assert.NotNil(t, resourcesUpdateCommand.CobraCommand())
}
