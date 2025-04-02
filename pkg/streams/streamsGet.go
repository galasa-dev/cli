/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package streams

import (
	"context"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/streamsformatter"
)

var (
	formatters = createFormatters()
)

func GetStreams(streamName string, format string, apiClient *galasaapi.APIClient, console spi.Console) error {

	var err error
	var streamData []galasaapi.Stream

	if streamName != "" {
		streamName, err = validateStreamName(streamName)
	}

	if err == nil {
		streamData, err = getStreamsFromRestApi(streamName, apiClient)
		if err == nil {
			err = formatFetchedStreamsAndWriteToConsole(streamData, console, format)
		}
	}

	return err

}

func formatFetchedStreamsAndWriteToConsole(streams []galasaapi.Stream, console spi.Console, outputFormatString string) error {

	var formattedOuptut string
	chosenFormatter, err := validateFormatFlag(outputFormatString)

	if err == nil {
		formattedOuptut, err = chosenFormatter.FormatStreams(streams)
		if err == nil {
			console.WriteString(formattedOuptut)
		}
	}

	return err

}

func getStreamsFromRestApi(
	streamName string,
	apiClient *galasaapi.APIClient,
) ([]galasaapi.Stream, error) {

	var streams []galasaapi.Stream
	var err error
	var restApiVersion string

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if err == nil {
		if streamName != "" {
			streams, err = getStreamByName(streamName, apiClient, restApiVersion)
		} else {
			streams, err = getAllStreams(apiClient, restApiVersion)
		}
	}

	return streams, err

}

func createFormatters() map[string]streamsformatter.StreamsFormatter {

	formatters := make(map[string]streamsformatter.StreamsFormatter, 0)
	summaryFormatter := streamsformatter.NewStreamsSummaryFormatter()
	yamlFormatter := streamsformatter.NewStreamsYamlFormatter()

	formatters[summaryFormatter.GetName()] = summaryFormatter
	formatters[yamlFormatter.GetName()] = yamlFormatter

	return formatters

}

func GetFormatterNamesAsString() string {
	names := make([]string, 0, len(formatters))
	for name := range formatters {
		names = append(names, name)
	}
	sort.Strings(names)
	formatterNames := strings.Builder{}

	for index, formatterName := range names {

		if index != 0 {
			formatterNames.WriteString(", ")
		}
		formatterNames.WriteString("'" + formatterName + "'")
	}

	return formatterNames.String()
}

func getStreamByName(
	streamName string,
	apiClient *galasaapi.APIClient,
	restApiVersion string,
) ([]galasaapi.Stream, error) {

	var err error
	var streamIn *galasaapi.Stream
	var resp *http.Response
	var context context.Context = nil
	var streamsToReturn []galasaapi.Stream

	apiCall := apiClient.StreamsAPIApi.GetStreamByName(context, streamName).ClientApiVersion(restApiVersion)

	streamIn, resp, err = apiCall.Execute()

	if err == nil {
		var statusCode int
		if resp != nil {
			defer resp.Body.Close()
			statusCode = resp.StatusCode
		}

		if statusCode != 200 {
			log.Println("getStreamsFromRestApi - Failed to retrieve list of test streams from API server")
			err = galasaErrors.NewGalasaErrorWithHttpStatusCode(statusCode, galasaErrors.GALASA_ERROR_RETRIEVING_STREAMS_FROM_API_SERVER, err.Error())
		}
		streamsToReturn = []galasaapi.Stream{*streamIn}
	}

	return streamsToReturn, err

}

func getAllStreams(
	apiClient *galasaapi.APIClient,
	restApiVersion string,
) ([]galasaapi.Stream, error) {

	var err error
	var streams []galasaapi.Stream
	var resp *http.Response
	var context context.Context = nil

	apiCall := apiClient.StreamsAPIApi.GetStreams(context).ClientApiVersion(restApiVersion)

	streams, resp, err = apiCall.Execute()

	if err == nil {

		var statusCode int
		if resp != nil {
			defer resp.Body.Close()
			statusCode = resp.StatusCode
		}

		if statusCode == 200 {
			log.Printf("getUserDataFromRestApi - %v test streams collected", len(streams))
		} else {
			log.Println("getStreamsFromRestApi - Failed to retrieve list of test streams from API server")
			err = galasaErrors.NewGalasaErrorWithHttpStatusCode(statusCode, galasaErrors.GALASA_ERROR_RETRIEVING_STREAMS_FROM_API_SERVER, err.Error())
			statusCode = 500
		}
	}

	return streams, err
}

func validateFormatFlag(outputFormatString string) (streamsformatter.StreamsFormatter, error) {
	var err error

	chosenFormatter, isPresent := formatters[outputFormatString]

	if !isPresent {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_INVALID_OUTPUT_FORMAT, outputFormatString, GetFormatterNamesAsString())
	}

	return chosenFormatter, err

}
