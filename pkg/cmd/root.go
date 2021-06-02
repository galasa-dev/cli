//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package cmd

import (
	"github.com/galasa.dev/cli/pkg/cli"

	"github.com/spf13/cobra"
)

func Root(p cli.Params) *cobra.Command {

	var cmd = &cobra.Command{
		Use:          "galasactl",
		Short:        "CLI for Galasa",
		Long:         "",
		SilenceUsage: true,
	}

	return cmd
}
