/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"log"
	"strings"

	"github.com/galasa.dev/cli/pkg/api"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
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

	prepareSelectionFlags = utils.TestSelectionFlags{}
)

func init() {
	runsAssembleCmd.Flags().StringVarP(&portfolioFilename, "portfolio", "p", "", "portfolio to add tests to")
	prepareFlagOverrides = runsAssembleCmd.Flags().StringSlice("override", make([]string, 0), "overrides to be sent with the tests (overrides in the portfolio will take precedence)")
	prepareAppend = runsAssembleCmd.Flags().Bool("append", false, "Append tests to existing portfolio")
	runsAssembleCmd.MarkFlagRequired("portfolio")
	utils.AddCommandFlags(runsAssembleCmd, &prepareSelectionFlags)

	runsCmd.AddCommand(runsAssembleCmd)
}

func executeAssemble(cmd *cobra.Command, args []string) {

	utils.CaptureLog(logFileName)

	log.Println("Galasa CLI - Assemble tests")

	// Operations on the file system will all be relative to the current folder.
	fileSystem := utils.NewOSFileSystem()

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

	apiClient, err := api.InitialiseAPI(bootstrap)
	if err != nil {
		panic(err)
	}

	testSelection := utils.SelectTests(apiClient, &prepareSelectionFlags)
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

	var portfolio *utils.Portfolio
	if *prepareAppend {
		portfolio, err = utils.LoadPortfolio(fileSystem, portfolioFilename)
		if err != nil {
			panic(err)
		}
	} else {
		portfolio = utils.NewPortfolio()
	}

	utils.CreatePortfolio(&testSelection, &testOverrides, portfolio)

	err = utils.WritePortfolio(fileSystem, portfolioFilename, portfolio)
	if err != nil {
		panic(err)
	}

	if *prepareAppend {
		log.Println("Portfolio appended")
	} else {
		log.Println("Portfolio created")
	}
}
