/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package launcher

import "github.com/galasa.dev/cli/pkg/galasaapi"

type MockLauncher struct {
}

func NewMockLauncher() *MockLauncher {
	return new(MockLauncher)
}

//-------------------------------------------------------------------
// Implementation of the launcher interface.
//-------------------------------------------------------------------

// GetRunsByGroup gets the lust of test runs for this groupName
func (launcher *MockLauncher) GetRunsByGroup(groupName string) (*galasaapi.TestRuns, error) {
	return newEmptyTestRun(), nil
}

func newEmptyTestRun() *galasaapi.TestRuns{
	isCompleteValue := false
	testRuns := new(galasaapi.TestRuns)
	testRuns.Complete = &isCompleteValue
	testRuns.Runs = []galasaapi.TestRun{}
	return testRuns
}

// SubmitTestRuns launch the test runs
func (launcher *MockLauncher) SubmitTestRuns(
	groupName string,
	classNames []string,
	requestType string,
	requestor string,
	stream string,
	isTraceEnabled bool,
	overrides map[string]interface{},
) (*galasaapi.TestRuns, error) {
	return newEmptyTestRun(), nil
}

// GetRunsById gets the Run information for the run with a specific run identifier
func (launcher *MockLauncher) GetRunsById(runId string) (*galasaapi.Run, error) {
	return &galasaapi.Run{}, nil
}

// GetStreams gets a list of streams available on this launcher
func (launcher *MockLauncher) GetStreams() ([]string, error) {
	return make([]string, 0), nil
}

// GetTestCatalog gets the test catalog for a given stream.
func (launcher *MockLauncher) GetTestCatalog(stream string) (TestCatalog, error) {
	return nil, nil
}
