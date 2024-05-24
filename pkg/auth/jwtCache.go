/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package auth

import (
	"log"
	"net/url"
	"time"

	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/golang-jwt/jwt/v5"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
)

const (
	// If a JWT only has less than a certain time left before it expires, don't use it.
	// This sets the boundary where tokens about-to-expire are never returned from the cache.
	TOKEN_EXPIRY_BUFFER_MINUTES = 10
)

type JwtCache interface {

	// Adds a jwt to the cache.
	Put(serverApiUrl string, galasaToken string, jwt string) error

	// Returns the jwt if we have one, or "" if not.
	Get(serverApiUrl string, galasaToken string) (jwt string, err error)

	// Clears the jwt if we have one in the cache. ie: Logs out.
	Clear(serverApiUrl string, galasaToken string) error

	// Clears all the cache content.
	ClearAll() error
}

type fileBasedJwtCache struct {
	fileSystem  spi.FileSystem
	galasaHome  spi.GalasaHome
	timeService spi.TimeService
}

func NewJwtCache(
	fileSystem spi.FileSystem,
	galasaHome spi.GalasaHome,
	timeService spi.TimeService,
) JwtCache {
	cache := new(fileBasedJwtCache)

	cache.fileSystem = fileSystem
	cache.galasaHome = galasaHome
	cache.timeService = timeService

	return cache
}

func (cache *fileBasedJwtCache) Put(serverApiUrl string, galasaToken string, jwt string) (err error) {
	file := utils.NewBearerTokenFile(cache.fileSystem, cache.galasaHome, cache.urlToFileName(serverApiUrl)+".json", cache.timeService)
	err = file.WriteJwt(jwt)
	return err
}

func (cache *fileBasedJwtCache) Clear(serverApiUrl string, galasaToken string) error {
	var err error
	file := utils.NewBearerTokenFile(cache.fileSystem, cache.galasaHome, cache.urlToFileName(serverApiUrl)+".json", cache.timeService)
	file.DeleteJwt()
	return err
}

func (cache *fileBasedJwtCache) ClearAll() error {
	err := utils.DeleteAllBearerTokenFiles(cache.fileSystem, cache.galasaHome)
	return err
}

// Gets the jwt from the cache, or returns a string if it's not present.
// Only jwts which are valid (not expired) are returned.
func (cache *fileBasedJwtCache) Get(serverApiUrl string, galasaToken string) (jwt string, err error) {
	var possiblyInvalidJwt string

	file := utils.NewBearerTokenFile(cache.fileSystem, cache.galasaHome, cache.urlToFileName(serverApiUrl)+".json", cache.timeService)

	var isExists bool
	isExists, err = file.Exists()
	if err == nil {

		if isExists {
			possiblyInvalidJwt, err = file.ReadJwt()

			if err == nil {
				var isValid bool
				isValid, err = cache.isBearerTokenValid(possiblyInvalidJwt, cache.timeService)
				if err == nil {
					if isValid {
						jwt = possiblyInvalidJwt
					}
				}
			}
		}
	}
	return jwt, err
}

// Converts an arbitrary URL to a name we can use as a filename.
func (cache *fileBasedJwtCache) urlToFileName(urlToConvert string) string {
	filename := url.QueryEscape(urlToConvert)
	return filename
}

// Checks whether a given bearer token is valid or not, returning true if it is valid and false otherwise
func (cache *fileBasedJwtCache) isBearerTokenValid(bearerTokenString string, timeService spi.TimeService) (bool, error) {
	var err error
	var bearerToken *jwt.Token
	var isValid bool = false

	// Decode the bearer token without verifying its signature
	bearerToken, _, err = jwt.NewParser().ParseUnverified(bearerTokenString, jwt.MapClaims{})
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_JWT_CANNOT_BE_PARSED, err.Error())
	} else {
		var tokenExpiry *jwt.NumericDate
		tokenExpiry, err = bearerToken.Claims.GetExpirationTime()
		if err != nil {
			err = galasaErrors.NewGalasaError(galasaErrors.GALASA_JWT_HAS_NO_EXPIRATION_DATETIME, err.Error())
		} else {
			// Add a buffer to the current time to make sure the bearer token does not expire within
			// this buffer (e.g. if the buffer is 10 mins, make sure the token doesn't expire within 10 mins)
			acceptableExpiryTime := timeService.Now().Add(time.Duration(TOKEN_EXPIRY_BUFFER_MINUTES) * time.Minute)
			if (tokenExpiry.Time).After(acceptableExpiryTime) {
				isValid = true
			} else {
				log.Printf("JWT token is valid, but due to expire shortly, so ignoring it.\n")
			}
		}
	}
	return isValid, err
}
