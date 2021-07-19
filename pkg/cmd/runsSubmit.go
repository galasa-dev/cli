//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package cmd

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	satori "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/galasa.dev/cli/pkg/api"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

var (
    runsSubmitCmd = &cobra.Command{
            Use:   "submit",
            Short: "submit a list of tests to the ecosystem",
            Long:  "Submit a list of tests to the ecosystem, monitor them and wait for them to complete",
            Args: cobra.NoArgs,
            Run:   executeSubmit,
    }

    groupName            string
    throttle             *int 
    pollFlag             *int64
    progressFlag         *int
    submitFlagOverrides  *[]string
    trace                *bool
    requestor            string 
    reportYamlFilename   string 
    reportJsonFilename   string 
    reportJunitFilename  string 

    submitSelectionFlags = utils.TestSelectionFlags{}
)

type TestReport struct {
    Tests    []TestRun          `yaml:"tests" json:"tests"`
}

type TestRun struct {
    Name      string            `yaml:"name" json:"name"`
	Bundle    string            `yaml:"bundle" json:"bundle"`
	Class     string            `yaml:"class" json:"class"`
	Stream    string            `yaml:"stream" json:"stream"`
    Status    string            `yaml:"status" json:"status"`
    Result    string            `yaml:"result" json:"result"`
    Overrides map[string]string `yaml:"overrides" json:"overrides"`
    Tests     []TestMethod      `yaml:"tests" json:"tests"`
}

type TestMethod struct {
    Method    string     `yaml:"name" json:"name"`
    Result    string     `yaml:"result" json:"result"`
}

type JunitTestSuites struct {
    XMLName  xml.Name            `xml:"testsuites"`
    ID       string              `xml:"id,attr"`
    Name     string              `xml:"name,attr"`
    Tests    int                 `xml:"tests,attr"`
    Failures int                 `xml:"failures,attr"`
    Time     int                 `xml:"time,attr"`
    Testsuite []JunitTestSuite   `xml:"testsuite"`
} 

type JunitTestSuite struct {
    ID       string            `xml:"id,attr"`
    Name     string            `xml:"name,attr"`
    Tests    int               `xml:"tests,attr"`
    Failures int               `xml:"failures,attr"`
    Time     int               `xml:"time,attr"`
    TestCase []JunitTestCase   `xml:"testcase"`
}

type JunitTestCase struct {
    ID       string          `xml:"id,attr"`
    Name     string          `xml:"name,attr"`
    Time     int             `xml:"time,attr"`
    Failure  *JunitFailure   `xml:"failure"`
}

type JunitFailure struct {
    Message    string   `xml:"message,attr"`
    Type       string   `xml:"type,attr"`
}

func init() {
    runsSubmitCmd.Flags().StringVarP(&portfolioFilename, "portfolio", "p", "", "portfolio containing the tests to run")
    runsSubmitCmd.Flags().StringVar(&reportYamlFilename, "reportYaml", "", "yaml file to record the final results in")
    runsSubmitCmd.Flags().StringVar(&reportJsonFilename, "reportJson", "", "json file to record the final results in")
    runsSubmitCmd.Flags().StringVar(&reportJunitFilename, "reportJunit", "", "junit xml file to record the final results in")
    runsSubmitCmd.Flags().StringVarP(&groupName, "group", "g", "", "the group name to assign the test runs to, if not provided, a psuedo unique id will be generated")
    runsSubmitCmd.Flags().StringVar(&requestor, "requestor", "cli", "(temporary until authentication is enabled on the ecosystem) the requestor id to be associated with the test runs")
    pollFlag = runsSubmitCmd.Flags().Int64("poll", 30, "in seconds, how often the cli will poll the ecosystem for the status of the test runs")
    progressFlag = runsSubmitCmd.Flags().Int("progress", 5, "in minutes, how often the cli will report the overall progress of the test runs, -1 or less will disable progress reports")
    throttle = runsSubmitCmd.Flags().Int("throttle", 3, "how many test runs can be submitted in parallel, 0 or less will disable throttling")
	submitFlagOverrides = runsSubmitCmd.Flags().StringSlice("override", make([]string, 0), "overrides to be sent with the tests (overrides in the portfolio will take precedence)")
    trace = runsSubmitCmd.Flags().Bool("trace", false, "Trace to be enabled on the test runs")
    utils.AddCommandFlags(runsSubmitCmd, &submitSelectionFlags)

    runsCmd.AddCommand(runsSubmitCmd)
}

