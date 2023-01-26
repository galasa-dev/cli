/*
 * Copyright contributors to the Galasa project
 */
package runs

// import (
// 	"os"
// 	"testing"

// 	"github.com/galasa.dev/cli/pkg/utils"
// 	"github.com/stretchr/testify/assert"
// )

// func TestJvmGetsLaunchedWithCorrectSyntax(t *testing.T) {

// 	// Given...
// 	mockFileSystem := utils.NewOSFileSystem()

// 	var testObrs []MavenCoordinates = []MavenCoordinates{
// 		{
// 			GroupId:    "dev.galasa.example.banking",
// 			ArtifactId: "dev.galasa.example.banking.obr",
// 			Version:    "0.0.1-SNAPSHOT",
// 		},
// 	}

// 	var testLocation TestLocation = TestLocation{
// 		OSGiBundleName: "dev.galasa.example.banking.payee",
// 		Class: JavaClassDef{
// 			PackageName: "dev.galasa.example.banking.payee",
// 			ClassName:   "TestPayee",
// 		},
// 	}

// 	javaHome := os.Getenv("JAVA_HOME")

// 	// remoteMaven := "https://repo.maven.apache.org/maven2"
// 	remoteMaven := "https://development.galasa.dev/main/maven-repo/obr/"

// 	// When
// 	err := executeTestInJVM(mockFileSystem, javaHome, testObrs, testLocation, remoteMaven)

// 	// Then
// 	if err != nil {
// 		assert.Fail(t, "Expecting no errors but there was one."+err.Error())
// 	}
// }
