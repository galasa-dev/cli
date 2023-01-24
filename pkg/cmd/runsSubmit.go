/*
 * Copyright contributors to the Galasa project
 */
package cmd

import (
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	randomGenerator "github.com/satori/go.uuid"
	"github.com/spf13/cobra"

	"github.com/galasa.dev/cli/pkg/api"
	"github.com/galasa.dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/runs"
	"github.com/galasa.dev/cli/pkg/utils"
)

// RunsSubmitCmdParameters - Holds variables set by cobra's command-line parsing.
// We collect the parameters here so that our unit tests can feed in different values
// easily.
type RunsSubmitCmdParameters struct {
	pollIntervalSeconds           int
	noExitCodeOnTestFailures      bool
	reportYamlFilename            string
	reportJsonFilename            string
	reportJunitFilename           string
	groupName                     string
	progressReportIntervalMinutes int
	throttle                      int
	overrides                     []string
	trace                         bool
	requestor                     string
	requestType                   string
	throttleFileName              string
	portfolioFileName             string
	isLocal                       bool // Is the runs submit to a local JVM ?
}

var (
	runsSubmitCmd = &cobra.Command{
		Use:   "submit",
		Short: "submit a list of tests to the ecosystem",
		Long:  "Submit a list of tests to the ecosystem, monitor them and wait for them to complete",
		Args:  cobra.NoArgs,
		Run:   executeSubmit,
	}

	// Variables set by cobra's command-line parsing.
	runsSubmitCmdParams RunsSubmitCmdParameters

	submitSelectionFlags = utils.TestSelectionFlags{}
)

const (
	DEFAULT_POLL_INTERVAL_SECONDS            int = 30
	MAX_INT                                  int = int(^uint(0) >> 1)
	DEFAULT_PROGRESS_REPORT_INTERVAL_MINUTES int = 5
	DEFAULT_THROTTLE_TESTS_AT_ONCE           int = 3
)

