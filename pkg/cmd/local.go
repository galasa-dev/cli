/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

func createLocalCmd(parentCmd *cobra.Command, rootCmdValues *RootCmdValues) (*cobra.Command, error) {

	localCmd := &cobra.Command{
		Use:   "local",
		Short: "Manipulate local system",
		Long:  "Manipulate local system",
	}
	parentCmd.AddCommand(localCmd)

	err := createLocalCmdChildren(localCmd, rootCmdValues)
	return localCmd, err
}

func createLocalCmdChildren(localCmd *cobra.Command, rootCmdValues *RootCmdValues) error {
	_, err := createLocalInitCmd(localCmd, rootCmdValues)
	return err
}
