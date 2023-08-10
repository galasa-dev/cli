/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package launcher

import (
	"log"
	"strconv"
	"testing"

	"github.com/galasa.dev/cli/pkg/api"
	"github.com/galasa.dev/cli/pkg/files"
	"github.com/galasa.dev/cli/pkg/props"

	"github.com/galasa.dev/cli/pkg/embedded"
	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCanCreateAJVMLauncher(t *testing.T) {

	env := utils.NewMockEnv()
	env.EnvVars["JAVA_HOME"] = "/java"

	fs := files.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fs, "/java")

	galasaHome, _ := utils.NewGalasaHome(fs, env, "")

	jvmLaunchParams := getBasicJvmLaunchParams()
	timeService := utils.NewMockTimeService()

	mockProcess := NewMockProcess()
	mockProcessFactory := NewMockProcessFactory(mockProcess)

	bootstrapProps := getBasicBootstrapProperties()

	launcher, err := NewJVMLauncher(
		bootstrapProps, env, fs, embedded.GetReadOnlyFileSystem(),
		jvmLaunchParams, timeService, mockProcessFactory, galasaHome)
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

func getBasicBootstrapProperties() props.JavaProperties {
	props := props.JavaProperties{}
	props[api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_OPTIONS] = "-Xmx80m"
	return props
}

func TestCantCreateAJVMLauncherIfJVMHomeNotSet(t *testing.T) {

	env := utils.NewMockEnv()
	// env.EnvVars["JAVA_HOME"] = "/java"

	fs := files.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fs, "/java")

	galasaHome, _ := utils.NewGalasaHome(fs, env, "")

	bootstrapProps := getBasicBootstrapProperties()

	jvmLaunchParams := getBasicJvmLaunchParams()
	timeService := utils.NewMockTimeService()

	mockProcess := NewMockProcess()
	mockProcessFactory := NewMockProcessFactory(mockProcess)

	launcher, err := NewJVMLauncher(
		bootstrapProps, env, fs, embedded.GetReadOnlyFileSystem(),
		jvmLaunchParams, timeService, mockProcessFactory, galasaHome)
	if err == nil {
		assert.Fail(t, "Constructor should have failed but it did not.")
	}
	assert.Nil(t, launcher, "Launcher reference was not nil")
	assert.Contains(t, err.Error(), "GAL1050E")
}

func TestCanCreateJvmLauncher(t *testing.T) {
	env := utils.NewMockEnv()
	env.EnvVars["JAVA_HOME"] = "/java"

	fs := files.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fs, "/java")

	jvmLaunchParams := getBasicJvmLaunchParams()
	timeService := utils.NewMockTimeService()
	mockProcess := NewMockProcess()
	mockProcessFactory := NewMockProcessFactory(mockProcess)
	galasaHome, _ := utils.NewGalasaHome(fs, env, "")

	bootstrapProps := getBasicBootstrapProperties()

	launcher, err := NewJVMLauncher(
		bootstrapProps, env, fs, embedded.GetReadOnlyFileSystem(),
		jvmLaunchParams, timeService, mockProcessFactory, galasaHome)

	if err != nil {
		assert.Fail(t, "JVM launcher should have been creatable.")
	}
	assert.NotNil(t, launcher, "Launcher returned is nil!")
}

