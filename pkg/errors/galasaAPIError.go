/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package errors

import (
	"encoding/json"
	"log"
)

type GalasaAPIError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
}

// This function reads a galasa API Error into a structure so that it can be displayed as
// the reason for the failure.
// NOTE: when this function is called ensure that the calling function has the  `defer resp.Body.Close()`
// called in order to ensure that the response body is closed when the function completes
func GetApiErrorFromResponse(body []byte) (*GalasaAPIError, error){
	var err error

	apiError := new(GalasaAPIError)

	err = json.Unmarshal(body, &apiError)

	if err != nil {
		log.Printf("GetApiErrorFromResponse FAIL - %v", err)
		err = NewGalasaError(GALASA_ERROR_UNABLE_TO_READ_RESPONSE_BODY, err.Error())
	}
	return apiError, err
}
