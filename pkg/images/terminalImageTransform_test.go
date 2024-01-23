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
	"strconv"
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

func createTextField(row int, column int, text string, textColor string) TerminalField {
    fieldContents := FieldContents{ Text: text }

    return TerminalField{
        Row: row,
        Column: column,
        Unformatted: false,
        FieldProtected: false,
        FieldNumeric: false,
        FieldDisplay: true,
        FieldIntenseDisplay: false,
        FieldSelectorPen: false,
        FieldModified: false,
        Contents: []FieldContents{ fieldContents },
        ForegroundColor: textColor,
    }
}

//----------------------------------------------------
// Tests
//----------------------------------------------------
func TestWritePngImageToDiskShouldCreateAPngFile(t *testing.T) {
    // Given...
    fs := files.NewMockFileSystem()
    tempDir, _ := fs.MkTempDir()

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
    err := WritePngImageToDisk(terminalImage, image, fs, tempDir)
    assert.Nil(t, err, "Should have successfully created a .png file")

    // Then...
    expectedPngFilePath := filepath.Join(tempDir, imageId + ".png")
    pngExists, _ := fs.Exists(expectedPngFilePath)
    assert.True(t, pngExists, "PNG file should have been created at '" + expectedPngFilePath + "'")

}

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

    terminalImage := TerminalImage{
        Id: imageId,
        Sequence: 1,
        Inbound: true,
        ImageSize: terminalSize,
        CursorRow: 0,
        CursorColumn: 0,
        Fields: []TerminalField{ createTextField(10, 13, "single text field in the middle", "d") },
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}

func TestRenderTerminalWithSmallerSizeRendersOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    terminalSize := TerminalSize{
        Rows: 18,
        Columns: 66,
    }

    terminalImage := TerminalImage{
        Id: imageId,
        Sequence: 1,
        Inbound: true,
        ImageSize: terminalSize,
        CursorRow: 0,
        CursorColumn: 0,
        Fields: []TerminalField{ createTextField(9, 15, "this terminal should be 66x18", "d") },
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

    terminalImage := TerminalImage{
        Id: imageId,
        Sequence: 1,
        Inbound: true,
        ImageSize: terminalSize,
        CursorRow: 0,
        CursorColumn: 0,
        Fields: []TerminalField{ createTextField(0, 0, "^ this is the origin (top left)", "d") },
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

    textField := createTextField(10, 20, "The '^' should be at the top right", "d")
    topRightField := createTextField(0, 79, "^", "d")

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

    textField := createTextField(10, 20, "The 'v' should be at the bottom left", "d")
    bottomLeftField := createTextField(26, 0, "v", "d")

    terminalImage := TerminalImage{
        Id: imageId,
        Sequence: 1,
        Inbound: false,
        Aid: "my-aid",
        ImageSize: terminalSize,
        CursorRow: 0,
        CursorColumn: 0,
        Fields: []TerminalField{ bottomLeftField, textField },
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

    textField := createTextField(10, 20, "The 'v' should be at the bottom right", "d")
    bottomRightField := createTextField(26, 79, "v", "d")

    terminalImage := TerminalImage{
        Id: imageId,
        Sequence: 1,
        Inbound: true,
        ImageSize: terminalSize,
        CursorRow: 0,
        CursorColumn: 0,
        Fields: []TerminalField{ bottomRightField, textField },
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

    textField := createTextField(10, 0, "0          1          2          3          4          5          6          7", "d")
    fullRowField := createTextField(11, 0, "01234567890123456789012345678901234567890123456789012345678901234567890123456789", "d")

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

func TestRenderTerminalWithFullColumnRendersOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    rows := 26
    terminalSize := TerminalSize{
        Rows: rows,
        Columns: 80,
    }

    terminalFields := make([]TerminalField, 0)
    for i := 0; i < rows; i++ {
        terminalFields = append(terminalFields, createTextField(i, 0, strconv.Itoa(i), "d"))
    }
    terminalFields = append(terminalFields, createTextField(10, 20, "Each of the 26 rows should have a number in", "d"))

    terminalImage := TerminalImage{
        Id: imageId,
        Sequence: 1,
        Inbound: true,
        ImageSize: terminalSize,
        CursorRow: 0,
        CursorColumn: 0,
        Fields: terminalFields,
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

    textField := createTextField(10, 0, "The next row should wrap around and continue on the row below it", "d")
    wrappedField := createTextField(11, 20, "0123456789012345678901234567890123456789012345678901234567890123456789", "d")

    terminalImage := TerminalImage{
        Id: imageId,
        Sequence: 1,
        Inbound: true,
        ImageSize: terminalSize,
        CursorRow: 0,
        CursorColumn: 0,
        Fields: []TerminalField{ wrappedField, textField },
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}

func TestRenderTerminaColorsRenderOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    terminalSize := TerminalSize{
        Rows: 26,
        Columns: 80,
    }

    defaultField := createTextField(10, 20, "This is the default color", "d")
    neutralField := createTextField(11, 20, "This is the neutral color", "n")
    redField := createTextField(12, 20, "This is red", "r")
    greenField := createTextField(13, 20, "This is green", "g")
    blueField := createTextField(14, 20, "This is blue", "b")
    pinkField := createTextField(15, 20, "This is pink", "p")
    turquoiseField := createTextField(16, 20, "This is turquoise", "t")
    yellowField := createTextField(17, 20, "This is yellow", "y")
    unknownColorField := createTextField(18, 20, "This is unknown, should render using the default color", "blah")

    terminalImage := TerminalImage{
        Id: imageId,
        Sequence: 1,
        Inbound: true,
        ImageSize: terminalSize,
        CursorRow: 0,
        CursorColumn: 0,
        Fields: []TerminalField{
            defaultField,
            neutralField,
            redField,
            greenField,
            blueField,
            pinkField,
            turquoiseField,
            yellowField,
            unknownColorField,
        },
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}

func TestRenderTerminaUnicodeTextRendersOk(t *testing.T) {
    // Given...
    fs := files.NewOSFileSystem()

    imageId := t.Name()
    terminalSize := TerminalSize{
        Rows: 26,
        Columns: 80,
    }

    symbolField := createTextField(9, 20, "Symbols: © Ø ® ß ◊ ¥ Ô º ™ € ¢ ∞ § Ω`", "d")
    greekField := createTextField(10, 20, "Greek: Χαίρετε", "d")
    japaneseField := createTextField(11, 20, "Japanese: こんにちは", "d")
    russianField := createTextField(12, 20, "Russian: Здравствуйте", "d")
    chineseField := createTextField(13, 20, "Chinese: 你好", "d")
    koreanField := createTextField(14, 20, "Korean: 여보세요", "d")

    terminalImage := TerminalImage{
        Id: imageId,
        Sequence: 1,
        Inbound: true,
        ImageSize: terminalSize,
        CursorRow: 0,
        CursorColumn: 0,
        Fields: []TerminalField{
            symbolField,
            greekField,
            japaneseField,
            russianField,
            chineseField,
            koreanField,
        },
    }

    // When...
    image := RenderTerminalImage(terminalImage)

    // Then...
    assertTerminalImageMatchesExpectedSnapshot(t, fs, image, terminalImage)
}
