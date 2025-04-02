/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package streamsformatter

import (
	"strings"

	"github.com/galasa-dev/cli/pkg/galasaapi"
	"gopkg.in/yaml.v3"
)

const (
	YAML_FORMATTER_NAME = "yaml"
)

type StreamsYamlFormatter struct {
}

func NewStreamsYamlFormatter() StreamsFormatter {
	return new(StreamsYamlFormatter)
}

func (*StreamsYamlFormatter) GetName() string {
	return YAML_FORMATTER_NAME
}

func (*StreamsYamlFormatter) FormatStreams(streams []galasaapi.Stream) (string, error) {
	var err error
	buff := strings.Builder{}

	for index, stream := range streams {
		if index > 0 {
			buff.WriteString("---\n")
		}

		type ObrInfo struct {
			MavenGroupID    string `yaml:"maven-group-id"`
			MavenArtifactID string `yaml:"maven-artifact-id"`
			MavenVersion    string `yaml:"maven-version"`
		}

		type RepositoryURL struct {
			URL string `yaml:"url"`
		}

		type StreamYAML struct {
			APIVersion string `yaml:"apiVersion"`
			Kind       string `yaml:"kind"`
			Metadata   struct {
				Name        string `yaml:"name"`
				Description string `yaml:"description"`
			} `yaml:"metadata"`
			Data struct {
				Repository  []RepositoryURL `yaml:"repository"`
				TestCatalog []RepositoryURL `yaml:"testCatalog"`
				Obrs        []ObrInfo       `yaml:"obrs"`
			} `yaml:"data"`
		}

		streamYAML := StreamYAML{
			APIVersion: "galasa-dev/v1alpha1",
			Kind:       "GalasaStream",
		}
		streamYAML.Metadata.Name = *stream.Metadata.Name
		streamYAML.Metadata.Description = *stream.Metadata.Description

		// Add repository and test catalog
		streamYAML.Data.Repository = []RepositoryURL{{URL: *stream.Data.Repository.Url}}
		streamYAML.Data.TestCatalog = []RepositoryURL{{URL: *stream.Data.TestCatalog.Url}}

		// Add the obrs section
		streamYAML.Data.Obrs = []ObrInfo{
			{
				MavenGroupID:    *stream.Data.Obrs[0].GroupId,
				MavenArtifactID: *stream.Data.Obrs[0].ArtifactId,
				MavenVersion:    *stream.Data.Obrs[0].Version,
			},
		}

		yamlBytes, err := yaml.Marshal(streamYAML)
		if err == nil {
			buff.Write(yamlBytes)
		}
	}

	return buff.String(), err
}
