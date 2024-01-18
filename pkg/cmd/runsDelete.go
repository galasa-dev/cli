/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import "github.com/spf13/cobra"

// Objective: Allow the user to do this:
//    runs delete --name U123
// And then galasactl deletes the run by abandoning it.

type RunsDeleteCommand struct {
	values       *RunsDeleteCmdValues
	cobraCommand *cobra.Command
}

type RunsDeleteCmdValues struct {
	runName string
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewRunsDeleteCommand(factory Factory, runsCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(RunsDeleteCommand)
	err := cmd.init(factory, runsCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsDeleteCommand) Name() string {
	return COMMAND_NAME_RUNS_DELETE
}

func (cmd *RunsDeleteCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsDeleteCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsDeleteCommand) init(factory Factory, runsCommand GalasaCommand, rootCommand GalasaCommand) error {
	var err error
	cmd.values = &RunsDeleteCmdValues{}
	cmd.cobraCommand, err = cmd.createRunsDeleteCobraCmd(
		factory,
		runsCommand,
		rootCommand.Values().(*RootCmdValues),
	)
	return err
}

func (cmd *RunsDeleteCommand) createRunsDeleteCobraCmd(factory Factory,
	runsCommand GalasaCommand,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {

	var err error = nil
	runsCmdValues := runsCommand.Values().(*RunsCmdValues)

	runsDeleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete an active run in the ecosystem",
		Long:    "Delete an active test run in the ecosystem if it is stuck or looping.",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs delete"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.executeDelete(factory, runsCmdValues, rootCmdValues)
		},
	}

	runsDeleteCmd.PersistentFlags().StringVar(&cmd.values.runName, "name", "", "the name of the test run to delete")

	runsDeleteCmd.MarkPersistentFlagRequired("name")

	runsCommand.CobraCommand().AddCommand(runsDeleteCmd)

	return runsDeleteCmd, err
}

func (cmd *RunsDeleteCommand) executeDelete(
	factory Factory,
	runsCmdValues *RunsCmdValues,
	rootCmdValues *RootCmdValues,
) error {

	var err error

	return err
}
