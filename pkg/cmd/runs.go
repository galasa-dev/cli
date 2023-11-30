/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

type RunsCmdValues struct {
	bootstrap string
}

type RunsCommand struct {
	cobraCommand *cobra.Command
	values       *RunsCmdValues
}

func NewRunsCmd(factory Factory, rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(RunsCommand)
	err := cmd.init(factory, rootCommand)
	return cmd, err
}

func (cmd *RunsCommand) GetName() string {
	return COMMAND_NAME_RUNS
}

func (cmd *RunsCommand) GetCobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsCommand) GetValues() interface{} {
	return cmd.values
}

func (cmd *RunsCommand) init(factory Factory, rootCommand GalasaCommand) error {

	var err error

	runsCmdValues := &RunsCmdValues{}

	runsCobraCmd := &cobra.Command{
		Use:   "runs",
		Short: "Manage test runs in the ecosystem",
		Long:  "Assembles, submits and monitors test runs in Galasa Ecosystem",
	}

	addBootstrapFlag(runsCobraCmd, &runsCmdValues.bootstrap)

	rootCommand.GetCobraCommand().AddCommand(runsCobraCmd)

	cmd.cobraCommand = runsCobraCmd
	cmd.values = runsCmdValues

	return err
}
