/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package errors

import (
	"encoding/json"
	"io"
	"net/http"
)


 type GalasaAPIError struct {
	Code    int `json:"error_code"`
	Message string `json:"error_message"`
}

/* This function reads a galasa API Error into a structure so that it can be displayed as 
 * the reason for the failure.
 * NOTE: when this function is called ensure that the calling function has the  `defer resp.Body.Close()`
 * called in order to ensure that the response body is closed when the function completes 
 */
func (apiError *GalasaAPIError) UnmarshalApiError( response *http.Response) (error){
	var err error
	var body []byte
	body, err = io.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &apiError)
	}
	return err
}