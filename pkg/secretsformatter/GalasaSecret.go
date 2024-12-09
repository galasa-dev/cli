/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package secretsformatter

import (
	"time"

	"github.com/galasa-dev/cli/pkg/galasaapi"
)

// The auto-generated OpenAPI structs don't include `yaml` annotations, which causes
// issues when it comes to marshalling data into GalasaSecret structs in order to display
// secrets in YAML format. This is a manually-maintained struct that includes `yaml` annotations.
type GalasaSecret struct {
	ApiVersion *string             `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind *string                   `json:"kind,omitempty"       yaml:"kind,omitempty"`
	Metadata *GalasaSecretMetadata `json:"metadata,omitempty"   yaml:"metadata,omitempty"`
	Data *GalasaSecretData         `json:"data,omitempty"       yaml:"data,omitempty"`
}

type GalasaSecretMetadata struct {
	Name *string                     `json:"name,omitempty"            yaml:"name,omitempty"`
	Description *string              `json:"description,omitempty"     yaml:"description,omitempty"`
	LastUpdatedTime *time.Time       `json:"lastUpdatedTime,omitempty" yaml:"lastUpdatedTime,omitempty"`
	LastUpdatedBy *string            `json:"lastUpdatedBy,omitempty"   yaml:"lastUpdatedBy,omitempty"`
	Encoding *string                 `json:"encoding,omitempty"        yaml:"encoding,omitempty"`
	Type *galasaapi.GalasaSecretType `json:"type,omitempty"            yaml:"type,omitempty"`
}

type GalasaSecretData struct {
	Username *string `json:"username,omitempty" yaml:"username,omitempty"`
	Password *string `json:"password,omitempty" yaml:"password,omitempty"`
	Token *string    `json:"token,omitempty"    yaml:"token,omitempty"`
}

func NewGalasaSecret(secret galasaapi.GalasaSecret) *GalasaSecret {
	return &GalasaSecret{
		ApiVersion: secret.ApiVersion,
		Kind: secret.Kind,
		Metadata: NewGalasaSecretMetadata(secret.Metadata),
		Data: NewGalasaSecretData(secret.Data),
	}
}

func NewGalasaSecretMetadata(metadata *galasaapi.GalasaSecretMetadata) *GalasaSecretMetadata {
	return &GalasaSecretMetadata{
		Name: metadata.Name,
		Description: metadata.Description,
		LastUpdatedTime: metadata.LastUpdatedTime,
		LastUpdatedBy: metadata.LastUpdatedBy,
		Encoding: metadata.Encoding,
		Type: metadata.Type,
	}
}

func NewGalasaSecretData(data *galasaapi.GalasaSecretData) *GalasaSecretData {
	return &GalasaSecretData{
		Username: data.Username,
		Password: data.Password,
		Token: data.Token,
	}
}