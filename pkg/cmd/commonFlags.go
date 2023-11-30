/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

// ------------------------------------------------------------------------------------------------
// Objectives
//   Functions which add a flag to a cobra command in a different way,
//   depending on the command it is being added to.
// ------------------------------------------------------------------------------------------------

func addBootstrapFlag(cobraCommand *cobra.Command, parsedValueLocation *string) {
	cobraCommand.PersistentFlags().StringVarP(parsedValueLocation, "bootstrap", "b", "",
		"Bootstrap URL. Should start with 'http://' or 'file://'. "+
			"If it starts with neither, it is assumed to be a fully-qualified path. "+
			"If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. "+
			"Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties")
}
