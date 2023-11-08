/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

func createProjectCmd(parentCmd *cobra.Command, rootCmdValues *RootCmdValues) (*cobra.Command, error) {

	var err error = nil

	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Manipulate local project source code",
		Long:  "Creates and manipulates Galasa test project source code",
	}

	parentCmd.AddCommand(projectCmd)

	err = createProjectCmdChildren(projectCmd, rootCmdValues)

	return projectCmd, err
}

func createProjectCmdChildren(projectCmd *cobra.Command, rootCmdValues *RootCmdValues) error {
	_, err := createProjectCreateCmd(projectCmd, rootCmdValues)
	return err
}
