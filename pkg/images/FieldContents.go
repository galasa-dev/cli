/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

type FieldContents struct {
    Characters []rune `json:"chars"`
    Text       string `json:"text"`
}

func (fieldContents *FieldContents) getCharacters() []rune {
    var contents []rune
    if (fieldContents.Characters != nil) {
        contents = fieldContents.Characters
    } else {
        contents = []rune(fieldContents.Text)
    }
    return contents
}