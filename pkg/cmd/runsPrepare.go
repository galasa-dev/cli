/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"
	"strings"

	"github.com/galasa-dev/cli/pkg/api"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/launcher"
	"github.com/galasa-dev/cli/pkg/runs"
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

type RunsPrepareCmdValues struct {
	portfolioFilename    string
	prepareFlagOverrides *[]string
	prepareAppend        *bool

	prepareSelectionFlags *utils.TestSelectionFlagValues
}

type RunsPrepareCommand struct {
	values       *RunsPrepareCmdValues
	cobraCommand *cobra.Command
}

func NewRunsPrepareCommand(factory spi.Factory, runsCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) (spi.GalasaCommand, error) {
	cmd := new(RunsPrepareCommand)
	err := cmd.init(factory, runsCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsPrepareCommand) Name() string {
	return COMMAND_NAME_RUNS_PREPARE
}

func (cmd *RunsPrepareCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsPrepareCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *RunsPrepareCommand) init(factory spi.Factory, runsCommand spi.GalasaCommand, rootCommand spi.GalasaCommand) error {
	var err error
	cmd.values = &RunsPrepareCmdValues{}
	cmd.cobraCommand, err = cmd.createCobraCommand(
		factory,
		runsCommand,
		rootCommand.Values().(*RootCmdValues),
	)
	return err
}
func (cmd *RunsPrepareCommand) createCobraCommand(
	factory spi.Factory,
	runsCommand spi.GalasaCommand,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {
	var err error

	cmd.values.prepareSelectionFlags = runs.NewTestSelectionFlagValues()

	runsPrepareCobraCmd := &cobra.Command{
		Use:     "prepare",
		Short:   "prepares a list of tests",
		Long:    "Prepares a list of tests from a test catalog providing specific overrides if required",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs prepare"},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.executeAssemble(factory, runsCommand.Values().(*RunsCmdValues), rootCmdValues)
		},
	}

	runsPrepareCobraCmd.Flags().StringVarP(&cmd.values.portfolioFilename, "portfolio", "p", "", "portfolio to add tests to")
	cmd.values.prepareFlagOverrides = runsPrepareCobraCmd.Flags().StringSlice("override", make([]string, 0), "overrides to be sent with the tests (overrides in the portfolio will take precedence)")
	cmd.values.prepareAppend = runsPrepareCobraCmd.Flags().Bool("append", false, "Append tests to existing portfolio")
	runsPrepareCobraCmd.MarkFlagRequired("portfolio")

	runs.AddCommandFlags(runsPrepareCobraCmd, cmd.values.prepareSelectionFlags)

	runsCommand.CobraCommand().AddCommand(runsPrepareCobraCmd)

	return runsPrepareCobraCmd, err
}

func (cmd *RunsPrepareCommand) executeAssemble(
	factory spi.Factory,
	runsCmdValues *RunsCmdValues,
	rootCmdValues *RootCmdValues,
) error {
	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err == nil {
		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Assemble tests")

		// Get the ability to query environment variables.
		env := factory.GetEnvironment()

		var galasaHome spi.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Convert overrides to a map
			testOverrides := make(map[string]string)
			for _, override := range *cmd.values.prepareFlagOverrides {
				pos := strings.Index(override, "=")
				if pos < 1 {
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PREPARE_INVALID_OVERRIDE, override)
					break
				}
				key := override[:pos]
				value := override[pos+1:]
				if value == "" {
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PREPARE_INVALID_OVERRIDE, override)
					break
				}

				testOverrides[key] = value
			}

			if err == nil {

				// Load the bootstrap properties.
				var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
				var bootstrapData *api.BootstrapData
				bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, runsCmdValues.bootstrap, urlService)
				if err == nil {

					// Create an API client
					var apiClient *galasaapi.APIClient
					apiServerUrl := bootstrapData.ApiServerURL
					authenticator := factory.GetAuthenticator(
						apiServerUrl,
						galasaHome,
					)
					apiClient, err = authenticator.GetAuthenticatedAPIClient()
					if err == nil {
						launcher := launcher.NewRemoteLauncher(apiServerUrl, apiClient)

						validator := runs.NewStreamBasedValidator()
						err = validator.Validate(cmd.values.prepareSelectionFlags)
						if err == nil {

							var testSelection runs.TestSelection
							testSelection, err = runs.SelectTests(launcher, cmd.values.prepareSelectionFlags)
							if err == nil {

								count := len(testSelection.Classes)
								if count < 1 {
									err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_NO_TESTS_SELECTED)
								} else {
									if count == 1 {
										log.Println("1 test was selected")
									} else {
										log.Printf("%v tests were selected", count)
									}
								}

								if err == nil {

									var portfolio *runs.Portfolio
									if *cmd.values.prepareAppend {
										portfolio, err = runs.ReadPortfolio(fileSystem, cmd.values.portfolioFilename)
									} else {
										portfolio = runs.NewPortfolio()
									}

									if err == nil {
										runs.AddClassesToPortfolio(&testSelection, &testOverrides, portfolio)

										err = runs.WritePortfolio(fileSystem, cmd.values.portfolioFilename, portfolio)
										if err == nil {
											if *cmd.values.prepareAppend {
												log.Println("Portfolio appended")
											} else {
												log.Println("Portfolio created")
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return err
}
