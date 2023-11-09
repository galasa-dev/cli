/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/utils"
)

type MockFactory struct {
	finalWordHandler FinalWordHandler
	fileSystem       files.FileSystem
	env              utils.Environment
	console          utils.Console
}

func NewMockFactory() Factory {
	return &MockFactory{}
}

func (this *MockFactory) GetFileSystem() files.FileSystem {
	if this.fileSystem == nil {
		this.fileSystem = files.NewMockFileSystem()
	}
	return this.fileSystem
}

func (this *MockFactory) GetEnvironment() utils.Environment {
	if this.env == nil {
		this.env = utils.NewMockEnv()
	}
	return this.env
}

func (this *MockFactory) GetFinalWordHandler() FinalWordHandler {
	if this.finalWordHandler == nil {
		this.finalWordHandler = NewMockFinalWordHandler()
	}
	return this.finalWordHandler
}

func (this *MockFactory) GetConsole() utils.Console {
	if this.console == nil {
		this.console = utils.NewMockConsole()
	}
	return this.console
}
