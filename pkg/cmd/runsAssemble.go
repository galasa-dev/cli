//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package cmd

import (
	"fmt"
	"os"

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

    portfolioFilename string
    stream string

)


func init() {
    runsAssembleCmd.Flags().StringVarP(&portfolioFilename, "portfolio", "p", "", "portfolio to add tests to")
    runsAssembleCmd.Flags().StringVarP(&stream, "stream", "s", "", "test stream to extract the tests from")
    runsAssembleCmd.MarkFlagRequired("portfolio")
    runsAssembleCmd.MarkFlagRequired("stream")

    utils.AddCommandFlags(runsAssembleCmd)

    runsCmd.AddCommand(runsAssembleCmd)
}

func executeAssemble(cmd *cobra.Command, args []string) {
    fmt.Println("Galasa CLI - Assemble tests")

    apiClient := api.InitialiseAPI(bootstrap)

    availableStreams := utils.FetchTestStreams(apiClient)

    err := utils.ValidateStream(availableStreams, stream)
    if err != nil {
        fmt.Printf("%v", err)
        os.Exit(1)
    }

    testCatalog, err := utils.FetchTestCatalog(apiClient, stream)
    if err != nil {
        panic(err)
    }

    fmt.Println("Test catalog retrieved")

    testSelection := utils.SelectTests(testCatalog, stream)

    portfolio := utils.CreatePortfolio(&testSelection)

    utils.WritePortfolio(portfolio, portfolioFilename)

    fmt.Println("Portfolio created")
}

