/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"log"
	"path/filepath"

	"github.com/galasa-dev/cli/pkg/files"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/props"
	"github.com/galasa-dev/cli/pkg/utils"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

const (
	CLIENT_ID_PROPERTY    = "GALASA_CLIENT_ID"
	SECRET_PROPERTY       = "GALASA_SECRET"
	ACCESS_TOKEN_PROPERTY = "GALASA_ACCESS_TOKEN"
)

// Gets authentication properties from the user's galasactl.properties file or from the environment or a mixture.
func GetAuthProperties(fileSystem files.FileSystem, galasaHome utils.GalasaHome, env utils.Environment) (galasaapi.AuthProperties, error) {
	var err error = nil

	// Work out which file we we want to draw properties from.
	galasactlPropertiesFilePath := filepath.Join(galasaHome.GetNativeFolderPath(), "galasactl.properties")

	// Get the file-based properties if we can
	authProperties, fileAccessErr := getAuthPropertiesFromFile(fileSystem, galasactlPropertiesFilePath, env)
	if fileAccessErr != nil {
		authProperties = *galasaapi.NewAuthProperties()
	}

	// We now have a structure which may be filled-in with values from the file.
	// Over-write those values if there is an environment variable set to do that.
	authProperties.SetClientId(getPropertyWithOverride(env, authProperties.GetClientId(), galasactlPropertiesFilePath, CLIENT_ID_PROPERTY))
	authProperties.SetRefreshToken(getPropertyWithOverride(env, authProperties.GetRefreshToken(), galasactlPropertiesFilePath, ACCESS_TOKEN_PROPERTY))
	authProperties.SetSecret(getPropertyWithOverride(env, authProperties.GetSecret(), galasactlPropertiesFilePath, SECRET_PROPERTY))

	// Make sure all the properties have values that we need.
	err = checkPropertyIsSet(authProperties.GetClientId(), CLIENT_ID_PROPERTY, galasactlPropertiesFilePath, fileAccessErr)
	if err == nil {
		err = checkPropertyIsSet(authProperties.GetRefreshToken(), ACCESS_TOKEN_PROPERTY, galasactlPropertiesFilePath, fileAccessErr)
		if err == nil {
			err = checkPropertyIsSet(authProperties.GetSecret(), SECRET_PROPERTY, galasactlPropertiesFilePath, fileAccessErr)
		}
	}

	return authProperties, err
}

func checkPropertyIsSet(propertyValue string, propertyName string, galasactlPropertiesFilePath string, fileAccessErr error) error {
	var err error
	if propertyValue == "" {
		// Property has not been set.
		if fileAccessErr != nil {
			// Being unable to read the file was the cause of this.
			err = fileAccessErr
		} else {
			log.Printf("Error: Auth property '%s' has not been set", propertyName)
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_AUTH_PROPERTY_NOT_AVAILABLE, propertyName, galasactlPropertiesFilePath)
		}
	}
	return err
}

func getPropertyWithOverride(env utils.Environment, valueFromFile string, filePathGatheredFrom string, propertyName string) string {
	value := env.GetEnv(propertyName)
	if value != "" {
		// env var has been set.
		if valueFromFile == "" {
			log.Printf("environment variable '%s' over-rides a value from file '%s'", propertyName, filePathGatheredFrom)
		} else {
			log.Printf("environment variable '%s' used to control authentication.", propertyName)
		}
	} else {
		value = valueFromFile
	}
	return value
}

// Gets authentication properties from the user's galasactl.properties file
func getAuthPropertiesFromFile(fileSystem files.FileSystem, galasactlPropertiesFilePath string, env utils.Environment) (galasaapi.AuthProperties, error) {
	var err error = nil
	authProperties := galasaapi.NewAuthProperties()

	var galasactlProperties props.JavaProperties
	galasactlProperties, err = props.ReadPropertiesFile(fileSystem, galasactlPropertiesFilePath)
	if err == nil {

		if err == nil {
			authProperties.SetClientId(galasactlProperties[CLIENT_ID_PROPERTY])
			authProperties.SetSecret(galasactlProperties[SECRET_PROPERTY])
			authProperties.SetRefreshToken(galasactlProperties[ACCESS_TOKEN_PROPERTY])
		}

	} else {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_FAILED_TO_READ_FILE, galasactlPropertiesFilePath, err.Error())
	}

	return *authProperties, err
}