func executeSubmit(cmd *cobra.Command, args []string) {
    fmt.Println("Galasa CLI - Submit tests")

    // Set the poll time
    if *pollFlag < 1 {
        *pollFlag = 30
    }
    poll := time.Second * time.Duration(*pollFlag)

    // Set the progress time
    if *progressFlag < 0 {
        *progressFlag = int(^uint(0) >> 1) // set to maximum size of the int
    } else if *progressFlag == 0 {
        *progressFlag = 5
    }
    progress := time.Minute * time.Duration(*progressFlag)

    // Set the throttle
    if *throttle <= 0 {
        *throttle = int(^uint(0) >> 1) // set to maximum size of the int
    } 

    apiClient := api.InitialiseAPI(bootstrap)


    //  Dont mix portfolio and test selection on the same command

    if portfolioFilename != "" {
        if utils.AreSelectionFlagsProvided(&submitSelectionFlags) {
            fmt.Println("The submit command does not support mixing of the test selection flags and a portfolio")
            os.Exit(1)
        }
    } else {
        if !utils.AreSelectionFlagsProvided(&submitSelectionFlags) {
            fmt.Println("The submit command requires either test selection flags or a portfolio")
            os.Exit(1)
        }
    }

    // Convert overrides to a map
    runOverrides := make(map[string]string)
    for _, override := range *submitFlagOverrides {
        pos := strings.Index(override, "=")
        if (pos < 1) {
            fmt.Printf("Invalid override '%v'",override)
            os.Exit(1)
        }
        key := override[:pos]
        value := override[pos+1:]
        if value == "" {
            fmt.Printf("Invalid override '%v'",override)
            os.Exit(1)
        }
        runOverrides[key] = value
    }

    // Do we need to ask the RAS for the test structure
    fetchRas := false
    if reportYamlFilename != "" {
        fetchRas = true
    }
    if reportJsonFilename != "" {
        fetchRas = true
    }
    if reportJunitFilename != "" {
        fetchRas = true
    }

    // Load the portfolio of tests
    var portfolio utils.Portfolio
    if portfolioFilename != "" {
        portfolio = utils.LoadPortfolio(portfolioFilename)  
    } else {
        testSelection := utils.SelectTests(apiClient, &submitSelectionFlags)

        testOverrides := make(map[string]string)
        portfolio = utils.NewPortfolio()
        utils.CreatePortfolio(&testSelection, &testOverrides, &portfolio)
    }

    if portfolio.Classes == nil || len(portfolio.Classes) < 1 {
        fmt.Println("There are no tests in the test porfolio")
        os.Exit(1)
    }

    // generate a group name if required
    if groupName == "" {
        groupName = satori.NewV4().String()
    }

    fmt.Printf("Using group name '%v' for test run submission\n", groupName)

    // Just check if it is already in use,  which is perfectly valid for custom group names
    uuidCheck, _, err := apiClient.RunsAPIApi.GetRunsGroup(nil, groupName).Execute()
    if err != nil {
        panic(err)
    }

    if uuidCheck.Runs != nil && len(*uuidCheck.Runs) > 0 {
        fmt.Printf("Group name '%v' is aleady in use\n", groupName)
    }

    // Build list of runs to submit

    readyRuns := make([]TestRun, 0, len(portfolio.Classes))

    for _, portfolioTest := range portfolio.Classes {
        newTestrun := TestRun {
            Bundle: portfolioTest.Bundle,
            Class: portfolioTest.Class,
            Stream: portfolioTest.Stream,
            Status: "queued",
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

        fmt.Printf("Added test %v/%v/%v to the ready queue\n", newTestrun.Stream, newTestrun.Bundle, newTestrun.Class)
    }


    // 
    //
    // Main submit loop
    //
    //

    submittedRuns  := make(map[string]*TestRun)
    rerunRuns      := make(map[string]*TestRun)
    finishedRuns   := make(map[string]*TestRun)
    lostRuns       := make(map[string]*TestRun)

    nextProgressReport := time.Now().Add(progress)
    for (len(readyRuns) > 0 || len(submittedRuns) > 0 || len(rerunRuns) > 0) { // Loop whilst there are runs to submit or are running
        for (len(submittedRuns) < *throttle && len(readyRuns) > 0) {
            readyRuns = submitRun(apiClient, groupName, readyRuns, submittedRuns, lostRuns, &runOverrides)
        }

        now := time.Now()
        if now.After(nextProgressReport) {
            reportProgress(readyRuns, submittedRuns, finishedRuns, lostRuns)
            nextProgressReport = now.Add(progress)
        } 

        time.Sleep(poll)

        runsFetchCurrentStatus(apiClient, groupName, readyRuns, submittedRuns, finishedRuns, lostRuns, fetchRas)
    }


    runOk := report(finishedRuns, lostRuns)

    if reportYamlFilename != "" {
        reportYaml(finishedRuns, lostRuns)
    }
    if reportJsonFilename != "" {
        reportJSON(finishedRuns, lostRuns)
    }
    if reportJunitFilename != "" {
        reportJunit(finishedRuns, lostRuns)
    }

    if !runOk {
        fmt.Println("Not all runs passed, exiting with code 1")
        os.Exit(1)
    }
}

func submitRun(apiClient *galasaapi.APIClient, groupName string, readyRuns []TestRun, submittedRuns map[string]*TestRun, lostRuns map[string]*TestRun, runOverrides *map[string]string) []TestRun {

    if len(readyRuns) < 1 {
        return readyRuns
    }

    nextRun := readyRuns[0]
    readyRuns = readyRuns[1:]
    
    className := nextRun.Bundle + "/" + nextRun.Class
    classNames := []string{className}
    
    submitOverrides := make(map[string]interface{})
    
    for key,value := range nextRun.Overrides {
        submitOverrides[key] = value
    }

    testRunRequest := galasaapi.NewTestRunRequest()
    testRunRequest.SetClassNames(classNames)
    testRunRequest.SetRequestorType("CLI")
    testRunRequest.SetRequestor(requestor)
    testRunRequest.SetTestStream(nextRun.Stream)
    testRunRequest.SetTrace(*trace)
    testRunRequest.SetOverrides(submitOverrides)
 
    resultGroup, _, err := apiClient.RunsAPIApi.PostSubmitTestRuns(nil, groupName).TestRunRequest(*testRunRequest).Execute()
    if err != nil {
        fmt.Printf("Failed to submit test %v/%v - %v\n", nextRun.Bundle, nextRun.Class, err)
        lostRuns[className] = &nextRun
        return readyRuns
    }

    if len(resultGroup.GetRuns()) < 1 {
        fmt.Printf("Lost the run attempting to submit test %v/%v\n", nextRun.Bundle, nextRun.Class)
        lostRuns[className] = &nextRun
        return readyRuns
    }

    submittedRun := resultGroup.GetRuns()[0]
    nextRun.Name = *submittedRun.Name

    submittedRuns[nextRun.Name] = &nextRun

    fmt.Printf("Run %v submitted - %v/%v/%v\n", nextRun.Name, nextRun.Stream, nextRun.Bundle, nextRun.Class)

    return readyRuns
}

func runsFetchCurrentStatus(apiClient *galasaapi.APIClient, groupName string, readyRuns []TestRun, submittedRuns map[string]*TestRun, finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun, fetchRas bool) {
    currentGroup, _, err := apiClient.RunsAPIApi.GetRunsGroup(nil, groupName).Execute()
    if err != nil {
        fmt.Printf("Received error from group request - %v\n", err)
        return
    }
    
    // a copy to find lost runs
    checkRuns := copyTestRuns(submittedRuns)
  
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
                        fmt.Printf("Failed to retrieve RAS run for %v - %v\n", checkRun.Name, err)
                    } else {
                        checkRun.Tests = make([]TestMethod, 0)

                        testStructure := rasRun.GetTestStructure()
                        for _, testMethod := range testStructure.GetMethods() {
                            test := TestMethod {
                                Method: testMethod.GetMethodName(),
                                Result: testMethod.GetResult(),
                            }

                            checkRun.Tests = append(checkRun.Tests, test)
                        }
                    }
                }

                fmt.Printf("Run %v has finished(%v) - %v/%v/%v\n", runName, result, checkRun.Stream, checkRun.Bundle, checkRun.Class)
            } else {
                // Check to see if there was a status change
                if checkRun.Status != currentRun.GetStatus() {
                    checkRun.Status = currentRun.GetStatus()
                    fmt.Printf("    Run %v status is now '%v' - %v/%v/%v\n", runName, checkRun.Status, checkRun.Stream, checkRun.Bundle, checkRun.Class)
                }
            }
        }
    }

    // Now deal with the lost runs
    for runName, lostRun := range checkRuns {
        lostRuns[runName] = lostRun
        delete(submittedRuns, runName)
        fmt.Printf("Run %v was lost - %v/%v/%v\n", runName, lostRun.Stream, lostRun.Bundle, lostRun.Class)
    }
    
}

