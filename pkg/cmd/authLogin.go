/*
 * Copyright contributors to the Galasa project
 */
 package cmd

 import (
	 "log"
 
	 "github.com/galasa.dev/cli/pkg/api"
	 "github.com/galasa.dev/cli/pkg/files"
	 "github.com/galasa.dev/cli/pkg/utils"
	 "github.com/spf13/cobra"
 )
 
 // Objective: Allow the user to do this:
 //    run get --runname 12345
 // And then show the results in a human-readable form.
 
 var (
	 authLoginCmd = &cobra.Command{
		 Use:   "login",
		 Short: "Authenticate against the galasa ecosystem.",
		 Long:  "Authenticate against the galasa ecosystem.",
		 Args:  cobra.NoArgs,
		 Run:   executeAuthLogin,
	 }
 
	 // Variables set by cobra's command-line parsing.
	 host		            string
	 token           		string

 )
 
 func init() {
	authLoginCmd.PersistentFlags().StringVar(&host, "host", "", "The host name of the galasa ecosystem you want to authenticate against."+
	 					 "If this is not specified then the value in the bootstrap identified in $GALASA_BOOTSTRAP will be used")
	authLoginCmd.PersistentFlags().StringVar(&token, "token", "", "The authentication token ")
	parentCommand := authCmd
	parentCommand.AddCommand(authLoginCmd)
 }
 
 func executeAuthLogin(cmd *cobra.Command, args []string) {
 
	 var err error
 
	 // Operations on the file system will all be relative to the current folder.
	 fileSystem := files.NewOSFileSystem()
 
	 err = utils.CaptureLog(fileSystem, logFileName)
	 if err != nil {
		 panic(err)
	 }
	 isCapturingLogs = true
 
	 log.Println("Galasa CLI - Get info about a run")
 
	 // Get the ability to query environment variables.
	 env := utils.NewEnvironment()
 
	 galasaHome, err := utils.NewGalasaHome(fileSystem, env, CmdParamGalasaHomePath)
	 if err != nil {
		 panic(err)
	 }
 
	 // Read the bootstrap properties.
	 var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
	 var bootstrapData *api.BootstrapData
	 bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, bootstrap, urlService)
	 if err != nil {
		 panic(err)
	 }
 
	 var console = utils.NewRealConsole()
 
	 apiServerUrl := bootstrapData.ApiServerURL
	 log.Printf("The API sever is at '%s'\n", apiServerUrl)
 
	 timeService := utils.NewRealTimeService()
 
	 // Call to process the command in a unit-testable way.
	 err = null
	 if err != nil {
		 panic(err)
	 }
 }
 