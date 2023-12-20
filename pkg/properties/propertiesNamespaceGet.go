/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"context"
	"log"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/propertiesformatter"
	"github.com/galasa-dev/cli/pkg/utils"
)

var (
	namespaceHasYamlFormat   = false
	validNamespaceFormatters = CreateFormatters(namespaceHasYamlFormat)
)

// GetPropertiesNamespaces - performs all the logic to implement the `galasactl properties namespace get` command
func GetPropertiesNamespaces(
	apiClient *galasaapi.APIClient,
	namespaceOutputFormat string,
	console utils.Console,
) error {
	var err error
	var chosenFormatter propertiesformatter.PropertyFormatter
	var context context.Context = nil

	chosenFormatter, err = validateOutputFormatFlagValue(namespaceOutputFormat, validNamespaceFormatters)
	if err == nil {
		var namespaces []galasaapi.Namespace
		namespaces, _, err = apiClient.ConfigurationPropertyStoreAPIApi.GetAllCpsNamespaces(context).Execute()
		log.Printf("GetPropertiesNamespaces -  namespaces collected: %v", namespaces)

		if err == nil {
			var outputText string

			outputText, err = chosenFormatter.FormatNamespaces(namespaces)

			if err == nil {
				console.WriteString(outputText)
			}

		} else {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_CPS_FAILED, err.Error())
		}
	}

	return err
}
