/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

func createProjectCmd(parentCmd *cobra.Command) (*cobra.Command, error) {

	var err error = nil

	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Manipulate local project source code",
		Long:  "Creates and manipulates Galasa test project source code",
	}

	parentCmd.AddCommand(projectCmd)

	err = createProjectCmdChildren(projectCmd)

	return projectCmd, err
}

func createProjectCmdChildren(projectCmd *cobra.Command) error {
	_, err := createProjectCreateCmd(projectCmd)
	return err
}
