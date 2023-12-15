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

func TestResourcesCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	resourcesCommand, err := commands.GetCommand(COMMAND_NAME_RESOURCES)
	assert.NotNil(t, resourcesCommand)
	assert.Equal(t, COMMAND_NAME_RESOURCES, resourcesCommand.Name())
	assert.NotNil(t, resourcesCommand.Values())
	assert.IsType(t, &ResourcesCmdValues{}, resourcesCommand.Values())
	assert.NotNil(t, resourcesCommand.CobraCommand())
	assert.Nil(t, err)
}
