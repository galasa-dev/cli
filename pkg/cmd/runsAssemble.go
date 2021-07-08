//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/galasa.dev/cli/pkg/api"
	"github.com/galasa.dev/cli/pkg/utils"
)

var (
        runsAssembleCmd = &cobra.Command{
            Use:   "assemble",
            Short: "assembles a list of tests",
            Long:  "Assembles a list of tests from a test catalog providing specific overrides if required",
            Args: cobra.NoArgs,
            Run:   executeAssemble,
    }

    portfolioFilename      string
    assembleFlagOverrides  *[]string
    assembleAppend         *bool

    assembleSelectionFlags = utils.TestSelectionFlags{}
)


func init() {
    runsAssembleCmd.Flags().StringVarP(&portfolioFilename, "portfolio", "p", "", "portfolio to add tests to")
    assembleFlagOverrides = runsAssembleCmd.Flags().StringSlice("override", make([]string, 0), "overrides to be sent with the tests (overrides in the portfolio will take precedence)")
    assembleAppend = runsAssembleCmd.Flags().Bool("append", false, "Append tests to existing portfolio")
    runsAssembleCmd.MarkFlagRequired("portfolio")

    utils.AddCommandFlags(runsAssembleCmd, &assembleSelectionFlags)

    runsCmd.AddCommand(runsAssembleCmd)
}

func executeAssemble(cmd *cobra.Command, args []string) {
    fmt.Println("Galasa CLI - Assemble tests")

    // Convert overrides to a map
    testOverrides := make(map[string]string)
    for _, override := range *assembleFlagOverrides {
        pos := strings.Index(override, "=")
        if (pos < 1) {
            fmt.Printf("Invalid override '%v'",override)
            os.Exit(1)
        }
        key := override[:pos]
        value := override[pos+1:]
        if value == "" {
            fmt.Printf("Invalid override '%v'",override)
            os.Exit(1)
        }

        testOverrides[key] = value
    }

    apiClient := api.InitialiseAPI(bootstrap)

    testSelection := utils.SelectTests(apiClient, &assembleSelectionFlags)
    count := len(testSelection.Classes)
    if count < 1 {
        fmt.Println("No tests were selected")
        os.Exit(1)
    } else {
        if count == 1 {
            fmt.Println("1 test was selected")
        } else {
            fmt.Printf("%v tests were selected", count)
        }
    }

    var portfolio utils.Portfolio
    if *assembleAppend {
        portfolio = utils.LoadPortfolio(portfolioFilename)
    } else {
        portfolio = utils.NewPortfolio()
    }

    utils.CreatePortfolio(&testSelection, &testOverrides, &portfolio)

    utils.WritePortfolio(portfolio, portfolioFilename)

    if *assembleAppend {
        fmt.Println("Portfolio appended")
    } else {
        fmt.Println("Portfolio created")
    }    
}

