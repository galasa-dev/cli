/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import "github.com/spf13/cobra"

// Objective: Allow the user to do this:
//    runs reset --name U123
// And then galasactl resets the run by requeuing it.

type RunsResetCommand struct {
	values       *RunsResetCmdValues
	cobraCommand *cobra.Command
}

type RunsResetCmdValues struct {
	runName string
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewRunsResetCommand(factory Factory, runsCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(RunsResetCommand)
	err := cmd.init(factory, runsCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsResetCommand) Name() string {
	return COMMAND_NAME_RUNS_RESET
}

func (cmd *RunsResetCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsResetCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsResetCommand) init(factory Factory, runsCommand GalasaCommand, rootCommand GalasaCommand) error {
	var err error
	cmd.values = &RunsResetCmdValues{}
	cmd.cobraCommand, err = cmd.createRunsResetCobraCmd(
		factory,
		runsCommand,
		rootCommand.Values().(*RootCmdValues),
	)
	return err
}

func (cmd *RunsResetCommand) createRunsResetCobraCmd(factory Factory,
	runsCommand GalasaCommand,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {

	var err error = nil
	runsCmdValues := runsCommand.Values().(*RunsCmdValues)

	runsResetCmd := &cobra.Command{
		Use:     "reset",
		Short:   "reset an active run in the ecosystem",
		Long:    "Reset an active test run in the ecosystem if it is stuck or looping.",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs reset"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.executeReset(factory, runsCmdValues, rootCmdValues)
		},
	}

	runsResetCmd.PersistentFlags().StringVar(&cmd.values.runName, "name", "", "the name of the test run to reset")

	runsResetCmd.MarkPersistentFlagRequired("name")

	runsCommand.CobraCommand().AddCommand(runsResetCmd)

	return runsResetCmd, err
}

func (cmd *RunsResetCommand) executeReset(
	factory Factory,
	runsCmdValues *RunsCmdValues,
	rootCmdValues *RootCmdValues,
) error {

	var err error

	return err
}
