/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import "github.com/galasa-dev/cli/pkg/utils"

type MockFinalWordHandler struct {
	ReportedObject interface{}
}

func NewMockFinalWordHandler() utils.FinalWordHandler {
	return new(MockFinalWordHandler)
}

func (handler *MockFinalWordHandler) FinalWord(rootCmd utils.GalasaCommand, errorToExctractFrom interface{}) {
	// Capture the final word object to see what was sent.
	handler.ReportedObject = errorToExctractFrom
}
