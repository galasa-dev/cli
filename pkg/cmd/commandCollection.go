/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/utils"
)

type CommandCollection interface {
	// name - One of the COMMAND_NAME_* constants.
	GetCommand(name string) (utils.GalasaCommand, error)

	GetRootCommand() utils.GalasaCommand

	Execute(args []string) error
}

type commandCollectionImpl struct {
	rootCommand utils.GalasaCommand

	commandMap map[string]utils.GalasaCommand
}

const (
	COMMAND_NAME_ROOT                     = "galasactl"
	COMMAND_NAME_AUTH                     = "auth"
	COMMAND_NAME_AUTH_LOGIN               = "auth login"
	COMMAND_NAME_AUTH_LOGOUT              = "auth logout"
	COMMAND_NAME_PROJECT                  = "project"
	COMMAND_NAME_PROJECT_CREATE           = "project create"
	COMMAND_NAME_LOCAL                    = "local"
	COMMAND_NAME_LOCAL_INIT               = "local init"
	COMMAND_NAME_PROPERTIES               = "properties"
	COMMAND_NAME_PROPERTIES_GET           = "properties get"
	COMMAND_NAME_PROPERTIES_SET           = "properties set"
	COMMAND_NAME_PROPERTIES_DELETE        = "properties delete"
	COMMAND_NAME_PROPERTIES_NAMESPACE     = "properties namespaces"
	COMMAND_NAME_PROPERTIES_NAMESPACE_GET = "properties namespaces get"
	COMMAND_NAME_RUNS                     = "runs"
	COMMAND_NAME_RUNS_DOWNLOAD            = "runs download"
	COMMAND_NAME_RUNS_GET                 = "runs get"
	COMMAND_NAME_RUNS_PREPARE             = "runs prepare"
	COMMAND_NAME_RUNS_SUBMIT              = "runs submit"
	COMMAND_NAME_RUNS_SUBMIT_LOCAL        = "runs submit local"
	COMMAND_NAME_RUNS_RESET               = "runs reset"
	COMMAND_NAME_RUNS_CANCEL              = "runs cancel"
	COMMAND_NAME_RESOURCES                = "resources"
	COMMAND_NAME_RESOURCES_APPLY          = "resources apply"
	COMMAND_NAME_RESOURCES_CREATE         = "resources create"
	COMMAND_NAME_RESOURCES_UPDATE         = "resources update"
	COMMAND_NAME_RESOURCES_DELETE         = "resources delete"
)

// -----------------------------------------------------------------
// Constructors
// -----------------------------------------------------------------
func NewCommandCollection(factory utils.Factory) (CommandCollection, error) {

	commands := new(commandCollectionImpl)

	err := commands.init(factory)

	return commands, err
}

// -----------------------------------------------------------------
// Public functions
// -----------------------------------------------------------------

func (commands *commandCollectionImpl) GetRootCommand() utils.GalasaCommand {
	cmd, _ := commands.GetCommand(COMMAND_NAME_ROOT)
	return cmd
}

func (commands *commandCollectionImpl) GetCommand(name string) (utils.GalasaCommand, error) {
	var err error
	cmd := commands.commandMap[name]
	if cmd == nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_COMMAND_NOT_FOUND_IN_CMD_COLLECTION)
		log.Printf("Caller tried to lookup %s in the command collection and it was not found.\n", name)
	}
	return cmd, err
}

func (commands *commandCollectionImpl) Execute(args []string) error {

	rootCmd := commands.GetRootCommand().CobraCommand()
	rootCmd.SetArgs(args)

	// Execute the command
	err := rootCmd.Execute()

	return err
}

// -----------------------------------------------------------------
// Private functions.
// -----------------------------------------------------------------
func (commands *commandCollectionImpl) init(factory utils.Factory) error {

	commands.commandMap = make(map[string]utils.GalasaCommand)

	rootCommand, err := NewRootCommand(factory)
	if err == nil {
		commands.rootCommand = rootCommand
		commands.commandMap[rootCommand.Name()] = rootCommand
	}

	if err == nil {
		err = commands.addAuthCommands(factory, rootCommand)
	}

	if err == nil {
		err = commands.addLocalCommands(factory, rootCommand)
	}

	if err == nil {
		err = commands.addProjectCommands(factory, rootCommand)
	}

	if err == nil {
		err = commands.addPropertiesCommands(factory, rootCommand)
	}

	if err == nil {
		err = commands.addRunsCommands(factory, rootCommand)
	}

	if err == nil {
		err = commands.addResourcesCommands(factory, rootCommand)
	}

	if err == nil {
		commands.setHelpFlags()
	}

	return err
}

