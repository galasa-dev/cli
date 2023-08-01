/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package launcher

import (
	"github.com/galasa.dev/cli/pkg/galasaapi"
)

type TestCatalog map[string]interface{}

// ----------------------------------------------------------------------------------
// Launcher something which launches and monitors tests in some environment.
type Launcher interface {

	// GetRunsByGroup gets the lust of test runs for this groupName
	GetRunsByGroup(groupName string) (*galasaapi.TestRuns, error)

	// SubmitTestRuns launch the test runs
	SubmitTestRuns(
		groupName string,
		classNames []string,
		requestType string,
		requestor string,
		stream string,
		isTraceEnabled bool,
		overrides map[string]interface{},
	) (*galasaapi.TestRuns, error)

	// GetRunsById gets the Run information for the run with a specific run identifier
	GetRunsById(runId string) (*galasaapi.Run, error)

	// GetStreams gets a list of streams available on this launcher
	GetStreams() ([]string, error)

	// GetTestCatalog gets the test catalog for a given stream.
	GetTestCatalog(stream string) (TestCatalog, error)
}
