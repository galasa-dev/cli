/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"log"
	"os/user"
	"strconv"
	"strings"
	"time"

	randomGenerator "github.com/satori/go.uuid"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/launcher"
	"github.com/galasa.dev/cli/pkg/props"
	"github.com/galasa.dev/cli/pkg/utils"
)

func ExecuteSubmitRuns(
	galasaHome utils.GalasaHome,
	fileSystem files.FileSystem,
	params utils.RunsSubmitCmdParameters,
	launcher launcher.Launcher,
	timeService utils.TimeService,
	testSelectionFlags *TestSelectionFlags,
) error {

	var err error = nil

	err = validateAndCorrectParams(galasaHome, fileSystem, &params, launcher, testSelectionFlags)
	if err != nil {
		return err
	}

	runOverrides, err := buildOverrideMap(fileSystem, params)
	if err != nil {
		return err
	}

	var portfolio *Portfolio
	portfolio, err = getPortfolio(fileSystem, params.PortfolioFileName, launcher, testSelectionFlags)
	if err != nil {
		return err
	}

	err = validatePortfolio(portfolio, params.PortfolioFileName)
	if err != nil {
		return err
	}

	// Build list of runs to submit
	readyRuns := buildListOfRunsToSubmit(portfolio, runOverrides)

	// Run all the tests
	var finishedRuns map[string]*TestRun
	var lostRuns map[string]*TestRun
	finishedRuns, lostRuns, err = executeSubmitRuns(
		fileSystem, params, launcher, timeService, readyRuns, runOverrides)

	// Report on the results.
	if err == nil {
		// Generate all the reports summarising the end-results.
		err = createReports(fileSystem, params, finishedRuns, lostRuns)
		if err == nil {

			// Fail the command if tests failed, and the user wanted us to fail if tests fail.
			failureCount := CountTotalFailedRuns(finishedRuns, lostRuns)
			if failureCount > 0 && !params.NoExitCodeOnTestFailures {
				// Not all runs passed
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_TESTS_FAILED, failureCount)
			}
		}
	}

	return err
}

func executeSubmitRuns(fileSystem files.FileSystem,
	params utils.RunsSubmitCmdParameters,
	launcher launcher.Launcher,
	timeService utils.TimeService,
	readyRuns []TestRun,
	runOverrides map[string]string) (map[string]*TestRun, map[string]*TestRun, error) {

	var err error = nil

	submittedRuns := make(map[string]*TestRun)
	rerunRuns := make(map[string]*TestRun)
	finishedRuns := make(map[string]*TestRun)
	lostRuns := make(map[string]*TestRun)

	progressReportInterval := time.Minute * time.Duration(params.ProgressReportIntervalMinutes)
	throttle := params.Throttle
	fetchRas := isRasDetailNeededForReports(params)
	pollInterval := time.Second * time.Duration(params.PollIntervalSeconds)

	err = writeThrottleFile(fileSystem, params.ThrottleFileName, throttle)
	if err != nil {
		return nil, nil, err
	}

	//
	// Main submit loop
	//
	nextProgressReport := timeService.Now().Add(progressReportInterval)
	isThrottleFileLost := false
	for len(readyRuns) > 0 || len(submittedRuns) > 0 || len(rerunRuns) > 0 { // Loop whilst there are runs to submit or are running

		for len(submittedRuns) < throttle && len(readyRuns) > 0 {
			readyRuns = submitRun(launcher, params.GroupName, readyRuns, submittedRuns,
				lostRuns, &runOverrides, params.Trace, params.Requestor, params.RequestType)
		}

		// Only do progress reporting if the user didn't disable it.
		if params.ProgressReportIntervalMinutes > 0 {
			now := timeService.Now()
			if now.After(nextProgressReport) {
				InterrimProgressReport(readyRuns, submittedRuns, finishedRuns, lostRuns, throttle)
				nextProgressReport = now.Add(progressReportInterval)
			}
		}

		throttle, isThrottleFileLost = updateThrottleFromFileIfDifferent(fileSystem, params.ThrottleFileName, throttle, isThrottleFileLost)

		runsFetchCurrentStatus(launcher, params.GroupName, readyRuns, submittedRuns, finishedRuns, lostRuns, fetchRas)

		// Only sleep if there are runs in progress but not yet finished.
		if len(submittedRuns) > 0 || len(rerunRuns) > 0 {
			log.Printf("Sleeping for the poll interval of %v seconds\n", params.PollIntervalSeconds)
			timeService.Sleep(pollInterval)
			log.Printf("Awake from poll interval sleep of %v seconds\n", params.PollIntervalSeconds)
		}
	}

	return finishedRuns, lostRuns, err
}

func writeThrottleFile(fileSystem files.FileSystem, throttleFileName string, throttle int) error {
	var err error = nil
	if throttleFileName != "" {
		// Throttle filename was specified. Lets use a throttle file.
		err = fileSystem.WriteTextFile(throttleFileName, strconv.Itoa(throttle))
		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_THROTTLE_FILE_WRITE, throttleFileName, err.Error())
		}
	}
	return err
}