func (commands *commandCollectionImpl) addAuthCommands(factory utils.Factory, rootCommand utils.GalasaCommand) error {
	var err error
	var authCommand utils.GalasaCommand
	var authLoginCommand utils.GalasaCommand
	var authLogoutCommand utils.GalasaCommand

	authCommand, err = NewAuthCommand(rootCommand)
	if err == nil {
		authLoginCommand, err = NewAuthLoginCommand(factory, authCommand, rootCommand)
		if err == nil {
			authLogoutCommand, err = NewAuthLogoutCommand(factory, authCommand, rootCommand)
		}
	}

	if err == nil {
		commands.commandMap[authCommand.Name()] = authCommand
		commands.commandMap[authLoginCommand.Name()] = authLoginCommand
		commands.commandMap[authLogoutCommand.Name()] = authLogoutCommand
	}

	return err
}

func (commands *commandCollectionImpl) addLocalCommands(factory utils.Factory, rootCommand utils.GalasaCommand) error {
	var err error
	var localCommand utils.GalasaCommand
	var localInitCommand utils.GalasaCommand

	localCommand, err = NewLocalCommand(rootCommand)
	if err == nil {
		localInitCommand, err = NewLocalInitCommand(factory, localCommand, rootCommand)
	}

	if err == nil {
		commands.commandMap[localCommand.Name()] = localCommand
		commands.commandMap[localInitCommand.Name()] = localInitCommand
	}
	return err
}

func (commands *commandCollectionImpl) addProjectCommands(factory utils.Factory, rootCommand utils.GalasaCommand) error {
	var err error

	var projectCommand utils.GalasaCommand

	projectCommand, err = NewProjectCmd(rootCommand)
	if err == nil {
		commands.commandMap[projectCommand.Name()] = projectCommand
	}

	if err == nil {
		var projectCreateCommand utils.GalasaCommand
		projectCreateCommand, err = NewProjectCreateCmd(factory, projectCommand, rootCommand)
		if err == nil {
			commands.commandMap[projectCreateCommand.Name()] = projectCreateCommand
		}
	}
	return err
}

func (commands *commandCollectionImpl) addPropertiesCommands(factory utils.Factory, rootCommand utils.GalasaCommand) error {
	var err error
	var propertiesCommand utils.GalasaCommand
	var propertiesGetCommand utils.GalasaCommand
	var propertiesDeleteCommand utils.GalasaCommand
	var propertiesSetCommand utils.GalasaCommand

	propertiesCommand, err = NewPropertiesCommand(rootCommand)
	if err == nil {
		propertiesGetCommand, err = NewPropertiesGetCommand(factory, propertiesCommand, rootCommand)
		if err == nil {
			propertiesSetCommand, err = NewPropertiesSetCommand(factory, propertiesCommand, rootCommand)
			if err == nil {
				propertiesDeleteCommand, err = NewPropertiesDeleteCommand(factory, propertiesCommand, rootCommand)
				if err == nil {
					err = commands.addPropertiesNamespaceCommands(factory, rootCommand, propertiesCommand)
				}
			}
		}
	}

	if err == nil {
		commands.commandMap[propertiesCommand.Name()] = propertiesCommand
		commands.commandMap[propertiesGetCommand.Name()] = propertiesGetCommand
		commands.commandMap[propertiesSetCommand.Name()] = propertiesSetCommand
		commands.commandMap[propertiesDeleteCommand.Name()] = propertiesDeleteCommand
	}

	return err
}

func (commands *commandCollectionImpl) addPropertiesNamespaceCommands(factory utils.Factory, rootCommand utils.GalasaCommand, propertiesCommand utils.GalasaCommand) error {
	var err error
	var propertiesNamespaceCommand utils.GalasaCommand
	var propertiesNamespaceGetCommand utils.GalasaCommand

	propertiesNamespaceCommand, err = NewPropertiesNamespaceCommand(propertiesCommand, rootCommand)
	if err == nil {
		propertiesNamespaceGetCommand, err = NewPropertiesNamespaceGetCommand(factory, propertiesNamespaceCommand, propertiesCommand, rootCommand)
	}

	if err == nil {
		commands.commandMap[propertiesNamespaceCommand.Name()] = propertiesNamespaceCommand
		commands.commandMap[propertiesNamespaceGetCommand.Name()] = propertiesNamespaceGetCommand
	}
	return err
}

