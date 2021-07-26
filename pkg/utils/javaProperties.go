/*
*  Licensed Materials - Property of IBM
*
* (c) Copyright IBM Corp. 2021.
*/

package utils

import (
	"bufio"
	"strings"
)

type JavaProperties map[string]string

func ReadProperties(propertyString string) (JavaProperties) {
	properties := JavaProperties{}

	scanner := bufio.NewScanner(strings.NewReader(propertyString))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		equalsPos := strings.Index(line, "=")
		if (equalsPos == -1) {
			continue
		}

		key := strings.TrimSpace(line[:equalsPos])
		if (key == "") {
			continue
		}

		value := strings.TrimSpace(line[equalsPos+1:])

		properties[key] = value
	}


	return properties
}