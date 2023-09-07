/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"os"
)

// Environment is a thin interface layer above the os package which can be mocked out
type Environment interface {
	GetEnv(propertyName string) string
}

//------------------------------------------------------------------------------------
// The implementation of the real os-delegating variant of the interface
//------------------------------------------------------------------------------------

type OSEnvironment struct {
}

// NewOSFileSystem creates an implementation of the thin file system layer which delegates
// to the real os package calls.
func NewOSEnvironment() *OSEnvironment {
	env := new(OSEnvironment)
	return env
}

func NewEnvironment() Environment {
	return NewOSEnvironment()
}

//------------------------------------------------------------------------------------
// Interface methods...
//------------------------------------------------------------------------------------

func (osEnv OSEnvironment) GetEnv(propertyName string) string {
	return os.Getenv(propertyName)
}
