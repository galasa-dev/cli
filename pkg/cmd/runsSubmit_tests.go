/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
*/
package cmd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRunsSubmitWithoutFlagsErrors(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"runs", "submit", "local"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Should throw an error asking for flags to be set
	assert.NotNil(t, err, "err should have been set!")
	assert.Contains(t, err.Error(), "GAL1009E: The submit command requires either test selection flags (eg: --stream, --class, --bundle, --package, --tag, --regex, --test) or --portfolio flag to be specified. Use the --help flag for more details.")
}

func TestRunsSubmitExecutesWithPortfolio(t *testing.T) {
	// Given...
	factory := NewMockFactory()
	var args []string = []string{"runs", "submit", "local", "--class", "osgi.bundle/class.path"}

	// When...
	err := Execute(factory, args)

	// Then...
	// Should throw an error asking for flags to be set
	assert.NotNil(t, err, "err should have been set!")
	assert.Contains(t, err.Error(), "required flag(s) \"obr\" not set")
}