/*
 * Copyright contributors to the Galasa project
 */
package launcher

import (
	"log"
	"testing"

	"github.com/galasa.dev/cli/pkg/embedded"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateAJVMLauncher(t *testing.T) {

	env := utils.NewMockEnv()
	env.EnvVars["JAVA_HOME"] = "/java"

	fileSystem := utils.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fileSystem, "/java")

	jvmLaunchParams := getBasicJvmLaunchParams()
	timeService := utils.NewMockTimeService()

	mockProcess := NewMockProcess()
	mockProcessFactory := NewMockProcessFactory(mockProcess)

	launcher, err := NewJVMLauncher(
		env, fileSystem, embedded.GetEmbeddedFileSystem(),
		jvmLaunchParams, timeService, mockProcessFactory)
	if err != nil {
		assert.Fail(t, "Constructor should not have failed but it did. error:%s", err.Error())
	}
	assert.NotNil(t, launcher, "Launcher reference was nil, shouldn't have been.")
}

func getBasicJvmLaunchParams() RunsSubmitLocalCmdParameters {
	return RunsSubmitLocalCmdParameters{
		Obrs:                nil,
		RemoteMaven:         "",
		TargetGalasaVersion: "",
	}
}

func TestCantCreateAJVMLauncherIfJVMHomeNotSet(t *testing.T) {

	env := utils.NewMockEnv()
	// env.EnvVars["JAVA_HOME"] = "/java"

	fileSystem := utils.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fileSystem, "/java")

	jvmLaunchParams := getBasicJvmLaunchParams()
	timeService := utils.NewMockTimeService()

	mockProcess := NewMockProcess()
	mockProcessFactory := NewMockProcessFactory(mockProcess)

	launcher, err := NewJVMLauncher(
		env, fileSystem, embedded.GetEmbeddedFileSystem(),
		jvmLaunchParams, timeService, mockProcessFactory)
	if err == nil {
		assert.Fail(t, "Constructor should have failed but it did not.")
	}
	assert.Nil(t, launcher, "Launcher reference was not nil")
	assert.Contains(t, err.Error(), "GAL1050E")
}

func TestCanCreateJvmLauncher(t *testing.T) {
	env := utils.NewMockEnv()
	env.EnvVars["JAVA_HOME"] = "/java"

	fileSystem := utils.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fileSystem, "/java")

	jvmLaunchParams := getBasicJvmLaunchParams()
	timeService := utils.NewMockTimeService()
	mockProcess := NewMockProcess()
	mockProcessFactory := NewMockProcessFactory(mockProcess)

	launcher, err := NewJVMLauncher(
		env, fileSystem, embedded.GetEmbeddedFileSystem(),
		jvmLaunchParams, timeService, mockProcessFactory)

	if err != nil {
		assert.Fail(t, "JVM launcher should have been creatable.")
	}
	assert.NotNil(t, launcher, "Launcher returned is nil!")
}

func TestCanLaunchLocalJvmTest(t *testing.T) {
	// Given...
	env := utils.NewMockEnv()
	env.EnvVars["JAVA_HOME"] = "/java"

	fileSystem := utils.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fileSystem, "/java")

	jvmLaunchParams := getBasicJvmLaunchParams()
	timeService := utils.NewMockTimeService()

	mockProcess := NewMockProcess()
	mockProcessFactory := NewMockProcessFactory(mockProcess)

	launcher, err := NewJVMLauncher(
		env, fileSystem, embedded.GetEmbeddedFileSystem(),
		jvmLaunchParams, timeService, mockProcessFactory)

	if err != nil {
		assert.Fail(t, "JVM launcher should have been creatable.")
	}
	assert.NotNil(t, launcher, "Launcher returned is nil!")

	isTraceEnabled := true
	var overrides map[string]interface{} = make(map[string]interface{})

	// When...
	testRuns, err := launcher.SubmitTestRuns(
		"myGroup",
		[]string{"galasa.dev.example.banking.account/galasa.dev.example.banking.account.TestAccount"},
		"myRequestType-UnitTest",
		"myRequestor",
		"unitTestStream",
		isTraceEnabled,
		overrides,
	)
	if err != nil {
		assert.Fail(t, "Launcher should have launched command OK")
	}
	assert.NotNil(t, testRuns, "Returned test runs is nil, should have been an object.")

	assert.Len(t, testRuns.Runs, 1, "Returned test runs array doesn't contain correct number of tests launched.")
	assert.False(t, *testRuns.Complete, "Returned test runs should not already be complete")
}