func init() {

	currentUserName := getCurrentUserName()

	runsSubmitCmd.Flags().StringVarP(&runsSubmitCmdParams.portfolioFileName, "portfolio", "p", "", "portfolio containing the tests to run")
	runsSubmitCmd.Flags().StringVar(&runsSubmitCmdParams.reportYamlFilename, "reportyaml", "", "yaml file to record the final results in")
	runsSubmitCmd.Flags().StringVar(&runsSubmitCmdParams.reportJsonFilename, "reportjson", "", "json file to record the final results in")
	runsSubmitCmd.Flags().StringVar(&runsSubmitCmdParams.reportJunitFilename, "reportjunit", "", "junit xml file to record the final results in")
	runsSubmitCmd.Flags().StringVarP(&runsSubmitCmdParams.groupName, "group", "g", "", "the group name to assign the test runs to, if not provided, a psuedo unique id will be generated")
	runsSubmitCmd.Flags().StringVar(&runsSubmitCmdParams.requestor, "requestor", currentUserName, "the requestor id to be associated with the test runs. Defaults to the current user id")
	runsSubmitCmd.Flags().StringVar(&runsSubmitCmdParams.requestType, "requesttype", "CLI", "the type of request, used to allocate a run name. Defaults to CLI.")

	runsSubmitCmd.Flags().StringVar(&runsSubmitCmdParams.throttleFileName, "throttlefile", "",
		"a file where the current throttle is stored. Periodically the throttle value is read from the file used. "+
			"Someone with edit access to the file can change it which dynamically takes effect. "+
			"Long-running large portfolios can be throttled back to nothing (paused) using this mechanism (if throttle is set to 0). "+
			"And they can be resumed (un-paused) if the value is set back. "+
			"This facility can allow the tests to not show a failure when the system under test is taken out of service for maintainence.")

	runsSubmitCmd.Flags().IntVar(&runsSubmitCmdParams.pollIntervalSeconds, "poll", DEFAULT_POLL_INTERVAL_SECONDS,
		"Optional. The interval time in seconds between successive polls of the ecosystem for the status of the test runs. "+
			"Defaults to "+strconv.Itoa(DEFAULT_POLL_INTERVAL_SECONDS)+" seconds. "+
			"If less than 1, then default value is used.")

	runsSubmitCmd.Flags().IntVar(&runsSubmitCmdParams.progressReportIntervalMinutes, "progress", DEFAULT_PROGRESS_REPORT_INTERVAL_MINUTES,
		"in minutes, how often the cli will report the overall progress of the test runs, -1 or less will disable progress reports. "+
			"Defaults to "+strconv.Itoa(DEFAULT_PROGRESS_REPORT_INTERVAL_MINUTES)+" minutes. "+
			"If less than 1, then default value is used.")

	runsSubmitCmd.Flags().IntVar(&runsSubmitCmdParams.throttle, "throttle", DEFAULT_THROTTLE_TESTS_AT_ONCE,
		"how many test runs can be submitted in parallel, 0 or less will disable throttling. Default is "+
			strconv.Itoa(DEFAULT_THROTTLE_TESTS_AT_ONCE))

	runsSubmitCmd.Flags().StringSliceVar(&runsSubmitCmdParams.overrides, "override", make([]string, 0),
		"overrides to be sent with the tests (overrides in the portfolio will take precedence). "+
			"Each override is of the form 'name=value'. Multiple instances of this flag can be used. "+
			"For example --override=prop1=val1 --override=prop2=val2")

	// The trace flag defaults to 'false' if you don't use it.
	// If you say '--trace' on it's own, it defaults to 'true'
	// If you say --trace=false or --trace=true you can set the value explicitly.
	runsSubmitCmd.Flags().BoolVar(&runsSubmitCmdParams.trace, "trace", false, "Trace to be enabled on the test runs")
	runsSubmitCmd.Flags().Lookup("trace").NoOptDefVal = "true"

	runsSubmitCmd.Flags().BoolVar(&(runsSubmitCmdParams.noExitCodeOnTestFailures), "noexitcodeontestfailures", false, "set to true if you don't want an exit code to be returned from galasactl if a test fails")

	// The local flag defaults to 'false' if you don't use it.
	// If you say '--local' on it's own, it defaults to 'true'
	// If you say --local=false or --local=true you can set the value explicitly.
	runsSubmitCmd.Flags().BoolVar(&(runsSubmitCmdParams.isLocal), "local", false, "set to true if you don't want an exit code to be returned from galasactl if a test fails")
	localFlag := runsSubmitCmd.Flags().Lookup("trace")
	localFlag.NoOptDefVal = "true"
	localFlag.Hidden = true // Currently this flag is hidden from user view until it works.

	utils.AddCommandFlags(runsSubmitCmd, &submitSelectionFlags)

	runsCmd.AddCommand(runsSubmitCmd)
}

func executeSubmit(cmd *cobra.Command, args []string) {

	var err error
	utils.CaptureLog(logFileName)

	// Operations on the file system will all be relative to the current folder.
	fileSystem := utils.NewOSFileSystem()

	if runsSubmitCmdParams.isLocal {

		javaHome := os.Getenv("JAVA_HOME")

		embeddedFileSystem := embedded.GetEmbeddedFileSystem()

		err = executeSubmitLocal(fileSystem, embeddedFileSystem, runsSubmitCmdParams, javaHome)

	} else {
		// the submit is targetting an ecosysystem to run the command.

		// An HTTP client which can communicate with the api server in an ecosystem.
		var apiClient *galasaapi.APIClient
		apiClient, err = api.InitialiseAPI(bootstrap)
		if err != nil {
			panic(err)
		}

		timeService := utils.NewRealTimeService()

		err = executeSubmitRemote(fileSystem, runsSubmitCmdParams, apiClient, timeService)
	}

	if err != nil {
		// Panic. If we could pass an error back we would.
		// The panic is recovered from in the root command, where
		// the error is logged/displayed before program exit.
		panic(err)
	}
}

