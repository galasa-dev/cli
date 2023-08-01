/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

import (
	"os"
	"strings"
)

// -------------------------------------------------
// A console where things can be written to so the user
// can see them.
type Console interface {
	WriteString(text string) error
}

// -------------------------------------------------
// A real implementation which writes text to stdout.
type RealConsole struct {
}

func NewRealConsole() *RealConsole {
	return new(RealConsole)
}

func (*RealConsole) WriteString(text string) error {
	_, err := os.Stdout.WriteString(text)
	return err
}

// -------------------------------------------------
// A mock implementation which writes text to a buffer
// Useful for unit testing.
type MockConsole struct {
	text strings.Builder
}

func NewMockConsole() *MockConsole {
	console := new(MockConsole)
	console.text = strings.Builder{}
	return console
}

func (data *MockConsole) WriteString(text string) error {
	_, err := data.text.WriteString(text)
	return err
}

func (data *MockConsole) ReadText() string {
	return data.text.String()
}
