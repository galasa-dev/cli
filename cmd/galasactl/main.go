//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package main

import (
	"fmt"

	"github.com/galasa.dev/cli/pkg/cli"
	"github.com/galasa.dev/cli/pkg/cmd"
)

func main() {
	gp := &cli.GalasaParams{}
	galasa := cmd.Root(gp)
	_ = galasa

	fmt.Println("Boo")
}
