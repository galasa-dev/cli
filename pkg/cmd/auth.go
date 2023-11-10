/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

func createAuthCmd(factory Factory, parentCmd *cobra.Command, rootCmdValues *RootCmdValues) (*cobra.Command, error) {

	var err error = nil

	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manages the authentication of users with a Galasa ecosystem",
		Long:  "Manages the authentication of users with a Galasa ecosystem",
	}

	parentCmd.AddCommand(authCmd)

	err = createAuthCmdChildren(factory, authCmd, rootCmdValues)

	return authCmd, err
}

func createAuthCmdChildren(factory Factory, authCmd *cobra.Command, rootCmdValues *RootCmdValues) error {
	_, err := createAuthLoginCmd(factory, authCmd, rootCmdValues)
	if err == nil {
		_, err = createAuthLogoutCmd(factory, authCmd, rootCmdValues)
	}
	return err
}