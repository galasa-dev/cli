/*
 * Copyright contributors to the Galasa project
 */
package utils

func AddJavaRuntimeToMock(fileSystem FileSystem, baseJavaFolderName string) {
	separator := fileSystem.GetFilePathSeparator()
	fileSystem.MkdirAll(baseJavaFolderName + separator + "bin")
	fileSystem.WriteBinaryFile(baseJavaFolderName+separator+
		"bin"+separator+"java"+fileSystem.GetExecutableExtension(),
		[]byte("some random content pretending to be a JRE program"),
	)
}
