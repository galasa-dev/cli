/*
 * Copyright contributors to the Galasa project
 */
package embedded

import "embed"

// Embed all the template files into the go executable, so there are no extra files
// we need to ship/install/locate on the target machine.
// We can access the "embedded" file system as if they are normal files.
//
//go:embed templates/*
var embeddedFileSystem embed.FS

func GetEmbeddedFileSystem() embed.FS {
	return embeddedFileSystem
}

func GetGalasaVersion() string {
	// Ideally, the build process would create something which go embeds, and we can read and return here.
	return "0.25.0"
}

func GetBootJarVersion() string {
	// Ideally, the build process would create something which go embeds, and we can read and return here.
	return "0.24.0"
}
