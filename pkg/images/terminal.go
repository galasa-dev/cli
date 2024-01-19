/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

type Terminal struct {
    Id          string 			`json:"id"`
    RunId       string 			`json:"runId"`
    Sequence    int 			`json:"sequence"`
    Images      []TerminalImage `json:"images"`
    DefaultSize TerminalSize    `json:"defaultSize"`
}
