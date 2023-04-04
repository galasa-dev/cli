/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"log"
	"strings"
)

type GalasaHome interface {
	GetNativeFolderPath() string
	GetUrlFolderPath() string
}

type galasaHomeImpl struct {
	path string
	fs   FileSystem
	env  Environment
}

func NewGalasaHome(fs FileSystem, env Environment) (GalasaHome, error) {
	var err error = nil
	var homeData *galasaHomeImpl = nil

	galasaHomePath := env.GetEnv("GALASA_HOME")
	if galasaHomePath == "" {
		var userHome string
		userHome, err = fs.GetUserHomeDirPath()
		if err == nil {
			galasaHomePath = userHome + fs.GetFilePathSeparator() + ".galasa"
		}
	} else {
		err = validateUserHomeDir(galasaHomePath, fs)
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

func validateUserHomeDir(path string, fs FileSystem) error {
	var err error = nil

	return err
}

func (homeData *galasaHomeImpl) GetNativeFolderPath() string {
	return homeData.path
}

func (homeData *galasaHomeImpl) GetUrlFolderPath() string {
	return strings.ReplaceAll(homeData.path, "\\", "/")
}
