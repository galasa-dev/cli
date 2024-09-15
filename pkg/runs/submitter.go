/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"fmt"
	"log"
	"os/user"
	"strconv"
	"strings"
	"time"

	randomGenerator "github.com/google/uuid"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/images"
	"github.com/galasa-dev/cli/pkg/launcher"
	"github.com/galasa-dev/cli/pkg/props"
	"github.com/galasa-dev/cli/pkg/runsformatter"
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/utils"
)

type Submitter struct {
	galasaHome  spi.GalasaHome
	fileSystem  spi.FileSystem
	launcher    launcher.Launcher
	timeService spi.TimeService
	env         spi.Environment
	console     spi.Console
	expander    images.ImageExpander
}

func NewSubmitter(
	galasaHome spi.GalasaHome,
	fileSystem spi.FileSystem,
	launcher launcher.Launcher,
	timeService spi.TimeService,
	env spi.Environment,
	console spi.Console,
	expander images.ImageExpander,
) *Submitter {
	instance := new(Submitter)
	instance.galasaHome = galasaHome
	instance.fileSystem = fileSystem
	instance.launcher = launcher
	instance.timeService = timeService
	instance.env = env
	instance.console = console
	instance.expander = expander
	return instance
}

func (submitter *Submitter) ExecuteSubmitRuns(
	params *utils.RunsSubmitCmdValues,
	TestSelectionFlagValues *utils.TestSelectionFlagValues,

) error {

	var err error

	err = submitter.validateAndCorrectParams(params, TestSelectionFlagValues)
	if err == nil {
		var runOverrides map[string]string
		runOverrides, err = submitter.buildOverrideMap(*params)
		if err == nil {
			var portfolio *Portfolio
			portfolio, err = submitter.getPortfolio(params.PortfolioFileName, TestSelectionFlagValues)
			if err == nil {
				err = submitter.validatePortfolio(portfolio, params.PortfolioFileName)
				if err == nil {
					err = submitter.executePortfolio(portfolio, runOverrides, *params)
				}
			}
		}
	}

	return err
}

func (submitter *Submitter) executePortfolio(portfolio *Portfolio,
	runOverrides map[string]string,
	params utils.RunsSubmitCmdValues,
) error {

	var err error

	// Build list of runs to submit
	readyRuns := submitter.buildListOfRunsToSubmit(portfolio, runOverrides)

	// Run all the tests
	var finishedRuns map[string]*TestRun
	var lostRuns map[string]*TestRun
	finishedRuns, lostRuns, err = submitter.executeSubmitRuns(
		params, readyRuns, runOverrides)

	// Report on the results.
	if err == nil {
		// Generate all the reports summarising the end-results.
		err = submitter.createReports(params, finishedRuns, lostRuns)
		if err == nil {

			err = reportRendedImages(finishedRuns, submitter)

			if err == nil {

				// Fail the command if tests failed, and the user wanted us to fail if tests fail.
				failureCount := CountTotalFailedRuns(finishedRuns, lostRuns)
				if failureCount > 0 && !params.NoExitCodeOnTestFailures {
					// Not all runs passed
					err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_TESTS_FAILED, failureCount)
				}
			}
		}

	}

	return err
}

func reportRendedImages(finishedRuns map[string]*TestRun, submitter *Submitter) error {
	var err error

	for runName := range finishedRuns {

		folderToScan := submitter.galasaHome.GetNativeFolderPath() + "/ras/" + runName
		err = submitter.expander.ExpandImages(folderToScan)
		if err != nil {
			break
		}
	}

	return err
}

