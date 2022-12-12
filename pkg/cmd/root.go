/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:     "galasactl",
		Short:   "CLI for Galasa",
		Long:    `A tool for controlling Galasa resources using the command-line.`,
		Version: "unknowncliversion-unknowngithash",
	}

	bootstrap   string
	logFileName string
)

func Execute() {

	// Catch execution if a panic happens.
	defer func() {
		errobj := recover()
		if errobj != nil {
			fmt.Fprintln(os.Stderr, errobj)
			log.Println(errobj)
			os.Exit(1)
		}
	}()

	// Execute the command
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&logFileName, "log", "l", "", "File to which log information will be sent")
	RootCmd.PersistentFlags().StringVarP(&bootstrap, "bootstrap", "b", "", "Bootstrap URL")
	RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
