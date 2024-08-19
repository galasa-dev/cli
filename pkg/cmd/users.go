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

type UsersCmdValues struct {
	ecosystemBootstrap string
	name               string
}

type UsersCommand struct {
	values       *UsersCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewUsersCommand(rootCmd spi.GalasaCommand) (spi.GalasaCommand, error) {

	cmd := new(UsersCommand)
	err := cmd.init(rootCmd)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *UsersCommand) Name() string {
	return COMMAND_NAME_USERS
}

func (cmd *UsersCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *UsersCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *UsersCommand) init(rootCmd spi.GalasaCommand) error {

	var err error

	cmd.values = &UsersCmdValues{}
	cmd.cobraCommand = cmd.createCobraCommand(rootCmd)

	return err
}

func (cmd *UsersCommand) createCobraCommand(
	rootCommand spi.GalasaCommand,
) *cobra.Command {

	usersCobraCmd := &cobra.Command{
		Use:   "users",
		Short: "Manages users in an ecosystem",
		Long:  "Allows interaction with the user servlet to return information about users.",
	}

	addBootstrapFlag(usersCobraCmd, &cmd.values.ecosystemBootstrap)

	rootCommand.CobraCommand().AddCommand(usersCobraCmd)

	return usersCobraCmd
}

func addLoginIdFlag(cmd *cobra.Command, isMandatory bool, userCmdValues *UsersCmdValues) {

	flagName := "id"
	var description string
	if isMandatory {
		description = "A mandatory flag that is required to return the currently logged in user."
	} else {
		description = "An optional flag that is required to return the currently logged in user."
	}
	description += "The input must be a string"

	cmd.PersistentFlags().StringVarP(&userCmdValues.name, flagName, "i", "", description)

	if isMandatory {
		cmd.MarkPersistentFlagRequired(flagName)
	}
}
