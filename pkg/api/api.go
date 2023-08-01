/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package api

import (
	"github.com/galasa.dev/cli/pkg/galasaapi"
)

func InitialiseAPI(apiServerUrl string) *galasaapi.APIClient {
	// Calculate the bootstrap for this execution

	var apiClient *galasaapi.APIClient = nil

	cfg := galasaapi.NewConfiguration()
	cfg.Debug = false
	cfg.Servers = galasaapi.ServerConfigurations{{URL: apiServerUrl}}
	apiClient = galasaapi.NewAPIClient(cfg)

	return apiClient
}
