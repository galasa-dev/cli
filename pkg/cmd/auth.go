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
	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Manages the authentication of users with a Galasa ecosystem",
		Long:  "Manages the authentication of users with a Galasa ecosystem",
	}
)

func init() {
	RootCmd.AddCommand(authCmd)
}
