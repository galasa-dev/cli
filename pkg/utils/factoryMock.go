/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"github.com/galasa-dev/cli/pkg/files"
)

type MockFactory struct {
	FinalWordHandler FinalWordHandler
	FileSystem       files.FileSystem
	Env              Environment
	StdOutConsole    Console
	StdErrConsole    Console
	TimeService      TimeService
	Authenticator    Authenticator
}

func NewMockFactory() *MockFactory {
	return new(MockFactory)
}

func (factory *MockFactory) GetFileSystem() files.FileSystem {
	if factory.FileSystem == nil {
		factory.FileSystem = files.NewMockFileSystem()
	}
	return factory.FileSystem
}

func (factory *MockFactory) GetEnvironment() Environment {
	if factory.Env == nil {
		factory.Env = NewMockEnv()
	}
	return factory.Env
}

func (factory *MockFactory) GetFinalWordHandler() FinalWordHandler {
	if factory.FinalWordHandler == nil {
		factory.FinalWordHandler = NewMockFinalWordHandler()
	}
	return factory.FinalWordHandler
}

func (factory *MockFactory) GetStdOutConsole() Console {
	if factory.StdOutConsole == nil {
		factory.StdOutConsole = NewMockConsole()
	}
	return factory.StdOutConsole
}

func (factory *MockFactory) GetStdErrConsole() Console {
	if factory.StdErrConsole == nil {
		factory.StdErrConsole = NewMockConsole()
	}
	return factory.StdErrConsole
}

func (factory *MockFactory) GetTimeService() TimeService {
	if factory.TimeService == nil {
		factory.TimeService = NewMockTimeService()
	}
	return factory.TimeService
}

func (factory *MockFactory) GetAuthenticator(apiServerUrl string, galasaHome GalasaHome) Authenticator {
	if factory.Authenticator == nil {
		factory.Authenticator = NewMockAuthenticator()
	}
	return factory.Authenticator
}
