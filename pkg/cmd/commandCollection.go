/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

type CommandCollection interface {
	// name - One of the COMMAND_NAME_* constants.
	GetCommand(name string) GalasaCommand

	GetRootCommand() GalasaCommand

	Execute(args []string) error
}

type CommandCollectionImpl struct {
	rootCommand GalasaCommand

	commandMap map[string]GalasaCommand
}

const (
	COMMAND_NAME_ROOT           = "galasactl"
	COMMAND_NAME_AUTH           = "auth"
	COMMAND_NAME_PROJECT        = "project"
	COMMAND_NAME_PROJECT_CREATE = "project create"
	COMMAND_NAME_LOCAL          = "local"
	COMMAND_NAME_LOCAL_INIT     = "local init"
	COMMAND_NAME_RUNS           = "runs"
)

// -----------------------------------------------------------------
// Public functions.
// -----------------------------------------------------------------
func NewCommandCollection(factory Factory) (CommandCollection, error) {

	commands := new(CommandCollectionImpl)

	err := commands.init(factory)

	return commands, err
}

// The main entry point into the cmd package.
func Execute(factory Factory, args []string) error {
	var err error

	finalWordHandler := factory.GetFinalWordHandler()

	var commands CommandCollection
	commands, err = NewCommandCollection(factory)

	if err == nil {

		// Catch execution if a panic happens.
		defer func() {
			err := recover()

			// Display the error and exit.
			finalWordHandler.FinalWord(commands.GetRootCommand(), err)
		}()

		// Execute the command
		err = commands.Execute(args)
	}
	finalWordHandler.FinalWord(commands.GetRootCommand(), err)
	return err
}

func (commands *CommandCollectionImpl) GetRootCommand() GalasaCommand {
	return commands.GetCommand(COMMAND_NAME_ROOT)
}

func (commands *CommandCollectionImpl) GetCommand(name string) GalasaCommand {
	cmd, _ := commands.commandMap[name]
	return cmd
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

	commands.commandMap = make(map[string]GalasaCommand)

	rootCommand, err := NewRootCommand(factory)
	if err == nil {
		commands.rootCommand = rootCommand
		commands.commandMap[rootCommand.GetName()] = rootCommand
	}

	var authCommand GalasaCommand
	if err == nil {
		authCommand, err = NewAuthCommand(factory, rootCommand)
		if err == nil {
			commands.commandMap[authCommand.GetName()] = authCommand
		}
	}

	var localCommand GalasaCommand
	if err == nil {
		localCommand, err = NewLocalCommand(factory, rootCommand)
		if err == nil {
			commands.commandMap[localCommand.GetName()] = localCommand
		}
	}

	var localInitCommand GalasaCommand
	if err == nil {
		localInitCommand, err = NewLocalInitCommand(factory, localCommand, rootCommand)
		if err == nil {
			commands.commandMap[localInitCommand.GetName()] = localInitCommand
		}
	}

	var projectCommand GalasaCommand
	if err == nil {
		projectCommand, err = NewProjectCmd(factory, rootCommand)
		if err == nil {
			commands.commandMap[projectCommand.GetName()] = projectCommand
		}
	}

	var projectCreateCommand GalasaCommand
	if err == nil {
		projectCreateCommand, err = NewProjectCreateCmd(factory, rootCommand, projectCommand)
		if err == nil {
			commands.commandMap[projectCreateCommand.GetName()] = projectCreateCommand
		}
	}

	var runsCommand GalasaCommand
	if err == nil {
		runsCommand, err = NewRunsCmd(factory, rootCommand)
		if err == nil {
			commands.commandMap[runsCommand.GetName()] = runsCommand
		}
	}

	// var runsDownloadCommand GalasaCommand
	// if err == nil {
	// 	runsDownloadCommand, err = NewRunsDownloadCmd(factory, runsCommand, rootCommand)
	// 	if err == nil {
	// 		commands.commandMap[runsDownloadCommand.GetName()] = runsDownloadCommand
	// 	}
	// }

	// if err == nil {
	// 	_, err = createRunsGetCmd(factory, runsCmd, runsCmdValues, rootCmdValues)
	// }

	// if err == nil {
	// 	_, err = createRunsPrepareCmd(factory, runsCmd, runsCmdValues, rootCmdValues)
	// }
	// if err == nil {
	// 	_, err = createRunsSubmitCmd(factory, runsCmd, runsCmdValues, rootCmdValues)
	// }

	// 	if err == nil {
	// 		_, err = createRunsCmd(factory, rootCmd, rootCmdValues)
	// 	}

	// 	if err == nil {
	// 		_, err = createPropertiesCmd(factory, rootCmd, rootCmdValues)
	// 	}

	// 	if err == nil {
	// 		_, err = createAuthCmd(factory, rootCmd, rootCmdValues)
	// 	}
	// 	return err
	// }

	if err == nil {
		sanitiseCommandHelpDescriptions(rootCommand.GetCobraCommand())
	}

	return err
}

// TODO: Make this an object method.
func sanitiseCommandHelpDescriptions(rootCmd *cobra.Command) {
	setHelpFlagForAllCommands(rootCmd, func(cobra *cobra.Command) {
		alias := cobra.NameAndAliases()
		//if the command has an alias,
		//the format would be cobra.Name, cobra.Aliases
		//otherwise it is just cobra.Name
		nameAndAliases := strings.Split(alias, ", ")
		if len(nameAndAliases) > 1 {
			alias = nameAndAliases[1]
		}

		cobra.Flags().BoolP("help", "h", false, "Displays the options for the "+alias+" command.")
	})
}

// TODO: Make this an object method.
func setHelpFlagForAllCommands(command *cobra.Command, setHelpFlag func(*cobra.Command)) {
	setHelpFlag(command)

	//for all the commands eg properties get, set etc
	for _, cobraCommand := range command.Commands() {
		setHelpFlagForAllCommands(cobraCommand, setHelpFlag)
	}
}
