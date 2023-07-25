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
	"gopkg.in/yaml.v3"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
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
	authBootstrap string
)

func init() {
	authLoginCmd.PersistentFlags().StringVarP(&authBootstrap, "bootstrap", "b", "", "Bootstrap URL")
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

	apiServerUrl := bootstrapData.ApiServerURL
	log.Printf("The API server is at '%s'\n", apiServerUrl)

	// Call to process the command in a unit-testable way.
	err = Login(
		apiServerUrl,
		fileSystem,
		galasaHome,
	)

	if err != nil {
		panic(err)
	}
}

func Login(apiServerUrl string, fileSystem files.FileSystem, galasaHome utils.GalasaHome) error {

	galasactlYamlFilePath := galasaHome.GetNativeFolderPath() + "/galasactl.yaml"
	galasactlYamlFile, err := fileSystem.ReadTextFile(galasactlYamlFilePath)
	if err == nil {
		var auth Auth
		err = yaml.Unmarshal([]byte(galasactlYamlFile), &auth)
		if err == nil {
			// To do: Pass these to the API to generate a JWT using the /token endpoint
			// clientId := auth.ClientId
			// secret := auth.Secret
			// accessToken := auth.AccessToken

			jwt := "jwt"

			bearerTokenFilePath := galasaHome.GetNativeFolderPath() + "/bearer-token.json"
			fileSystem.Create(bearerTokenFilePath)
			fileSystem.WriteTextFile(bearerTokenFilePath, jwt)
		} else {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_UNMARSHAL_GALASACTL_YAML_FILE)
		}
	} else {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_READ_GALASACTL_YAML_FILE)
	}

	return err
}

type Auth struct {
	ClientId    string
	Secret      string
	AccessToken string
}
