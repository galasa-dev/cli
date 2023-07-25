/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"fmt"
	"log"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/spf13/cobra"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
)

var (
	authLogoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "Log out from a Galasa ecosystem",
		Long:  "Log out from a Galasa ecosystem",
		Args:  cobra.NoArgs,
		Run:   executeAuthLogout,
	}
)

func init() {
	authCmd.AddCommand(authLogoutCmd)
}

func executeAuthLogout(cmd *cobra.Command, args []string) {

	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, logFileName)
	if err != nil {
		panic(err)
	}
	isCapturingLogs = true

	log.Println("Galasa CLI - Log out of an ecosystem")

	// Get the ability to query environment variables.
	env := utils.NewEnvironment()

	galasaHome, err := utils.NewGalasaHome(fileSystem, env, CmdParamGalasaHomePath)
	if err != nil {
		panic(err)
	}

	console := utils.NewRealConsole()

	// Call to process the command in a unit-testable way.
	err = Logout(
		fileSystem,
		console,
		env,
		galasaHome,
	)

	if err != nil {
		panic(err)
	}
}

func Logout(fileSystem files.FileSystem, console utils.Console, env utils.Environment, galasaHome utils.GalasaHome) error {

	var err error
	bearerTokenFile := fmt.Sprintf("%s/%s", galasaHome.GetNativeFolderPath(), "bearer-token.json")
	if _, err := fileSystem.Exists(bearerTokenFile); err == nil {
		fileSystem.DeleteFile(bearerTokenFile)
	}

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNABLE_TO_DELETE_BEARER_TOKEN_FILE)
	}

	return err
}