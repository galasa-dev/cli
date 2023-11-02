/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"context"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/propertiesformatter"
	"github.com/galasa-dev/cli/pkg/utils"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

// GetNamespaceProperties - performs all the logic to implement the `galasactl properties namespace get` command
func GetNamespaceProperties(
	apiServerUrl string,
	console utils.Console,
) error {
	var err error
	var chosenFormatter propertiesformatter.PropertyFormatter
	var context context.Context = nil

	//only format so far is the default format, summary
	chosenFormatter, err = validateOutputFormatFlagValue("summary", validFormatters)
	if err == nil {
		var namespaces []galasaapi.Namespace
		restClient := api.InitialiseAPI(apiServerUrl)
		namespaces, _, err = restClient.ConfigurationPropertyStoreAPIApi.GetAllCpsNamespaces(context).Execute()

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
