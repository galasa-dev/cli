/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/galasa.dev/cli/pkg/embedded"
	"github.com/galasa.dev/cli/pkg/files"
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

	var err error = nil

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, logFileName)
	if err != nil {
		panic(err)
	}
	isCapturingLogs = true

	env := utils.NewEnvironment()

	err = localEnvInit(fileSystem, env, CmdParamGalasaHomePath, isDevelopmentLocalInit)
	if err != nil {
		panic(err)
	}
}

func localEnvInit(
	fileSystem files.FileSystem,
	env utils.Environment,
	cmdFlagGalasaHome string,
	isDevelopment bool,
) error {

	galasaHome, err := utils.NewGalasaHome(fileSystem, env, cmdFlagGalasaHome)
	if err == nil {
		embeddedFileSystem := embedded.GetReadOnlyFileSystem()
		err = utils.InitialiseGalasaHomeFolder(galasaHome, fileSystem, embeddedFileSystem)
		if err == nil {
			err = utils.InitialiseM2Folder(fileSystem, embeddedFileSystem, isDevelopment)
		}
	}
	return err
}
