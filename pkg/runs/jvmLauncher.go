/*
 * Copyright contributors to the Galasa project
 */
package runs

import (
	"log"
	"os/exec"
	"strings"

	"github.com/galasa.dev/cli/pkg/embedded"
	"github.com/galasa.dev/cli/pkg/utils"
)

type TestLocation struct {
	OSGiBundleName string
	Class          JavaClassDef
}

type TestResultsSummary struct {
	MethodPasses int
	MethodFails  int
}

func executeTestInJVM(fileSystem utils.FileSystem, javaHome string, testObrs []MavenCoordinates, testLocation TestLocation, remoteMaven string) error {
	cmd, args, err := getCommandSyntax(fileSystem, javaHome, testObrs, testLocation, remoteMaven)
	if err == nil {
		log.Printf("Launching command '%s' '%v'\n", cmd, args)
		jvmProcess := exec.Command(cmd, args...)

		// Wait for it to complete.
		jvmOutput, err := jvmProcess.Output()
		if err != nil {
			log.Printf("Failed to launch the JVM. %s\n", err.Error())
			log.Printf("Failing command is %s %v\n", cmd, args)
		} else {
			log.Printf("jvm output is %s\n", jvmOutput)

			results := summariseJVMOutput(jvmOutput)
			log.Printf("Results: Method passes: %d, fails: %d", results.MethodPasses, results.MethodFails)
		}
	}
	return err
}

func summariseJVMOutput(jvmOutput []byte) TestResultsSummary {
	jvmOutputStr := string(jvmOutput[:])

	results := TestResultsSummary{MethodPasses: 0, MethodFails: 0}

	results.MethodPasses = strings.Count(jvmOutputStr, "*** Passed - Test method ")
	results.MethodFails = strings.Count(jvmOutputStr, "*** Failed - Test method ")

	return results
}

// getCommandSyntax From the parameters we aim to build a command-line incantation which would launch the test in a JVM...
// For example:
// java -jar ${BOOT_JAR_PATH} \
// --localmaven file:${M2_PATH}/repository/ \
// --remotemaven $REMOTE_MAVEN \
// --bootstrap file:${HOME}/.galasa/bootstrap.properties \
// --overrides file:${HOME}/.galasa/overrides.properties \
// --obr mvn:dev.galasa/dev.galasa.uber.obr/${OBR_VERSION}/obr \
// --obr mvn:${TEST_OBR_GROUP_ID}/${TEST_OBR_ARTIFACT_ID}/${TEST_OBR_VERSION}/obr \
// --test ${TEST_BUNDLE}/${TEST_JAVA_CLASS} | tee jvm-log.txt | grep "[*][*][*]" | grep -v "[*][*][*][*]" | sed -e "s/[--]*//g"
//
// For example:
//
//	java -jar /Users/mcobbett/builds/galasa/code/external/galasa-dev/cli/pkg/embedded/templates/galasahome/lib/galasa-boot-0.24.0.jar \
//	    --localmaven file:/Users/mcobbett/.m2/repository/ \
//	    --remotemaven https://development.galasa.dev/main/maven-repo/obr/ \
//	    --bootstrap file:/Users/mcobbett/.galasa/bootstrap.properties \
//	    --overrides file:/Users/mcobbett/.galasa/overrides.properties \
//	    --obr mvn:dev.galasa/dev.galasa.uber.obr/0.25.0/obr \
//	    --obr mvn:dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOT/obr \
//	    --test dev.galasa.example.banking.payee/dev.galasa.example.banking.payee.TestPayee
func getCommandSyntax(fileSystem utils.FileSystem, javaHome string, testObrs []MavenCoordinates, testLocation TestLocation, remoteMaven string) (string, []string, error) {

	var cmd string = ""
	var args []string = make([]string, 0)

	bootJarPath, err := utils.GetGalasaBootJarPath(fileSystem)
	if err == nil {

		cmd = javaHome + utils.FILE_SYSTEM_PATH_SEPARATOR + "bin" +
			utils.FILE_SYSTEM_PATH_SEPARATOR + "java"

		args = append(args, "-jar")
		args = append(args, bootJarPath)

		var userHome string
		userHome, err = fileSystem.GetUserHomeDir()

		// --localmaven file:${M2_PATH}/repository/
		args = append(args, "--localmaven")
		localMavenPath := "file:" + userHome + utils.FILE_SYSTEM_PATH_SEPARATOR +
			".m2" + utils.FILE_SYSTEM_PATH_SEPARATOR + "repository"
		args = append(args, localMavenPath)

		// --remotemaven $REMOTE_MAVEN
		args = append(args, "--remotemaven")
		args = append(args, remoteMaven)

		// --bootstrap file:${HOME}/.galasa/bootstrap.properties
		args = append(args, "--bootstrap")
		bootstrapPath := "file:" + userHome + utils.FILE_SYSTEM_PATH_SEPARATOR +
			".galasa" + utils.FILE_SYSTEM_PATH_SEPARATOR + "bootstrap.properties"
		args = append(args, bootstrapPath)

		// --overrides file:${HOME}/.galasa/overrides.properties
		args = append(args, "--overrides")
		overridesPath := "file:" + userHome + utils.FILE_SYSTEM_PATH_SEPARATOR +
			".galasa" + utils.FILE_SYSTEM_PATH_SEPARATOR + "overrides.properties"
		args = append(args, overridesPath)

		for _, obrCoordinate := range testObrs {
			// We are aiming for this:
			// mvn:${TEST_OBR_GROUP_ID}/${TEST_OBR_ARTIFACT_ID}/${TEST_OBR_VERSION}/obr
			args = append(args, "--obr")
			obrMvnPath := "mvn:" + obrCoordinate.GroupId + "/" +
				obrCoordinate.ArtifactId + "/" + obrCoordinate.Version + "/obr"
			args = append(args, obrMvnPath)
		}

		// --obr mvn:dev.galasa/dev.galasa.uber.obr/${OBR_VERSION}/obr
		args = append(args, "--obr")
		galasaUberObrPath := "mvn:dev.galasa/dev.galasa.uber.obr/" + embedded.GetGalasaVersion() + "/obr"
		args = append(args, galasaUberObrPath)

		// --test ${TEST_BUNDLE}/${TEST_JAVA_CLASS}
		args = append(args, "--test")
		args = append(args, testLocation.OSGiBundleName+"/"+testLocation.Class.PackageName+"."+testLocation.Class.ClassName)

	}

	return cmd, args, err
}
