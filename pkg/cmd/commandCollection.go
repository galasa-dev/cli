/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/spi"
)

type CommandCollection interface {
	// name - One of the COMMAND_NAME_* constants.
	GetCommand(name string) (spi.GalasaCommand, error)

	GetRootCommand() spi.GalasaCommand

	Execute(args []string) error
}

type commandCollectionImpl struct {
	rootCommand spi.GalasaCommand

	commandMap map[string]spi.GalasaCommand
}

const (
	COMMAND_NAME_ROOT                     = "galasactl"
	COMMAND_NAME_AUTH                     = "auth" //This is a command, not a secret //pragma: allowlist secret
	COMMAND_NAME_AUTH_LOGIN               = "auth login"
	COMMAND_NAME_AUTH_LOGOUT              = "auth logout"
	COMMAND_NAME_AUTH_TOKENS              = "auth tokens"
	COMMAND_NAME_AUTH_TOKENS_GET          = "auth tokens get"
	COMMAND_NAME_AUTH_TOKENS_DELETE       = "auth tokens delete"
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
	COMMAND_NAME_RUNS_DELETE              = "runs delete"
	COMMAND_NAME_RESOURCES                = "resources"
	COMMAND_NAME_RESOURCES_APPLY          = "resources apply"
	COMMAND_NAME_RESOURCES_CREATE         = "resources create"
	COMMAND_NAME_RESOURCES_UPDATE         = "resources update"
	COMMAND_NAME_RESOURCES_DELETE         = "resources delete"
	COMMAND_NAME_USERS                    = "users"
	COMMAND_NAME_USERS_GET                = "users get"
)

// -----------------------------------------------------------------
// Constructors
// -----------------------------------------------------------------
func NewCommandCollection(factory spi.Factory) (CommandCollection, error) {

	commands := new(commandCollectionImpl)

	err := commands.init(factory)

	return commands, err
}

// -----------------------------------------------------------------
// Public functions
// -----------------------------------------------------------------

func (commands *commandCollectionImpl) GetRootCommand() spi.GalasaCommand {
	cmd, _ := commands.GetCommand(COMMAND_NAME_ROOT)
	return cmd
}

func (commands *commandCollectionImpl) GetCommand(name string) (spi.GalasaCommand, error) {
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
func (commands *commandCollectionImpl) init(factory spi.Factory) error {

	commands.commandMap = make(map[string]spi.GalasaCommand)

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
		err = commands.addUsersCommands(factory, rootCommand)
	}

	if err == nil {
		commands.setHelpFlags()
	}

	return err
}

func (commands *commandCollectionImpl) addAuthCommands(factory spi.Factory, rootCommand spi.GalasaCommand) error {
	var err error
	var authCommand spi.GalasaCommand
	var authLoginCommand spi.GalasaCommand
	var authLogoutCommand spi.GalasaCommand

	authCommand, err = NewAuthCommand(rootCommand)
	if err == nil {
		authLoginCommand, err = NewAuthLoginCommand(factory, authCommand, rootCommand)
		if err == nil {
			authLogoutCommand, err = NewAuthLogoutCommand(factory, authCommand, rootCommand)
			if err == nil {
				err = commands.addAuthTokensCommands(factory, authCommand, rootCommand)
			}
		}
	}

	if err == nil {
		commands.commandMap[authCommand.Name()] = authCommand
		commands.commandMap[authLoginCommand.Name()] = authLoginCommand
		commands.commandMap[authLogoutCommand.Name()] = authLogoutCommand
	}

	return err
}

func (commands *commandCollectionImpl) addAuthTokensCommands(factory spi.Factory, authCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) error {
	var err error
	var authTokensCommand spi.GalasaCommand
	var authTokensGetCommand spi.GalasaCommand
	var authTokensDeleteCommand spi.GalasaCommand

	authTokensCommand, err = NewAuthTokensCommand(authCommand, rootCommand)
	if err == nil {
		authTokensGetCommand, err = NewAuthTokensGetCommand(factory, authTokensCommand, rootCommand)
		if err == nil {
			authTokensDeleteCommand, err = NewAuthTokensDeleteCommand(factory, authTokensCommand, rootCommand)
		}
	}

	if err == nil {
		commands.commandMap[authTokensCommand.Name()] = authTokensCommand
		commands.commandMap[authTokensGetCommand.Name()] = authTokensGetCommand
		commands.commandMap[authTokensDeleteCommand.Name()] = authTokensDeleteCommand
	}

	return err
}

