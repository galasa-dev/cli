/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

func CountTotalFailedRuns(finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun) int {

	totalFailed := len(lostRuns)

	for _, run := range finishedRuns {
		// Anything which didn't pass failed by definition.
		if !strings.HasPrefix(run.Result, "Passed") {
			totalFailed = totalFailed + 1
		}
	}

	return totalFailed
}

// FinalHumanReadableReport - Creates a human readable report of how it went.
func FinalHumanReadableReport(finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun) {
	report := FinalHumanReadableReportAsString(finishedRuns, lostRuns)
	log.Print(report)
	fmt.Fprint(os.Stdout, report)
}

func FinalHumanReadableReportAsString(finishedRuns map[string]*TestRun, lostRuns map[string]*TestRun) string {

	resultCounts := make(map[string]int, 0)

	resultCounts["Passed"] = 0
	resultCounts["Failed"] = 0
	resultCounts["Passed With Defects"] = 0
	resultCounts["Failed With Defects"] = 0

	for _, run := range finishedRuns {
		c, ok := resultCounts[run.Result]
		if !ok {
			resultCounts[run.Result] = 1
		} else {
			resultCounts[run.Result] = c + 1
		}
	}

	resultCounts["Lost"] = len(lostRuns)

	var buff bytes.Buffer

	fmt.Fprintln(&buff, "***")
	fmt.Fprintln(&buff, "*** Final report")
	fmt.Fprintln(&buff, "*** ---------------")
	fmt.Fprintln(&buff, "***")
	fmt.Fprintln(&buff, "*** Passed test runs:-")
	found := false
	for runName, run := range finishedRuns {
		if strings.HasPrefix(run.Result, "Passed") && !strings.HasPrefix(run.Result, "Passed With Defects") {
			fmt.Fprintf(&buff, "***     Run %v - %v/%v/%v\n", runName, run.Stream, run.Bundle, run.Class)
			found = true
		}
	}
	if !found {
		fmt.Fprintln(&buff, "***     None")
	}

	fmt.Fprintln(&buff, "***")
	fmt.Fprintln(&buff, "*** Failed test runs:-")
	found = false
	for runName, run := range finishedRuns {
		if strings.HasPrefix(run.Result, "Failed") && !strings.HasPrefix(run.Result, "Failed With Defects") {
			fmt.Fprintf(&buff, "***     Run %v - %v/%v/%v\n", runName, run.Stream, run.Bundle, run.Class)
			found = true
		}
	}
	if !found {
		fmt.Fprintln(&buff, "***     None")
	}

	fmt.Fprintln(&buff, "***")
	fmt.Fprintln(&buff, "*** Passed With Defects test runs:-")
	found = false
	for runName, run := range finishedRuns {
		if strings.HasPrefix(run.Result, "Passed With Defects") {
			fmt.Fprintf(&buff, "***     Run %v - %v/%v/%v\n", runName, run.Stream, run.Bundle, run.Class)
			found = true
		}
	}
	if !found {
		fmt.Fprintln(&buff, "***     None")
	}

	fmt.Fprintln(&buff, "***")
	fmt.Fprintln(&buff, "*** Failed With Defects test runs:-")
	found = false
	for runName, run := range finishedRuns {
		if strings.HasPrefix(run.Result, "Failed With Defects") {
			log.Printf("***     Run %v - %v/%v/%v\n", runName, run.Stream, run.Bundle, run.Class)
			found = true
		}
	}
	if !found {
		fmt.Fprintln(&buff, "***     None")
	}

	fmt.Fprintln(&buff, "***")
	fmt.Fprintln(&buff, "*** Other test runs:-")
	found = false
	for runName, run := range finishedRuns {
		if !strings.HasPrefix(run.Result, "Passed") && !strings.HasPrefix(run.Result, "Failed") {
			fmt.Fprintf(&buff, "***     Run %v(%v) - %v/%v/%v\n", runName, run.Result, run.Stream, run.Bundle, run.Class)
			found = true
		}
	}
	if !found {
		fmt.Fprintln(&buff, "***     None")
	}
	fmt.Fprintln(&buff, "***")
	fmt.Fprintln(&buff, "*** results")
	resultsSoFar := "*** results "
	for result, count := range resultCounts {
		resultsSoFar = resultsSoFar + fmt.Sprintf(", %v=%v", result, count)
	}
	fmt.Fprintln(&buff, resultsSoFar)
	return buff.String()
}

func InterrimProgressReport(
	readyRuns []TestRun,
	submittedRuns map[string]*TestRun,
	finishedRuns map[string]*TestRun,
	lostRuns map[string]*TestRun,
	throttle int) {

	report := InterrimProgressReportAsString(readyRuns, submittedRuns, finishedRuns, lostRuns, throttle)
	log.Print(report)
}

func InterrimProgressReportAsString(
	readyRuns []TestRun,
	submittedRuns map[string]*TestRun,
	finishedRuns map[string]*TestRun,
	lostRuns map[string]*TestRun,
	throttle int) string {

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

	var buff bytes.Buffer

	fmt.Fprintln(&buff, "***")
	fmt.Fprintln(&buff, "*** Progress report")
	fmt.Fprintln(&buff, "*** ---------------")
	for runName, run := range submittedRuns {
		log.Printf("***     Run %v is currently %v - %v/%v/%v\n", runName, run.Status, run.Stream, run.Bundle, run.Class)
	}
	fmt.Fprintln(&buff, "*** ----------------------------------------------------------------------------")
	fmt.Fprintf(&buff, "*** run status, ready=%v, submitted=%v, finished=%v, lost=%v\n", ready, submitted, finished, lost)
	fmt.Fprintf(&buff, "*** throttle=%v\n", throttle)
	if len(resultCounts) > 0 {
		resultsSoFar := "*** results so far"
		for result, count := range resultCounts {
			resultsSoFar = resultsSoFar + fmt.Sprintf(", %v=%v", result, count)
		}
		fmt.Fprintln(&buff, resultsSoFar)
	}
	fmt.Fprintln(&buff, "***")
	return buff.String()
}
