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

func (apiError *GalasaAPIError) SetGalasaAPIError( response *http.Response) (error){
	var err error
	var body []byte
	body, err = io.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &apiError)
	}
	return err
}