func executeSubmitRemote(
	fileSystem utils.FileSystem,
	params RunsSubmitCmdParameters,
	apiClient *galasaapi.APIClient,
	timeService utils.TimeService) error {

	log.Println("Galasa CLI - Submit tests (Remote)")

	var err error = nil

	err = validateAndCorrectParams(&params, apiClient)
	if err != nil {
		return err
	}

	runOverrides, err := buildOverrideMap(params)
	if err != nil {
		return err
	}

	var portfolio *utils.Portfolio
	portfolio, err = getPortfolio(fileSystem, params.portfolioFileName, apiClient)
	if err != nil {
		return err
	}

	err = validatePortfolio(portfolio)
	if err != nil {
		return err
	}

	// Build list of runs to submit
	readyRuns := buildListOfRunsToSubmit(portfolio, runOverrides)

	submittedRuns := make(map[string]*runs.TestRun)
	rerunRuns := make(map[string]*runs.TestRun)
	finishedRuns := make(map[string]*runs.TestRun)
	lostRuns := make(map[string]*runs.TestRun)

	progressReportInterval := time.Minute * time.Duration(params.progressReportIntervalMinutes)
	throttle := params.throttle
	fetchRas := isRasDetailNeededForReports(params)
	pollInterval := time.Second * time.Duration(params.pollIntervalSeconds)

	err = writeThrottleFile(fileSystem, params.throttleFileName, throttle)
	if err != nil {
		return err
	}

	//
	// Main submit loop
	//
	nextProgressReport := timeService.Now().Add(progressReportInterval)
	isThrottleFileLost := false
	for len(readyRuns) > 0 || len(submittedRuns) > 0 || len(rerunRuns) > 0 { // Loop whilst there are runs to submit or are running

		for len(submittedRuns) < throttle && len(readyRuns) > 0 {
			readyRuns = submitRun(apiClient, params.groupName, readyRuns, submittedRuns,
				lostRuns, &runOverrides, params.trace, params.requestor, params.requestType)
		}

		now := timeService.Now()
		if now.After(nextProgressReport) {
			runs.InterrimProgressReport(readyRuns, submittedRuns, finishedRuns, lostRuns, throttle)
			nextProgressReport = now.Add(progressReportInterval)
		}

		timeService.Sleep(pollInterval)

		throttle, isThrottleFileLost = updateThrottleFromFileIfDifferent(fileSystem, params.throttleFileName, throttle, isThrottleFileLost)

		runsFetchCurrentStatus(apiClient, params.groupName, readyRuns, submittedRuns, finishedRuns, lostRuns, fetchRas)
	}

	// Generate all the reports summarising the end-results.
	err = createReports(fileSystem, params, finishedRuns, lostRuns)
	if err == nil {

		// Fail the command if tests failed, and the user wanted us to fail if tests fail.
		failureCount := runs.CountTotalFailedRuns(finishedRuns, lostRuns)
		if failureCount > 0 && !params.noExitCodeOnTestFailures {
			// Not all runs passed
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_TESTS_FAILED)
		}
	}

	return err
}

func writeThrottleFile(fileSystem utils.FileSystem, throttleFileName string, throttle int) error {
	err := fileSystem.WriteTextFile(throttleFileName, strconv.Itoa(throttle))
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_THROTTLE_FILE_WRITE, throttleFileName, err.Error())
	}
	return err
}

func updateThrottleFromFileIfDifferent(
	fileSystem utils.FileSystem,
	throttleFileName string,
	currentThrottle int,
	wasThrottleFileLostAlready bool) (int, bool) {

	var newThrottle int = currentThrottle
	var isThrottleFileLost bool = false

	savedThrottle, err := readThrottleFile(fileSystem, throttleFileName)
	if err != nil {
		if wasThrottleFileLostAlready {
			// Don't log it, as we logged it when it was first lost.
		} else {
			// We just lost it, so log the fact.
			log.Printf("Error with throttle file %v\n", err)
		}
		// Don't report the throttle file as being lost until after it's been found again.
		isThrottleFileLost = true
	} else {
		isThrottleFileLost = false
		// Only log something if we are changing the throttle value.
		if savedThrottle != currentThrottle {
			log.Printf("Changing throttle from %v to %v\n", currentThrottle, newThrottle)
		}
		newThrottle = savedThrottle
	}
	return newThrottle, isThrottleFileLost
}

