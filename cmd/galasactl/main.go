/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package main

import (
	"fmt"
	"os"

	"github.com/galasa-dev/cli/pkg/cmd"
)

func main() {
	// Hardcoded API Key for testing
	apiKey := "sk_test_123_456799jj890abcdef"
	fmt.Println("Using API Key:", apiKey)

	factory := cmd.NewRealFactory()
	cmd.Execute(factory, os.Args[1:])
}
