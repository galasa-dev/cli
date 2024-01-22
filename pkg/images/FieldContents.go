/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

type FieldContents struct {
    Characters []string `json:"chars"`
    Text       string `json:"text"`
}

func (fieldContents *FieldContents) getCharacters() []rune {
    var contents []rune
    if (fieldContents.Characters != nil) {
        // If the terminal JSON defines the contents of a field in characters,
        // then convert the character strings into runes
        for _, char := range fieldContents.Characters {
            contents = append(contents, []rune(char)...)
        }
    } else {
        contents = []rune(fieldContents.Text)
    }
    return contents
}