const (
	INT_TYPE_VARIANT_PLAIN_INT = 0
	INT_TYPE_VARIANT_INT8      = 8
	INT_TYPE_VARIANT_INT16     = 16
	INT_TYPE_VARIANT_INT32     = 32
	INT_TYPE_VARIANT_INT64     = 64
)

func readThrottleFile(fileSystem utils.FileSystem, throttleFileName string) (int, error) {
	var savedThrottle int = 0
	var intermediateThrottle int64
	contents, err := fileSystem.ReadTextFile(throttleFileName)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_THROTTLE_FILE_READ, throttleFileName, err.Error())
	} else {
		intermediateThrottle, err = strconv.ParseInt(contents, 10, INT_TYPE_VARIANT_PLAIN_INT)
		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_THROTTLE_FILE_INVALID, contents, throttleFileName, err.Error())
			log.Printf("Throttle file %s contains bad integer value %s returning %s\n", throttleFileName, contents, err.Error())
		} else {
			savedThrottle = int(intermediateThrottle)
		}
	}
	return savedThrottle, err
}

func submitRun(
	apiClient *galasaapi.APIClient,
	groupName string,
	readyRuns []runs.TestRun,
	submittedRuns map[string]*runs.TestRun,
	lostRuns map[string]*runs.TestRun,
	runOverrides *map[string]string,
	trace bool,
	requestor string,
	requestType string) []runs.TestRun {

	if len(readyRuns) < 1 {
		return readyRuns
	}

	nextRun := readyRuns[0]
	readyRuns = readyRuns[1:]

	className := nextRun.Bundle + "/" + nextRun.Class
	classNames := []string{className}

	submitOverrides := make(map[string]interface{})

	for key, value := range nextRun.Overrides {
		submitOverrides[key] = value
	}

	testRunRequest := galasaapi.NewTestRunRequest()
	testRunRequest.SetClassNames(classNames)
	testRunRequest.SetRequestorType(requestType)
	testRunRequest.SetRequestor(requestor)
	testRunRequest.SetTestStream(nextRun.Stream)
	testRunRequest.SetTrace(trace)
	testRunRequest.SetOverrides(submitOverrides)

	resultGroup, _, err := apiClient.RunsAPIApi.PostSubmitTestRuns(nil, groupName).TestRunRequest(*testRunRequest).Execute()
	if err != nil {
		log.Printf("Failed to submit test %v/%v - %v\n", nextRun.Bundle, nextRun.Class, err)
		lostRuns[className] = &nextRun
		return readyRuns
	}

	if len(resultGroup.GetRuns()) < 1 {
		log.Printf("Lost the run attempting to submit test %v/%v\n", nextRun.Bundle, nextRun.Class)
		lostRuns[className] = &nextRun
		return readyRuns
	}

	submittedRun := resultGroup.GetRuns()[0]
	nextRun.Name = *submittedRun.Name

	submittedRuns[nextRun.Name] = &nextRun

	log.Printf("Run %v submitted - %v/%v/%v\n", nextRun.Name, nextRun.Stream, nextRun.Bundle, nextRun.Class)

	return readyRuns
}

