/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

import (
	"log"
	"strings"

	"github.com/galasa-dev/cli/pkg/files"
)

// Given a root folder, we scan for .gz files which need expansion into images.

type ImageExpander interface {
	ExpandImages(rootFolderPath string) error
	GetExpandedImageFileCount() int
}

type ImageExpanderImpl struct {
	fs                  files.FileSystem
	renderer            ImageRenderer
	expandedFileCounter int
}

func NewImageExpander(fs files.FileSystem, renderer ImageRenderer) ImageExpander {
	expander := new(ImageExpanderImpl)
	expander.fs = fs
	expander.renderer = renderer
	expander.expandedFileCounter = 0
	return expander
}

func (expander *ImageExpanderImpl) GetExpandedImageFileCount() int {
	return expander.expandedFileCounter
}

func (expander *ImageExpanderImpl) ExpandImages(rootFolderPath string) error {
	var err error

	var paths []string
	paths, err = expander.fs.GetAllFilePaths(rootFolderPath)

	if err == nil {
		for _, filePath := range paths {
			if strings.HasSuffix(filePath, ".gz") {
				// It's a gz file we might like to expand.

				err = expander.expandGzFile(filePath)
				if err != nil {
					break
				}
			}
		}
	}

	return err
}

func (expander *ImageExpanderImpl) expandGzFile(gzFilePath string) error {
	var err error
	var targetImageFilePath string
	var targetImageFolderPath string
	targetImageFilePath, targetImageFolderPath, err = expander.calculateTargetImagePaths(gzFilePath)

	gzip := files.NewGzipFile(expander.fs, gzFilePath)
	var binaryContent []byte
	binaryContent, err = gzip.ReadBytes()

	if err == nil {
		var imageBytes []byte
		imageBytes, err = expander.renderer.RenderJsonBytesToImageBytes(binaryContent)

		if err == nil {
			err = expander.fs.MkdirAll(targetImageFolderPath)
			if err == nil {
				var isExistsAlready bool
				isExistsAlready, err = expander.fs.Exists(targetImageFilePath)
				if err == nil {
					if isExistsAlready {
						log.Printf("File %s already exists.\n", targetImageFilePath)
					} else {
						err = expander.fs.WriteBinaryFile(targetImageFilePath, imageBytes)
						if err == nil {
							expander.expandedFileCounter = expander.expandedFileCounter + 1
						}
					}
				}
			}
		}
	}

	return err
}

func (expander *ImageExpanderImpl) calculateTargetImagePaths(gzFilePath string) (string, string, error) {
	var err error
	var desiredImageFilePath string
	var desiredImageFolderPath string

	// Figure out the file path of the image we want to create.
	// Into a folder called "images" next to the .gz folder
	separator := expander.fs.GetFilePathSeparator()
	filePathParts := strings.Split(gzFilePath, separator)

	if len(filePathParts)-3 < 0 {
		log.Printf("gz file %s found, but it's not in a 'terminals/termXXX' folder so ignoring.\n", gzFilePath)
	} else {
		// The json descriptions of the panels appear in zos3270/terminals/term1/term1-0001.gz
		// So we want images to appear in zos3270/images/term1/term1-0001.png
		if filePathParts[len(filePathParts)-3] != "terminals" {
			// It's not in the correct structure, so ignoring.
			log.Printf("gz file %s found, but it's not in a 'terminals' folder so ignoring.\n", gzFilePath)
		} else {

			// Replace the .gz file extension with .png
			simpleFileName := filePathParts[len(filePathParts)-1]
			indexOfExtension := strings.LastIndex(simpleFileName, ".gz")
			simpleTargetFileName := simpleFileName[:indexOfExtension] + ".png"
			filePathParts[len(filePathParts)-1] = simpleTargetFileName

			// Replace the "terminals" part of the path with "images"
			filePathParts[len(filePathParts)-3] = "images"

			// Roll together all the parts of the files to get the folder and the file
			desiredImageFilePath = strings.Join(filePathParts, separator)
			desiredImageFolderPath = strings.Join(filePathParts[:len(filePathParts)-1], separator)

		}
	}
	return desiredImageFilePath, desiredImageFolderPath, err
}