func report(finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun) bool {

    resultCounts := make(map[string]int, 0)

    resultCounts["Passed"] = 0
    resultCounts["Failed"] = 0

    for _, run := range finishedRuns {
        c, ok := resultCounts[run.Result]
        if !ok {
            resultCounts[run.Result] = 1
        } else {
            resultCounts[run.Result] = c + 1
        }
    }

    resultCounts["Lost"] = len(lostRuns)

    totalFailed := len(lostRuns)

    fmt.Println("***")
    fmt.Println("*** Final report")
    fmt.Println("*** ---------------")
    fmt.Println("***")
    fmt.Println("*** Passed test runs:-")
    found := false
    for runName, run := range finishedRuns {
        if strings.HasPrefix(run.Result, "Passed") {
            fmt.Printf("***     Run %v - %v/%v/%v\n", runName, run.Stream, run.Bundle, run.Class)
            found = true
        }
    }
    if !found {
        fmt.Println("***     None")
    }

    fmt.Println("***")
    fmt.Println("*** Failed test runs:-")
    found = false
    for runName, run := range finishedRuns {
        if strings.HasPrefix(run.Result, "Failed") {
            fmt.Printf("***     Run %v - %v/%v/%v\n", runName, run.Stream, run.Bundle, run.Class)
            found = true
            totalFailed = totalFailed + 1
        }
    }
    if !found {
        fmt.Println("***     None")
    }

    fmt.Println("***")
    fmt.Println("*** Other test runs:-")
    found = false
    for runName, run := range finishedRuns {
        if !strings.HasPrefix(run.Result, "Passed") && !strings.HasPrefix(run.Result, "Failed") {
            fmt.Printf("***     Run %v(%v) - %v/%v/%v\n", runName, run.Result, run.Stream, run.Bundle, run.Class)
            found = true
            totalFailed = totalFailed + 1
        }
    }
    if !found {
        fmt.Println("***     None")
    }
    fmt.Println("***")
    fmt.Print("*** results")
    for result, count := range resultCounts {
        fmt.Printf(", %v=%v", result, count)
    }
    fmt.Print("\n")


    if totalFailed > 0 {
        return false
    }

    return true
}