func (submitter *Submitter) executeSubmitRuns(
	params utils.RunsSubmitCmdValues,
	readyRuns []TestRun,
	runOverrides map[string]string,
) (map[string]*TestRun, map[string]*TestRun, error) {

	var err error

	submittedRuns := make(map[string]*TestRun)
	rerunRuns := make(map[string]*TestRun)
	finishedRuns := make(map[string]*TestRun)
	lostRuns := make(map[string]*TestRun)

	progressReportInterval := time.Minute * time.Duration(params.ProgressReportIntervalMinutes)
	throttle := params.Throttle
	fetchRas := submitter.isRasDetailNeededForReports(params)
	pollInterval := time.Second * time.Duration(params.PollIntervalSeconds)

	err = submitter.writeThrottleFile(params.ThrottleFileName, throttle)
	if err != nil {
		return nil, nil, err
	}

	currentUser := submitter.GetCurrentUserName()
	//
	// Main submit loop
	//
	nextProgressReport := submitter.timeService.Now().Add(progressReportInterval)
	isThrottleFileLost := false

	for len(readyRuns) > 0 || len(submittedRuns) > 0 || len(rerunRuns) > 0 { // Loop whilst there are runs to submit or are running

		for len(submittedRuns) < throttle && len(readyRuns) > 0 {

			readyRuns, err = submitter.submitRun(params.GroupName, readyRuns, submittedRuns,
				lostRuns, &runOverrides, params.Trace, currentUser, params.RequestType)

			if err != nil {
				// Ignore the error and continue to process the list of available runs.
				submitter.console.WriteString(fmt.Sprintf("%s\n", err.Error()))
			}
		}

		// Only do progress reporting if the user didn't disable it.
		if params.ProgressReportIntervalMinutes > 0 {
			now := submitter.timeService.Now()
			if now.After(nextProgressReport) {
				//convert TestRun
				submitter.displayInterrimProgressReport(readyRuns, submittedRuns, finishedRuns, lostRuns, throttle)
				nextProgressReport = now.Add(progressReportInterval)
			}
		}

		throttle, isThrottleFileLost = submitter.updateThrottleFromFileIfDifferent(params.ThrottleFileName, throttle, isThrottleFileLost)

		submitter.runsFetchCurrentStatus(params.GroupName, submittedRuns, finishedRuns, lostRuns, fetchRas)

		// Only sleep if there are runs in progress but not yet finished.
		if len(submittedRuns) > 0 || len(rerunRuns) > 0 {
			log.Printf("Sleeping for the poll interval of %v seconds\n", params.PollIntervalSeconds)
			submitter.timeService.Sleep(pollInterval)
			log.Printf("Awake from poll interval sleep of %v seconds\n", params.PollIntervalSeconds)
		}
	}

	return finishedRuns, lostRuns, err
}

func (submitter *Submitter) displayInterrimProgressReport(readyRuns []TestRun,
	submittedRuns map[string]*TestRun,
	finishedRuns map[string]*TestRun,
	lostRuns map[string]*TestRun,
	throttle int) {

	ready := len(readyRuns)
	submitted := len(submittedRuns)
	finished := len(finishedRuns)
	lost := len(lostRuns)

	log.Println("Progress report")
	for runName, run := range submittedRuns {
		log.Printf("***     Run %v is currently %v - %v/%v/%v\n", runName, run.Status, run.Stream, run.Bundle, run.Class)
	}
	log.Println("----------------------------------------------------------------------------")
	log.Printf("Run status: Ready=%v, Submitted=%v, Finished=%v, Lost=%v\n", ready, submitted, finished, lost)
	log.Printf("Throttle=%v\n", throttle)

	if finished > 0 {
		submitter.displayTestRunResults(finishedRuns, lostRuns)
	}
}

func (submitter *Submitter) writeThrottleFile(throttleFileName string, throttle int) error {
	var err error
	if throttleFileName != "" {
		// Throttle filename was specified. Lets use a throttle file.
		err = submitter.fileSystem.WriteTextFile(throttleFileName, strconv.Itoa(throttle))
		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_THROTTLE_FILE_WRITE, throttleFileName, err.Error())
		}
	}
	return err
}

