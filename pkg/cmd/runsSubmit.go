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
	"github.com/galasa-dev/cli/pkg/auth"
	"github.com/galasa-dev/cli/pkg/launcher"
	"github.com/galasa-dev/cli/pkg/runs"
	"github.com/galasa-dev/cli/pkg/utils"
)

type RunsSubmitCommand struct {
	values       *utils.RunsSubmitCmdValues
	cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors
// ------------------------------------------------------------------------------------------------
func NewRunsSubmitCommand(factory Factory, runsCommand GalasaCommand, rootCommand GalasaCommand) (GalasaCommand, error) {
	cmd := new(RunsSubmitCommand)
	err := cmd.init(factory, runsCommand, rootCommand)
	return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *RunsSubmitCommand) Name() string {
	return COMMAND_NAME_RUNS_SUBMIT
}

func (cmd *RunsSubmitCommand) CobraCommand() *cobra.Command {
	return cmd.cobraCommand
}

func (cmd *RunsSubmitCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------

func (cmd *RunsSubmitCommand) init(factory Factory, runsCommand GalasaCommand, rootCommand GalasaCommand) error {
	var err error
	cmd.values = &utils.RunsSubmitCmdValues{}
	cmd.cobraCommand, err = cmd.createRunsSubmitCobraCmd(
		factory,
		cmd.values,
		runsCommand.CobraCommand(),
		runsCommand.Values().(*RunsCmdValues),
		rootCommand.Values().(*RootCmdValues),
	)
	return err
}

func (cmd *RunsSubmitCommand) createRunsSubmitCobraCmd(factory Factory,
	runsSubmitCmdValues *utils.RunsSubmitCmdValues,
	parentCmd *cobra.Command,
	runsCmdValues *RunsCmdValues,
	rootCmdValues *RootCmdValues) (*cobra.Command, error) {

	var err error = nil

	submitSelectionFlags := runs.NewTestSelectionFlagValues()
	runsSubmitCmdValues.TestSelectionFlagValues = submitSelectionFlags

	runsSubmitCmd := &cobra.Command{
		Use:     "submit",
		Short:   "submit a list of tests to the ecosystem",
		Long:    "Submit a list of tests to the ecosystem, monitor them and wait for them to complete",
		Args:    cobra.NoArgs,
		Aliases: []string{"runs submit"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeSubmit(factory, runsSubmitCmdValues, runsCmdValues, rootCmdValues)
		},
	}

	runsSubmitCmd.Flags().StringVarP(&runsSubmitCmdValues.PortfolioFileName, "portfolio", "p", "", "portfolio containing the tests to run")

	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdValues.ReportYamlFilename, "reportyaml", "", "yaml file to record the final results in")
	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdValues.ReportJsonFilename, "reportjson", "", "json file to record the final results in")
	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdValues.ReportJunitFilename, "reportjunit", "", "junit xml file to record the final results in")
	runsSubmitCmd.PersistentFlags().StringVarP(&runsSubmitCmdValues.GroupName, "group", "g", "", "the group name to assign the test runs to, if not provided, a psuedo unique id will be generated")
	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdValues.RequestType, "requesttype", "CLI", "the type of request, used to allocate a run name. Defaults to CLI.")

	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdValues.ThrottleFileName, "throttlefile", "",
		"a file where the current throttle is stored. Periodically the throttle value is read from the file used. "+
			"Someone with edit access to the file can change it which dynamically takes effect. "+
			"Long-running large portfolios can be throttled back to nothing (paused) using this mechanism (if throttle is set to 0). "+
			"And they can be resumed (un-paused) if the value is set back. "+
			"This facility can allow the tests to not show a failure when the system under test is taken out of service for maintainence."+
			"Optional. If not specified, no throttle file is used.",
	)

	runsSubmitCmd.PersistentFlags().IntVar(&runsSubmitCmdValues.PollIntervalSeconds, "poll", runs.DEFAULT_POLL_INTERVAL_SECONDS,
		"Optional. The interval time in seconds between successive polls of the test runs status. "+
			"Defaults to "+strconv.Itoa(runs.DEFAULT_POLL_INTERVAL_SECONDS)+" seconds. "+
			"If less than 1, then default value is used.")

	runsSubmitCmd.PersistentFlags().IntVar(&runsSubmitCmdValues.ProgressReportIntervalMinutes, "progress", runs.DEFAULT_PROGRESS_REPORT_INTERVAL_MINUTES,
		"in minutes, how often the cli will report the overall progress of the test runs. A value of 0 or less disables progress reporting.")

	runsSubmitCmd.PersistentFlags().IntVar(&runsSubmitCmdValues.Throttle, "throttle", runs.DEFAULT_THROTTLE_TESTS_AT_ONCE,
		"how many test runs can be submitted in parallel, 0 or less will disable throttling. 1 causes tests to be run sequentially.")

	runsSubmitCmd.PersistentFlags().StringVar(&runsSubmitCmdValues.OverrideFilePath, "overridefile", "",
		"path to a properties file containing override properties. Defaults to overrides.properties in galasa home folder if that file exists. "+
			"Overrides from --override options will take precedence over properties in this property file. "+
			"A file path of '-' disables reading any properties file.")

	runsSubmitCmd.PersistentFlags().StringSliceVar(&runsSubmitCmdValues.Overrides, "override", make([]string, 0),
		"overrides to be sent with the tests (overrides in the portfolio will take precedence). "+
			"Each override is of the form 'name=value'. Multiple instances of this flag can be used. "+
			"For example --override=prop1=val1 --override=prop2=val2")

	// The trace flag defaults to 'false' if you don't use it.
	// If you say '--trace' on it's own, it defaults to 'true'
	// If you say --trace=false or --trace=true you can set the value explicitly.
	runsSubmitCmd.PersistentFlags().BoolVar(&runsSubmitCmdValues.Trace, "trace", false, "Trace to be enabled on the test runs")
	runsSubmitCmd.PersistentFlags().Lookup("trace").NoOptDefVal = "true"

	runsSubmitCmd.PersistentFlags().BoolVar(&(runsSubmitCmdValues.NoExitCodeOnTestFailures), "noexitcodeontestfailures", false, "set to true if you don't want an exit code to be returned from galasactl if a test fails")

	runs.AddCommandFlags(runsSubmitCmd, submitSelectionFlags)

	parentCmd.AddCommand(runsSubmitCmd)

	return runsSubmitCmd, err
}

