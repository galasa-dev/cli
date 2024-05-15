/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"context"
	"log"
	"sort"
	"strings"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/tokensformatter"
	"github.com/galasa-dev/cli/pkg/utils"
)

var (
	validTokenFormatters = CreateFormatters()
)

// GetTokens - performs all the logic to implement the `galasactl auth tokens get` command
func GetTokens(
	apiClient *galasaapi.APIClient,
	tokensOutputFormat string,
	console utils.Console,
) error {
	log.Println("GetTokens - ENTEREDr")
	var err error = nil
	var context context.Context = nil
	var authTokens *galasaapi.AuthTokens

	authTokens, _, err = apiClient.AuthenticationAPIApi.GetTokens(context).Execute()
	if err != nil {
		log.Println("GetTokens - Failed to retrieve list of tokens from API server")
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_TOKEN_LIST_FROM_API_SERVER, err.Error())
	} else {
		log.Printf("GetTokens -  tokens collected: %v", authTokens)
		var chosenFormatter tokensformatter.TokenFormatter

		chosenFormatter, err = validateOutputFormatFlagValue(tokensOutputFormat, validTokenFormatters)
		if err == nil {

			var outputText string
			outputText, err = chosenFormatter.FormatTokens(authTokens.GetTokens())

			if err == nil {
				console.WriteString(outputText)
			}
		}

	}

	return err
}

// Ensures the user has provided a valid output format as part of the "runs get" command.
func validateOutputFormatFlagValue(tokensOutputFormat string, validFormatters map[string]tokensformatter.TokenFormatter) (tokensformatter.TokenFormatter, error) {
	var err error

	chosenFormatter, isPresent := validFormatters[tokensOutputFormat]

	if !isPresent {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OUTPUT_FORMAT, tokensOutputFormat, GetFormatterNamesString(validFormatters))
	}

	return chosenFormatter, err
}

// GetFormatterNamesString builds a string of comma separated, quoted formatter names
func GetFormatterNamesString(validFormatters map[string]tokensformatter.TokenFormatter) string {
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

func CreateFormatters() map[string]tokensformatter.TokenFormatter {
	validFormatters := make(map[string]tokensformatter.TokenFormatter, 0)
	summaryFormatter := tokensformatter.NewTokenSummaryFormatter()

	validFormatters[summaryFormatter.GetName()] = summaryFormatter

	return validFormatters
}
