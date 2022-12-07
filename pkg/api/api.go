/*
 * Copyright contributors to the Galasa project
 */
package api

import (
	"fmt"
	"io"
	"log"
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

func InitialiseAPI(providedBootstrap string) *galasaapi.APIClient {
	// Calculate the bootstrap for this execution
	bootstrap = providedBootstrap

	if bootstrap == "" {
		bootstrap = os.Getenv("GALASA_BOOTSTRAP")
	}

	if bootstrap == "" {
		bootstrap = "~/.galasa/bootstrap"
	}

	loadBootstrap()

	cfg := galasaapi.NewConfiguration()
	cfg.Debug = false
	cfg.Servers = galasaapi.ServerConfigurations{{URL: baseURL}}
	apiClient := galasaapi.NewAPIClient(cfg)

	return apiClient
}

func loadBootstrap() {
	//	fmt.Printf("using bootstrap %v\n", bootstrap)

	bootstrapString := new(strings.Builder)

	baseURL = "http://127.0.0.1"

	if strings.HasPrefix(bootstrap, "http:") || strings.HasPrefix(bootstrap, "https:") {

		if strings.HasSuffix(bootstrap, "/bootstrap") {
			baseURL = bootstrap[:len(bootstrap)-10]
		} else {
			msg := fmt.Sprintf(galasaErrors.GALASA_ERROR_BOOTSTRAP_URL_BAD_ENDING.Template, bootstrap)
			log.Println(msg)
			panic(msg)
		}

		resp, err := http.Get(bootstrap)
		if err != nil {
			msg := fmt.Sprintf(galasaErrors.GALASA_ERROR_FAILED_TO_GET_BOOTSTRAP.Template, bootstrap, err.Error())
			log.Println(msg)
			panic(msg)
		}
		defer resp.Body.Close()

		_, err = io.Copy(bootstrapString, resp.Body)
		if err != nil {
			msg := fmt.Sprintf(galasaErrors.GALASA_ERROR_BAD_BOOTSTRAP_CONTENT.Template, bootstrap, err.Error())
			log.Println(msg)
			panic(msg)
		}

		//		fmt.Printf("base=%v\n", baseURL)
	} else { // assume file
		msg := fmt.Sprintf(galasaErrors.GALASA_ERROR_UNSUPPORTED_BOOTSTRAP_URL.Template, bootstrap)
		log.Println(msg)
		panic(msg)
	}

	// read the lines and extract the properties
	//	fmt.Printf("bootstrap contents:-\n%v\n", bootstrapString.String())
	bootstrapProperties = utils.ReadProperties(bootstrapString.String())
}
