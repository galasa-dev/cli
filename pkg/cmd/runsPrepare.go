/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"
	"strings"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/launcher"
	"github.com/galasa.dev/cli/pkg/runs"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	runsAssembleCmd = &cobra.Command{
		Use:   "prepare",
		Short: "prepares a list of tests",
		Long:  "Prepares a list of tests from a test catalog providing specific overrides if required",
		Args:  cobra.NoArgs,
		Run:   executeAssemble,
	}

	portfolioFilename    string
	prepareFlagOverrides *[]string
	prepareAppend        *bool

	prepareSelectionFlags = runs.TestSelectionFlags{}
)

func init() {
	runsAssembleCmd.Flags().StringVarP(&portfolioFilename, "portfolio", "p", "", "portfolio to add tests to")
	prepareFlagOverrides = runsAssembleCmd.Flags().StringSlice("override", make([]string, 0), "overrides to be sent with the tests (overrides in the portfolio will take precedence)")
	prepareAppend = runsAssembleCmd.Flags().Bool("append", false, "Append tests to existing portfolio")
	runsAssembleCmd.MarkFlagRequired("portfolio")
	runs.AddCommandFlags(runsAssembleCmd, &prepareSelectionFlags)

	runsCmd.AddCommand(runsAssembleCmd)
}

func executeAssemble(cmd *cobra.Command, args []string) {
	var err error = nil

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, logFileName)
	if err != nil {
		panic(err)
	}
	isCapturingLogs = true

	log.Println("Galasa CLI - Assemble tests")

	// Get the ability to query environment variables.
	env := utils.NewEnvironment()

	galasaHome, err := utils.NewGalasaHome(fileSystem, env, CmdParamGalasaHomePath)
	if err != nil {
		panic(err)
	}

	// Convert overrides to a map
	testOverrides := make(map[string]string)
	for _, override := range *prepareFlagOverrides {
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
		galasaHome, fileSystem, env, bootstrap, urlService)
	if err != nil {
		panic(err)
	}

	// Create an API client
	launcher := launcher.NewRemoteLauncher(bootstrapData.ApiServerURL)

	testSelection, err := runs.SelectTests(launcher, &prepareSelectionFlags)
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
	if *prepareAppend {
		portfolio, err = runs.ReadPortfolio(fileSystem, portfolioFilename)
		if err != nil {
			panic(err)
		}
	} else {
		portfolio = runs.NewPortfolio()
	}

	runs.AddClassesToPortfolio(&testSelection, &testOverrides, portfolio)

	err = runs.WritePortfolio(fileSystem, portfolioFilename, portfolio)
	if err != nil {
		panic(err)
	}

	if *prepareAppend {
		log.Println("Portfolio appended")
	} else {
		log.Println("Portfolio created")
	}
}
