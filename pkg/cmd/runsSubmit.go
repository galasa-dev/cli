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
	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/launcher"
	"github.com/galasa.dev/cli/pkg/runs"
	"github.com/galasa.dev/cli/pkg/utils"
)

var (
	runsSubmitCmd = &cobra.Command{
		Use:   "submit",
		Short: "submit a list of tests to the ecosystem",
		Long:  "Submit a list of tests to the ecosystem, monitor them and wait for them to complete",
		Args:  cobra.NoArgs,
		Run:   executeSubmit,
	}

	// Variables set by cobra's command-line parsing.
	runsSubmitCmdParams utils.RunsSubmitCmdParameters

	submitSelectionFlags = runs.TestSelectionFlags{}
)

func init() {

	runsSubmitCmd.PersistentFlags().StringVarP(&runsSubmitCmdParams.PortfolioFileName, "portfolio", "p", "", "portfolio containing the tests to run")
	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdParams.ReportYamlFilename, "reportyaml", "", "yaml file to record the final results in")
	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdParams.ReportJsonFilename, "reportjson", "", "json file to record the final results in")
	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdParams.ReportJunitFilename, "reportjunit", "", "junit xml file to record the final results in")
	runsSubmitCmd.PersistentFlags().StringVarP(&runsSubmitCmdParams.GroupName, "group", "g", "", "the group name to assign the test runs to, if not provided, a psuedo unique id will be generated")
	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdParams.Requestor, "requestor", "", "the requestor id to be associated with the test runs. Defaults to the current user id.")
	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdParams.RequestType, "requesttype", "CLI", "the type of request, used to allocate a run name. Defaults to CLI.")

	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdParams.ThrottleFileName, "throttlefile", "",
		"a file where the current throttle is stored. Periodically the throttle value is read from the file used. "+
			"Someone with edit access to the file can change it which dynamically takes effect. "+
			"Long-running large portfolios can be throttled back to nothing (paused) using this mechanism (if throttle is set to 0). "+
			"And they can be resumed (un-paused) if the value is set back. "+
			"This facility can allow the tests to not show a failure when the system under test is taken out of service for maintainence."+
			"Optional. If not specified, no throttle file is used.",
	)

	runsSubmitCmd.PersistentFlags().IntVar(&runsSubmitCmdParams.PollIntervalSeconds, "poll", runs.DEFAULT_POLL_INTERVAL_SECONDS,
		"Optional. The interval time in seconds between successive polls of the test runs status. "+
			"Defaults to "+strconv.Itoa(runs.DEFAULT_POLL_INTERVAL_SECONDS)+" seconds. "+
			"If less than 1, then default value is used.")

	runsSubmitCmd.PersistentFlags().IntVar(&runsSubmitCmdParams.ProgressReportIntervalMinutes, "progress", runs.DEFAULT_PROGRESS_REPORT_INTERVAL_MINUTES,
		"in minutes, how often the cli will report the overall progress of the test runs. A value of 0 or less disables progress reporting.")

	runsSubmitCmd.PersistentFlags().IntVar(&runsSubmitCmdParams.Throttle, "throttle", runs.DEFAULT_THROTTLE_TESTS_AT_ONCE,
		"how many test runs can be submitted in parallel, 0 or less will disable throttling. 1 causes tests to be run sequentially.")

	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdParams.OverrideFilePath, "overridefile", "",
		"path to a properties file containing override properties. Defaults to overrides.properties in galasa home folder if that file exists. "+
			"Overrides from --override options will take precedence over properties in this property file. "+
			"A file path of '-' disables reading any properties file.")

	runsSubmitCmd.PersistentFlags().StringSliceVar(&runsSubmitCmdParams.Overrides, "override", make([]string, 0),
		"overrides to be sent with the tests (overrides in the portfolio will take precedence). "+
			"Each override is of the form 'name=value'. Multiple instances of this flag can be used. "+
			"For example --override=prop1=val1 --override=prop2=val2")

	// The trace flag defaults to 'false' if you don't use it.
	// If you say '--trace' on it's own, it defaults to 'true'
	// If you say --trace=false or --trace=true you can set the value explicitly.
	runsSubmitCmd.PersistentFlags().BoolVar(&runsSubmitCmdParams.Trace, "trace", false, "Trace to be enabled on the test runs")
	runsSubmitCmd.PersistentFlags().Lookup("trace").NoOptDefVal = "true"

	runsSubmitCmd.PersistentFlags().BoolVar(&(runsSubmitCmdParams.NoExitCodeOnTestFailures), "noexitcodeontestfailures", false, "set to true if you don't want an exit code to be returned from galasactl if a test fails")

	runs.AddCommandFlags(runsSubmitCmd, &submitSelectionFlags)

	runsCmd.AddCommand(runsSubmitCmd)
}

func executeSubmit(cmd *cobra.Command, args []string) {

	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := files.NewOSFileSystem()

	err = utils.CaptureLog(fileSystem, logFileName)
	if err != nil {
		panic(err)
	}
	isCapturingLogs = true

	log.Println("Galasa CLI - Submit tests (Remote)")

	// Get the ability to query environment variables.
	env := utils.NewEnvironment()

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
	var launcherInstance launcher.Launcher = nil

	// The launcher we are going to use to start/monitor tests.
	launcherInstance = launcher.NewRemoteLauncher(bootstrapData.ApiServerURL)

	if err == nil {
		err = runs.ExecuteSubmitRuns(galasaHome, fileSystem, runsSubmitCmdParams, launcherInstance, timeService, &submitSelectionFlags)
	}

	if err != nil {
		// Panic. If we could pass an error back we would.
		// The panic is recovered from in the root command, where
		// the error is logged/displayed before program exit.
		panic(err)
	}
}
