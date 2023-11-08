/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

func createProjectCmd(factory Factory, parentCmd *cobra.Command, rootCmdValues *RootCmdValues) (*cobra.Command, error) {

	var err error = nil

	projectCmd := &cobra.Command{
		Use:   "project",
		Short: "Manipulate local project source code",
		Long:  "Creates and manipulates Galasa test project source code",
	}

	parentCmd.AddCommand(projectCmd)

	err = createProjectCmdChildren(factory, projectCmd, rootCmdValues)

	return projectCmd, err
}

func createProjectCmdChildren(factory Factory, projectCmd *cobra.Command, rootCmdValues *RootCmdValues) error {
	_, err := createProjectCreateCmd(factory, projectCmd, rootCmdValues)
	return err
}
