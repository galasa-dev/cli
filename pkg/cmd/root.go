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
		log.Printf("Error : %s. exit code 1.", err.Error())
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	log.Printf("OK. Exit code 0")
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&logFileName, "log", "l", "",
		"File to which log information will be sent. Any folder referred to must exist. "+
			"An existing file will be over-written. "+
			"Specify \"-\" to log to stderr. "+
			"Defaults to not logging.")
	RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
