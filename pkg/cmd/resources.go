/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package cmd

import (
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/spf13/cobra"
)

type ResourcesCmdValues struct {
	bootstrap string
	filePath  string
}

type ResourcesCommand struct {
	cobraCommand *cobra.Command
	values       *ResourcesCmdValues
}

// ------------------------------------------------------------------------------------------------
// Constructors
// ------------------------------------------------------------------------------------------------

func NewResourcesCmd(rootCommand spi.GalasaCommand) (spi.GalasaCommand, error) {
	cmd := new(ResourcesCommand)
	err := cmd.init(rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public functions
// ------------------------------------------------------------------------------------------------

func (cmd *ResourcesCommand) Name() string {
	return COMMAND_NAME_RESOURCES
}

func (cmd *ResourcesCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *ResourcesCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private functions
// ------------------------------------------------------------------------------------------------

func (cmd *ResourcesCommand) init(rootCmd spi.GalasaCommand) error {

	var err error

	cmd.values = &ResourcesCmdValues{}
	cmd.cobraCommand, err = cmd.createCobraCommand(rootCmd)

	return err
}

func (cmd *ResourcesCommand) createCobraCommand(rootCommand spi.GalasaCommand) (*cobra.Command, error) {

	var err error

	resourcesCobraCmd := &cobra.Command{
		Use:   "resources",
		Short: "Manages resources in an ecosystem",
		Long:  "Allows interaction with the Resources endpoint to create and maintain resources in the Galasa Ecosystem",
	}

	addBootstrapFlag(resourcesCobraCmd, &cmd.values.bootstrap)
	resourcesCobraCmd.PersistentFlags().StringVarP(&cmd.values.filePath, "file", "f", "",
		"The file containing yaml definitions of resources to be applied manipulated by this command. "+
			"This can be a fully-qualified path or path relative to the current directory."+
			"Example: my_resources.yaml")
	resourcesCobraCmd.MarkPersistentFlagRequired("file")

	rootCommand.CobraCommand().AddCommand(resourcesCobraCmd)

	return resourcesCobraCmd, err
}
