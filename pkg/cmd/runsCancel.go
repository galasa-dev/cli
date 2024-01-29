/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import "github.com/spf13/cobra"

// Objective: Allow the user to do this:
//    runs cancel --name U123
// And then galasactl cancels the run by abandoning it.

type RunsCancelCommand struct {
	values       *RunsCancelCmdValues
	cobraCommand *cobra.Command
}

type RunsCancelCmdValues struct {
	runName string
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewRunsCancelCommand(factory Factory, runsCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(RunsCancelCommand)
	err := cmd.init(factory, runsCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsCancelCommand) Name() string {
	return COMMAND_NAME_RUNS_CANCEL
}

func (cmd *RunsCancelCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsCancelCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsCancelCommand) init(factory Factory, runsCommand GalasaCommand, rootCommand GalasaCommand) error {
	var err error
	cmd.values = &RunsCancelCmdValues{}
	cmd.cobraCommand, err = cmd.createRunsCancelCobraCmd(
		factory,
		runsCommand,
		rootCommand.Values().(*RootCmdValues),
	)
	return err
}

func (cmd *RunsCancelCommand) createRunsCancelCobraCmd(factory Factory,
	runsCommand GalasaCommand,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {

	var err error = nil
	runsCmdValues := runsCommand.Values().(*RunsCmdValues)

	runsCancelCmd := &cobra.Command{
		Use:     "cancel",
		Short:   "cancel an active run in the ecosystem",
		Long:    "Cancel an active test run in the ecosystem if it is stuck or looping.",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs cancel"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.executeCancel(factory, runsCmdValues, rootCmdValues)
		},
	}

	runsCancelCmd.PersistentFlags().StringVar(&cmd.values.runName, "name", "", "the name of the test run to cancel")

	runsCancelCmd.MarkPersistentFlagRequired("name")

	runsCommand.CobraCommand().AddCommand(runsCancelCmd)

	return runsCancelCmd, err
}

func (cmd *RunsCancelCommand) executeCancel(
	factory Factory,
	runsCmdValues *RunsCmdValues,
	rootCmdValues *RootCmdValues,
) error {

	var err error

	return err
}
