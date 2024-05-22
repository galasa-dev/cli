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
	finalWordHandler utils.FinalWordHandler
	fileSystem       files.FileSystem
	env              utils.Environment
	stdOutConsole    utils.Console
	stdErrConsole    utils.Console
	timeService      utils.TimeService
}

func NewMockFactory() utils.Factory {
	return &MockFactory{}
}

func (factory *MockFactory) GetFileSystem() files.FileSystem {
	if factory.fileSystem == nil {
		factory.fileSystem = files.NewMockFileSystem()
	}
	return factory.fileSystem
}

func (factory *MockFactory) GetEnvironment() utils.Environment {
	if factory.env == nil {
		factory.env = utils.NewMockEnv()
	}
	return factory.env
}

func (factory *MockFactory) GetFinalWordHandler() utils.FinalWordHandler {
	if factory.finalWordHandler == nil {
		factory.finalWordHandler = NewMockFinalWordHandler()
	}
	return factory.finalWordHandler
}

func (factory *MockFactory) GetStdOutConsole() utils.Console {
	if factory.stdOutConsole == nil {
		factory.stdOutConsole = utils.NewMockConsole()
	}
	return factory.stdOutConsole
}

func (factory *MockFactory) GetStdErrConsole() utils.Console {
	if factory.stdErrConsole == nil {
		factory.stdErrConsole = utils.NewMockConsole()
	}
	return factory.stdErrConsole
}

func (factory *MockFactory) GetTimeService() utils.TimeService {
	if factory.timeService == nil {
		factory.timeService = utils.NewMockTimeService()
	}
	return factory.timeService
}
