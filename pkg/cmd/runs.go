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

	runsCobraCmd.PersistentFlags().StringVarP(&runsCmdValues.bootstrap, "bootstrap", "b", "",
		"Bootstrap URL. Should start with 'http://' or 'file://'. "+
			"If it starts with neither, it is assumed to be a fully-qualified path. "+
			"If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. "+
			"Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties")

	rootCommand.GetCobraCommand().AddCommand(runsCobraCmd)

	cmd.cobraCommand = runsCobraCmd
	cmd.values = runsCmdValues

	return err
}
