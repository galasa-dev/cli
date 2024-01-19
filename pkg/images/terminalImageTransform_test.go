/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"path/filepath"
	"testing"

	"github.com/galasa-dev/cli/pkg/files"
	"github.com/stretchr/testify/assert"
)

//----------------------------------------------------
// Utility functions
//----------------------------------------------------
func assertTerminalImageMatchesExpectedSnapshot(t *testing.T, fs files.FileSystem, actualImage *image.RGBA, terminalImage TerminalImage) {
    buf := new(bytes.Buffer)
    err := png.Encode(buf, actualImage)
    actualImageBytes := buf.Bytes()
    assert.Nil(t, err, "Image should successfully be encoded into PNG format")
    assert.NotEmpty(t, actualImageBytes, "Image data should not be empty")

    pngImageToCompareAgainst, err := filepath.Glob(filepath.Join("testdata", terminalImage.Id + ".png"))
    if err != nil || len(pngImageToCompareAgainst) == 0 {
        writeRenderedImageToTempDir(t, fs, terminalImage, actualImage)
        t.Fatalf("Failed to find expected image to compare against")
    }

    expectedFileBytes, err := fs.ReadBinaryFile(pngImageToCompareAgainst[0])
    assert.Nil(t, err, "Failed to read expected image file")

    actualImageSize := len(actualImageBytes)
    expectedImageSize := len(expectedFileBytes)
    if actualImageSize != expectedImageSize {
        writeRenderedImageToTempDir(t, fs, terminalImage, actualImage)
        t.Fatalf("Rendered image size '%d' does not match the expected image size '%d' ", actualImageSize, expectedImageSize)
    }

    for i, actualByte := range actualImageBytes {
        expectedByte := expectedFileBytes[i]
        if actualByte != expectedByte {
            writeRenderedImageToTempDir(t, fs, terminalImage, actualImage)
            t.Fatalf("Rendered image byte '%s' does not match expected image byte '%s'", string(actualByte), string(expectedByte))
        }
    }
}

func writeRenderedImageToTempDir(t *testing.T, fs files.FileSystem, terminalImage TerminalImage, actualImage *image.RGBA) {
    outputDirectory, err := fs.MkTempDir()
    if err == nil {
        err = WritePngImageToDisk(terminalImage, actualImage, fs, outputDirectory)
    }

    if err != nil {
        t.Log("Failed to write the rendered image to a temporary directory")
    } else {
        fmt.Printf("Rendered image written to: %s", filepath.Join(outputDirectory, terminalImage.Id + ".png"))
    }
}


//----------------------------------------------------
// Test functions
//----------------------------------------------------
func TestRenderEmptyTerminalRendersOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    terminalSize := TerminalSize{
    Rows: 26,
    Columns: 80,
    }
    terminalImage := TerminalImage{
    Id: imageId,
    Sequence: 1,
    Inbound: true,
    ImageSize: terminalSize,
    CursorRow: 0,
    CursorColumn: 0,
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}

func TestRenderTerminalWithFieldRendersOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    terminalSize := TerminalSize{
    Rows: 26,
    Columns: 80,
    }

    fieldContents := FieldContents{
        Text: "single text field in the middle",
    }
    terminalField := TerminalField{
        Row: 10,
        Column: 13,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ fieldContents },
    }

    terminalImage := TerminalImage{
    Id: imageId,
    Sequence: 1,
    Inbound: true,
    ImageSize: terminalSize,
    CursorRow: 0,
    CursorColumn: 0,
    Fields: []TerminalField{ terminalField },
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}

func TestRenderTerminalWithFieldAtOriginRendersOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    terminalSize := TerminalSize{
    Rows: 26,
    Columns: 80,
    }

    fieldContents := FieldContents{
        Text: "^ this is the origin (top left)",
    }
    terminalField := TerminalField{
        Row: 0,
        Column: 0,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ fieldContents },
    }

    terminalImage := TerminalImage{
    Id: imageId,
    Sequence: 1,
    Inbound: true,
    ImageSize: terminalSize,
    CursorRow: 0,
    CursorColumn: 0,
    Fields: []TerminalField{ terminalField },
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}

