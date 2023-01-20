/*
 * Copyright contributors to the Galasa project
 */
package api

import (
	"io"
	"net/http"
	"os"
	"strings"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/utils"
)

var (
	bootstrap string
	baseURL   string

	bootstrapProperties utils.JavaProperties
)

func InitialiseAPI(providedBootstrap string) (*galasaapi.APIClient, error) {
	// Calculate the bootstrap for this execution
	bootstrap = providedBootstrap
	var apiClient *galasaapi.APIClient = nil

	if bootstrap == "" {
		bootstrap = os.Getenv("GALASA_BOOTSTRAP")
	}

	if bootstrap == "" {
		bootstrap = "~/.galasa/bootstrap"
	}

	err := loadBootstrap()

	if err == nil {
		cfg := galasaapi.NewConfiguration()
		cfg.Debug = false
		cfg.Servers = galasaapi.ServerConfigurations{{URL: baseURL}}
		apiClient = galasaapi.NewAPIClient(cfg)
	}

	return apiClient, err
}

func loadBootstrap() error {
	//	fmt.Printf("using bootstrap %v\n", bootstrap)

	bootstrapString := new(strings.Builder)

	baseURL = "http://127.0.0.1"

	if strings.HasPrefix(bootstrap, "http:") || strings.HasPrefix(bootstrap, "https:") {

		if strings.HasSuffix(bootstrap, "/bootstrap") {
			baseURL = bootstrap[:len(bootstrap)-10]
		} else {
			err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_BOOTSTRAP_URL_BAD_ENDING, bootstrap)
			return err
		}

		resp, err := http.Get(bootstrap)
		if err != nil {
			err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_GET_BOOTSTRAP, bootstrap, err.Error())
			return err
		}
		defer resp.Body.Close()

		_, err = io.Copy(bootstrapString, resp.Body)
		if err != nil {
			err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_BAD_BOOTSTRAP_CONTENT, bootstrap, err.Error())
			return err
		}

		//		fmt.Printf("base=%v\n", baseURL)
	} else { // assume file
		err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_UNSUPPORTED_BOOTSTRAP_URL, bootstrap)
		return err
	}

	// read the lines and extract the properties
	//	fmt.Printf("bootstrap contents:-\n%v\n", bootstrapString.String())
	bootstrapProperties = utils.ReadProperties(bootstrapString.String())

	return nil
}