func TestCanLaunchLocalJvmTest(t *testing.T) {
	// Given...
	env := utils.NewMockEnv()
	env.EnvVars["JAVA_HOME"] = "/java"

	fs := files.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fs, "/java")

	galasaHome, _ := utils.NewGalasaHome(fs, env, "")

	jvmLaunchParams := getBasicJvmLaunchParams()
	timeService := utils.NewMockTimeService()

	mockProcess := NewMockProcess()
	mockProcessFactory := NewMockProcessFactory(mockProcess)

	bootstrapProps := getBasicBootstrapProperties()

	launcher, err := NewJVMLauncher(
		bootstrapProps, env, fs, embedded.GetReadOnlyFileSystem(),
		jvmLaunchParams, timeService, mockProcessFactory, galasaHome)

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

	fs := files.NewMockFileSystem()
	utils.AddJavaRuntimeToMock(fs, "/java")

	galasaHome, _ := utils.NewGalasaHome(fs, env, "")

	jvmLaunchParams := getBasicJvmLaunchParams()
	timeService := utils.NewMockTimeService()

	mockProcess := NewMockProcess()
	mockProcessFactory := NewMockProcessFactory(mockProcess)

	bootstrapProps := getBasicBootstrapProperties()

	launcher, err := NewJVMLauncher(
		bootstrapProps, env, fs, embedded.GetReadOnlyFileSystem(),
		jvmLaunchParams, timeService, mockProcessFactory, galasaHome)
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
	fs.WriteTextFile("/temp/ras/L12345/structure.json", structureJsonContent)

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
	fs := files.NewMockFileSystem()
	env := utils.NewMockEnv()
	galasaHome, _ := utils.NewGalasaHome(fs, env, "")

	overridesGotBack := addStandardOverrideProperties(galasaHome, fs, overrides)
	assert.Contains(t, overridesGotBack, "framework.resultarchive.store")
}

func TestJvmLauncherSets3270TerminalOutputFormatProperty(t *testing.T) {
	overrides := make(map[string]interface{})
	fs := files.NewMockFileSystem()
	env := utils.NewMockEnv()
	galasaHome, _ := utils.NewGalasaHome(fs, env, "")

	overridesGotBack := addStandardOverrideProperties(galasaHome, fs, overrides)

	assert.Contains(t, overridesGotBack, "zos3270.terminal.output")
	assert.Contains(t, overridesGotBack["zos3270.terminal.output"], "png")
	assert.Contains(t, overridesGotBack["zos3270.terminal.output"], "json")
}

func TestCanCreateTempPropsFile(t *testing.T) {
	overrides := make(map[string]interface{})
	fs := files.NewMockFileSystem()
	env := utils.NewMockEnv()
	galasaHome, _ := utils.NewGalasaHome(fs, env, "")

	// When
	tempFolder, tempPropsFile, err := prepareTempFiles(galasaHome, fs, overrides)

	// Then the temp folder should exist.
	assert.Nil(t, err)
	assert.NotEmpty(t, tempFolder)
	exists, err := fs.DirExists(tempFolder)
	assert.Nil(t, err)
	assert.True(t, exists)

	// The temp property file should exist
	assert.NotEmpty(t, tempPropsFile)
	overridesGotBack, err := props.ReadPropertiesFile(fs, tempPropsFile)
	assert.Nil(t, err)
	assert.Contains(t, overridesGotBack, "framework.resultarchive.store")
	assert.Contains(t, overridesGotBack, "framework.request.type.LOCAL.prefix")
}

func getDefaultCommandSyntaxTestParameters() (
	props.JavaProperties,
	utils.Environment,
	utils.GalasaHome,
	*files.MockFileSystem,
	string,
	[]utils.MavenCoordinates,
	TestLocation,
	string,
	string,
	string,
	bool,
) {
	bootstrapProps := getBasicBootstrapProperties()
	fs := files.NewOverridableMockFileSystem()
	javaHome := "my_java_home"
	testObrs := make([]utils.MavenCoordinates, 0)
	testObrs = append(
		testObrs,
		utils.MavenCoordinates{
			GroupId:    "myGroup",
			ArtifactId: "myArtifact",
			Version:    "0.2",
			Classifier: "myClassifier",
		},
	)
	testLocation := TestLocation{
		OSGiBundleName:         "myTestBundle",
		QualifiedJavaClassName: "myClass",
	}
	remoteMaven := "myRemoteMaven"
	galasaVersionToRun := "0.99.0"
	overridesFilePath := "C:/myFolder/myOverrides.props"
	isTraceEnabled := true

	env := utils.NewMockEnv()
	galasaHome, _ := utils.NewGalasaHome(fs, env, "")

	return bootstrapProps, env, galasaHome, fs, javaHome, testObrs, testLocation,
		remoteMaven, galasaVersionToRun, overridesFilePath, isTraceEnabled
}

