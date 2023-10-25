/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
*/
package utils

func ConvertToYaml(kind string, cpsnamespace string, name string, value string) string {
	returnYaml := `Kind: ` + kind + `
metadata:
	cpsnamespace: ` + cpsnamespace + `
	name: ` + name + `
data:
	value: ` + value + `
`

	return returnYaml
}