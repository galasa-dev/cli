/*
 * Copyright contributors to the Galasa project
 */
 package cmd

 import (
	 "github.com/spf13/cobra"
 )
 
 var (
	 authCmd = &cobra.Command{
		 Use:   "auth",
		 Short: "Manages authentication of users to the ecosystem",
		 Long:  "Manages authentication of users to the ecosystem",
	 }
	 bootstrapFile string
 )
 
 func init() {
	 cmd := runsCmd
	 parentCmd := RootCmd
 
	 cmd.PersistentFlags().StringVarP(&bootstrap, "bootstrap", "b", "", "Bootstrap URL")
 
	 parentCmd.AddCommand(runsCmd)
 }
 