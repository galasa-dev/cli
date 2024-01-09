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
type ResourcesDeleteCmdValues struct {
}

type ResourcesDeleteCommand struct {
	values       *ResourcesDeleteCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewResourcesDeleteCommand(factory Factory, resourcesCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {

	cmd := new(ResourcesDeleteCommand)
	err := cmd.init(factory, resourcesCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *ResourcesDeleteCommand) Name() string {
	return COMMAND_NAME_RESOURCES_DELETE
}

func (cmd *ResourcesDeleteCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *ResourcesDeleteCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *ResourcesDeleteCommand) init(factory Factory, resourcesCommand GalasaCommand, rootCommand GalasaCommand) error {

	var err error = nil

	cmd.values = &ResourcesDeleteCmdValues{}
	cmd.cobraCommand = cmd.createCobraCommand(factory, resourcesCommand, rootCommand.Values().(*RootCmdValues))

	return err
}

func (cmd *ResourcesDeleteCommand) createCobraCommand(
	factory Factory,
	resourcesCommand GalasaCommand,
	rootCommandValues *RootCmdValues,
) *cobra.Command {

	resourcesDeleteCommandValues := resourcesCommand.Values().(*ResourcesCmdValues)
	resourcesDeleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete Galasa Ecosystem resources.",
		Long:    "Delete Galasa Ecosystem resources in a file.",
		Args:    cobra.NoArgs,
		Aliases: []string{"resources delete"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeResourcesDelete(factory,
				cmd, args, resourcesDeleteCommandValues, rootCommandValues)
		},
	}

	resourcesCommand.CobraCommand().AddCommand(resourcesDeleteCmd)

	return resourcesDeleteCmd
}

func executeResourcesDelete(factory Factory,
	resourcesDeleteCmd *cobra.Command,
	args []string,
	resourcesCmdValues *ResourcesCmdValues,
	rootCmdValues *RootCmdValues,
) error {
	action := "delete"

	err := loadAndPassDataIntoResourcesApi(action, factory, resourcesCmdValues, rootCmdValues)

	return err
}
