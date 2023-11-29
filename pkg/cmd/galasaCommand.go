/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/spf13/cobra"
)

type GalasaCommand interface {
	GetName() string
	GetCobraCommand() *cobra.Command
}
