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
    commsFlagSet GalasaFlagSet,
) (spi.GalasaCommand, error) {

    cmd := new(SecretsDeleteCommand)

    err := cmd.init(factory, secretsDeleteCommand, commsFlagSet)
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
func (cmd *SecretsDeleteCommand) init(factory spi.Factory, secretsCommand spi.GalasaCommand, commsFlagSet GalasaFlagSet) error {
    var err error

    cmd.cobraCommand, err = cmd.createCobraCmd(factory, secretsCommand, commsFlagSet.Values().(*CommsFlagSetValues))

    return err
}

func (cmd *SecretsDeleteCommand) createCobraCmd(
    factory spi.Factory,
    secretsCommand spi.GalasaCommand,
    commsFlagSetValues *CommsFlagSetValues,
) (*cobra.Command, error) {

    var err error

    secretsCommandValues := secretsCommand.Values().(*SecretsCmdValues)
    secretsDeleteCobraCmd := &cobra.Command{
        Use:     "delete",
        Short:   "Deletes a secret from the credentials store",
        Long:    "Deletes a secret from the credentials store",
        Aliases: []string{COMMAND_NAME_SECRETS_DELETE},
        RunE: func(cobraCommand *cobra.Command, args []string) error {
			return cmd.executeSecretsDelete(factory, secretsCommand.Values().(*SecretsCmdValues), commsFlagSetValues)
        },
    }

    addSecretNameFlag(secretsDeleteCobraCmd, true, secretsCommandValues)

    secretsCommand.CobraCommand().AddCommand(secretsDeleteCobraCmd)

    return secretsDeleteCobraCmd, err
}

func (cmd *SecretsDeleteCommand) executeSecretsDelete(
    factory spi.Factory,
    secretsCmdValues *SecretsCmdValues,
    commsFlagSetValues *CommsFlagSetValues,
) error {

    var err error
    // Operations on the file system will all be relative to the current folder.
    fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, commsFlagSetValues.logFileName)
	if err == nil {
		commsFlagSetValues.isCapturingLogs = true
	
		log.Println("Galasa CLI - Delete a secret from the credentials store")
	
		env := factory.GetEnvironment()
	
		var galasaHome spi.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, commsFlagSetValues.CmdParamGalasaHomePath)
		if err == nil {

			timeService := factory.GetTimeService()
			commsRetrier := api.NewCommsRetrier(commsFlagSetValues.maxRetries, commsFlagSetValues.retryBackoffSeconds, timeService)

			// Read the bootstrap properties, retrying if a rate limit has been exceeded
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			loadBootstrapWithRetriesFunc := func() error {
				bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, commsFlagSetValues.bootstrap, urlService)
				return err
			}

			err = commsRetrier.ExecuteCommandWithRateLimitRetries(loadBootstrapWithRetriesFunc)
			if err == nil {
	
				var console = factory.GetStdOutConsole()
	
				apiServerUrl := bootstrapData.ApiServerURL
				log.Printf("The API server is at '%s'\n", apiServerUrl)
	
				authenticator := factory.GetAuthenticator(
					apiServerUrl,
					galasaHome,
				)

				commsRetrier, err = api.NewCommsRetrierWithAPIClient(
					commsFlagSetValues.maxRetries,
					commsFlagSetValues.retryBackoffSeconds,
					timeService,
					authenticator,
				)

				byteReader := factory.GetByteReader()
	
				if err == nil {
					deleteSecretFunc := func(apiClient *galasaapi.APIClient) error {
						return secrets.DeleteSecret(secretsCmdValues.name, console, apiClient, byteReader)
					}
					err = commsRetrier.ExecuteCommandWithRetries(deleteSecretFunc)
				}
			}
		}
	}

    return err
}
