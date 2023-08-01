/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package runs

import (
	"testing"

	"github.com/galasa.dev/cli/pkg/files"
	"github.com/stretchr/testify/assert"
)

func TestCanWriteAndReadAPortfolio(t *testing.T) {
	portfolio := NewPortfolio()

	testSelection := new(TestSelection)
	testSelection.Classes = make([]TestClass, 0)

	testSelection.Classes = append(testSelection.Classes, TestClass{
		Bundle: "myBundle",
		Class:  "myClass",
		Stream: "myStream",
	})

	testOverrides := make(map[string]string)

	AddClassesToPortfolio(testSelection, &testOverrides, portfolio)

	fs := files.NewMockFileSystem()
	err := WritePortfolio(fs, "my.portfolio", portfolio)
	assert.Nil(t, err)

	portfolioGotBack, err := ReadPortfolio(fs, "my.portfolio")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(portfolio.Classes))
	assert.Equal(t, "myBundle", portfolioGotBack.Classes[0].Bundle)
	assert.Equal(t, "myClass", portfolioGotBack.Classes[0].Class)
	assert.Equal(t, "myStream", portfolioGotBack.Classes[0].Stream)

}
