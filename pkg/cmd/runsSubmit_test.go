/*
 * Copyright contributors to the Galasa project
 */
package cmd

/*
func TestCanSubmitSmallPortfolio(t *testing.T) {
	// Given...
	mockFileSystem := utils.NewMockFileSystem()

	portfolioText := ` Some portfolio text. tbd.
	`

	mockFileSystem.WriteTextFile("small.portfolio", portfolioText)

	params := RunsSubmitCmdParameters{
		pollIntervalSeconds:           1,
		noExitCodeOnTestFailures:      true,
		progressReportIntervalMinutes: 1,
		throttle:                      1,
		trace:                         false,
		reportYamlFilename:            "a.yaml",
		reportJsonFilename:            "a.json",
		reportJunitFilename:           "a.junit.xml",
		groupName:                     "babe",
		portfolioFileName:             "small.portfolio",
	}

	mockTimeService := utils.NewMockTimeServiceAsMock()

	apiClient := &galasaapi.APIClient{}

	// When ...
	err := executeSubmitRemote(
		mockFileSystem,
		params,
		apiClient,
		mockTimeService)

	// Then...
	if err != nil {
		assert.Fail(t, "Not expecting error "+err.Error())
	}
	assert.Fail(t, "Exiting.")
}
*/
