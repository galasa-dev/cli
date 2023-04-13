/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/galasa.dev/cli/pkg/embedded"
	"github.com/galasa.dev/cli/pkg/launcher"
	"github.com/galasa.dev/cli/pkg/runs"
	"github.com/galasa.dev/cli/pkg/utils"
)

var (
	runsSubmitLocalCmd = &cobra.Command{
		Use:   "local",
		Short: "submit a list of tests to be run on a local java virtual machine (JVM)",
		Long:  "Submit a list of tests to a local JVM, monitor them and wait for them to complete",
		Args:  cobra.NoArgs,
		Run:   executeSubmitLocal,
	}

	// Variables set by cobra's command-line parsing.
	runsSubmitLocalCmdParams launcher.RunsSubmitLocalCmdParameters

	submitLocalSelectionFlags = runs.TestSelectionFlags{}
)

func init() {

	// currentUserName := runs.GetCurrentUserName()

	runsSubmitLocalCmd.Flags().StringVar(&runsSubmitLocalCmdParams.RemoteMaven, "remoteMaven",
		"https://repo.maven.apache.org/maven2",
		"the url of the remote maven where galasa bundles can be loaded from. "+
			"Defaults to maven central.")

	currentGalasaVersion := embedded.GetGalasaVersion()
	runsSubmitLocalCmd.Flags().StringVar(&runsSubmitLocalCmdParams.TargetGalasaVersion, "galasaVersion",
		currentGalasaVersion,
		"the version of galasa you want to use to run your tests. "+
			"This should match the version of the galasa obr you built your test bundles against.")

	runsSubmitLocalCmd.Flags().StringSliceVar(&runsSubmitLocalCmdParams.Obrs, "obr", make([]string, 0),
		"The maven coordinates of the obr bundle(s) which refer to your test bundles. "+
			"The format of this parameter is 'mvn:${TEST_OBR_GROUP_ID}/${TEST_OBR_ARTIFACT_ID}/${TEST_OBR_VERSION}/obr' "+
			"Multiple instances of this flag can be used to describe multiple obr bundles.")

	runs.AddCommandFlags(runsSubmitLocalCmd, &submitLocalSelectionFlags)

	runsSubmitCmd.AddCommand(runsSubmitLocalCmd)
}

func executeSubmitLocal(cmd *cobra.Command, args []string) {

	var err error

	utils.CaptureLog(logFileName)
	isCapturingLogs = true

	log.Println("Galasa CLI - Submit tests (Local)")

	// Operations on the file system will all be relative to the current folder.
	fileSystem := utils.NewOSFileSystem()

	// Get the ability to query environment variables.
	env := utils.NewEnvironment()

	// Work out where galasa home is, only once.
	galasaHome, err := utils.NewGalasaHome(fileSystem, env, CmdParamGalasaHomePath)
	if err != nil {
		panic(err)
	}

	timeService := utils.NewRealTimeService()

	// the submit is targetting a local JVM
	embeddedFileSystem := embedded.GetEmbeddedFileSystem()

	// Something which can kick off new operating system processes
	processFactory := launcher.NewRealProcessFactory()

	var launcherInstance launcher.Launcher
	launcherInstance, err = launcher.NewJVMLauncher(
		env, fileSystem, embeddedFileSystem,
		runsSubmitLocalCmdParams, timeService,
		processFactory, galasaHome)

	if err == nil {
		err = runs.ExecuteSubmitRuns(
			galasaHome,
			fileSystem,
			runsSubmitCmdParams,
			launcherInstance,
			timeService,
			&submitLocalSelectionFlags,
		)
	}

	if err != nil {
		// Panic. If we could pass an error back we would.
		// The panic is recovered from in the root command, where
		// the error is logged/displayed before program exit.
		panic(err)
	}
}
