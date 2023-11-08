/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

type RunsCmdValues struct {
	bootstrap string
}

func createRunsCmd(parentCmd *cobra.Command, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	var err error = nil

	runsCmdValues := &RunsCmdValues{}

	runsCmd := &cobra.Command{
		Use:   "runs",
		Short: "Manage test runs in the ecosystem",
		Long:  "Assembles, submits and monitors test runs in Galasa Ecosystem",
	}

	runsCmd.PersistentFlags().StringVarP(&runsCmdValues.bootstrap, "bootstrap", "b", "",
		"Bootstrap URL. Should start with 'http://' or 'file://'. "+
			"If it starts with neither, it is assumed to be a fully-qualified path. "+
			"If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. "+
			"Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties")

	parentCmd.AddCommand(runsCmd)

	err = createRunsCmdChildren(runsCmd, runsCmdValues, rootCmdValues)

	return runsCmd, err
}

func createRunsCmdChildren(runsCmd *cobra.Command, runsCmdValues *RunsCmdValues, rootCmdValues *RootCmdValues) error {

	_, err := createRunsDownloadCmd(runsCmd, runsCmdValues, rootCmdValues)
	if err == nil {
		_, err = createRunsGetCmd(runsCmd, runsCmdValues, rootCmdValues)
	}
	if err == nil {
		_, err = createRunsPrepareCmd(runsCmd, runsCmdValues, rootCmdValues)
	}
	if err == nil {
		_, err = createRunsSubmitCmd(runsCmd, runsCmdValues, rootCmdValues)
	}
	return err
}
