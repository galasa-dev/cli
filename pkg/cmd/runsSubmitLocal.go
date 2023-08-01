/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/galasa.dev/cli/pkg/api"
	"github.com/galasa.dev/cli/pkg/embedded"
	"github.com/galasa.dev/cli/pkg/files"
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

	currentGalasaVersion, _ := embedded.GetGalasaVersion()
	runsSubmitLocalCmd.Flags().StringVar(&runsSubmitLocalCmdParams.TargetGalasaVersion, "galasaVersion",
		currentGalasaVersion,
		"the version of galasa you want to use to run your tests. "+
			"This should match the version of the galasa obr you built your test bundles against.")

	runsSubmitLocalCmd.Flags().StringSliceVar(&runsSubmitLocalCmdParams.Obrs, "obr", make([]string, 0),
		"The maven coordinates of the obr bundle(s) which refer to your test bundles. "+
			"The format of this parameter is 'mvn:${TEST_OBR_GROUP_ID}/${TEST_OBR_ARTIFACT_ID}/${TEST_OBR_VERSION}/obr' "+
			"Multiple instances of this flag can be used to describe multiple obr bundles.")

	runsSubmitLocalCmd.Flags().Uint32Var(&runsSubmitLocalCmdParams.DebugPort, "debugPort", 0,
		"The port to use when the --debug option causes the testcase to connect to a java debugger. "+
			"The default value used is "+strconv.FormatUint(uint64(launcher.DEBUG_PORT_DEFAULT), 10)+" which can be "+
			"overridden by the '"+api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_PORT+"' property in the bootstrap file, "+
			"which in turn can be overridden by this explicit parameter on the galasactl command.",
	)

	runsSubmitLocalCmd.Flags().StringVar(&runsSubmitLocalCmdParams.DebugMode, "debugMode", "",
		"The mode to use when the --debug option causes the testcase to connect to a Java debugger. "+
			"Valid values are 'listen' or 'attach'. "+
			"'listen' means the testcase JVM will pause on startup, waiting for the Java debugger to connect to the debug port "+
			"(see the --debugPort option). "+
			"'attach' means the testcase JVM will pause on startup, trying to attach to a java debugger which is listening on the debug port. "+
			"The default value is 'listen' but can be overridden by the '"+api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_MODE+"' property in the bootstrap file, "+
			"which in turn can be overridden by this explicit parameter on the galasactl command.",
	)

	runsSubmitLocalCmd.Flags().BoolVar(&runsSubmitLocalCmdParams.IsDebugEnabled, "debug", false,
		"When set (or true) the debugger pauses on startup and tries to connect to a Java debugger. "+
			"The connection is established using the --debugMode and --debugPort values.",
	)

	runs.AddCommandFlags(runsSubmitLocalCmd, &submitLocalSelectionFlags)

	runsSubmitCmd.AddCommand(runsSubmitLocalCmd)
}

func executeSubmitLocal(cmd *cobra.Command, args []string) {

	var err error = nil

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, logFileName)
	if err != nil {
		panic(err)
	}
	isCapturingLogs = true

	log.Println("Galasa CLI - Submit tests (Local)")

	// Get the ability to query environment variables.
	env := utils.NewEnvironment()

	// Work out where galasa home is, only once.
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

	timeService := utils.NewRealTimeService()

	// the submit is targetting a local JVM
	embeddedFileSystem := embedded.GetReadOnlyFileSystem()

	// Something which can kick off new operating system processes
	processFactory := launcher.NewRealProcessFactory()

	var launcherInstance launcher.Launcher
	launcherInstance, err = launcher.NewJVMLauncher(
		bootstrapData.Properties, env, fileSystem, embeddedFileSystem,
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
