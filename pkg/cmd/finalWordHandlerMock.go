/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

type MockFinalWordHandler struct {
	ReportedObject interface{}
}

func NewMockFinalWordHandler() FinalWordHandler {
	return new(MockFinalWordHandler)
}

func (this *MockFinalWordHandler) FinalWord(rootCmd GalasaCommand, o interface{}) {
	// Capture the final word object to see what was sent.
	this.ReportedObject = o
}
