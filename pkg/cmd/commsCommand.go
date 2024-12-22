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

type CommsCmdValues struct {
	bootstrap string
    maxRetries int
    retryBackoffSeconds float64
	*RootCmdValues
}

type CommsCommand struct {
	cobraCommand *cobra.Command
	values       *CommsCmdValues
}

// ------------------------------------------------------------------------------------------------
// Constructors
// ------------------------------------------------------------------------------------------------

func NewCommsCommand(rootCommand spi.GalasaCommand) (spi.GalasaCommand, error) {
	cmd := new(CommsCommand)
	err := cmd.init(rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public functions
// ------------------------------------------------------------------------------------------------

func (cmd *CommsCommand) Name() string {
    // This is a hidden command, so no assume no name
	return ""
}

func (cmd *CommsCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *CommsCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private functions
// ------------------------------------------------------------------------------------------------

func (cmd *CommsCommand) init(rootCmd spi.GalasaCommand) error {

	var err error

	cmd.values = &CommsCmdValues{
		RootCmdValues: rootCmd.Values().(*RootCmdValues),
	}
	cmd.cobraCommand, err = cmd.createCobraCommand(rootCmd)

	return err
}

func (cmd *CommsCommand) createCobraCommand(rootCommand spi.GalasaCommand) (*cobra.Command, error) {

	var err error

	commsCobraCmd := &cobra.Command{
		Hidden: true,
		SilenceUsage: true,
	}

	addBootstrapFlag(commsCobraCmd, &cmd.values.bootstrap)
    addRateLimitRetryFlags(commsCobraCmd, &cmd.values.maxRetries, &cmd.values.retryBackoffSeconds)

	rootCommand.CobraCommand().AddCommand(commsCobraCmd)

	return commsCobraCmd, err
}