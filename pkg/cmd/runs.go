/*
 * Copyright contributors to the Galasa project
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
	bootstrap string
)

func init() {
	cmd := runsCmd
	parentCmd := RootCmd

	cmd.PersistentFlags().StringVarP(&bootstrap, "bootstrap", "b", "", "Bootstrap URL")

	parentCmd.AddCommand(runsCmd)
}
