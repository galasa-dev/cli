/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

import (
	"log"

	"github.com/galasa-dev/cli/pkg/files"
)

type ImageFileWriter interface {
	WriteImageFile(simpleFileName string, imageBytes []byte) error
	GetImageFilesWrittenCount() int
}

type ImageFileWriterImpl struct {
	fs                     files.FileSystem
	imageFolderPath        string
	imageFilesWrittenCount int
}

func NewImageFileWriter(fs files.FileSystem, imageFolderPath string) ImageFileWriter {
	writer := new(ImageFileWriterImpl)
	writer.fs = fs
	writer.imageFolderPath = imageFolderPath
	writer.imageFilesWrittenCount = 0
	return writer
}

func (writer *ImageFileWriterImpl) GetImageFilesWrittenCount() int {
	return writer.imageFilesWrittenCount
}

func (writer *ImageFileWriterImpl) WriteImageFile(simpleFileName string, imageBytes []byte) error {
	var err error
	targetImageFilePath := writer.imageFolderPath + writer.fs.GetFilePathSeparator() + simpleFileName
	if err == nil {
		var isExistsAlready bool
		isExistsAlready, err = writer.fs.Exists(targetImageFilePath)
		if err == nil {
			if isExistsAlready {
				log.Printf("File %s already exists. So not over-writing it.\n", targetImageFilePath)
			} else {
				err = writer.fs.WriteBinaryFile(targetImageFilePath, imageBytes)
				if err == nil {
					writer.imageFilesWrittenCount = writer.imageFilesWrittenCount + 1
					log.Printf("Image file has been created: %s\n", targetImageFilePath)
				}
			}
		}
	}
	return err
}
