/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package props

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/galasa.dev/cli/pkg/files"
)

type JavaProperties map[string]string

func ReadPropertiesFile(fs files.FileSystem, filePath string) (JavaProperties, error) {
	var properties JavaProperties
	contents, err := fs.ReadTextFile(filePath)
	if err == nil {
		properties = ReadProperties(contents)
	}
	return properties, err
}

func ReadProperties(propertyString string) JavaProperties {
	properties := JavaProperties{}

	scanner := bufio.NewScanner(strings.NewReader(propertyString))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		equalsPos := strings.Index(line, "=")
		if equalsPos == -1 {
			continue
		}

		key := strings.TrimSpace(line[:equalsPos])
		if key == "" {
			continue
		}

		value := strings.TrimSpace(line[equalsPos+1:])

		properties[key] = value
	}

	return properties
}

func WritePropertiesFile(fs files.FileSystem, path string, properties map[string]interface{}) error {
	var err error = nil

	buff := new(bytes.Buffer)

	// Extract all the property keys
	keys := make([]string, 0)
	for k, _ := range properties {
		keys = append(keys, k)
	}

	// Sort the property keys into sort order.
	sort.Strings(keys)

	// Write out the properties
	for _, key := range keys {
		buff.WriteString(fmt.Sprintf("%s=%v\n", key, properties[key]))
	}

	// Write it all out to a file.
	contents := buff.String()
	err = fs.WriteTextFile(path, contents)

	log.Printf("Properties file %s written containing this:\n%s", path, contents)

	return err
}