func (commands *commandCollectionImpl) addRunsCommands(factory utils.Factory, rootCommand utils.GalasaCommand) error {

	var err error
	var runsCommand utils.GalasaCommand
	var runsDownloadCommand utils.GalasaCommand
	var runsGetCommand utils.GalasaCommand
	var runsPrepareCommand utils.GalasaCommand
	var runsSubmitCommand utils.GalasaCommand
	var runsSubmitLocalCommand utils.GalasaCommand
	var runsResetCommand utils.GalasaCommand
	var runsCancelCommand utils.GalasaCommand

	runsCommand, err = NewRunsCmd(rootCommand)
	if err == nil {
		runsDownloadCommand, err = NewRunsDownloadCommand(factory, runsCommand, rootCommand)
		if err == nil {
			runsGetCommand, err = NewRunsGetCommand(factory, runsCommand, rootCommand)
			if err == nil {
				runsPrepareCommand, err = NewRunsPrepareCommand(factory, runsCommand, rootCommand)
				if err == nil {
					runsSubmitCommand, err = NewRunsSubmitCommand(factory, runsCommand, rootCommand)
					if err == nil {
						runsSubmitLocalCommand, err = NewRunsSubmitLocalCommand(factory, runsSubmitCommand, runsCommand, rootCommand)
						if err == nil {
							runsResetCommand, err = NewRunsResetCommand(factory, runsCommand, rootCommand)
							if err == nil {
								runsCancelCommand, err = NewRunsCancelCommand(factory, runsCommand, rootCommand)
							}
						}
					}
				}
			}
		}
	}

	if err == nil {
		commands.commandMap[runsCommand.Name()] = runsCommand
		commands.commandMap[runsDownloadCommand.Name()] = runsDownloadCommand
		commands.commandMap[runsGetCommand.Name()] = runsGetCommand
		commands.commandMap[runsPrepareCommand.Name()] = runsPrepareCommand
		commands.commandMap[runsSubmitCommand.Name()] = runsSubmitCommand
		commands.commandMap[runsSubmitLocalCommand.Name()] = runsSubmitLocalCommand
		commands.commandMap[runsResetCommand.Name()] = runsResetCommand
		commands.commandMap[runsCancelCommand.Name()] = runsCancelCommand
	}

	return err
}

func (commands *commandCollectionImpl) addResourcesCommands(factory utils.Factory, rootCommand utils.GalasaCommand) error {

	var err error
	var resourcesCommand utils.GalasaCommand
	var resourcesApplyCommand utils.GalasaCommand
	var resourcesCreateCommand utils.GalasaCommand
	var resourcesUpdateCommand utils.GalasaCommand
	var resourcesDeleteCommand utils.GalasaCommand

	resourcesCommand, err = NewResourcesCmd(rootCommand)
	if err == nil {
		resourcesApplyCommand, err = NewResourcesApplyCommand(factory, resourcesCommand, rootCommand)
		if err == nil {
			resourcesCreateCommand, err = NewResourcesCreateCommand(factory, resourcesCommand, rootCommand)
			if err == nil {
				resourcesUpdateCommand, err = NewResourcesUpdateCommand(factory, resourcesCommand, rootCommand)
				if err == nil {
					resourcesDeleteCommand, err = NewResourcesDeleteCommand(factory, resourcesCommand, rootCommand)
				}
			}
		}
	}

	if err == nil {
		commands.commandMap[resourcesCommand.Name()] = resourcesCommand
		commands.commandMap[resourcesApplyCommand.Name()] = resourcesApplyCommand
		commands.commandMap[resourcesCreateCommand.Name()] = resourcesCreateCommand
		commands.commandMap[resourcesUpdateCommand.Name()] = resourcesUpdateCommand
		commands.commandMap[resourcesDeleteCommand.Name()] = resourcesDeleteCommand
	}

	return err
}

func (commands *commandCollectionImpl) setHelpFlags() {
	for _, command := range commands.commandMap {
		command.CobraCommand().Flags().BoolP("help", "h", false, "Displays the options for the '"+command.Name()+"' command.")
	}
}
