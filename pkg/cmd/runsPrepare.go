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
	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/launcher"
	"github.com/galasa-dev/cli/pkg/runs"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

type RunsPrepareCmdValues struct {
	portfolioFilename    string
	prepareFlagOverrides *[]string
	prepareAppend        *bool

	prepareSelectionFlags *utils.TestSelectionFlagValues
}

func createRunsPrepareCmd(parentCmd *cobra.Command, runsCmdValues *RunsCmdValues, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	var err error = nil

	runsPrepareCmdValues := &RunsPrepareCmdValues{}

	runsPrepareCmdValues.prepareSelectionFlags = runs.NewTestSelectionFlagValues()

	runsPrepareCmd := &cobra.Command{
		Use:     "prepare",
		Short:   "prepares a list of tests",
		Long:    "Prepares a list of tests from a test catalog providing specific overrides if required",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs prepare"},
		Run: func(cmd *cobra.Command, args []string) {
			executeAssemble(cmd, args, runsPrepareCmdValues, runsCmdValues, rootCmdValues)
		},
	}

	runsPrepareCmd.Flags().StringVarP(&runsPrepareCmdValues.portfolioFilename, "portfolio", "p", "", "portfolio to add tests to")
	runsPrepareCmdValues.prepareFlagOverrides = runsPrepareCmd.Flags().StringSlice("override", make([]string, 0), "overrides to be sent with the tests (overrides in the portfolio will take precedence)")
	runsPrepareCmdValues.prepareAppend = runsPrepareCmd.Flags().Bool("append", false, "Append tests to existing portfolio")
	runsPrepareCmd.MarkFlagRequired("portfolio")

	runs.AddCommandFlags(runsPrepareCmd, runsPrepareCmdValues.prepareSelectionFlags)

	parentCmd.AddCommand(runsPrepareCmd)

	// There are no sub-command children to add to the command tree.

	return runsPrepareCmd, err
}

func executeAssemble(cmd *cobra.Command, args []string, runsPrepareCmdValues *RunsPrepareCmdValues, runsCmdValues *RunsCmdValues, rootCmdValues *RootCmdValues) {
	var err error = nil

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err != nil {
		panic(err)
	}
	rootCmdValues.isCapturingLogs = true

	log.Println("Galasa CLI - Assemble tests")

	// Get the ability to query environment variables.
	env := utils.NewEnvironment()

	galasaHome, err := utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
	if err != nil {
		panic(err)
	}

	// Convert overrides to a map
	testOverrides := make(map[string]string)
	for _, override := range *runsPrepareCmdValues.prepareFlagOverrides {
		pos := strings.Index(override, "=")
		if pos < 1 {
			err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PREPARE_INVALID_OVERRIDE, override)
			panic(err)
		}
		key := override[:pos]
		value := override[pos+1:]
		if value == "" {
			err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PREPARE_INVALID_OVERRIDE, override)
			panic(err)
		}

		testOverrides[key] = value
	}

	// Load the bootstrap properties.
	var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
	bootstrapData, err := api.LoadBootstrap(
		galasaHome, fileSystem, env, runsCmdValues.bootstrap, urlService)
	if err != nil {
		panic(err)
	}

	// Create an API client
	launcher := launcher.NewRemoteLauncher(bootstrapData.ApiServerURL)

	validator := runs.NewStreamBasedValidator()
	err = validator.Validate(runsPrepareCmdValues.prepareSelectionFlags)
	if err != nil {
		panic(err)
	}

	testSelection, err := runs.SelectTests(launcher, runsPrepareCmdValues.prepareSelectionFlags)
	if err != nil {
		panic(err)
	}

	count := len(testSelection.Classes)
	if count < 1 {
		err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_NO_TESTS_SELECTED)
		panic(err)
	} else {
		if count == 1 {
			log.Println("1 test was selected")
		} else {
			log.Printf("%v tests were selected", count)
		}
	}

	var portfolio *runs.Portfolio
	if *runsPrepareCmdValues.prepareAppend {
		portfolio, err = runs.ReadPortfolio(fileSystem, runsPrepareCmdValues.portfolioFilename)
		if err != nil {
			panic(err)
		}
	} else {
		portfolio = runs.NewPortfolio()
	}

	runs.AddClassesToPortfolio(&testSelection, &testOverrides, portfolio)

	err = runs.WritePortfolio(fileSystem, runsPrepareCmdValues.portfolioFilename, portfolio)
	if err != nil {
		panic(err)
	}

	if *runsPrepareCmdValues.prepareAppend {
		log.Println("Portfolio appended")
	} else {
		log.Println("Portfolio created")
	}
}
