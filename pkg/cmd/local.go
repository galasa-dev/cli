/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

type LocalCommand struct {
	cobraCommand *cobra.Command
}

func NewLocalCommand(factory Factory, rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(LocalCommand)
	err := cmd.init(factory, rootCommand)
	return cmd, err
}

func (cmd *LocalCommand) init(factory Factory, rootCommand GalasaCommand) error {

	var err error

	localCobraCmd := &cobra.Command{
		Use:   "local",
		Short: "Manipulate local system",
		Long:  "Manipulate local system",
	}
	rootCommand.GetCobraCommand().AddCommand(localCobraCmd)

	cmd.cobraCommand = localCobraCmd

	return err
}

func (cmd *LocalCommand) GetName() string {
	return COMMAND_NAME_LOCAL
}

func (cmd *LocalCommand) GetCobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *LocalCommand) GetValues() interface{} {
	return nil
}
