/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

// The main entry point into the cmd package.
func Execute(factory Factory, args []string) error {
	var err error

	finalWordHandler := factory.GetFinalWordHandler()

	var commands CommandCollection
	commands, err = NewCommandCollection(factory)

	if err == nil {

		// Catch execution if a panic happens.
		defer func() {
			err := recover()

			// Display the error and exit.
			finalWordHandler.FinalWord(commands.GetRootCommand(), err)
		}()

		// Execute the command
		err = commands.Execute(args)
	}
	finalWordHandler.FinalWord(commands.GetRootCommand(), err)
	return err
}