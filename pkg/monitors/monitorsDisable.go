/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package monitors

import (
	"log"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
)

func DisableMonitor(
	monitorName string,
	apiClient *galasaapi.APIClient,
	byteReader spi.ByteReader,
) error {
	var err error

	desiredEnabledState := false
	err = setMonitorIsEnabledState(monitorName, desiredEnabledState, apiClient, byteReader)

	log.Printf("DisableMonitor exiting. err is %v\n", err)
	return err
}
