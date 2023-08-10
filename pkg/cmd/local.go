/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	localCmd = &cobra.Command{
		Use:   "local",
		Short: "Manipulate local system",
		Long:  "Manipulate local system",
	}
)

func init() {
	RootCmd.AddCommand(localCmd)
}
