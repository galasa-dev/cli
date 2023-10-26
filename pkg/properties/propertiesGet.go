/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"context"
	"sort"
	"strings"

	"github.com/galasa-dev/cli/pkg/api"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/propertiesformatter"
	"github.com/galasa-dev/cli/pkg/utils"
)

var (
	validFormatters = CreateFormatters()
)

// GetProperties - performs all the logic to implement the `galasactl properties get` command,
// but in a unit-testable manner.
func GetProperties(
	namespace string,
	name string,
	prefix string,
	suffix string,
	infix string,
	apiServerUrl string,
	propertiesOutputFormat string,
	console utils.Console,
) error {
	var err error
	err = checkNameNotUsedWithPrefixSuffixInfix(name, prefix, suffix, infix)
	if (err == nil) {
		var chosenFormatter propertiesformatter.PropertyFormatter

		chosenFormatter, err = validateOutputFormatFlagValue(propertiesOutputFormat, validFormatters)
		if err == nil {
			var cpsProperty []galasaapi.CpsProperty
			cpsProperty, err = getCpsPropertiesFromRestApi(namespace, name, prefix, suffix, infix, apiServerUrl, console)
			if err == nil {
				var outputText string

			//convert galasaapi.CpsProperty into formattable data
			outputText, err = chosenFormatter.FormatProperties(cpsProperty)

				if err == nil {
					console.WriteString(outputText)
				}

			}
		}
	}

	return err
}

// Retrieves properties from the ecosystem API that match a given namespace.
// Multiple properties can be returned as the namespace is not unique.
func getCpsPropertiesFromRestApi(
	namespace string,
	name string,
	prefix string,
	suffix string,
	infix string,
	apiServerUrl string,
	console utils.Console,
) ([]galasaapi.CpsProperty, error) {

	var err error = nil

	var context context.Context = nil

	// An HTTP client which can communicate with the api server in an ecosystem.
	restClient := api.InitialiseAPI(apiServerUrl)

	var cpsProperties = make([]galasaapi.CpsProperty, 0)

	if name == "" {
		apicall := restClient.ConfigurationPropertyStoreAPIApi.QueryCpsNamespaceProperties(context, namespace)
		if prefix != "" {
			apicall = apicall.Prefix(prefix)
		}
		if suffix != "" {
			apicall = apicall.Suffix(suffix)
		}
		if infix != "" {
			apicall = apicall.Infix(infix)
		}
		cpsProperties, _, err = apicall.Execute()
	} else {
		apicall := restClient.ConfigurationPropertyStoreAPIApi.GetCpsProperty(context, namespace, name)
		cpsProperties, _, err = apicall.Execute()
	}

	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_QUERY_NAMESPACE_FAILED, err.Error())
	}

	return cpsProperties, err
}

func CreateFormatters() map[string]propertiesformatter.PropertyFormatter {
	validFormatters := make(map[string]propertiesformatter.PropertyFormatter, 0)
	summaryFormatter := propertiesformatter.NewPropertySummaryFormatter()
	rawFormatter := propertiesformatter.NewPropertyRawFormatter()
	yamlFormatter := propertiesformatter.NewPropertyYamlFormatter()

	validFormatters[summaryFormatter.GetName()] = summaryFormatter
	validFormatters[rawFormatter.GetName()] = rawFormatter
	validFormatters[yamlFormatter.GetName()] = yamlFormatter

	return validFormatters
}

// Ensures the user has provided a valid output format as part of the "runs get" command.
func validateOutputFormatFlagValue(propertiesOutputFormat string, validFormatters map[string]propertiesformatter.PropertyFormatter) (propertiesformatter.PropertyFormatter, error) {
	var err error

	chosenFormatter, isPresent := validFormatters[propertiesOutputFormat]

	if !isPresent {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OUTPUT_FORMAT, propertiesOutputFormat, GetFormatterNamesString(validFormatters))
	}

	return chosenFormatter, err
}

// GetFormatterNamesString builds a string of comma separated, quoted formatter names
func GetFormatterNamesString(validFormatters map[string]propertiesformatter.PropertyFormatter) string {
	// extract names into a sorted slice
	names := make([]string, 0, len(validFormatters))
	for name := range validFormatters {
		names = append(names, name)
	}
	sort.Strings(names)

	// render list of sorted names into string
	formatterNames := strings.Builder{}

	for count, formatterName := range names {

		if count != 0 {
			formatterNames.WriteString(", ")
		}
		formatterNames.WriteString("'" + formatterName + "'")
	}

	return formatterNames.String()
}

func checkNameNotUsedWithPrefixSuffixInfix(name string, prefix string, suffix string, infix string) error {
	var err error
	if name != "" && (prefix != "" || suffix != "" || infix != "") {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_PROPERTIES_FLAG_COMBINATION)
	}
	return err
}
