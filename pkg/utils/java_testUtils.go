/*
 * Copyright contributors to the Galasa project
 */
package utils

func AddJavaRuntimeToMock(fileSystem FileSystem, baseJavaFolderName string) {
	fileSystem.MkdirAll(baseJavaFolderName + "/bin")
	fileSystem.WriteBinaryFile(baseJavaFolderName+"/bin/java", []byte("some random content pretending to be a JRE program"))
}
