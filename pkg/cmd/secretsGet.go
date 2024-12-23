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

type SecretsGetCmdValues struct {
	outputFormat string
}

type SecretsGetCommand struct {
	values *SecretsGetCmdValues
    cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewSecretsGetCommand(
    factory spi.Factory,
    secretsGetCommand spi.GalasaCommand,
    commsCmd spi.GalasaCommand,
) (spi.GalasaCommand, error) {

    cmd := new(SecretsGetCommand)

    err := cmd.init(factory, secretsGetCommand, commsCmd)
    return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *SecretsGetCommand) Name() string {
    return COMMAND_NAME_SECRETS_GET
}

func (cmd *SecretsGetCommand) CobraCommand() *cobra.Command {
    return cmd.cobraCommand
}

func (cmd *SecretsGetCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *SecretsGetCommand) init(factory spi.Factory, secretsCommand spi.GalasaCommand, commsCmd spi.GalasaCommand) error {
    var err error

	cmd.values = &SecretsGetCmdValues{}
    cmd.cobraCommand, err = cmd.createCobraCmd(factory, secretsCommand, commsCmd.Values().(*CommsCmdValues))

    return err
}

func (cmd *SecretsGetCommand) createCobraCmd(
    factory spi.Factory,
    secretsCommand spi.GalasaCommand,
    commsCommandValues *CommsCmdValues,
) (*cobra.Command, error) {

    var err error

    secretsCommandValues := secretsCommand.Values().(*SecretsCmdValues)
    secretsGetCobraCmd := &cobra.Command{
        Use:     "get",
        Short:   "Get secrets from the credentials store",
        Long:    "Get a list of secrets or a specific secret from the credentials store",
        Aliases: []string{COMMAND_NAME_SECRETS_GET},
        RunE: func(cobraCommand *cobra.Command, args []string) error {
			executionFunc := func() error {
            	return cmd.executeSecretsGet(factory, secretsCommand.Values().(*SecretsCmdValues), commsCommandValues)
			}
			return executeCommandWithRetries(factory, commsCommandValues, executionFunc)
        },
    }

    addSecretNameFlag(secretsGetCobraCmd, false, secretsCommandValues)

	formatters := secrets.GetFormatterNamesAsString()
	secretsGetCobraCmd.Flags().StringVar(&cmd.values.outputFormat, "format", "summary", "the output format of the returned secrets. Supported formats are: "+formatters+".")

    secretsCommand.CobraCommand().AddCommand(secretsGetCobraCmd)

    return secretsGetCobraCmd, err
}

func (cmd *SecretsGetCommand) executeSecretsGet(
    factory spi.Factory,
    secretsCmdValues *SecretsCmdValues,
    commsCmdValues *CommsCmdValues,
) error {

    var err error
    // Operations on the file system will all be relative to the current folder.
    fileSystem := factory.GetFileSystem()

	commsCmdValues.isCapturingLogs = true

	log.Println("Galasa CLI - Get secrets from the ecosystem")

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
				err = secrets.GetSecrets(secretsCmdValues.name, cmd.values.outputFormat, console, apiClient, byteReader)
			}
		}
	}

    return err
}
