//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

var (
        runsAssembleCmd = &cobra.Command{
            Use:   "assemble",
            Short: "assembles a list of tests",
            Long:  "Assembles a list of tests from a test catalog providing specific overrides if required",
            Run:   execute,
    }

    assembleOutputFile string

)


func init() {
    runsAssembleCmd.PersistentFlags().StringVarP(&assembleOutputFile, "output", "o", "", "output file to add tests to")

    runsCmd.AddCommand(runsAssembleCmd)
}

func execute(cmd *cobra.Command, args []string) {
    fmt.Println("Galasa CLI - Assemble tests")
}