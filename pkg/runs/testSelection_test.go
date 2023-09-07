/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------
// Functions
func TestAreSelectionFlagsProvidedDoesntPanicWhenFlagsAreNotSet(t *testing.T) {
	flags := NewTestSelectionFlags()
	areTheyProvided := AreSelectionFlagsProvided(flags)
	assert.False(t, areTheyProvided)
}

func TestAllocatingNewFlagsStructureHasEmptyArrays(t *testing.T) {
	flags := NewTestSelectionFlags()
	assert.NotNil(t, flags)
	assert.NotNil(t, flags.bundles)
	assert.Equal(t, len(*flags.bundles), 0)

	assert.NotNil(t, flags.packages)
	assert.Equal(t, len(*flags.packages), 0)

	assert.NotNil(t, flags.tests)
	assert.Equal(t, len(*flags.tests), 0)

	assert.NotNil(t, flags.tags)
	assert.Equal(t, len(*flags.tags), 0)

	assert.NotNil(t, flags.classes)
	assert.Equal(t, len(*flags.classes), 0)

	assert.NotNil(t, flags.regexSelect)
}
