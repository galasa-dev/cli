/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"encoding/xml"
	"log"
	"strings"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/files"
)

type JunitTestSuites struct {
	XMLName   xml.Name         `xml:"testsuites"`
	ID        string           `xml:"id,attr"`
	Name      string           `xml:"name,attr"`
	Tests     int              `xml:"tests,attr"`
	Failures  int              `xml:"failures,attr"`
	Time      int              `xml:"time,attr"`
	Testsuite []JunitTestSuite `xml:"testsuite"`
}

type JunitTestSuite struct {
	ID       string          `xml:"id,attr"`
	Name     string          `xml:"name,attr"`
	Tests    int             `xml:"tests,attr"`
	Failures int             `xml:"failures,attr"`
	Time     int             `xml:"time,attr"`
	TestCase []JunitTestCase `xml:"testcase"`
}

type JunitTestCase struct {
	ID      string        `xml:"id,attr"`
	Name    string        `xml:"name,attr"`
	Time    int           `xml:"time,attr"`
	Failure *JunitFailure `xml:"failure"`
}

type JunitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
}

func ReportJunit(
	fileSystem files.FileSystem,
	reportJunitFilename string,
	groupName string,
	finishedRuns map[string]*TestRun,
	lostRuns map[string]*TestRun) error {

	var testSuites JunitTestSuites
	testSuites.ID = groupName
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
				failure.Type = "Unknown"

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
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_REPORT_JUNIT_PREPARE, reportJunitFilename, err.Error())
	} else {

		prologue := "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n" + string(data)

		err = fileSystem.WriteTextFile(reportJunitFilename, prologue)
		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SUBMIT_REPORT_JUNIT_WRITE_FAIL, reportJunitFilename, err.Error())
		} else {
			log.Printf("Junit XML test report written to %v\n", reportJunitFilename)
		}
	}
	return err
}