func (commands *commandCollectionImpl) addLocalCommands(factory spi.Factory, rootCommand spi.GalasaCommand) error {
	var err error
	var localCommand spi.GalasaCommand
	var localInitCommand spi.GalasaCommand

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

func (commands *commandCollectionImpl) addProjectCommands(factory spi.Factory, rootCommand spi.GalasaCommand) error {
	var err error

	var projectCommand spi.GalasaCommand

	projectCommand, err = NewProjectCmd(rootCommand)
	if err == nil {
		commands.commandMap[projectCommand.Name()] = projectCommand
	}

	if err == nil {
		var projectCreateCommand spi.GalasaCommand
		projectCreateCommand, err = NewProjectCreateCmd(factory, projectCommand, rootCommand)
		if err == nil {
			commands.commandMap[projectCreateCommand.Name()] = projectCreateCommand
		}
	}
	return err
}

func (commands *commandCollectionImpl) addPropertiesCommands(factory spi.Factory, rootCommand spi.GalasaCommand) error {
	var err error
	var propertiesCommand spi.GalasaCommand
	var propertiesGetCommand spi.GalasaCommand
	var propertiesDeleteCommand spi.GalasaCommand
	var propertiesSetCommand spi.GalasaCommand

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

func (commands *commandCollectionImpl) addPropertiesNamespaceCommands(factory spi.Factory, rootCommand spi.GalasaCommand, propertiesCommand spi.GalasaCommand) error {
	var err error
	var propertiesNamespaceCommand spi.GalasaCommand
	var propertiesNamespaceGetCommand spi.GalasaCommand

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

func (commands *commandCollectionImpl) addRunsCommands(factory spi.Factory, rootCommand spi.GalasaCommand) error {

	var err error
	var runsCommand spi.GalasaCommand
	var runsDownloadCommand spi.GalasaCommand
	var runsGetCommand spi.GalasaCommand
	var runsPrepareCommand spi.GalasaCommand
	var runsSubmitCommand spi.GalasaCommand
	var runsSubmitLocalCommand spi.GalasaCommand
	var runsResetCommand spi.GalasaCommand
	var runsCancelCommand spi.GalasaCommand
	var runsDeleteCommand spi.GalasaCommand

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
								if err == nil {
									runsDeleteCommand, err = NewRunsDeleteCommand(factory, runsCommand, rootCommand)
								}
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
		commands.commandMap[runsDeleteCommand.Name()] = runsDeleteCommand
	}

	return err
}

func (commands *commandCollectionImpl) addResourcesCommands(factory spi.Factory, rootCommand spi.GalasaCommand) error {

	var err error
	var resourcesCommand spi.GalasaCommand
	var resourcesApplyCommand spi.GalasaCommand
	var resourcesCreateCommand spi.GalasaCommand
	var resourcesUpdateCommand spi.GalasaCommand
	var resourcesDeleteCommand spi.GalasaCommand

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

func (commands *commandCollectionImpl) addUsersCommands(factory spi.Factory, rootCommand spi.GalasaCommand) error {

	var err error
	var usersCommand spi.GalasaCommand
	var usersGetCommand spi.GalasaCommand

	usersCommand, err = NewUsersCommand(rootCommand)

	if err == nil {
		usersGetCommand, err = NewUsersGetCommand(factory, usersCommand, rootCommand)
	}

	if err == nil {
		commands.commandMap[usersCommand.Name()] = usersCommand
		commands.commandMap[usersGetCommand.Name()] = usersGetCommand
	}

	return err
}

func (commands *commandCollectionImpl) setHelpFlags() {
	for _, command := range commands.commandMap {
		command.CobraCommand().Flags().BoolP("help", "h", false, "Displays the options for the '"+command.Name()+"' command.")
	}
}