func TestCommandIncludesTraceWhenTraceIsEnabled(t *testing.T) {

	bootstrapProps, _, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := true
	isDebugEnabled := false
	var debugPort uint32 = 0
	debugMode := ""

	cmd, args, err := getCommandSyntax(
		bootstrapProps,
		galasaHome,
		fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, cmd)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.Contains(t, args, "--trace")
}

func TestCommandDoesNotIncludeTraceWhenTraceIsDisabled(t *testing.T) {
	bootstrapProps, _, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := false
	isDebugEnabled := false
	var debugPort uint32 = 0
	debugMode := ""

	cmd, args, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, cmd)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.NotContains(t, args, "--trace")
}

func TestCommandSyntaxContainsJavaHomeUnixSlashes(t *testing.T) {
	bootstrapProps, _, galasaHome, fs,
		_,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled := getDefaultCommandSyntaxTestParameters()

	javaHome := "myJavaHome"
	fs.SetFilePathSeparator("/")
	isDebugEnabled := false
	var debugPort uint32 = 0
	debugMode := ""

	cmd, args, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, cmd)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.Equal(t, cmd, "myJavaHome/bin/java")
}

func TestCommandSyntaxContainsJavaHomeWindowsSlashes(t *testing.T) {
	bootstrapProps, _, galasaHome, fs,
		_,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled := getDefaultCommandSyntaxTestParameters()

	javaHome := "myJavaHome"
	fs.SetFilePathSeparator("\\")
	fs.SetExecutableExtension(".exe")

	isDebugEnabled := false
	var debugPort uint32 = 0
	debugMode := ""

	cmd, args, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, cmd)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.Equal(t, cmd, "myJavaHome\\bin\\java")
}

func TestSingleValidObrIsValid(t *testing.T) {
	obrInputs := []string{"mvn:dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOT/obr"}
	mavenCoordinates, err := validateObrs(obrInputs)
	assert.Nil(t, err)
	assert.Len(t, mavenCoordinates, 1)
	assert.NotNil(t, mavenCoordinates)
	assert.Equal(t, mavenCoordinates[0].ArtifactId, "dev.galasa.example.banking.obr")
	assert.Equal(t, mavenCoordinates[0].Classifier, "obr")
	assert.Equal(t, mavenCoordinates[0].GroupId, "dev.galasa.example.banking")
	assert.Equal(t, mavenCoordinates[0].Version, "0.0.1-SNAPSHOT")
}

func TestSingleObrIsInvalidTooFewPartsWithSlashSeparator(t *testing.T) {
	obrInputs := []string{"mvn:dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOTobr"}
	mavenCoordinates, err := validateObrs(obrInputs)
	assert.NotNil(t, err)
	assert.NotNil(t, mavenCoordinates)
	assert.Len(t, mavenCoordinates, 0)
	assert.Contains(t, err.Error(), "GAL1060E")
}

func TestSingleObrIsInvalidTooManyPartsWithSlashSeparator(t *testing.T) {
	obrInputs := []string{"mvn:dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOT//obr"}
	mavenCoordinates, err := validateObrs(obrInputs)
	assert.NotNil(t, err)
	assert.NotNil(t, mavenCoordinates)
	assert.Len(t, mavenCoordinates, 0)
	assert.Contains(t, err.Error(), "GAL1061E")
}

func TestSingleObrIsInvalidTooManyPartsWithMissingMvnPrefix(t *testing.T) {
	obrInputs := []string{"dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOT/obr"}
	mavenCoordinates, err := validateObrs(obrInputs)
	assert.NotNil(t, err)
	assert.NotNil(t, mavenCoordinates)
	assert.Len(t, mavenCoordinates, 0)
	assert.Contains(t, err.Error(), "GAL1062E")
}