func runsFetchCurrentStatus(
	apiClient *galasaapi.APIClient,
	groupName string,
	readyRuns []runs.TestRun,
	submittedRuns map[string]*runs.TestRun,
	finishedRuns map[string]*runs.TestRun,
	lostRuns map[string]*runs.TestRun,
	fetchRas bool) {

	currentGroup, _, err := apiClient.RunsAPIApi.GetRunsGroup(nil, groupName).Execute()
	if err != nil {
		log.Printf("Received error from group request - %v\n", err)
		return
	}

	// a copy to find lost runs
	checkRuns := runs.DeepClone(submittedRuns)

	for _, currentRun := range currentGroup.GetRuns() {
		runName := currentRun.GetName()

		checkRun, ok := submittedRuns[runName]
		if ok {
			//First remove from the checkRuns as we know it still exists in the ecosystem
			delete(checkRuns, runName)

			// now check to see if it is finished
			if currentRun.GetStatus() == "finished" {
				finishedRuns[runName] = checkRun
				delete(submittedRuns, runName)

				result := "unknown"
				if currentRun.HasResult() {
					result = currentRun.GetResult()
				}
				checkRun.Result = result

				// Extract the ras run result to get the method names if a report is requested
				rasRunID := currentRun.RasRunId
				if fetchRas && rasRunID != nil {
					rasRun, _, err := apiClient.ResultArchiveStoreAPIApi.GetRasRunById(nil, *rasRunID).Execute()
					if err != nil {
						log.Printf("Failed to retrieve RAS run for %v - %v\n", checkRun.Name, err)
					} else {
						checkRun.Tests = make([]runs.TestMethod, 0)

						testStructure := rasRun.GetTestStructure()
						for _, testMethod := range testStructure.GetMethods() {
							test := runs.TestMethod{
								Method: testMethod.GetMethodName(),
								Result: testMethod.GetResult(),
							}

							checkRun.Tests = append(checkRun.Tests, test)
						}
					}
				}

				log.Printf("Run %v has finished(%v) - %v/%v/%v\n", runName, result, checkRun.Stream, checkRun.Bundle, checkRun.Class)
			} else {
				// Check to see if there was a status change
				if checkRun.Status != currentRun.GetStatus() {
					checkRun.Status = currentRun.GetStatus()
					log.Printf("    Run %v status is now '%v' - %v/%v/%v\n", runName, checkRun.Status, checkRun.Stream, checkRun.Bundle, checkRun.Class)
				}
			}
		}
	}

	// Now deal with the lost runs
	for runName, lostRun := range checkRuns {
		lostRuns[runName] = lostRun
		delete(submittedRuns, runName)
		log.Printf("Run %v was lost - %v/%v/%v\n", runName, lostRun.Stream, lostRun.Bundle, lostRun.Class)
	}

}

func createReports(fileSystem utils.FileSystem, params RunsSubmitCmdParameters, finishedRuns map[string]*runs.TestRun, lostRuns map[string]*runs.TestRun) error {
	runs.FinalHumanReadableReport(finishedRuns, lostRuns)

	var err error = nil

	if params.reportYamlFilename != "" {
		err = runs.ReportYaml(fileSystem, params.reportYamlFilename, finishedRuns, lostRuns)
	}

	if err == nil {
		if params.reportJsonFilename != "" {
			err = runs.ReportJSON(fileSystem, params.reportJsonFilename, finishedRuns, lostRuns)
		}
	}

	if err == nil {
		if params.reportJunitFilename != "" {
			err = runs.ReportJunit(fileSystem, params.reportJunitFilename, params.groupName, finishedRuns, lostRuns)
		}
	}

	return err
}

func isRasDetailNeededForReports(params RunsSubmitCmdParameters) bool {

	// Do we need to ask the RAS for the test structure
	isRasDetailNeeded := false
	if params.reportYamlFilename != "" {
		isRasDetailNeeded = true
	}
	if params.reportJsonFilename != "" {
		isRasDetailNeeded = true
	}
	if params.reportJunitFilename != "" {
		isRasDetailNeeded = true
	}

	return isRasDetailNeeded
}

func buildListOfRunsToSubmit(portfolio *utils.Portfolio, runOverrides map[string]string) []runs.TestRun {
	readyRuns := make([]runs.TestRun, 0, len(portfolio.Classes))

	for _, portfolioTest := range portfolio.Classes {
		newTestrun := runs.TestRun{
			Bundle:    portfolioTest.Bundle,
			Class:     portfolioTest.Class,
			Stream:    portfolioTest.Stream,
			Status:    "queued",
			Overrides: make(map[string]string, 0),
		}

		// load the run overrides
		for key, value := range runOverrides {
			newTestrun.Overrides[key] = value
		}

		// load the assemble overrides, they take precedence on the run overrides
		for key, value := range portfolioTest.Overrides {
			newTestrun.Overrides[key] = value
		}

		readyRuns = append(readyRuns, newTestrun)

		log.Printf("Added test %v/%v/%v to the ready queue\n", newTestrun.Stream, newTestrun.Bundle, newTestrun.Class)
	}

	return readyRuns
}

