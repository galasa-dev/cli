/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/galasa-dev/cli/pkg/embedded"
	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

type LocalInitCmdValues struct {
	isDevelopmentLocalInit bool
}

func createLocalInitCmd(parentCmd *cobra.Command, rootCmdValues *RootCmdValues) (*cobra.Command, error) {
	var err error = nil

	localInitCmdValues := &LocalInitCmdValues{}

	localInitCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialises Galasa home folder",
		Long:  "Initialises Galasa home folder in home directory with all the properties files",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			executeEnvInit(cmd, args, localInitCmdValues, rootCmdValues)
		},
	}

	localInitCmd.Flags().BoolVar(&localInitCmdValues.isDevelopmentLocalInit, "development", false, "Use bleeding-edge galasa versions and repositories.")

	parentCmd.AddCommand(localInitCmd)

	// There are no children commands to add to the command tree from here.

	return localInitCmd, err
}

func executeEnvInit(cmd *cobra.Command, args []string, localInitCmdValues *LocalInitCmdValues, rootCmdValues *RootCmdValues) {

	var err error = nil

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err != nil {
		panic(err)
	}
	rootCmdValues.isCapturingLogs = true

	env := utils.NewEnvironment()

	err = localEnvInit(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath, localInitCmdValues.isDevelopmentLocalInit)
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
