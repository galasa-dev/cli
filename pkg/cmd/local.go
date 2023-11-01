/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

func createLocalCmd(rootCmd *cobra.Command) (*cobra.Command, error) {
	localCmd := &cobra.Command{
		Use:   "local",
		Short: "Manipulate local system",
		Long:  "Manipulate local system",
	}
	rootCmd.AddCommand(localCmd)

	err := createLocalCmdChildren(localCmd)
	return localCmd, err
}

func createLocalCmdChildren(localCmd *cobra.Command) error {
	_, err := createLocalInitCmd(localCmd)
	return err
}
