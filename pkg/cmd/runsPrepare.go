//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/galasa.dev/cli/pkg/api"
	"github.com/galasa.dev/cli/pkg/utils"
)

var (
        runsAssembleCmd = &cobra.Command{
            Use:   "prepare",
            Short: "prepares a list of tests",
            Long:  "Prepares a list of tests from a test catalog providing specific overrides if required",
            Args: cobra.NoArgs,
            Run:   executeAssemble,
    }

    portfolioFilename      string
    prepareFlagOverrides  *[]string
    prepareAppend         *bool

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
    log.Println("Galasa CLI - Assemble tests")

    // Convert overrides to a map
    testOverrides := make(map[string]string)
    for _, override := range *prepareFlagOverrides {
        pos := strings.Index(override, "=")
        if (pos < 1) {
            log.Printf("Invalid override '%v'",override)
            os.Exit(1)
        }
        key := override[:pos]
        value := override[pos+1:]
        if value == "" {
            log.Printf("Invalid override '%v'",override)
            os.Exit(1)
        }

        testOverrides[key] = value
    }

    apiClient := api.InitialiseAPI(bootstrap)

    testSelection := utils.SelectTests(apiClient, &prepareSelectionFlags)
    count := len(testSelection.Classes)
    if count < 1 {
        log.Println("No tests were selected")
        os.Exit(1)
    } else {
        if count == 1 {
            log.Println("1 test was selected")
        } else {
            log.Printf("%v tests were selected", count)
        }
    }

    var portfolio utils.Portfolio
    if *prepareAppend {
        portfolio = utils.LoadPortfolio(portfolioFilename)
    } else {
        portfolio = utils.NewPortfolio()
    }

    utils.CreatePortfolio(&testSelection, &testOverrides, &portfolio)

    utils.WritePortfolio(portfolio, portfolioFilename)

    if *prepareAppend {
        log.Println("Portfolio appended")
    } else {
        log.Println("Portfolio created")
    }    
}

