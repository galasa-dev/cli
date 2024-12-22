/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/secrets"
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

type SecretsDeleteCommand struct {
    cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewSecretsDeleteCommand(
    factory spi.Factory,
    secretsDeleteCommand spi.GalasaCommand,
    commsCmd spi.GalasaCommand,
) (spi.GalasaCommand, error) {

    cmd := new(SecretsDeleteCommand)

    err := cmd.init(factory, secretsDeleteCommand, commsCmd)
    return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *SecretsDeleteCommand) Name() string {
    return COMMAND_NAME_SECRETS_DELETE
}

func (cmd *SecretsDeleteCommand) CobraCommand() *cobra.Command {
    return cmd.cobraCommand
}

func (cmd *SecretsDeleteCommand) Values() interface{} {
    return nil
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *SecretsDeleteCommand) init(factory spi.Factory, secretsCommand spi.GalasaCommand, commsCmd spi.GalasaCommand) error {
    var err error

    cmd.cobraCommand, err = cmd.createCobraCmd(factory, secretsCommand, commsCmd.Values().(*CommsCmdValues))

    return err
}

func (cmd *SecretsDeleteCommand) createCobraCmd(
    factory spi.Factory,
    secretsCommand spi.GalasaCommand,
    commsCommandValues *CommsCmdValues,
) (*cobra.Command, error) {

    var err error

    secretsCommandValues := secretsCommand.Values().(*SecretsCmdValues)
    secretsDeleteCobraCmd := &cobra.Command{
        Use:     "delete",
        Short:   "Deletes a secret from the credentials store",
        Long:    "Deletes a secret from the credentials store",
        Aliases: []string{COMMAND_NAME_SECRETS_DELETE},
        RunE: func(cobraCommand *cobra.Command, args []string) error {
            return cmd.executeSecretsDelete(factory, secretsCommand.Values().(*SecretsCmdValues), commsCommandValues)
        },
    }

    addSecretNameFlag(secretsDeleteCobraCmd, true, secretsCommandValues)

    secretsCommand.CobraCommand().AddCommand(secretsDeleteCobraCmd)

    return secretsDeleteCobraCmd, err
}

func (cmd *SecretsDeleteCommand) executeSecretsDelete(
    factory spi.Factory,
    secretsCmdValues *SecretsCmdValues,
    commsCmdValues *CommsCmdValues,
) error {

    var err error
    // Operations on the file system will all be relative to the current folder.
    fileSystem := factory.GetFileSystem()

    err = utils.CaptureLog(fileSystem, commsCmdValues.logFileName)

    if err == nil {
        commsCmdValues.isCapturingLogs = true

        log.Println("Galasa CLI - Delete a secret from the credentials store")

        env := factory.GetEnvironment()

        var galasaHome spi.GalasaHome
        galasaHome, err = utils.NewGalasaHome(fileSystem, env, commsCmdValues.CmdParamGalasaHomePath)
        if err == nil {

            var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
            var bootstrapData *api.BootstrapData
            bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, commsCmdValues.bootstrap, urlService)
            if err == nil {

				var console = factory.GetStdOutConsole()

                apiServerUrl := bootstrapData.ApiServerURL
                log.Printf("The API server is at '%s'\n", apiServerUrl)

                authenticator := factory.GetAuthenticator(
                    apiServerUrl,
                    galasaHome,
                )

                var apiClient *galasaapi.APIClient
                apiClient, err = authenticator.GetAuthenticatedAPIClient()

				byteReader := factory.GetByteReader()

                if err == nil {
                    err = secrets.DeleteSecret(secretsCmdValues.name, console, apiClient, byteReader)
                }
            }
        }
    }
    return err
}
