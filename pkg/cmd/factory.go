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

// We use the factory to create instances of various classes.
// Some are cached, so you get the same one back each time
// Some are fresh objects created each time.
// We do this so we can have a real and a mock implementation
// to make unit testing easier.
type Factory interface {
	GetFileSystem() files.FileSystem
	GetEnvironment() utils.Environment
	GetFinalWordHandler() FinalWordHandler
	GetStdOutConsole() utils.Console
	GetStdErrConsole() utils.Console
	GetTimeService() utils.TimeService
}

// Allocates real objects with real implementations,
// none of which are generally great for unit testing.
// eg: A real file system can leave debris behind when a test runs.
type RealFactory struct {
	stdOutConsole utils.Console
	stdErrConsole utils.Console
}

func NewRealFactory() Factory {
	return &RealFactory{}
}

func (*RealFactory) GetFileSystem() files.FileSystem {
	return files.NewOSFileSystem()
}

func (*RealFactory) GetEnvironment() utils.Environment {
	return utils.NewEnvironment()
}

func (*RealFactory) GetFinalWordHandler() FinalWordHandler {
	return NewRealFinalWordHandler()
}

// We only ever expect there to be a single console object, which collects all the
// command output.
func (this *RealFactory) GetStdOutConsole() utils.Console {
	if this.stdOutConsole == nil {
		this.stdOutConsole = utils.NewRealConsole()
	}
	return this.stdOutConsole
}

func (this *RealFactory) GetStdErrConsole() utils.Console {
	if this.stdErrConsole == nil {
		this.stdErrConsole = utils.NewRealConsole()
	}
	return this.stdErrConsole
}

func (*RealFactory) GetTimeService() utils.TimeService {
	return utils.NewRealTimeService()
}
