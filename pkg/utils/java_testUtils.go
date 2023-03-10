/*
 * Copyright contributors to the Galasa project
 */
package utils

func AddJavaRuntimeToMock(fileSystem FileSystem, baseJavaFolderName string) {
	fileSystem.MkdirAll(baseJavaFolderName + FILE_SYSTEM_PATH_SEPARATOR + "bin")
	fileSystem.WriteBinaryFile(baseJavaFolderName+FILE_SYSTEM_PATH_SEPARATOR+
		"bin"+FILE_SYSTEM_PATH_SEPARATOR+"java"+fileSystem.GetExecutableExtension(),
		[]byte("some random content pretending to be a JRE program"),
	)
}