func reportProgress(readyRuns []TestRun, submittedRuns map[string]*TestRun, finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun) {
    ready := len(readyRuns)
    submitted := len(submittedRuns)
    finished := len(finishedRuns)
    lost := len(lostRuns)

    resultCounts := make(map[string]int, 0)

    for _, run := range finishedRuns {
        c, ok := resultCounts[run.Result]
        if !ok {
            resultCounts[run.Result] = 1
        } else {
            resultCounts[run.Result] = c + 1
        }
    }


    fmt.Println("***")
    fmt.Println("*** Progress report")
    fmt.Println("*** ---------------")
    for runName, run := range submittedRuns {
        fmt.Printf("***     Run %v is currently %v - %v/%v/%v\n", runName, run.Status, run.Stream, run.Bundle, run.Class)
    }
    fmt.Println("*** ----------------------------------------------------------------------------")
    fmt.Printf("*** run status, ready=%v, submitted=%v, finished=%v, lost=%v\n", ready, submitted, finished, lost)
    if len(resultCounts) > 0 {
        fmt.Print("*** results so far")
        for result, count := range resultCounts {
            fmt.Printf(", %v=%v", result, count)
        }
        fmt.Print("\n")
    }   
    fmt.Println("***")
}


func copyTestRuns(original map[string]*TestRun) map[string]*TestRun {
    new := make(map[string]*TestRun)
    for k,v := range original {
        new[k] = v
    }

    return new
}


