/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
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

// Gets authentication properties from the user's galasactl.properties file
func GetAuthProperties(fileSystem files.FileSystem, galasaHome utils.GalasaHome) (galasaapi.AuthProperties, error) {
    var err error = nil
    authProperties := galasaapi.NewAuthProperties()

    galasactlPropertiesFilePath := filepath.Join(galasaHome.GetNativeFolderPath(), "galasactl.properties")
    
    var galasactlProperties props.JavaProperties
    galasactlProperties, err = props.ReadPropertiesFile(fileSystem, galasactlPropertiesFilePath)
    if err == nil {
        requiredAuthProperties := getAuthPropertiesList()
        err = validateRequiredGalasactlProperties(requiredAuthProperties, galasactlProperties)

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

// Ensures the provided galasactl properties contain values for the required properties, returning an error if a property is missing
func validateRequiredGalasactlProperties(requiredProperties []string, galasactlProperties props.JavaProperties) error {
    var err error = nil
    for _, property := range requiredProperties {
        if galasactlProperties[property] == "" {
            err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_MISSING_GALASACTL_PROPERTY, property)
            break
        }
    }
    return err
}

// Returns a list of auth properties
func getAuthPropertiesList() []string {
	return []string{CLIENT_ID_PROPERTY, SECRET_PROPERTY, ACCESS_TOKEN_PROPERTY}
}
