/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"os"
	"os/user"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

// Environment is a thin interface layer above the os package which can be mocked out
type Environment interface {
	GetEnv(propertyName string) string
	GetUserName() (string, error)
}

//------------------------------------------------------------------------------------
// The implementation of the real os-delegating variant of the interface
//------------------------------------------------------------------------------------

type OSEnvironment struct {
}

// NewOSEnvironment creates a real wrapper over the os environment
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

func (osEnv OSEnvironment) GetUserName() (string, error) {
	name, err := user.Current()
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_RETRIEVING_USERNAME_FAILED, err.Error())
	}
	return name.Username, err
}
