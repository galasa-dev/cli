/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package launcher

type SubmitRunStructure struct {
	requestor        string
	obrFromPortfolio string
	isTraceEnabled   bool
	overrides        map[string]interface{}
}

type JavaSubmitRunStructure struct {
	SubmitRunStructure
	groupName   string
	className   string
	requestType string
	stream      string
}

type GherkinSubmitRunStructure struct {
	SubmitRunStructure
	gherkinFile string
}
