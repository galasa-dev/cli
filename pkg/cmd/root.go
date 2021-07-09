//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "galasactl",
    Short: "CLI for Galasa",
    Long:  "",
    Version: "0.17.0-alpha",
}

var (
    bootstrap string
)

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().StringVarP(&bootstrap, "bootstrap", "b", "", "Bootstrap URL")
    rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
