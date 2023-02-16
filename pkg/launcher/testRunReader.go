/*
 * Copyright contributors to the Galasa project
 */
package launcher

import (
	"encoding/json"
	"net/url"

	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"

	"log"
)

// On disk, we have a file like this:
// {
// 	"runName": "U527",
// 	"bundle": "dev.galasa.example.banking.account",
// 	"testName": "dev.galasa.example.banking.account.TestAccount",
// 	"testShortName": "TestAccount",
// 	"requestor": "mcobbett",
// 	"status": "finished",
// 	"result": "Passed",
// 	"queued": "2023-02-14T14:16:42.854396Z",
// 	"startTime": "2023-02-14T14:16:42.881464Z",
// 	"endTime": "2023-02-14T14:16:43.082639Z",
// 	"methods": [
// 	  {
// 		"className": "dev.galasa.example.banking.account.TestAccount",
// 		"methodName": "simpleSampleTest",
// 		"type": "Test",
// 		"befores": [],
// 		"afters": [],
// 		"status": "finished",
// 		"result": "Passed",
// 		"runLogStart": 0,
// 		"runLogEnd": 0,
// 		"startTime": "2023-02-14T14:16:43.063018Z",
// 		"endTime": "2023-02-14T14:16:43.076501Z"
// 	  }
// 	]
// }
//
//
// Unfortunately, it doesn't serialise into the structure we need. Namely
// galasaapi.TestRun
// So we have to do some marshalling to get the TestRun object we want.
//

// Read a test run file from disk into the in-memory structure.
func readTestRunFromJsonFile(
	fileSystem utils.FileSystem,
	jsonFilePath string,
) (*galasaapi.TestRun, error) {

	testRunData := &galasaapi.TestRun{}
	url, err := url.Parse(jsonFilePath)
	if err == nil {
		var jsonContent string
		jsonContent, err = fileSystem.ReadTextFile(url.Path)
		if err == nil && len(jsonContent) > 0 {
			var f interface{}
			err = json.Unmarshal([]byte(jsonContent), &f)

			if err == nil {
				fields := f.(map[string]interface{})

				testRunData.Name = getStringField(fields, "runName")
				testRunData.Test = getStringField(fields, "testShortName")
				testRunData.BundleName = getStringField(fields, "bundle")
				testRunData.TestName = getStringField(fields, "testName")
				testRunData.Status = getStringField(fields, "status")
				testRunData.Result = getStringField(fields, "result")
				testRunData.Queued = getStringField(fields, "queued")
				testRunData.Requestor = getStringField(fields, "requestor")
				testRunData.RasRunId = getStringField(fields, "runName")

				log.Printf("Test %s status %s result %s\n", *testRunData.Name, *testRunData.Status, *testRunData.Result)
				// testRunData.Stream = stream
				// testRunData.Local = isLocal
				// testRunData.Trace = isTraceEnabled
				// testRunData.Type = testType
				// testRunData.Group = group
			}
		}
	}
	return testRunData, err
}

func getStringField(fields map[string]interface{}, fieldName string) *string {
	var strValue = ""
	value := fields[fieldName]
	if value != nil {
		strValue = value.(string)
	}
	return &strValue
}
