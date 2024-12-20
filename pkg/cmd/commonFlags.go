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

func addRateLimitRetryFlags(cobraCommand *cobra.Command, maxRetries *int, retryBackoffSeconds *float64) {
	cobraCommand.PersistentFlags().IntVar(maxRetries, "rate-limit-retries", 3,
		"The maximum number of retries that should be made when requests to the Galasa Service fail due to rate limits being exceeded. Must be a whole number. "+
			"Defaults to 3 retries")

	cobraCommand.PersistentFlags().Float64Var(retryBackoffSeconds, "rate-limit-retry-backoff", float64(1),
	"The amount of time in seconds to wait before retrying a command if it failed due to rate limits being exceeded. Defaults to 1 second.")
}
