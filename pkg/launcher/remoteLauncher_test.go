/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package launcher

import (
	"testing"

	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/stretchr/testify/assert"
)

func TestProcessingGoodPropertiesExtractsStreamsOk(t *testing.T) {

	var inputProperties []galasaapi.GalasaProperty = make([]galasaapi.GalasaProperty, 0)

	name1 := "thames"
	name1full := "test.stream." + name1 + ".repo"
	name2 := "avon"
	name2full := "test.stream." + name2 + ".repo"

	inputProperties = append(inputProperties, galasaapi.GalasaProperty{
		Metadata: &galasaapi.GalasaPropertyMetadata{
			Name: &name1full,
		},
	})

	inputProperties = append(inputProperties, galasaapi.GalasaProperty{
		Metadata: &galasaapi.GalasaPropertyMetadata{
			Name: &name2full,
		},
	})

	streams, err := getStreamNamesFromProperties(inputProperties)
	assert.Nil(t, err)
	assert.NotNil(t, streams)
	assert.Equal(t, 2, len(streams))

	assert.Equal(t, streams[0], name1)
	assert.Equal(t, streams[1], name2)
}

func TestProcessingEmptyPropertiesListExtractsZeroStreamsOk(t *testing.T) {

	var inputProperties []galasaapi.GalasaProperty = make([]galasaapi.GalasaProperty, 0)

	streams, err := getStreamNamesFromProperties(inputProperties)

	assert.Nil(t, err)
	assert.NotNil(t, streams)
	assert.Equal(t, 0, len(streams))
}


