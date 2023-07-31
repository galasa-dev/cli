/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"context"
	"log"
	"path/filepath"

	"github.com/galasa.dev/cli/pkg/api"
	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/galasaapi"
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

type AuthYaml struct {
	Auth AuthPropertiesYaml `yaml:"auth,omitempty"`
}

type AuthPropertiesYaml struct {
	ClientId    string `yaml:"client_id,omitempty"`
	Secret      string `yaml:"secret,omitempty"`
	AccessToken string `yaml:"access_token,omitempty"`
}

func Login(apiServerUrl string, fileSystem files.FileSystem, galasaHome utils.GalasaHome) error {

	var err error = nil
	var authProperties galasaapi.AuthProperties
	authProperties, err = getAuthProperties(fileSystem, galasaHome)
	if err == nil {
		var jwt string
		jwt, err = GetJwtFromRestApi(apiServerUrl, authProperties)
		if err == nil {
			err = writeBearerTokenJsonFile(fileSystem, galasaHome, jwt)
		}
	}

	return err
}

// Gets authentication properties from the user's galasactl.yaml file
func getAuthProperties(fileSystem files.FileSystem, galasaHome utils.GalasaHome) (galasaapi.AuthProperties, error) {
	var err error = nil
	var authParent AuthYaml

	galasactlYamlFilePath := filepath.Join(galasaHome.GetNativeFolderPath(), "galasactl.yaml")
	galasactlYamlFile, err := fileSystem.ReadTextFile(galasactlYamlFilePath)
	if err == nil {
		err = yaml.Unmarshal([]byte(galasactlYamlFile), &authParent)
		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_UNMARSHAL_GALASACTL_YAML_FILE)
		}
	} else {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_READ_GALASACTL_YAML_FILE)
	}

	// Convert the YAML representations of the auth properties into the OpenAPI-generated "AuthProperties" type
	authProperties := galasaapi.NewAuthProperties()
	authProperties.SetClientId(authParent.Auth.ClientId)
	authProperties.SetSecret(authParent.Auth.Secret)
	authProperties.SetAccessToken(authParent.Auth.AccessToken)

	return *authProperties, err
}

func GetJwtFromRestApi(apiServerUrl string, authProperties galasaapi.AuthProperties) (string, error) {
	var err error = nil
	var context context.Context = nil
	var jwtJsonStr string

	restClient := api.InitialiseAPI(apiServerUrl)

	tokenResponse, httpResponse, err := restClient.AuthenticationAPIApi.PostAuthenticate(context).
		AuthProperties(authProperties).
		Execute()
	defer httpResponse.Body.Close()

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_BEARER_TOKEN, err.Error())
		log.Printf("Failed to retrieve JWT from API server. %s", err.Error())
	} else {
		var tokenResponseJson []byte
		tokenResponseJson, err = tokenResponse.MarshalJSON()
		jwtJsonStr = string(tokenResponseJson)
		log.Println("JWT received from API server OK")

	}
	return jwtJsonStr, err
}

func writeBearerTokenJsonFile(fileSystem files.FileSystem, galasaHome utils.GalasaHome, jwt string) error {
	bearerTokenFilePath := filepath.Join(galasaHome.GetNativeFolderPath(), "bearer-token.json")

	log.Printf("Writing JWT to bearer token file '%s'", bearerTokenFilePath)
	err := fileSystem.WriteTextFile(bearerTokenFilePath, jwt)

	if err == nil {
		log.Printf("Written JWT to bearer token file '%s' OK", bearerTokenFilePath)
	} else {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_WRITE_FILE, bearerTokenFilePath, err.Error())
		log.Printf("Failed to write bearer token file '%s'. %s", bearerTokenFilePath, err.Error())
	}
	return err
}