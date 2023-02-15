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
	envInitCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialises Galasa home folder",
		Long:  "Initialises Galasa home folder in home directory with all the properties files",
		Args:  cobra.NoArgs,
		Run:   executeEnvInit,
	}
)

func init() {
	parentCommand := localCmd
	parentCommand.AddCommand(envInitCmd)
}

func executeEnvInit(cmd *cobra.Command, args []string) {
	utils.CaptureLog(logFileName)
	isCapturingLogs = true
	
	fileSystem := utils.NewOSFileSystem()
	err := envInit(fileSystem)
	if err != nil {
		panic(err)
	}
}

func envInit(fileSystem utils.FileSystem) error {
	embeddedFileSystem := embedded.GetEmbeddedFileSystem()
	err := utils.InitialiseGalasaHomeFolder(fileSystem, embeddedFileSystem)
	if err == nil {
		err = utils.InitialiseM2Folder(fileSystem, embeddedFileSystem)
	}
	return err
}
