/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package streams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStreamNameFlagReturnsError_EmptyOrNullValue(t *testing.T) {

	invalidStreamName := ""

	loginId, err := validateStreamName(invalidStreamName)

	assert.NotNil(t, err)
	assert.Empty(t, loginId)
}
