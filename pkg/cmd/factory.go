/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/galasa-dev/cli/pkg/auth"
	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/utils"
)

// Allocates real objects with real implementations,
// none of which are generally great for unit testing.
// eg: A real file system can leave debris behind when a test runs.
type RealFactory struct {
	stdOutConsole utils.Console
	stdErrConsole utils.Console
}

func NewRealFactory() utils.Factory {
	return &RealFactory{}
}

func (*RealFactory) GetFileSystem() files.FileSystem {
	return files.NewOSFileSystem()
}

func (*RealFactory) GetEnvironment() utils.Environment {
	return utils.NewEnvironment()
}

func (*RealFactory) GetFinalWordHandler() utils.FinalWordHandler {
	return NewRealFinalWordHandler()
}

// We only ever expect there to be a single console object, which collects all the
// command output.
func (factory *RealFactory) GetStdOutConsole() utils.Console {
	if factory.stdOutConsole == nil {
		factory.stdOutConsole = utils.NewRealConsole()
	}
	return factory.stdOutConsole
}

func (factory *RealFactory) GetStdErrConsole() utils.Console {
	if factory.stdErrConsole == nil {
		factory.stdErrConsole = utils.NewRealConsole()
	}
	return factory.stdErrConsole
}

func (*RealFactory) GetTimeService() utils.TimeService {
	return utils.NewRealTimeService()
}

func (factory *RealFactory) GetAuthenticator(apiServerUrl string, galasaHome utils.GalasaHome) utils.Authenticator {
	jwtCache := auth.NewJwtCache(factory.GetFileSystem(), galasaHome, factory.GetTimeService())
	return auth.NewAuthenticator(apiServerUrl, factory.GetFileSystem(), galasaHome, factory.GetTimeService(), factory.GetEnvironment(), jwtCache)
}
