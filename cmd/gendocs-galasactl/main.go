/*
 * Copyright contributors to the Galasa project
 */
package main

import (
	"log"
	"os"

	command "github.com/galasa.dev/cli/pkg/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd := command.RootCmd
	err := doc.GenMarkdownTree(cmd, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}
