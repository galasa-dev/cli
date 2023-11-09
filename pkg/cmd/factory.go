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

type Factory interface {
	GetFileSystem() files.FileSystem
	GetEnvironment() utils.Environment
	GetFinalWordHandler() FinalWordHandler
	GetConsole() utils.Console
}

type RealFactory struct {
	console utils.Console
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
func (this *RealFactory) GetConsole() utils.Console {
	if this.console == nil {
		this.console = utils.NewRealConsole()
	}
	return this.console
}
