/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"log"
	"strings"

	"github.com/galasa-dev/cli/pkg/files"
)

type GalasaHome interface {
	GetNativeFolderPath() string
	GetUrlFolderPath() string
}

type galasaHomeImpl struct {
	path string
	fs   files.FileSystem
	env  Environment
}

func NewGalasaHome(fs files.FileSystem, env Environment, cmdFlagGalasaHome string) (GalasaHome, error) {
	var err error
	var homeData *galasaHomeImpl = nil

	galasaHomePath := cmdFlagGalasaHome
	if galasaHomePath == "" {
		galasaHomePath = env.GetEnv("GALASA_HOME")
		if galasaHomePath == "" {
			var userHome string
			userHome, err = fs.GetUserHomeDirPath()
			if err == nil {
				galasaHomePath = userHome + fs.GetFilePathSeparator() + ".galasa"
			}
		}
	}

	if err == nil {
		// All is well, so lets allocate a structure to pack with data.
		homeData = new(galasaHomeImpl)
		homeData.fs = fs
		homeData.env = env
		homeData.path = galasaHomePath

		log.Printf("Galasa home is '%s'", galasaHomePath)
	}

	return homeData, err
}

func (homeData *galasaHomeImpl) GetNativeFolderPath() string {
	return homeData.path
}

func (homeData *galasaHomeImpl) GetUrlFolderPath() string {
	return strings.ReplaceAll(homeData.path, "\\", "/")
}
