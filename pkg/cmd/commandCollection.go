/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

type CommandCollection interface {
	// name - One of the COMMAND_NAME_* constants.
	GetCommand(name string) GalasaCommand

	GetRootCommand() GalasaCommand

	Execute(args []string) error
}

type CommandCollectionImpl struct {
	rootCommand GalasaCommand
}

const (
	COMMAND_NAME_ROOT = "galasactl"
)

// -----------------------------------------------------------------
// Public functions.
// -----------------------------------------------------------------
func NewCommandCollection(factory Factory) (CommandCollection, error) {

	commands := new(CommandCollectionImpl)

	err := commands.init(factory)

	return commands, err
}

func (commands *CommandCollectionImpl) GetRootCommand() GalasaCommand {
	return commands.GetCommand(COMMAND_NAME_ROOT)
}

func (commands *CommandCollectionImpl) GetCommand(name string) GalasaCommand {
	return commands.rootCommand
}

func (commands *CommandCollectionImpl) Execute(args []string) error {

	rootCmd := commands.GetRootCommand().GetCobraCommand()
	rootCmd.SetArgs(args)

	// Execute the command
	err := rootCmd.Execute()

	return err
}

// -----------------------------------------------------------------
// Private functions.
// -----------------------------------------------------------------
func (commands *CommandCollectionImpl) init(factory Factory) error {

	rootCommand, err := NewRootCommand(factory)

	if err == nil {
		commands.rootCommand = rootCommand
	}
	return err
}
