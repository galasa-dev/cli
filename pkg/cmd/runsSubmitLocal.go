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

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/embedded"
	"github.com/galasa-dev/cli/pkg/launcher"
	"github.com/galasa-dev/cli/pkg/runs"
	"github.com/galasa-dev/cli/pkg/utils"
)

var ()

// Variables set by cobra's command-line parsing.
type RunsSubmitLocalCmdValues struct {
	runsSubmitLocalCmdParams  *launcher.RunsSubmitLocalCmdParameters
	submitLocalSelectionFlags *utils.TestSelectionFlagValues
}

func createRunsSubmitLocalCmd(
	factory Factory,
	parentCmd *cobra.Command,
	runsSubmitCmdValues *utils.RunsSubmitCmdValues,
	runsCmdValues *RunsCmdValues,
	rootCmdValues *RootCmdValues,
) (*cobra.Command, error) {
	var err error = nil

	// Allocate storage to capture the parsed values.
	runsSubmitLocalCmdValues := &RunsSubmitLocalCmdValues{
		runsSubmitLocalCmdParams:  &launcher.RunsSubmitLocalCmdParameters{},
		submitLocalSelectionFlags: runs.NewTestSelectionFlagValues(),
	}

	runsSubmitLocalCmd := &cobra.Command{
		Use:     "local",
		Short:   "submit a list of tests to be run on a local java virtual machine (JVM)",
		Long:    "Submit a list of tests to a local JVM, monitor them and wait for them to complete",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs submit local"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeSubmitLocal(factory, cmd, args, runsSubmitLocalCmdValues, runsSubmitCmdValues, runsCmdValues, rootCmdValues)
		},
	}

	//currentUserName := runs.GetCurrentUserName()

	runsSubmitLocalCmd.Flags().StringVar(&runsSubmitLocalCmdValues.runsSubmitLocalCmdParams.RemoteMaven, "remoteMaven",
		"https://repo.maven.apache.org/maven2",
		"the url of the remote maven where galasa bundles can be loaded from. "+
			"Defaults to maven central.")

	runsSubmitLocalCmd.Flags().StringVar(&runsSubmitLocalCmdValues.runsSubmitLocalCmdParams.LocalMaven, "localMaven", "",
		"The url of a local maven repository are where galasa bundles can be loaded from on your local file system. Defaults to your home .m2/repository file. Please note that this should be in a URL form e.g. 'file:///Users/myuserid/.m2/repository', or 'file://C:/Users/myuserid/.m2/repository'")

	currentGalasaVersion, _ := embedded.GetGalasaVersion()
	runsSubmitLocalCmd.Flags().StringVar(&runsSubmitLocalCmdValues.runsSubmitLocalCmdParams.TargetGalasaVersion, "galasaVersion",
		currentGalasaVersion,
		"the version of galasa you want to use to run your tests. "+
			"This should match the version of the galasa obr you built your test bundles against.")

	runsSubmitLocalCmd.Flags().StringSliceVar(&runsSubmitLocalCmdValues.runsSubmitLocalCmdParams.Obrs, "obr", make([]string, 0),
		"The maven coordinates of the obr bundle(s) which refer to your test bundles. "+
			"The format of this parameter is 'mvn:${TEST_OBR_GROUP_ID}/${TEST_OBR_ARTIFACT_ID}/${TEST_OBR_VERSION}/obr' "+
			"Multiple instances of this flag can be used to describe multiple obr bundles.")
	runsSubmitLocalCmd.MarkFlagRequired("obr")

	runsSubmitLocalCmd.Flags().Uint32Var(&runsSubmitLocalCmdValues.runsSubmitLocalCmdParams.DebugPort, "debugPort", 0,
		"The port to use when the --debug option causes the testcase to connect to a java debugger. "+
			"The default value used is "+strconv.FormatUint(uint64(launcher.DEBUG_PORT_DEFAULT), 10)+" which can be "+
			"overridden by the '"+api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_PORT+"' property in the bootstrap file, "+
			"which in turn can be overridden by this explicit parameter on the galasactl command.",
	)

	runsSubmitLocalCmd.Flags().StringVar(&runsSubmitLocalCmdValues.runsSubmitLocalCmdParams.DebugMode, "debugMode", "",
		"The mode to use when the --debug option causes the testcase to connect to a Java debugger. "+
			"Valid values are 'listen' or 'attach'. "+
			"'listen' means the testcase JVM will pause on startup, waiting for the Java debugger to connect to the debug port "+
			"(see the --debugPort option). "+
			"'attach' means the testcase JVM will pause on startup, trying to attach to a java debugger which is listening on the debug port. "+
			"The default value is 'listen' but can be overridden by the '"+api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_MODE+"' property in the bootstrap file, "+
			"which in turn can be overridden by this explicit parameter on the galasactl command.",
	)

	runsSubmitLocalCmd.Flags().BoolVar(&runsSubmitLocalCmdValues.runsSubmitLocalCmdParams.IsDebugEnabled, "debug", false,
		"When set (or true) the debugger pauses on startup and tries to connect to a Java debugger. "+
			"The connection is established using the --debugMode and --debugPort values.",
	)

	runs.AddClassFlag(runsSubmitLocalCmd, runsSubmitLocalCmdValues.submitLocalSelectionFlags, true, "test class names."+
		" The format of each entry is osgi-bundle-name/java-class-name. Java class names are fully qualified. No .class suffix is needed.")

	parentCmd.AddCommand(runsSubmitLocalCmd)

	// There are no children of this commands.

	return runsSubmitLocalCmd, err
}

func executeSubmitLocal(
	factory Factory,
	cmd *cobra.Command,
	args []string,
	runsSubmitLocalCmdValues *RunsSubmitLocalCmdValues,
	runsSubmitCmdValues *utils.RunsSubmitCmdValues,
	runsCmdValues *RunsCmdValues,
	rootCmdValues *RootCmdValues,
) error {

	var err error = nil

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()
	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)

	if err == nil {
		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Submit tests (Local)")

		// Get the ability to query environment variables.
		env := factory.GetEnvironment()

		// Work out where galasa home is, only once.
		var galasaHome utils.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Read the bootstrap properties.
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, runsCmdValues.bootstrap, urlService)
			if err == nil {

				timeService := utils.NewRealTimeService()

				// the submit is targetting a local JVM
				embeddedFileSystem := embedded.GetReadOnlyFileSystem()

				// Something which can kick off new operating system processes
				processFactory := launcher.NewRealProcessFactory()

				// Validate the test selection parameters.
				validator := runs.NewObrBasedValidator()
				err = validator.Validate(runsSubmitLocalCmdValues.submitLocalSelectionFlags)
				if err == nil {

					// A launcher is needed to launch anythihng
					var launcherInstance launcher.Launcher
					launcherInstance, err = launcher.NewJVMLauncher(
						bootstrapData.Properties, env, fileSystem, embeddedFileSystem,
						runsSubmitLocalCmdValues.runsSubmitLocalCmdParams, timeService,
						processFactory, galasaHome)

					if err == nil {
						var console = factory.GetConsole()

						// Do the launching of the tests.
						submitter := runs.NewSubmitter(
							galasaHome,
							fileSystem,
							launcherInstance,
							timeService,
							env,
							console,
						)

						err = submitter.ExecuteSubmitRuns(
							runsSubmitCmdValues,
							runsSubmitLocalCmdValues.submitLocalSelectionFlags,
						)
					}
				}
			}
		}
	}

	return err
}
