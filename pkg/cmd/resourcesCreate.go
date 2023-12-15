/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"github.com/spf13/cobra"
)

// Variables set by cobra's command-line parsing.
type ResourcesCreateCmdValues struct {
}

type ResourcesCreateCommand struct {
	values       *ResourcesCreateCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewResourcesCreateCommand(factory Factory, resourcesCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {

	cmd := new(ResourcesCreateCommand)
	err := cmd.init(factory, resourcesCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *ResourcesCreateCommand) Name() string {
	return COMMAND_NAME_RESOURCES_CREATE
}

func (cmd *ResourcesCreateCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *ResourcesCreateCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *ResourcesCreateCommand) init(factory Factory, resourcesCommand GalasaCommand, rootCommand GalasaCommand) error {

	var err error = nil

	cmd.values = &ResourcesCreateCmdValues{}
	cmd.cobraCommand = cmd.createCobraCommand(factory, resourcesCommand, rootCommand.Values().(*RootCmdValues))

	return err
}

func (cmd *ResourcesCreateCommand) createCobraCommand(
	factory Factory,
	resourcesCommand GalasaCommand,
	rootCommandValues *RootCmdValues,
) *cobra.Command {

	resourcesCreateCommandValues := resourcesCommand.Values().(*ResourcesCmdValues)
	resourcesCreateCmd := &cobra.Command{
		Use:     "create",
		Short:   "Update Galasa Ecosystem resources.",
		Long:    "Create Galasa Ecosystem resources from definitions held in a file.",
		Args:    cobra.NoArgs,
		Aliases: []string{"resources create"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeResourcesCreate(factory,
				cmd, args, resourcesCreateCommandValues, rootCommandValues)
		},
	}

	resourcesCommand.CobraCommand().AddCommand(resourcesCreateCmd)

	return resourcesCreateCmd
}

func executeResourcesCreate(factory Factory,
	resourcesCreateCmd *cobra.Command,
	args []string,
	resourcesCmdValues *ResourcesCmdValues,
	rootCmdValues *RootCmdValues,
) error {
	action := "create"

	err := loadAndPassDataIntoResourcesApi(action, factory, resourcesCmdValues, rootCmdValues)

	return err
}