func TestCanGetRunGroupStatus(t *testing.T) {
	// Given...
	env := utils.NewMockEnv()
	env.EnvVars["JAVA_HOME"] = "/java"

	fileSystem := utils.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fileSystem, "/java")

	jvmLaunchParams := getBasicJvmLaunchParams()
	timeService := utils.NewMockTimeService()

	mockProcess := NewMockProcess()
	mockProcessFactory := NewMockProcessFactory(mockProcess)

	launcher, err := NewJVMLauncher(
		env, fileSystem, embedded.GetEmbeddedFileSystem(),
		jvmLaunchParams, timeService, mockProcessFactory)
	if err != nil {
		assert.Fail(t, "Launcher should have launched command OK")
	}

	isTraceEnabled := true
	var overrides map[string]interface{} = make(map[string]interface{})

	launcher.SubmitTestRuns(
		"myGroup",
		[]string{"galasa.dev.example.banking.account/galasa.dev.example.banking.account.TestAccount"},
		"myRequestType-UnitTest",
		"myRequestor",
		"unitTestStream",
		isTraceEnabled,
		overrides,
	)

	// Wait for the child process to complete...
	mockProcess.Wait()

	// Simulate the test writing some final status to disk.
	structureJsonContent := `
	{
		"runName": "U2",
		"bundle": "dev.galasa.example.banking.account",
		"testName": "dev.galasa.example.banking.account.TestAccountExtended",
		"testShortName": "TestAccountExtended",
		"requestor": "unknown",
		"status": "finished",
		"result": "Passed",
		"queued": "2023-02-17T16:24:52.041118Z",
		"startTime": "2023-02-17T16:24:52.068591Z",
		"endTime": "2023-02-17T16:24:52.268396Z",
		"methods": [
		  {
			"className": "dev.galasa.example.banking.account.TestAccountExtended",
			"methodName": "simpleSampleTest",
			"type": "Test",
			"befores": [],
			"afters": [],
			"status": "finished",
			"result": "Passed",
			"runLogStart": 0,
			"runLogEnd": 0,
			"startTime": "2023-02-17T16:24:52.238868Z",
			"endTime": "2023-02-17T16:24:52.263756Z"
		  },
		  {
			"className": "dev.galasa.example.banking.account.TestAccountExtended",
			"methodName": "testRetrieveBundleResourceFileAsStringMethod",
			"type": "Test",
			"befores": [],
			"afters": [],
			"status": "finished",
			"result": "Passed",
			"runLogStart": 0,
			"runLogEnd": 0,
			"startTime": "2023-02-17T16:24:52.264511Z",
			"endTime": "2023-02-17T16:24:52.265325Z"
		  }
		]
	}`
	fileSystem.WriteTextFile("/temp/ras/L12345/structure.json", structureJsonContent)

	// When...
	testRuns, err := launcher.GetRunsByGroup("myGroup")

	// Then...
	if err != nil {
		assert.Fail(t, "Launcher should have returned some test status")
	}
	assert.NotNil(t, testRuns, "Returned test runs status is nil, should have been an object.")

	assert.Len(t, testRuns.Runs, 1, "Returned test runs array doesn't contain correct number of tests launched.")
	log.Printf("getRunsByGroup returned *testRUns.Complete of %v", *testRuns.Complete)
	if !*testRuns.Complete {
		assert.Fail(t, "Returned test runs should all be marked as complete")
	}
}

func TestJvmLauncherSetsRASStoreOverride(t *testing.T) {
	overrides := make(map[string]interface{})
	fs := utils.NewMockFileSystem()
	overridesGotBack, err := addStandardProperties(fs, overrides)
	assert.Nil(t, err)
	assert.Contains(t, overridesGotBack, "framework.resultarchive.store")
}

func TestCanCreateTempPropsFile(t *testing.T) {
	overrides := make(map[string]interface{})
	fs := utils.NewMockFileSystem()

	// When
	tempFolder, tempPropsFile, err := prepareTempFiles(fs, overrides)

	// Then the temp folder should exist.
	assert.Nil(t, err)
	assert.NotEmpty(t, tempFolder)
	exists, err := fs.DirExists(tempFolder)
	assert.Nil(t, err)
	assert.True(t, exists)

	// The temp property file should exist
	assert.NotEmpty(t, tempPropsFile)
	overridesGotBack, err := utils.ReadPropertiesFile(fs, tempPropsFile)
	assert.Nil(t, err)
	assert.Contains(t, overridesGotBack, "framework.resultarchive.store")
}
