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

var (
	authLoginCmd = &cobra.Command{
		Use:   "login",
		Short: "Authenticate against a Galasa ecosystem",
		Long:  "Log in to a Galasa ecosystem using an existing access token",
		Args:  cobra.NoArgs,
		Run:   executeAuthLogin,
	}

	// Variables set by cobra's command-line parsing.
	authBootstrap  string
	token          string
)

func init() {
	authLoginCmd.PersistentFlags().StringVarP(&authBootstrap, "bootstrap", "b", "", "Bootstrap URL")
	authLoginCmd.PersistentFlags().StringVar(&token, "token", "", "The authentication token")
	authCmd.AddCommand(authLoginCmd)
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

	log.Println("Galasa CLI - Log in to an ecosystem")

	// Get the ability to query environment variables.
	env := utils.NewEnvironment()

	galasaHome, err := utils.NewGalasaHome(fileSystem, env, CmdParamGalasaHomePath)
	if err != nil {
		panic(err)
	}

	// Read the bootstrap properties.
	var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
	var bootstrapData *api.BootstrapData
	bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, authBootstrap, urlService)
	if err != nil {
		panic(err)
	}

	console := utils.NewRealConsole()

	apiServerUrl := bootstrapData.ApiServerURL
	log.Printf("The API server is at '%s'\n", apiServerUrl)

	// Call to process the command in a unit-testable way.
	err = Login(
		token,
		apiServerUrl,
		fileSystem,
		console,
	)

	if err != nil {
		panic(err)
	}
}

func Login(token string, apiServerUrl string, fileSystem files.FileSystem, console utils.Console) error {
	return nil
}