func validateAndCorrectParams(params *RunsSubmitCmdParameters, apiClient *galasaapi.APIClient) error {

	var err error = nil

	// Guard against the poll time being less than 1 second
	if params.pollIntervalSeconds < 1 {
		log.Printf("poll value is invalid. Less than 1. Defaulting value to %v seconds.\n", DEFAULT_POLL_INTERVAL_SECONDS)
		params.pollIntervalSeconds = DEFAULT_POLL_INTERVAL_SECONDS
	}

	// Set the progress time
	if params.progressReportIntervalMinutes < 0 {
		params.progressReportIntervalMinutes = MAX_INT
	} else if params.progressReportIntervalMinutes == 0 {
		params.progressReportIntervalMinutes = 5
	}

	// Set the throttle
	if params.throttle <= 0 {
		params.throttle = MAX_INT // set to maximum size of the int
	}

	//  Dont mix portfolio and test selection on the same command
	if params.portfolioFileName != "" {
		if utils.AreSelectionFlagsProvided(&submitSelectionFlags) {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_MIX_FLAGS_AND_PORTFOLIO)
		}
	} else {
		if !utils.AreSelectionFlagsProvided(&submitSelectionFlags) {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_MISSING_ACTION_FLAGS)
		}
	}

	if err == nil {
		// generate a group name if required
		if params.groupName == "" {
			params.groupName = randomGenerator.NewV4().String()
		}
		log.Printf("Using group name '%v' for test run submission\n", params.groupName)

		_, err = checkIfGroupAlreadyInUse(*apiClient, params.groupName)
	}

	return err
}

func buildOverrideMap(params RunsSubmitCmdParameters) (map[string]string, error) {
	var err error = nil

	// Convert overrides to a map
	runOverrides := make(map[string]string)
	for _, override := range params.overrides {
		pos := strings.Index(override, "=")
		if pos < 1 {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_INVALID_OVERRIDE, override)
			return nil, err
		}
		key := override[:pos]
		value := override[pos+1:]
		if value == "" {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_INVALID_OVERRIDE, override)
			return nil, err
		}
		runOverrides[key] = value
	}
	return runOverrides, nil
}

func getPortfolio(fileSystem utils.FileSystem, portfolioFileName string, apiClient *galasaapi.APIClient) (*utils.Portfolio, error) {
	// Load the portfolio of tests
	var portfolio *utils.Portfolio
	var err error = nil

	if portfolioFileName != "" {
		portfolio, err = utils.LoadPortfolio(fileSystem, portfolioFileName)
		if err != nil {
			return nil, err
		}
	} else {
		// There is no portfolio file, so create an in-memory portfolio 
		// from the tests we can find from the test selection.
		testSelection := utils.SelectTests(apiClient, &submitSelectionFlags)

		testOverrides := make(map[string]string)
		portfolio = utils.NewPortfolio()
		utils.CreatePortfolio(&testSelection, &testOverrides, portfolio)
	}
	return portfolio, nil
}

func getCurrentUserName() string {
	userName := "cli"
	currentUser, err := user.Current()
	if err == nil {
		userName = currentUser.Username
	}
	return userName
}

func validatePortfolio(portfolio *utils.Portfolio) error {
	var err error = nil
	if portfolio.Classes == nil || len(portfolio.Classes) < 1 {
		// Empty portfolio
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_EMPTY_PORTFOLIO, portfolioFilename)
	}
	return err
}

func checkIfGroupAlreadyInUse(apiClient galasaapi.APIClient, groupName string) (bool, error) {
	isInUse := false
	var err error = nil
	// Just check if it is already in use,  which is perfectly valid for custom group names
	uuidCheck, _, err := apiClient.RunsAPIApi.GetRunsGroup(nil, groupName).Execute()
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_RUNS_GROUP_CHECK, groupName, err.Error())
	} else {

		if uuidCheck.Runs != nil && len(uuidCheck.Runs) > 0 {
			log.Printf("Group name '%v' is aleady in use\n", groupName)
			isInUse = true
		}
	}
	return isInUse, err
}
