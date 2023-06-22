/*
 * Copyright contributors to the Galasa project
 */
package embedded

import (
	"embed"
	"log"

	"github.com/galasa.dev/cli/pkg/props"
)

const (
	PROPERTY_NAME_GALASACTL_VERSION        = "galasactl.version"
	PROPERTY_NAME_GALASA_BOOT_JAR_VERSION  = "galasa.boot.jar.version"
	PROPERTY_NAME_GALASA_FRAMEWORK_VERSION = "galasa.framework.version"
)

// Embed all the template files into the go executable, so there are no extra files
// we need to ship/install/locate on the target machine.
// We can access the "embedded" file system as if they are normal files.
//
//go:embed templates/*
var embeddedFileSystem embed.FS

// An instance of the ReadOnlyFileSystem interface, set once, used many times.
// It just delegates to teh embed.FS
var readOnlyFileSystem ReadOnlyFileSystem

type versions struct {
	galasaFrameworkVersion string
	galasaBootJarVersion   string
	galasactlVersion       string
}

var (
	versionsCache *versions = nil
	PropsFileName           = "version/build.properties"
)

func GetGalasaVersion() (string, error) {
	fs := GetReadOnlyFileSystem()
	cache, err := getVersionsCache(fs)
	var version string
	if err == nil {
		version = cache.galasaFrameworkVersion
	}
	return version, err
}

func GetBootJarVersion() (string, error) {
	fs := GetReadOnlyFileSystem()
	cache, err := getVersionsCache(fs)
	var version string
	if err == nil {
		version = cache.galasaBootJarVersion
	}
	return version, err
}

func GetGalasaCtlVersion() (string, error) {
	fs := GetReadOnlyFileSystem()
	cache, err := getVersionsCache(fs)
	var version string
	if err == nil {
		version = cache.galasactlVersion
	}
	return version, err
}

func GetReadOnlyFileSystem() ReadOnlyFileSystem {
	if readOnlyFileSystem == nil {
		readOnlyFileSystem = NewReadOnlyFileSystem()
	}
	return readOnlyFileSystem
}

func getVersionsCache(fs ReadOnlyFileSystem) (*versions, error) {
	var (
		err   error
		bytes []byte
	)
	if versionsCache == nil {
		versionsCache = &versions{
			galasaFrameworkVersion: "unknown",
			galasaBootJarVersion:   "unknown",
			galasactlVersion:       "unknown",
		}
		log.Printf("Loading the properties file '%s'...", PropsFileName)
		bytes, err = fs.ReadFile(PropsFileName)
		if err != nil {
			log.Printf("Failure. %s", err.Error())
		} else {
			propsFileContent := string(bytes)
			properties := props.ReadProperties(propsFileContent)

			versionsCache.galasaBootJarVersion = properties[PROPERTY_NAME_GALASA_BOOT_JAR_VERSION]
			versionsCache.galasaFrameworkVersion = properties[PROPERTY_NAME_GALASA_FRAMEWORK_VERSION]
			versionsCache.galasactlVersion = properties[PROPERTY_NAME_GALASACTL_VERSION]
		}
	}
	return versionsCache, err
}
