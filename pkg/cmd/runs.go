/*
*  Licensed Materials - Property of IBM
*
* (c) Copyright IBM Corp. 2021.
*/

package cmd

import (
    "github.com/spf13/cobra"
)

var (
        runsCmd = &cobra.Command{
        Use:   "runs",
        Short: "Manage test runs in the ecosystem",
        Long:  "Assembles, submits and monitors test runs in Galasa Ecosystem",
    }
)


func init() {
    rootCmd.AddCommand(runsCmd)
}