func reportYaml(finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun) {

    var testReport TestReport
    testReport.Tests = make([]TestRun, 0)

    for _, run := range finishedRuns {
        testReport.Tests = append(testReport.Tests, *run)
    }

    for _, run := range lostRuns {
        testReport.Tests = append(testReport.Tests, *run)
    }

	file,err := os.Create(reportYamlFilename)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(file)

	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)

	err = encoder.Encode(&testReport)
	if err != nil {
		panic(err)
	}
	w.Flush()
	encoder.Close()
	file.Close()  
    
    fmt.Printf("Yaml test report written to %v\n", reportYamlFilename)
}

func reportJSON(finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun) {

    var testReport TestReport
    testReport.Tests = make([]TestRun, 0)

    for _, run := range finishedRuns {
        testReport.Tests = append(testReport.Tests, *run)
    }

    for _, run := range lostRuns {
        testReport.Tests = append(testReport.Tests, *run)
    }

    data, err := json.MarshalIndent(&testReport, "", " ")
    if err != nil {
        panic(err)
    }

    err = ioutil.WriteFile(reportJsonFilename, data, 0644)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Json test report written to %v\n", reportJsonFilename)
}


func reportJunit(finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun) {

    var testSuites JunitTestSuites
    testSuites.ID   = groupName
    testSuites.Name = "Galasa test run"
    testSuites.Tests = 0
    testSuites.Failures = 0
    testSuites.Time = 0
    testSuites.Testsuite = make([]JunitTestSuite, 0)

    for _, run := range finishedRuns {
        var testSuite JunitTestSuite

        testSuite.ID = run.Name
        testSuite.Name = run.Stream + "/" + run.Bundle + "/" + run.Class
        testSuite.TestCase = make([]JunitTestCase, 0)

        for _, method := range run.Tests {
            var testCase JunitTestCase
            testCase.ID = method.Method
            testCase.Name = method.Method

            testSuites.Tests = testSuites.Tests + 1
            testSuite.Tests = testSuite.Tests + 1
            if !strings.HasPrefix(method.Result, "Passed") {
                testSuites.Failures = testSuites.Failures + 1
                testSuite.Failures = testSuite.Failures + 1

                var failure JunitFailure
                failure.Message = "Failure messages are unavailable at this time"
                failure.Type    = "Unknown"

                testCase.Failure = &failure
            }
            
            testSuite.TestCase = append(testSuite.TestCase, testCase)
        }

        testSuites.Testsuite = append(testSuites.Testsuite, testSuite)
    }

    for range lostRuns {
        testSuites.Tests = testSuites.Tests + 1
        testSuites.Failures = testSuites.Failures + 1
    }


    data, err := xml.MarshalIndent(&testSuites, "", "    ")
    if err != nil {
        panic(err)
    }

    prologue := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n" + string(data))

    err = ioutil.WriteFile(reportJunitFilename, prologue, 0644)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Junit XML test report written to %v\n", reportJunitFilename)
}