func executeSubmit(
	factory Factory,
	runsSubmitCmdValues *utils.RunsSubmitCmdValues,
	runsCmdValues *RunsCmdValues,
	rootCmdValues *RootCmdValues,
) error {

	var err error

	// Operations on the file system will all be relative to the current folder.
	fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, rootCmdValues.logFileName)
	if err == nil {

		rootCmdValues.isCapturingLogs = true

		log.Println("Galasa CLI - Submit tests (Remote)")

		// Get the ability to query environment variables.
		env := factory.GetEnvironment()

		var galasaHome utils.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, rootCmdValues.CmdParamGalasaHomePath)
		if err == nil {

			// Read the bootstrap properties.
			var urlService *api.RealUrlResolutionService = new(api.RealUrlResolutionService)
			var bootstrapData *api.BootstrapData
			bootstrapData, err = api.LoadBootstrap(galasaHome, fileSystem, env, runsCmdValues.bootstrap, urlService)
			if err == nil {

				timeService := factory.GetTimeService()
				var launcherInstance launcher.Launcher = nil

				// The launcher we are going to use to start/monitor tests.
				apiServerUrl := bootstrapData.ApiServerURL
				apiClient := auth.GetAuthenticatedAPIClient(apiServerUrl, fileSystem, galasaHome, timeService)
				launcherInstance = launcher.NewRemoteLauncher(apiServerUrl, apiClient)

				validator := runs.NewStreamBasedValidator()
				err = validator.Validate(runsSubmitCmdValues.TestSelectionFlagValues)
				if err == nil {

					var console = factory.GetStdOutConsole()

					submitter := runs.NewSubmitter(galasaHome, fileSystem, launcherInstance, timeService, env, console)

					if err == nil {
						err = submitter.ExecuteSubmitRuns(runsSubmitCmdValues, runsSubmitCmdValues.TestSelectionFlagValues)
					}
				}
			}
		}
	}

	return err
}
