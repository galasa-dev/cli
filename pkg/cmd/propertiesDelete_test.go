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

func TestPropertiesDeleteCommandInCommandCollectionHasName(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesDeleteCommand := commands.GetCommand(COMMAND_NAME_PROPERTIES_DELETE)
	assert.Equal(t, COMMAND_NAME_PROPERTIES_DELETE, propertiesDeleteCommand.Name())
	assert.Nil(t, propertiesDeleteCommand.Values())
	assert.NotNil(t, propertiesDeleteCommand.CobraCommand())
}
