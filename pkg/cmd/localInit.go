/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"github.com/galasa.dev/cli/pkg/embedded"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	localInitCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialises Galasa home folder",
		Long:  "Initialises Galasa home folder in home directory with all the properties files",
		Args:  cobra.NoArgs,
		Run:   executeEnvInit,
	}

	isDevelopmentLocalInit bool
)

func init() {
	localInitCmd.Flags().BoolVar(&isDevelopmentLocalInit, "development", false, "Use bleeding-edge galasa versions and repositories.")

	parentCommand := localCmd
	parentCommand.AddCommand(localInitCmd)
}

func executeEnvInit(cmd *cobra.Command, args []string) {
	utils.CaptureLog(logFileName)
	isCapturingLogs = true

	fileSystem := utils.NewOSFileSystem()
	env := utils.NewEnvironment()

	err := localEnvInit(fileSystem, env, CmdParamGalasaHomePath, isDevelopmentLocalInit)
	if err != nil {
		panic(err)
	}
}

func localEnvInit(
	fileSystem utils.FileSystem,
	env utils.Environment,
	cmdFlagGalasaHome string,
	isDevelopment bool,
) error {

	galasaHome, err := utils.NewGalasaHome(fileSystem, env, cmdFlagGalasaHome)
	if err == nil {
		embeddedFileSystem := embedded.GetEmbeddedFileSystem()
		err = utils.InitialiseGalasaHomeFolder(galasaHome, fileSystem, embeddedFileSystem)
		if err == nil {
			err = utils.InitialiseM2Folder(fileSystem, embeddedFileSystem, isDevelopment)
		}
	}
	return err
}
