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
	rootCmd = &cobra.Command{
		Use:     "galasactl",
		Short:   "CLI for Galasa",
		Long:    "",
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&logFileName, "log", "l", "", "File to which log information will be sent")
	rootCmd.PersistentFlags().StringVarP(&bootstrap, "bootstrap", "b", "", "Bootstrap URL")
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