func TestSingleObrIsInvalidTooManyPartsWithMissingObrSuffix(t *testing.T) {
	obrInputs := []string{"mvn:dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOT/mysuffix"}
	mavenCoordinates, err := validateObrs(obrInputs)
	assert.NotNil(t, err)
	assert.NotNil(t, mavenCoordinates)
	assert.Len(t, mavenCoordinates, 0)
	assert.Contains(t, err.Error(), "GAL1063E")
}

func TestValidClassInputGetsSplitCorrectly(t *testing.T) {
	userInput := "myBundle/myClass"
	testLocation, err := classNameUserInputToTestClassLocation(userInput)
	assert.NotNil(t, testLocation)
	assert.Nil(t, err)
	assert.Equal(t, testLocation.OSGiBundleName, "myBundle")
	assert.Equal(t, testLocation.QualifiedJavaClassName, "myClass")
}

func TestInvalidClassInputNoSlashGetsError(t *testing.T) {
	userInput := "myBundleNoSlashmyClass"
	testLocation, err := classNameUserInputToTestClassLocation(userInput)
	assert.Nil(t, testLocation)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1064E")
}

func TestInvalidClassInputManySlashesGetsError(t *testing.T) {
	userInput := "myBundle/With/More/Slashes/Class"
	testLocation, err := classNameUserInputToTestClassLocation(userInput)
	assert.Nil(t, testLocation)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1065E")
}

func TestInvalidClassInputWithClassSuffixGetsError(t *testing.T) {
	userInput := "myBundle/myClass.class"
	testLocation, err := classNameUserInputToTestClassLocation(userInput)
	assert.Nil(t, testLocation)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "GAL1066E")
}

func TestCommandIncludesGALASA_HOMESystemProperty(t *testing.T) {

	bootstrapProps, _, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := true
	isDebugEnabled := false
	var debugPort uint32 = 0
	debugMode := ""

	cmd, args, err := getCommandSyntax(
		bootstrapProps,
		galasaHome,
		fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled,
		debugPort,
		debugMode,
	)

	assert.NotNil(t, cmd)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.Contains(t, args, `-DGALASA_HOME="/User/Home/testuser/.galasa"`)
}

func TestCommandIncludesFlagsFromBootstrapProperties(t *testing.T) {
	bootstrapProps, _, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := false
	isDebugEnabled := false
	var debugPort uint32 = 0
	debugMode := ""

	cmd, args, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, cmd)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.Contains(t, args, "-Xmx80m")
}

func TestCommandIncludesTwoFlagsFromBootstrapProperties(t *testing.T) {
	bootstrapProps,
		_, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	bootstrapProps[api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_OPTIONS] = "-Xmx40m -Xms20m"
	isTraceEnabled := false
	isDebugEnabled := false
	var debugPort uint32 = 0
	debugMode := ""

	cmd, args, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, cmd)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.Contains(t, args, "-Xmx40m")
	assert.Contains(t, args, "-Xms20m")
}

func TestCommandIncludesDefaultDebugPortAndMode(t *testing.T) {
	bootstrapProps,
		_, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := false
	isDebugEnabled := true // <<<< Debug is turned on. No overrides to debugPort in either boostrap or explicit command option.
	var debugPort uint32 = 0
	debugMode := ""

	command, args, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, command)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.Contains(t, args, "-agentlib:jdwp=transport=dt_socket,address=*:"+strconv.FormatUint(uint64(DEBUG_PORT_DEFAULT), 10)+",server=y,suspend=y")
}

func TestCommandDrawsValidDebugPortFromBootstrap(t *testing.T) {
	bootstrapProps,
		_, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := false
	isDebugEnabled := true // <<<< Debug is turned on. No overrides to debugPort in either boostrap or explicit command option.
	var debugPort uint32 = 0
	debugMode := ""

	bootstrapProps[api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_PORT] = "345"

	command, args, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, command)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.Contains(t, args, "-agentlib:jdwp=transport=dt_socket,address=*:345,server=y,suspend=y")
}

