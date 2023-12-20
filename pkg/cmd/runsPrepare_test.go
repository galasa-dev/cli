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

func TestRunsPrepareCommandInCommandCollection(t *testing.T) {

	factory := NewMockFactory()
	commands, _ := NewCommandCollection(factory)

	cmd, err := commands.GetCommand(COMMAND_NAME_RUNS_PREPARE)
	assert.Nil(t, err)
	
	assert.Equal(t, COMMAND_NAME_RUNS_PREPARE, cmd.Name())
	assert.NotNil(t, cmd.Values())
	assert.IsType(t, &RunsPrepareCmdValues{}, cmd.Values())
	assert.NotNil(t, cmd.CobraCommand())
}


func TestRunsPrepareHelpFlagSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	
	var args []string = []string{"runs", "prepare", "--help"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("Displays the options for the 'runs prepare' command.", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPrepareNoParametersReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	
	var args []string = []string{"runs", "prepare"}

	// When...
	err := Execute(factory, args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "Error: required flag(s) \"portfolio\" not set", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "required flag(s) \"portfolio\" not set")
}

func TestRunsPreparePortfolioFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "portfolio.file"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioPortfolioNoParamsFlagReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "Error: flag needs an argument: --portfolio", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "flag needs an argument: --portfolio")
}

func TestRunsPreparePortfolioAppendFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--append"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioBundleNoParamsFlagReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--bundle"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "Error: flag needs an argument: --bundle", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "flag needs an argument: --bundle")
}

func TestRunsPreparePortfolioBundlesFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--bundle", "bundle.name"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioClassNoParamsFlagReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--class"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "Error: flag needs an argument: --class", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "flag needs an argument: --class")
}

func TestRunsPreparePortfolioClassFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--class", "class.stuff"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioOverrideNoParamsFlagReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--override"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "Error: flag needs an argument: --override", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "flag needs an argument: --override")
}

func TestRunsPreparePortfolioOverrideFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--override", "override string one"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioOverrideFlagMultipleInstancesReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--override", "override string one", "--override", "override string two"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioPackageNoParamsFlagReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--package"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "Error: flag needs an argument: --package", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "flag needs an argument: --package")
}

func TestRunsPreparePortfolioPackageFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--package", "packagethingy"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioPackageFlagMultipleInstancesReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--package", "packageName1", "--package", "packageName2"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioRegexFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--regex"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioStreamNoParamsFlagReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--stream"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "Error: flag needs an argument: --stream", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "flag needs an argument: --stream")
}

func TestRunsPreparePortfolioStreamFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--stream", "arlo.stream"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioTagNoParamsFlagReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--tag"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "Error: flag needs an argument: --tag", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "flag needs an argument: --tag")
}

func TestRunsPreparePortfolioTagFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--tag", "tag.stuff"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioTagFlagCommaSeperatedListValuesSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--tag", "comma,seperated,tags"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioTestNoParamsFlagReturnsError(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--test"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "Error: flag needs an argument: --test", "", factory, t)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "flag needs an argument: --test")
}

func TestRunsPreparePortfolioTestFlagReturnsOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--test", "tag.stuff"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}

func TestRunsPreparePortfolioTestFlagCommaSeperatedListValuesSetCorrectly(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--test", "comma,seperated,tests"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}
func TestRunsPrepareAllFlagsReturnOk(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	commandCollection := setupTestCommandCollection(COMMAND_NAME_RUNS_PREPARE, factory, t)
	
	var args []string = []string{"runs", "prepare", "--portfolio", "roo.yaml", "--test", "test,stuff", "--tag", "tag,list", 
	"--append", "--bundle", "comma,seperated,bundles", "--class", "class,list", "--override", "list,o,overrides", "--package", "package,list",
"--regex", "--stream", "stream,list"}

	// When...
	err := commandCollection.Execute(args)

	// Then...

	// Check what the user saw is reasonable.
	checkOutput("", "", "", factory, t)

	assert.Nil(t, err)
}