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

	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/streamsformatter"
)

func GetStreams(streamName string, apiClient *galasaapi.APIClient, console spi.Console) error {

	streamData, err := getStreamsFromRestApi(streamName, apiClient)
	if err == nil {
		err = formatFetchedStreamsAndWriteToConsole(streamData, console)
	}

	return err

}

func formatFetchedStreamsAndWriteToConsole(streams []galasaapi.Stream, console spi.Console) error {

	summaryFormatter := streamsformatter.NewStreamsSummaryFormatter()
	outputText, err := summaryFormatter.FormatStreams(streams)

	if err == nil {
		console.WriteString(outputText)
	}

	return err

}

func getStreamsFromRestApi(
	streamName string,
	apiClient *galasaapi.APIClient,
) ([]galasaapi.Stream, error) {

	var context context.Context = nil
	var streams []galasaapi.Stream
	var err error
	var restApiVersion string

	restApiVersion, err = embedded.GetGalasactlRestApiVersion()

	if streamName != "" {
		streamName, err = validateStreamName(streamName)

		if err == nil {
			streams, err = getStreamByName(streamName, apiClient, &context, restApiVersion)
		}

	} else {

		if err == nil {
			streams, err = getAllStreams(apiClient, &context, restApiVersion)
		}

	}

	return streams, err

}

func getStreamByName(
	streamName string,
	apiClient *galasaapi.APIClient,
	context *context.Context,
	restApiVersion string,
) ([]galasaapi.Stream, error) {

	var err error
	var streamIn *galasaapi.Stream
	var resp *http.Response

	apiCall := apiClient.StreamsAPIApi.GetStreamByName(*context, streamName).ClientApiVersion(restApiVersion)

	streamIn, resp, err = apiCall.Execute()

	var statusCode int
	if resp != nil {
		defer resp.Body.Close()
		statusCode = resp.StatusCode
	}

	if err != nil {
		log.Println("getStreamsFromRestApi - Failed to retrieve list of test streams from API server")
		err = galasaErrors.NewGalasaErrorWithHttpStatusCode(statusCode, galasaErrors.GALASA_ERROR_RETRIEVING_USER_LIST_FROM_API_SERVER, err.Error())
	}

	streamsToReturn := []galasaapi.Stream{*streamIn}
	return streamsToReturn, err

}

func getAllStreams(
	apiClient *galasaapi.APIClient,
	context *context.Context,
	restApiVersion string,
) ([]galasaapi.Stream, error) {

	var err error
	var streams []galasaapi.Stream
	var resp *http.Response

	apiCall := apiClient.StreamsAPIApi.GetStreams(*context).ClientApiVersion(restApiVersion)

	streams, resp, err = apiCall.Execute()

	var statusCode int
	if resp != nil {
		defer resp.Body.Close()
		statusCode = resp.StatusCode
	}

	if err != nil {
		log.Println("getStreamsFromRestApi - Failed to retrieve list of test streams from API server")
		err = galasaErrors.NewGalasaErrorWithHttpStatusCode(statusCode, galasaErrors.GALASA_ERROR_RETRIEVING_STREAMS_FROM_API_SERVER, err.Error())
	} else {
		log.Printf("getUserDataFromRestApi - %v test streams collected", len(streams))

	}

	return streams, err
}
