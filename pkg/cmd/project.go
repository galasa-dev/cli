/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

type ProjectCommand struct {
	cobraCommand *cobra.Command
}

func NewProjectCmd(factory Factory, rootCmd GalasaCommand) (GalasaCommand, error) {
	cmd := new(ProjectCommand)
	err := cmd.init(factory, rootCmd)
	return cmd, err
}

func (cmd *ProjectCommand) init(factory Factory, rootCmd GalasaCommand) error {

	var err error = nil

	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Manipulate local project source code",
		Long:  "Creates and manipulates Galasa test project source code",
	}

	cmd.cobraCommand = projectCmd
	rootCmd.CobraCommand().AddCommand(projectCmd)

	return err
}

func (cmd *ProjectCommand) Name() string {
	return COMMAND_NAME_PROJECT
}

func (cmd *ProjectCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *ProjectCommand) Values() interface{} {
	return nil
}
