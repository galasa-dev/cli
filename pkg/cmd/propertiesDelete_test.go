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

func TestPropertiesDeleteNoArgsReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"properties", "delete"}

	// When...
	err := Execute(factory, args)

	// Then...
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Usage:")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Contains(t, errText, "Error: required flag(s) \"name\", \"namespace\" not set")

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.NotNil(t, err)
}

func TestPropertiesDeleteWithoutName(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"properties", "delete", "--namespace", "jitters"}

	// When...
	err := Execute(factory, args)

	// Then...
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Usage:")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Contains(t, errText, "Error: required flag(s) \"name\" not set")

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.NotNil(t, err)
}

func TestPropertiesDeleteWithoutNamespace(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"properties", "delete", "--name", "jeepers"}

	// When...
	err := Execute(factory, args)

	// Then...
	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Contains(t, outText, "Usage:")

	stdErrConsole := factory.GetStdErrConsole().(*utils.MockConsole)
	errText := stdErrConsole.ReadText()
	assert.Contains(t, errText, "Error: required flag(s) \"namespace\" not set")

	// We expect an exit code of 1 for this command. But it seems that syntax errors caught by cobra still return no error.
	finalWordHandler := factory.GetFinalWordHandler().(*MockFinalWordHandler)
	o := finalWordHandler.ReportedObject
	assert.Nil(t, o)

	assert.NotNil(t, err)
}

func TestPropertiesDeleteWithNameAndNamespace(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	fs := factory.GetFileSystem()
	var args []string = []string{"properties", "delete", "--namespace", "gyro", "--name", "space.ball"}
	homeDir, _ := fs.GetUserHomeDirPath()
	galasaDir := homeDir + "/.galasa/"
	fs.WriteTextFile(galasaDir+"bootstrap.properties", "")

	// When...
	err := Execute(factory, args)

	// Then...
	assert.Nil(t, err)

	stdOutConsole := factory.GetStdOutConsole().(*utils.MockConsole)
	outText := stdOutConsole.ReadText()
	assert.Equal(t, outText, "")
}