func updateThrottleFromFileIfDifferent(
	fileSystem files.FileSystem,
	throttleFileName string,
	currentThrottle int,
	wasThrottleFileLostAlready bool) (int, bool) {

	var newThrottle int = currentThrottle
	var isThrottleFileLost bool = wasThrottleFileLostAlready

	// Only bother with anything here if there is a throttle file specified by the user.
	if throttleFileName != "" {

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

func readThrottleFile(fileSystem files.FileSystem, throttleFileName string) (int, error) {
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
	launcher launcher.Launcher,
	groupName string,
	readyRuns []TestRun,
	submittedRuns map[string]*TestRun,
	lostRuns map[string]*TestRun,
	runOverrides *map[string]string,
	trace bool,
	requestor string,
	requestType string) []TestRun {

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

	var resultGroup *galasaapi.TestRuns
	var err error
	resultGroup, err = launcher.SubmitTestRuns(groupName, classNames, requestType, requestor, nextRun.Stream, trace, submitOverrides)
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
	launcher launcher.Launcher,
	groupName string,
	readyRuns []TestRun,
	submittedRuns map[string]*TestRun,
	finishedRuns map[string]*TestRun,
	lostRuns map[string]*TestRun,
	fetchRas bool) {

	currentGroup, err := launcher.GetRunsByGroup(groupName)
	if err != nil {
		log.Printf("Received error from group request - %v\n", err)
		return
	}

	// a copy to find lost runs
	checkRuns := DeepClone(submittedRuns)

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

					rasRun, err := launcher.GetRunsById(*rasRunID)

					if err != nil {
						log.Printf("Failed to retrieve RAS run for %v - %v\n", checkRun.Name, err)
					} else {
						checkRun.Tests = make([]TestMethod, 0)

						testStructure := rasRun.GetTestStructure()
						for _, testMethod := range testStructure.GetMethods() {
							test := TestMethod{
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

func createReports(fileSystem files.FileSystem, params utils.RunsSubmitCmdParameters,
	finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun) error {

	FinalHumanReadableReport(finishedRuns, lostRuns)

	var err error = nil

	if params.ReportYamlFilename != "" {
		err = ReportYaml(fileSystem, params.ReportYamlFilename, finishedRuns, lostRuns)
	}

	if err == nil {
		if params.ReportJsonFilename != "" {
			err = ReportJSON(fileSystem, params.ReportJsonFilename, finishedRuns, lostRuns)
		}
	}

	if err == nil {
		if params.ReportJunitFilename != "" {
			err = ReportJunit(fileSystem, params.ReportJunitFilename, params.GroupName, finishedRuns, lostRuns)
		}
	}

	return err
}

func isRasDetailNeededForReports(params utils.RunsSubmitCmdParameters) bool {

	// Do we need to ask the RAS for the test structure
	isRasDetailNeeded := false
	if params.ReportYamlFilename != "" {
		isRasDetailNeeded = true
	}
	if params.ReportJsonFilename != "" {
		isRasDetailNeeded = true
	}
	if params.ReportJunitFilename != "" {
		isRasDetailNeeded = true
	}

	return isRasDetailNeeded
}

func buildListOfRunsToSubmit(portfolio *Portfolio, runOverrides map[string]string) []TestRun {
	readyRuns := make([]TestRun, 0, len(portfolio.Classes))

	for _, portfolioTest := range portfolio.Classes {
		newTestrun := TestRun{
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

func validateAndCorrectParams(
	galasaHome utils.GalasaHome,
	fs files.FileSystem,
	params *utils.RunsSubmitCmdParameters,
	launcher launcher.Launcher,
	submitSelectionFlags *TestSelectionFlags,
) error {

	var err error = nil

	// Guard against the poll time being less than 1 second
	if params.PollIntervalSeconds < 1 {
		log.Printf("poll value is invalid. Less than 1. Defaulting value to %v seconds.\n", DEFAULT_POLL_INTERVAL_SECONDS)
		params.PollIntervalSeconds = DEFAULT_POLL_INTERVAL_SECONDS
	}

	// Set the progress reporting interval time
	if params.ProgressReportIntervalMinutes <= 0 {
		params.ProgressReportIntervalMinutes = 0
	}

	// Set the throttle
	if params.Throttle <= 0 {
		params.Throttle = MAX_INT // set to maximum size of the int
	}

	//  Dont mix portfolio and test selection on the same command
	if params.PortfolioFileName != "" {
		if AreSelectionFlagsProvided(submitSelectionFlags) {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_MIX_FLAGS_AND_PORTFOLIO)
		}
	} else {
		if !AreSelectionFlagsProvided(submitSelectionFlags) {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_MISSING_ACTION_FLAGS)
		}
	}

	if err == nil {
		// generate a group name if required
		if params.GroupName == "" {
			params.GroupName = randomGenerator.NewV4().String()
		}
		log.Printf("Using group name '%v' for test run submission\n", params.GroupName)

		_, err = checkIfGroupAlreadyInUse(launcher, params.GroupName)
	}

	if err == nil {
		err = correctOverrideFilePathParameter(galasaHome, fs, params)
	}

	tildaExpandAllPaths(fs, params)

	return err
}

func correctOverrideFilePathParameter(
	galasaHome utils.GalasaHome,
	fs files.FileSystem,
	params *utils.RunsSubmitCmdParameters,
) error {
	var err error
	// Correct the default overrideFile path if it wasn't specified.
	if params.OverrideFilePath == "" {

		params.OverrideFilePath = galasaHome.GetUrlFolderPath() + "/overrides.properties"
		var isFileThere bool
		isFileThere, err = fs.Exists(params.OverrideFilePath)
		if err == nil {
			if !isFileThere {
				// The flag wasn't specified.
				// And we don't have an overrides file to read from the .galasa folder.
				// So treat this the same as the user not wanting to use an override file.
				// If the file existed, then we'd want to use it.
				params.OverrideFilePath = "-"
			}
		}
	}
	return err
}

func tildaExpandAllPaths(fs files.FileSystem, params *utils.RunsSubmitCmdParameters) error {
	var err error = nil

	if err == nil {
		params.OverrideFilePath, err = files.TildaExpansion(fs, params.OverrideFilePath)
	}

	if err == nil {
		params.PortfolioFileName, err = files.TildaExpansion(fs, params.PortfolioFileName)
	}

	if err == nil {
		params.ReportJsonFilename, err = files.TildaExpansion(fs, params.ReportJsonFilename)
	}

	if err == nil {
		params.ReportJunitFilename, err = files.TildaExpansion(fs, params.ReportJunitFilename)
	}

	if err == nil {
		params.ReportYamlFilename, err = files.TildaExpansion(fs, params.ReportYamlFilename)
	}

	if err == nil {
		params.ThrottleFileName, err = files.TildaExpansion(fs, params.ThrottleFileName)
	}
	return err
}

func buildOverrideMap(fileSystem files.FileSystem, commandParameters utils.RunsSubmitCmdParameters) (map[string]string, error) {

	path := commandParameters.OverrideFilePath
	runOverrides, err := loadOverrideFile(fileSystem, path)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_LOAD_OVERRIDES_FILE, path, err.Error())
	} else {
		runOverrides, err = addOverridesFromCmdLine(runOverrides, commandParameters.Overrides)
	}

	return runOverrides, err
}

func loadOverrideFile(fileSystem files.FileSystem, overrideFilePath string) (map[string]string, error) {

	var (
		overrides props.JavaProperties
		err       error = nil
	)

	if overrideFilePath == "-" {
		// Don't read properties from a file.
		overrides = make(map[string]string)
	} else {
		overrides, err = props.ReadPropertiesFile(fileSystem, overrideFilePath)
	}

	return overrides, err
}

func addOverridesFromCmdLine(overrides map[string]string, commandLineOverrides []string) (map[string]string, error) {
	var err error = nil

	// Convert overrides to a map
	for _, override := range commandLineOverrides {
		pos := strings.Index(override, "=")
		if pos < 1 {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_INVALID_OVERRIDE, override)
			break
		}
		key := override[:pos]
		value := override[pos+1:]
		if value == "" {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_INVALID_OVERRIDE, override)
			break
		}
		overrides[key] = value
	}

	// Discard overrides if there was an error.
	if err != nil {
		overrides = nil
	}
	return overrides, nil
}

func getPortfolio(fileSystem files.FileSystem, portfolioFileName string, launcher launcher.Launcher, submitSelectionFlags *TestSelectionFlags) (*Portfolio, error) {
	// Load the portfolio of tests
	var portfolio *Portfolio = nil
	var err error = nil

	if portfolioFileName != "" {
		portfolio, err = ReadPortfolio(fileSystem, portfolioFileName)
	} else {
		// There is no portfolio file, so create an in-memory portfolio
		// from the tests we can find from the test selection.
		var testSelection TestSelection
		testSelection, err = SelectTests(launcher, submitSelectionFlags)

		if err == nil {
			testOverrides := make(map[string]string)
			portfolio = NewPortfolio()
			AddClassesToPortfolio(&testSelection, &testOverrides, portfolio)
		}
	}
	return portfolio, err
}

func GetCurrentUserName() string {
	userName := "cli"
	currentUser, err := user.Current()
	if err == nil {
		userName = currentUser.Username
	}
	return userName
}

func validatePortfolio(portfolio *Portfolio, portfolioFilename string) error {
	var err error = nil
	if portfolio.Classes == nil || len(portfolio.Classes) < 1 {
		// Empty portfolio
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_EMPTY_PORTFOLIO, portfolioFilename)
	}
	return err
}

func checkIfGroupAlreadyInUse(launcher launcher.Launcher, groupName string) (bool, error) {
	isInUse := false
	var err error = nil

	// Just check if it is already in use,  which is perfectly valid for custom group names
	uuidCheck, err := launcher.GetRunsByGroup(groupName)
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
