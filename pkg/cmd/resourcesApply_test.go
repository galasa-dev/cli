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

func TestResourcesApplyCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	resourcesApplyCommand, err := commands.GetCommand(COMMAND_NAME_RESOURCES_APPLY)
	assert.NotNil(t, resourcesApplyCommand)
	assert.Equal(t, COMMAND_NAME_RESOURCES_APPLY, resourcesApplyCommand.Name())
	assert.NotNil(t, resourcesApplyCommand.Values())
	assert.IsType(t, &ResourcesApplyCmdValues{}, resourcesApplyCommand.Values())
	assert.NotNil(t, resourcesApplyCommand.CobraCommand())
	assert.Nil(t, err)
}
