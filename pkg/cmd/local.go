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

// ------------------------------------------------------------------------------------------------
// Constructors
// ------------------------------------------------------------------------------------------------
func NewLocalCommand(rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(LocalCommand)
	err := cmd.init(rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public functions
// ------------------------------------------------------------------------------------------------
func (cmd *LocalCommand) Name() string {
	return COMMAND_NAME_LOCAL
}

func (cmd *LocalCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *LocalCommand) Values() interface{} {
	return nil
}

// ------------------------------------------------------------------------------------------------
// Private functions
// ------------------------------------------------------------------------------------------------
func (cmd *LocalCommand) init(rootCommand GalasaCommand) error {
	var err error
	cmd.cobraCommand, err = cmd.createCobraCommand(rootCommand)
	return err
}

func (cmd *LocalCommand) createCobraCommand(rootCommand GalasaCommand) (*cobra.Command, error) {
	var err error
	localCobraCmd := &cobra.Command{
		Use:   "local",
		Short: "Manipulate local system",
		Long:  "Manipulate local system",
	}
	rootCommand.CobraCommand().AddCommand(localCobraCmd)
	return localCobraCmd, err
}
