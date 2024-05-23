/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"github.com/galasa-dev/cli/pkg/utils"
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
func NewResourcesDeleteCommand(factory utils.Factory, resourcesCommand utils.GalasaCommand, rootCommand utils.GalasaCommand) (utils.GalasaCommand, error) {

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

func (cmd *ResourcesDeleteCommand) init(factory utils.Factory, resourcesCommand utils.GalasaCommand, rootCommand utils.GalasaCommand) error {

	var err error

	cmd.values = &ResourcesDeleteCmdValues{}
	cmd.cobraCommand = cmd.createCobraCommand(factory, resourcesCommand, rootCommand.Values().(*RootCmdValues))

	return err
}

func (cmd *ResourcesDeleteCommand) createCobraCommand(
	factory utils.Factory,
	resourcesCommand utils.GalasaCommand,
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
				resourcesDeleteCommandValues, rootCommandValues)
		},
	}

	resourcesCommand.CobraCommand().AddCommand(resourcesDeleteCmd)

	return resourcesDeleteCmd
}

func executeResourcesDelete(factory utils.Factory,
	resourcesCmdValues *ResourcesCmdValues,
	rootCmdValues *RootCmdValues,
) error {
	action := "delete"

	err := loadAndPassDataIntoResourcesApi(action, factory, resourcesCmdValues, rootCmdValues)

	return err
}