func (submitter *Submitter) updateThrottleFromFileIfDifferent(
	throttleFileName string,
	currentThrottle int,
	wasThrottleFileLostAlready bool,
) (int, bool) {

	var newThrottle int = currentThrottle
	var isThrottleFileLost bool = wasThrottleFileLostAlready

	// Only bother with anything here if there is a throttle file specified by the user.
	if throttleFileName != "" {

		savedThrottle, err := submitter.readThrottleFile(throttleFileName)
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

func (submitter *Submitter) readThrottleFile(throttleFileName string) (int, error) {
	var savedThrottle int = 0
	var intermediateThrottle int64
	contents, err := submitter.fileSystem.ReadTextFile(throttleFileName)
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

func (submitter *Submitter) submitRun(
	groupName string,
	readyRuns []TestRun,
	submittedRuns map[string]*TestRun,
	lostRuns map[string]*TestRun,
	runOverrides *map[string]string, // This doesn't appear to be used. Why not ?
	trace bool,
	requestor string,
	requestType string,
) ([]TestRun, error) {

	var err error
	if len(readyRuns) >= 1 {

		nextRun := readyRuns[0]
		readyRuns = readyRuns[1:]

		className := nextRun.Bundle + "/" + nextRun.Class

		submitOverrides := make(map[string]interface{})

		for key, value := range nextRun.Overrides {
			submitOverrides[key] = value
		}

		var resultGroup *galasaapi.TestRuns
		log.Printf("submitRun - %s, %s", className, requestType)
		resultGroup, err = submitter.launcher.SubmitTestRun(groupName, className, requestType, requestor,
			nextRun.Stream, nextRun.Obr, trace, nextRun.GherkinUrl, nextRun.GherkinFeature, submitOverrides)
		if err != nil {
			log.Printf("Failed to submit test %v/%v - %v\n", nextRun.Bundle, nextRun.Class, err)
			lostRuns[className] = &nextRun
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_SUBMIT_TEST, nextRun.Bundle, nextRun.Class, err.Error())
		} else {
			if len(resultGroup.GetRuns()) < 1 {
				log.Printf("Lost the run attempting to submit test %v/%v\n", nextRun.Bundle, nextRun.Class)
				lostRuns[className] = &nextRun
				err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_TEST_NOT_IN_RUN_GROUP_LOST, nextRun.Bundle, nextRun.Class)
			}

			if err == nil {
				submittedRun := resultGroup.GetRuns()[0]
				nextRun.Name = *submittedRun.Name

				submittedRuns[nextRun.Name] = &nextRun

				if nextRun.GherkinUrl != "" {
					log.Printf("Run %v submitted - %v\n", nextRun.Name, nextRun.GherkinFeature)
				} else {
					log.Printf("Run %v submitted - %v/%v/%v\n", nextRun.Name, nextRun.Stream, nextRun.Bundle, nextRun.Class)
				}
			}
		}

	}
	return readyRuns, err
}

func (submitter *Submitter) runsFetchCurrentStatus(
	groupName string,
	submittedRuns map[string]*TestRun,
	finishedRuns map[string]*TestRun,
	lostRuns map[string]*TestRun,
	fetchRas bool) {

	currentGroup, err := submitter.launcher.GetRunsByGroup(groupName)
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

					var rasRun *galasaapi.Run
					rasRun, err = submitter.launcher.GetRunsById(*rasRunID)
					if err != nil {
						log.Printf("runsFetchCurrentStatus - Failed to retrieve RAS run for %v - %v\n", checkRun.Name, err)
					} else {
						checkRun.Tests = make([]TestMethod, 0)

						testStructure := rasRun.GetTestStructure()
						log.Printf("runsFetchCurrentStatus - testStructure- %v", testStructure)

						for _, testMethod := range testStructure.GetMethods() {
							test := TestMethod{
								Method: testMethod.GetMethodName(),
								Result: testMethod.GetResult(),
							}

							checkRun.Tests = append(checkRun.Tests, test)
						}
					}
				}

				if checkRun.GherkinUrl != "" {
					log.Printf("Run %v has finished(%v) - %v\n", runName, result, checkRun.GherkinFeature)
				} else {
					log.Printf("Run %v has finished(%v) - %v/%v/%v\n", runName, result, checkRun.Stream, checkRun.Bundle, checkRun.Class)
				}
			} else {
				// Check to see if there was a status change
				if checkRun.Status != currentRun.GetStatus() {
					checkRun.Status = currentRun.GetStatus()
					if checkRun.GherkinUrl != "" {
						log.Printf("    Run %v status is now '%v' - %v\n", runName, checkRun.Status, checkRun.GherkinFeature)
					} else {
						log.Printf("    Run %v status is now '%v' - %v/%v/%v\n", runName, checkRun.Status, checkRun.Stream, checkRun.Bundle, checkRun.Class)
					}
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

func (submitter *Submitter) createReports(params utils.RunsSubmitCmdValues,
	finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun) error {

	//convert TestRun tests into formattable data
	submitter.displayTestRunResults(finishedRuns, lostRuns)

	var err error
	if params.ReportYamlFilename != "" {
		err = ReportYaml(submitter.fileSystem, params.ReportYamlFilename, finishedRuns, lostRuns)
	}

	if err == nil {
		if params.ReportJsonFilename != "" {
			err = ReportJSON(submitter.fileSystem, params.ReportJsonFilename, finishedRuns, lostRuns)
		}
	}

	if err == nil {
		if params.ReportJunitFilename != "" {
			err = ReportJunit(submitter.fileSystem, params.ReportJunitFilename, params.GroupName, finishedRuns, lostRuns)
		}
	}

	return err
}

func (submitter *Submitter) displayTestRunResults(finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun) {
	var formatter = runsformatter.NewSummaryFormatter()
	var err error
	var outputText string

	formattableTest := FormattableTestFromTestRun(finishedRuns, lostRuns)
	outputText, err = formatter.FormatRuns(formattableTest)
	if err == nil {
		submitter.console.WriteString(outputText)
	}
}

func (submitter *Submitter) isRasDetailNeededForReports(params utils.RunsSubmitCmdValues) bool {

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

func (submitter *Submitter) buildListOfRunsToSubmit(portfolio *Portfolio, runOverrides map[string]string) []TestRun {
	log.Printf("buildListOfRunsToSubmit - portfolio %v, runOverrides %v", portfolio, runOverrides)
	readyRuns := make([]TestRun, 0, len(portfolio.Classes))
	currentUser := submitter.GetCurrentUserName()
	for _, portfolioTest := range portfolio.Classes {
		newTestrun := TestRun{
			Bundle:         portfolioTest.Bundle,
			Class:          portfolioTest.Class,
			Stream:         portfolioTest.Stream,
			Obr:            portfolioTest.Obr,
			QueuedTimeUTC:  submitter.timeService.Now().String(),
			Requestor:      currentUser,
			Status:         "queued",
			Overrides:      make(map[string]string, 0),
			GherkinUrl:     portfolioTest.GherkinUrl,
			GherkinFeature: submitter.getFeatureFromGherkinUrl(portfolioTest.GherkinUrl),
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
		if newTestrun.GherkinUrl == "" {
			log.Printf("Added test %v/%v/%v to the ready queue\n", newTestrun.Stream, newTestrun.Bundle, newTestrun.Class)
		} else {
			log.Printf("Added gherkin test %v to the ready queue\n", newTestrun.GherkinFeature)
		}
	}

	return readyRuns
}

func (submitter *Submitter) validateAndCorrectParams(
	params *utils.RunsSubmitCmdValues,
	submitSelectionFlags *utils.TestSelectionFlagValues,
) error {

	var err error

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
			params.GroupName = randomGenerator.NewString()
		}
		log.Printf("Using group name '%v' for test run submission\n", params.GroupName)

		_, err = submitter.checkIfGroupAlreadyInUse(params.GroupName)
	}

	if err == nil {
		err = submitter.correctOverrideFilePathParameter(params)
	}

	submitter.tildaExpandAllPaths(params)

	return err
}

func (submitter *Submitter) correctOverrideFilePathParameter(
	params *utils.RunsSubmitCmdValues,
) error {
	var err error
	// Correct the default overrideFile path if it wasn't specified.
	if params.OverrideFilePath == "" {

		params.OverrideFilePath = submitter.galasaHome.GetUrlFolderPath() + "/overrides.properties"
		var isFileThere bool
		isFileThere, err = submitter.fileSystem.Exists(params.OverrideFilePath)
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

func (submitter *Submitter) tildaExpandAllPaths(params *utils.RunsSubmitCmdValues) error {
	var err error

	params.OverrideFilePath, err = files.TildaExpansion(submitter.fileSystem, params.OverrideFilePath)

	if err == nil {
		params.PortfolioFileName, err = files.TildaExpansion(submitter.fileSystem, params.PortfolioFileName)
	}

	if err == nil {
		params.ReportJsonFilename, err = files.TildaExpansion(submitter.fileSystem, params.ReportJsonFilename)
	}

	if err == nil {
		params.ReportJunitFilename, err = files.TildaExpansion(submitter.fileSystem, params.ReportJunitFilename)
	}

	if err == nil {
		params.ReportYamlFilename, err = files.TildaExpansion(submitter.fileSystem, params.ReportYamlFilename)
	}

	if err == nil {
		params.ThrottleFileName, err = files.TildaExpansion(submitter.fileSystem, params.ThrottleFileName)
	}
	return err
}

func (submitter *Submitter) buildOverrideMap(commandParameters utils.RunsSubmitCmdValues) (map[string]string, error) {

	path := commandParameters.OverrideFilePath
	runOverrides, err := submitter.loadOverrideFile(path)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_LOAD_OVERRIDES_FILE, path, err.Error())
	} else {
		runOverrides, err = submitter.addOverridesFromCmdLine(runOverrides, commandParameters.Overrides)
	}

	return runOverrides, err
}

func (submitter *Submitter) loadOverrideFile(overrideFilePath string) (map[string]string, error) {

	var (
		overrides props.JavaProperties
		err       error = nil
	)

	if overrideFilePath == "-" {
		// Don't read properties from a file.
		overrides = make(map[string]string)
	} else {
		overrides, err = props.ReadPropertiesFile(submitter.fileSystem, overrideFilePath)
	}

	return overrides, err
}

func (submitter *Submitter) addOverridesFromCmdLine(overrides map[string]string, commandLineOverrides []string) (map[string]string, error) {
	var err error

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

func (submitter *Submitter) getPortfolio(portfolioFileName string, submitSelectionFlags *utils.TestSelectionFlagValues) (*Portfolio, error) {
	// Load the portfolio of tests
	var portfolio *Portfolio = nil
	var err error

	if portfolioFileName != "" {
		portfolio, err = ReadPortfolio(submitter.fileSystem, portfolioFileName)
	} else {
		// There is no portfolio file, so create an in-memory portfolio
		// from the tests we can find from the test selection.
		var testSelection TestSelection
		testSelection, err = SelectTests(submitter.launcher, submitSelectionFlags)

		if err == nil {
			testOverrides := make(map[string]string)
			portfolio = NewPortfolio()
			AddClassesToPortfolio(&testSelection, &testOverrides, portfolio)
		}
	}
	return portfolio, err
}

func (submitter *Submitter) GetCurrentUserName() string {
	userName := "cli"
	currentUser, err := user.Current()
	if err == nil {
		userName = currentUser.Username
	}
	return userName
}

func (submitter *Submitter) validatePortfolio(portfolio *Portfolio, portfolioFilename string) error {
	var err error
	if portfolio.Classes == nil || len(portfolio.Classes) < 1 {
		// Empty portfolio
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_EMPTY_PORTFOLIO, portfolioFilename)
	}
	return err
}

func (submitter *Submitter) checkIfGroupAlreadyInUse(groupName string) (bool, error) {
	isInUse := false
	var err error

	// Just check if it is already in use,  which is perfectly valid for custom group names
	var uuidCheck *galasaapi.TestRuns
	uuidCheck, err = submitter.launcher.GetRunsByGroup(groupName)
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

func (submitter *Submitter) getFeatureFromGherkinUrl(gherkinURL string) string {
	// split the Gherkin URL and select the last element from the array which should be the feature file name
	featureSlice := strings.Split(gherkinURL, "/")
	featureName := featureSlice[len(featureSlice)-1]
	// remove the .feature extension from the url
	featureName = strings.TrimSuffix(featureName, ".feature")
	return featureName
}
