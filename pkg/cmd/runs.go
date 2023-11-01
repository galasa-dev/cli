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
	bootstrap string
)

func init() {

}

func createRunsCmd(parentCmd *cobra.Command) (*cobra.Command, error) {
	var err error = nil

	runsCmd := &cobra.Command{
		Use:   "runs",
		Short: "Manage test runs in the ecosystem",
		Long:  "Assembles, submits and monitors test runs in Galasa Ecosystem",
	}

	cmd := runsCmd

	cmd.PersistentFlags().StringVarP(&bootstrap, "bootstrap", "b", "",
		"Bootstrap URL. Should start with 'http://' or 'file://'. "+
			"If it starts with neither, it is assumed to be a fully-qualified path. "+
			"If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. "+
			"Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties")

	parentCmd.AddCommand(runsCmd)

	err = createRunsCmdChildren(runsCmd)

	return runsCmd, err
}

func createRunsCmdChildren(runsCmd *cobra.Command) error {

	_, err := createRunsDownloadCmd(runsCmd)
	if err == nil {
		_, err = createRunsGetCmd(runsCmd)
	}
	if err == nil {
		_, err = createRunsPrepareCmd(runsCmd)
	}
	if err == nil {
		_, err = createRunsSubmitCmd(runsCmd)
	}
	return err
}
