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
	bootstrap string
}

type RunsCommand struct {
	cobraCommand *cobra.Command
	values       *RunsCmdValues
}

// ------------------------------------------------------------------------------------------------
// Constructors
// ------------------------------------------------------------------------------------------------

func NewRunsCmd(rootCommand spi.GalasaCommand) (spi.GalasaCommand, error) {
	cmd := new(RunsCommand)
	err := cmd.init(rootCommand)
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

func (cmd *RunsCommand) init(rootCmd spi.GalasaCommand) error {

	var err error

	cmd.values = &RunsCmdValues{}
	cmd.cobraCommand, err = cmd.createCobraCommand(rootCmd)

	return err
}

func (cmd *RunsCommand) createCobraCommand(rootCommand spi.GalasaCommand) (*cobra.Command, error) {

	var err error

	runsCobraCmd := &cobra.Command{
		Use:   "runs",
		Short: "Manage test runs in the ecosystem",
		Long:  "Assembles, submits and monitors test runs in Galasa Ecosystem",
	}

	addBootstrapFlag(runsCobraCmd, &cmd.values.bootstrap)

	rootCommand.CobraCommand().AddCommand(runsCobraCmd)

	return runsCobraCmd, err
}
