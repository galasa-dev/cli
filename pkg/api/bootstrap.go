/*
 * Copyright contributors to the Galasa project
 */
package api

import (
	"io"
	"net/http"
	"strings"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/utils"
)

const (
	BOOTSTRAP_PROPERTY_NAME_REMOTE_API_SERVER_URL string = "framework.api.server.url"
)

type BootstrapData struct {
	// Path - the raw path that a user has given us, either from the command-line
	// option or the GALASA_BOOTSTRAP environment variable.
	Path string

	// URL - The URL on which can be followed to reach the bootstrap contents.
	ApiServerURL string

	// Properties - The properties which are read from the bootstrap
	Properties utils.JavaProperties
}

type UrlResolutionService interface {
	Get(url string) (string, error)
}

type RealUrlResolutionService struct {
}

// get - Gets the string contents from a URL
func (*RealUrlResolutionService) Get(url string) (string, error) {
	var resp *http.Response
	var contents string = ""
	var err error = nil

	resp, err = http.Get(url)
	// Wrap the error inside a galasa error.
	if err == nil {

		// Make sure the http response is closed (eventually).
		defer resp.Body.Close()

		buffer := new(strings.Builder)
		_, err = io.Copy(buffer, resp.Body)
		if err == nil {
			contents = buffer.String()
		}
	}
	return contents, err
}

// getDefaultBootstrapPath - Work out where the boostrap file can normally be found.
func getDefaultBootstrapPath(fileSystem utils.FileSystem) (string, error) {
	var path string
	home, err := fileSystem.GetUserHomeDir()
	if err == nil {
		path = home + utils.FILE_SYSTEM_PATH_SEPARATOR + ".galasa" +
			utils.FILE_SYSTEM_PATH_SEPARATOR + "bootstrap.properties"
	}
	return path, err
}

// loadBootstrap - Loads the contents of a bootstrap file into memory.
// bootstrapPath - Where do we find the bootstrap contents from ? This can be a URL must end in /bootstrap
func LoadBootstrap(fileSystem utils.FileSystem, env utils.Environment,
	bootstrapPath string, urlResolutionService UrlResolutionService) (*BootstrapData, error) {

	var err error = nil

	path := bootstrapPath

	var bootstrap *BootstrapData = nil

	// If the --bootstrap flag wasn't specified by the user... default to the value in the
	// GALASA_BOOTSTRAP environment variable.
	if path == "" {
		path = env.GetEnv("GALASA_BOOTSTRAP")
	}

	// If it's still not clear, use the default bootstrap.properties in the ${HOME}/.galasa folder.
	if path == "" {
		path, err = getDefaultBootstrapPath(fileSystem)
	}

	if err == nil {

		// Default the API server to assume it's running locally, natively.
		defaultApiServerURL := "http://127.0.0.1"

		if strings.HasPrefix(path, "http:") || strings.HasPrefix(path, "https:") {
			// The path looks like a URL...
			bootstrap, err = loadBootstrapFromUrl(path, defaultApiServerURL, urlResolutionService)
		} else {
			// The path looks like a file...
			bootstrap, err = loadBootstrapFromFile(path, defaultApiServerURL, fileSystem)
		}

		if err == nil {
			// Now we have a collection of boot properties. Find the Api server URL if there is one
			// to populate the specified field in the structure.
			if bootstrap.Properties != nil {
				apiServerUrlFromPropsFile := bootstrap.Properties[BOOTSTRAP_PROPERTY_NAME_REMOTE_API_SERVER_URL]
				if apiServerUrlFromPropsFile != "" {
					bootstrap.ApiServerURL = apiServerUrlFromPropsFile
				}
			}
		} else {
			// Don't return any data if there was a failure.
			bootstrap = nil
		}
	}

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_LOAD_BOOTSTRAP_FILE, path, err.Error())
	}

	return bootstrap, err
}

func cleanPath(fileSystem utils.FileSystem, path string) (string, error) {
	var err error = nil
	if path != "" {
		path = removeLeadingFileColon(path)
		path, err = utils.TildaExpansion(fileSystem, path)
	}
	return path, err
}

func removeLeadingFileColon(path string) string {
	path = strings.TrimPrefix(path, "file:")
	return path
}

func loadBootstrapFromFile(path string, defaultApiServerURL string, fileSystem utils.FileSystem) (*BootstrapData, error) {
	bootstrap := new(BootstrapData)
	var content string
	var err error = nil

	bootstrap.ApiServerURL = defaultApiServerURL

	path, err = cleanPath(fileSystem, path)
	if err == nil {
		content, err = fileSystem.ReadTextFile(path)
	}

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_READ_BOOTSTRAP_FILE, path, err.Error())
	} else {
		// read the lines and extract the properties
		//	fmt.Printf("bootstrap contents:-\n%v\n", bootstrapString.String())
		bootstrap.Properties = utils.ReadProperties(content)
	}

	if err != nil {
		bootstrap = nil
	}
	return bootstrap, err
}

func loadBootstrapFromUrl(path string, defaultApiServerURL string,
	urlResolutionService UrlResolutionService) (*BootstrapData, error) {

	bootstrap := new(BootstrapData)
	bootstrap.Path = path
	var err error = nil

	bootstrap.ApiServerURL = defaultApiServerURL

	// Check that the provided url has /bootstrap at the end.
	// Then strip it off to form the URL we will use.

	if !strings.HasSuffix(path, "/bootstrap") {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_BOOTSTRAP_URL_BAD_ENDING, bootstrap.Path)
	} else {
		// strip off the /bootstrap to get the ApiServerURL
		bootstrap.ApiServerURL = path[:len(path)-10]

		// Use the bootstrap URL to get the bootstrap contents
		var bootstrapContents string
		bootstrapContents, err = urlResolutionService.Get(path)
		if err != nil {
			// Wrap any http error as a galasa error.
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_GET_BOOTSTRAP, path, err.Error())
		} else {
			// read the lines and extract the properties
			//	fmt.Printf("bootstrap contents:-\n%v\n", bootstrapString.String())
			bootstrap.Properties = utils.ReadProperties(bootstrapContents)
		}
	}

	return bootstrap, err
}
