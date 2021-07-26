/*
*  Licensed Materials - Property of IBM
*
* (c) Copyright IBM Corp. 2021.
*/

package api

import (
	"os"
	"io"
	"strings"
	"net/http"

	"github.com/galasa.dev/cli/pkg/utils"
	"github.com/galasa.dev/cli/pkg/galasaapi"
)

var (
	bootstrap string
	baseURL   string

	bootstrapProperties utils.JavaProperties
)

func InitialiseAPI(providedBootstrap string) *galasaapi.APIClient {
	// Calculate the bootstrap for this execution
	bootstrap = providedBootstrap

	if (bootstrap == "") {
		bootstrap = os.Getenv("GALASA_BOOTSTRAP")
	}

	if (bootstrap == "") {
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

	if (strings.HasPrefix(bootstrap, "http:") || strings.HasPrefix(bootstrap, "https:")) {
		resp, err := http.Get(bootstrap)
		if (err != nil) {
			panic(err)
		}
		defer resp.Body.Close()

		_, err = io.Copy(bootstrapString, resp.Body)
		if (err != nil) {
			panic(err)
		}

		if strings.HasSuffix(bootstrap, "/bootstrap") {
			baseURL = bootstrap[:len(bootstrap)-10]
		} else {
			panic("bootstrap url does not end in /bootstrap")
		}

//		fmt.Printf("base=%v\n", baseURL)
	} else { // assume file
		panic("unsupported bootstrap")
	}

	// read the lines and extract the properties
//	fmt.Printf("bootstrap contents:-\n%v\n", bootstrapString.String())
	bootstrapProperties = utils.ReadProperties(bootstrapString.String())
}