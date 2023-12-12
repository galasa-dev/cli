/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"testing"

	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestPropertiesNoCommandsProducesUsageReport(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"properties"}

	// When...
	Execute(factory, args)

	// Then...
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Usage:")
	assert.Contains(t, outText, "galasactl properties [command]")

	// We expect an exit code of 0 for this command.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)
}

func TestPropertiesCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	propertiesCommand := commands.GetCommand(COMMAND_NAME_PROPERTIES)
	assert.NotNil(t, propertiesCommand)
	assert.Equal(t, COMMAND_NAME_PROPERTIES, propertiesCommand.Name())
	assert.NotNil(t, propertiesCommand.Values())
	assert.IsType(t, &PropertiesCmdValues{}, propertiesCommand.Values())
	assert.NotNil(t, propertiesCommand.CobraCommand())
}