func TestCommandDrawsInvalidDebugPortFromBootstrap(t *testing.T) {
	bootstrapProps,
		_, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := false
	isDebugEnabled := true // <<<< Debug is turned on. No overrides to debugPort in either boostrap or explicit command option.
	var debugPort uint32 = 0
	debugMode := ""

	bootstrapProps[api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_PORT] = "-456"

	_, _, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, err)

	assert.Contains(t, err.Error(), "-456")
	assert.Contains(t, err.Error(), api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_PORT)
	assert.Contains(t, err.Error(), "GAL1072E")
}

func TestCommandDrawsValidDebugModeFromBootstrap(t *testing.T) {
	bootstrapProps,
		_, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := false
	isDebugEnabled := true // <<<< Debug is turned on. No overrides to debugPort in either boostrap or explicit command option.
	var debugPort uint32 = 0
	debugMode := ""

	bootstrapProps[api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_MODE] = "attach"

	command, args, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, command)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.Contains(
		t,
		args,
		"-agentlib:jdwp=transport=dt_socket,address=*:"+
			strconv.FormatUint(uint64(DEBUG_PORT_DEFAULT), 10)+
			",server=n,suspend=y",
	)
}

func TestCommandDrawsInvalidDebugModeFromBootstrap(t *testing.T) {
	bootstrapProps,
		_, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := false
	isDebugEnabled := true // <<<< Debug is turned on. No overrides to debugPort in either boostrap or explicit command option.
	var debugPort uint32 = 0
	debugMode := ""

	bootstrapProps[api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_MODE] = "shout" //  << Invalid !

	_, _, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, err)

	assert.Contains(t, err.Error(), "shout")
	assert.Contains(t, err.Error(), api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_MODE)
	assert.Contains(t, err.Error(), "GAL1070E")
}

func TestCommandDrawsValidDebugModeListenFromCommandLine(t *testing.T) {
	bootstrapProps,
		_, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := false
	isDebugEnabled := true // <<<< Debug is turned on. No overrides to debugPort in either boostrap or explicit command option.
	var debugPort uint32 = 0
	debugMode := "listen"

	command, args, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, command)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.Contains(
		t,
		args,
		"-agentlib:jdwp=transport=dt_socket,address=*:"+
			strconv.FormatUint(uint64(DEBUG_PORT_DEFAULT), 10)+
			",server=y,suspend=y",
	)
}

func TestCommandDrawsValidDebugModeAttachFromCommandLine(t *testing.T) {
	bootstrapProps,
		_, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := false
	isDebugEnabled := true // <<<< Debug is turned on. No overrides to debugPort in either boostrap or explicit command option.
	var debugPort uint32 = 0
	debugMode := "attach"

	command, args, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, command)
	assert.NotNil(t, args)
	assert.Nil(t, err)

	assert.Contains(
		t,
		args,
		"-agentlib:jdwp=transport=dt_socket,address=*:"+
			strconv.FormatUint(uint64(DEBUG_PORT_DEFAULT), 10)+
			",server=n,suspend=y",
	)
}

func TestCommandDrawsInvalidDebugModeFromCommandLine(t *testing.T) {
	bootstrapProps,
		_, galasaHome, fs,
		javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		_ := getDefaultCommandSyntaxTestParameters()

	isTraceEnabled := false
	isDebugEnabled := true // <<<< Debug is turned on. No overrides to debugPort in either boostrap or explicit command option.
	var debugPort uint32 = 0
	debugMode := "invalidMode"

	_, _, err := getCommandSyntax(
		bootstrapProps, galasaHome, fs, javaHome,
		testObrs,
		testLocation,
		remoteMaven,
		galasaVersionToRun,
		overridesFilePath,
		isTraceEnabled,
		isDebugEnabled, debugPort, debugMode,
	)

	assert.NotNil(t, err)

	assert.Contains(t, err.Error(), "invalidMode")
	assert.Contains(t, err.Error(), api.BOOTSTRAP_PROPERTY_NAME_LOCAL_JVM_LAUNCH_DEBUG_MODE)
	assert.Contains(t, err.Error(), "GAL1071E")
}
