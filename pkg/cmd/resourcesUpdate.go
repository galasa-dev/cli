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
type ResourcesUpdateCmdValues struct {
}

type ResourcesUpdateCommand struct {
	values       *ResourcesUpdateCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewResourcesUpdateCommand(factory Factory, resourcesCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {

	cmd := new(ResourcesUpdateCommand)
	err := cmd.init(factory, resourcesCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *ResourcesUpdateCommand) Name() string {
	return COMMAND_NAME_RESOURCES_UPDATE
}

func (cmd *ResourcesUpdateCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *ResourcesUpdateCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *ResourcesUpdateCommand) init(factory Factory, resourcesCommand GalasaCommand, rootCommand GalasaCommand) error {

	var err error = nil

	cmd.values = &ResourcesUpdateCmdValues{}
	cmd.cobraCommand = cmd.createCobraCommand(factory, resourcesCommand, rootCommand.Values().(*RootCmdValues))

	return err
}

func (cmd *ResourcesUpdateCommand) createCobraCommand(
	factory Factory,
	resourcesCommand GalasaCommand,
	rootCommandValues *RootCmdValues,
) *cobra.Command {

	resourcesUpdateCommandValues := resourcesCommand.Values().(*ResourcesCmdValues)
	resourcesUpdateCmd := &cobra.Command{
		Use:     "update",
		Short:   "Update Galasa Ecosystem resources.",
		Long:    "Update Galasa Ecosystem resources from definitions held in a file.",
		Args:    cobra.NoArgs,
		Aliases: []string{"resources update"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeResourcesUpdate(factory,
				resourcesUpdateCommandValues, rootCommandValues)
		},
	}

	resourcesCommand.CobraCommand().AddCommand(resourcesUpdateCmd)

	return resourcesUpdateCmd
}

func executeResourcesUpdate(factory Factory,
	resourcesCmdValues *ResourcesCmdValues,
	rootCmdValues *RootCmdValues,
) error {
	action := "update"

	err := loadAndPassDataIntoResourcesApi(action, factory, resourcesCmdValues, rootCmdValues)

	return err
}
