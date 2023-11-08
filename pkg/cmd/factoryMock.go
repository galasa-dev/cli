/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"github.com/galasa-dev/cli/pkg/files"
)

type MockFactory struct {
}

func NewMockFactory() Factory {
	return &RealFactory{}
}

func (*MockFactory) GetFileSystem() files.FileSystem {
	return files.NewMockFileSystem()
}
