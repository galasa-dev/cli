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

type RunsCmdValues struct {
}

type RunsCommand struct {
	cobraCommand *cobra.Command
	values       *RunsCmdValues
}

// ------------------------------------------------------------------------------------------------
// Constructors
// ------------------------------------------------------------------------------------------------

func NewRunsCmd(rootCommand spi.GalasaCommand, commsCommand spi.GalasaCommand) (spi.GalasaCommand, error) {
	cmd := new(RunsCommand)
	err := cmd.init(rootCommand, commsCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public functions
// ------------------------------------------------------------------------------------------------

func (cmd *RunsCommand) Name() string {
	return COMMAND_NAME_RUNS
}

func (cmd *RunsCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private functions
// ------------------------------------------------------------------------------------------------

func (cmd *RunsCommand) init(rootCmd spi.GalasaCommand, commsCommand spi.GalasaCommand) error {

	var err error

	cmd.values = &RunsCmdValues{}
	cmd.cobraCommand, err = cmd.createCobraCommand(rootCmd, commsCommand)

	return err
}

func (cmd *RunsCommand) createCobraCommand(rootCommand spi.GalasaCommand, commsCommand spi.GalasaCommand) (*cobra.Command, error) {

	var err error

	runsCobraCmd := &cobra.Command{
		Use:   "runs",
		Short: "Manage test runs in the ecosystem",
		Long:  "Assembles, submits and monitors test runs in Galasa Ecosystem",
	}

	runsCobraCmd.PersistentFlags().AddFlagSet(commsCommand.CobraCommand().PersistentFlags())
	rootCommand.CobraCommand().AddCommand(runsCobraCmd)

	return runsCobraCmd, err
}
