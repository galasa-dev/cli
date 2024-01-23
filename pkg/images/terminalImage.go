/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

type TerminalImage struct {
    Id           string 		 `json:"id"`
    Sequence     int 		     `json:"sequence"`
    Inbound      bool 		  	 `json:"inbound"`
    Type         string       	 `json:"type"`
    ImageSize    TerminalSize 	 `json:"imageSize"`
    CursorRow    int 		  	 `json:"cursorRow"`
    CursorColumn int 		  	 `json:"cursorColumn"`
    Aid 		 string 	     `json:"aid"`
    Fields 		 []TerminalField `json:"fields"`
}