func TestRenderTerminalWithFieldAtTopRightRendersOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    terminalSize := TerminalSize{
    Rows: 26,
    Columns: 80,
    }

    fieldContents := FieldContents{
        Text: "The '^' should be at the top right",
    }
    textField := TerminalField{
        Row: 10,
        Column: 20,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ fieldContents },
    }

    topRightFieldContents := FieldContents{
        Text: "^",
    }
    topRightField := TerminalField{
        Row: 0,
        Column: 79,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ topRightFieldContents },
    }

    terminalImage := TerminalImage{
    Id: imageId,
    Sequence: 1,
    Inbound: true,
    ImageSize: terminalSize,
    CursorRow: 0,
    CursorColumn: 0,
    Fields: []TerminalField{ topRightField, textField },
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}

func TestRenderTerminalWithFieldAtBottomLeftRendersOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    terminalSize := TerminalSize{
    Rows: 26,
    Columns: 80,
    }

    fieldContents := FieldContents{
        Text: "The 'v' should be at the bottom left",
    }
    textField := TerminalField{
        Row: 10,
        Column: 20,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ fieldContents },
    }

    topRightFieldContents := FieldContents{
        Text: "v",
    }
    topRightField := TerminalField{
        Row: 26,
        Column: 0,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ topRightFieldContents },
    }

    terminalImage := TerminalImage{
    Id: imageId,
    Sequence: 1,
    Inbound: false,
    Aid: "my-aid",
    ImageSize: terminalSize,
    CursorRow: 0,
    CursorColumn: 0,
    Fields: []TerminalField{ topRightField, textField },
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}

func TestRenderTerminalWithFieldAtBottomRightRendersOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    terminalSize := TerminalSize{
    Rows: 26,
    Columns: 80,
    }

    fieldContents := FieldContents{
        Text: "The 'v' should be at the bottom right",
    }
    textField := TerminalField{
        Row: 10,
        Column: 20,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ fieldContents },
    }

    topRightFieldContents := FieldContents{
        Text: "v",
    }
    topRightField := TerminalField{
        Row: 26,
        Column: 79,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ topRightFieldContents },
    }

    terminalImage := TerminalImage{
    Id: imageId,
    Sequence: 1,
    Inbound: true,
    ImageSize: terminalSize,
    CursorRow: 0,
    CursorColumn: 0,
    Fields: []TerminalField{ topRightField, textField },
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}

func TestRenderTerminalWithFullRowRendersOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    terminalSize := TerminalSize{
    Rows: 26,
    Columns: 80,
    }

    guideContents := FieldContents{
        Text: "0          1          2          3          4          5          6          7",
    }
    textField := TerminalField{
        Row: 10,
        Column: 0,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ guideContents },
    }

    fullRowContents := FieldContents{
        Text: "01234567890123456789012345678901234567890123456789012345678901234567890123456789",
    }
    fullRowField := TerminalField{
        Row: 11,
        Column: 0,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ fullRowContents },
    }

    terminalImage := TerminalImage{
    Id: imageId,
    Sequence: 1,
    Inbound: true,
    ImageSize: terminalSize,
    CursorRow: 0,
    CursorColumn: 0,
    Fields: []TerminalField{ fullRowField, textField },
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}

func TestRenderTerminalWithWrappingRowRendersOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    terminalSize := TerminalSize{
    Rows: 26,
    Columns: 80,
    }

    guideContents := FieldContents{
        Text: "The next row should wrap around and continue on the row below it",
    }
    textField := TerminalField{
        Row: 10,
        Column: 0,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ guideContents },
    }

    fullRowContents := FieldContents{
        Text: "0123456789012345678901234567890123456789012345678901234567890123456789",
    }
    fullRowField := TerminalField{
        Row: 11,
        Column: 20,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ fullRowContents },
    }

    terminalImage := TerminalImage{
    Id: imageId,
    Sequence: 1,
    Inbound: true,
    ImageSize: terminalSize,
    CursorRow: 0,
    CursorColumn: 0,
    Fields: []TerminalField{ fullRowField, textField },
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}
