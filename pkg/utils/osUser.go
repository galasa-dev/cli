/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"os/user"
)

type Username interface {
	GetUsername() string
}

type OSUsername struct {
}

func NewOSUsername() *OSUsername {
	name := new(OSUsername)
	return name
}

func NewUsername() Username {
	return NewOSUsername()
}

func (osUser OSUsername) GetUsername() string {
	name, _ := user.Current()
	return name.Username
}
