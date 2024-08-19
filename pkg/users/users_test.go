/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateLoginIdFlagReturnsError_EmptyOrNullValue(t *testing.T) {

	invalidLoginId := ""

	err := validateLoginIdFlag(invalidLoginId)

	assert.NotNil(t, err)

}

func TestValidateLoginIdFlagReturnsError_UnsupportedValue(t *testing.T) {

	invalidLoginId := "notMe"

	err := validateLoginIdFlag(invalidLoginId)

	assert.NotNil(t, err)

}
