/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/galasa-dev/cli/pkg/files"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
    FONT_WIDTH = 7
    FONT_HEIGHT = 13
)

// Renders an RGBA image representation of a 3270 terminal and returns the rendered image
func RenderTerminalImage(terminalImage TerminalImage) *image.RGBA {
    targetColumnCount := terminalImage.ImageSize.Columns
    targetRowCount := terminalImage.ImageSize.Rows + 2

    imagePixelWidth := targetColumnCount * FONT_WIDTH
    imagePixelHeight := targetRowCount * FONT_HEIGHT
    img := createImageBase(imagePixelWidth, imagePixelHeight)

    for _, field := range terminalImage.Fields {
        column := field.Column
        row := field.Row

        for _, contents := range field.Contents {
            // Field contents can sometimes span multiple rows, so draw each character individually,
            // adjusting the current row whenever the image column boundary is reached
            for _, char := range contents.getCharacters() {
                if column >= targetColumnCount {
                    column = 0
                    row++
                }
                drawString(img, column, row, string(char))
                column++
            }
        }
    }
    statusText := getStatusText(terminalImage, terminalImage.ImageSize.Columns, terminalImage.ImageSize.Rows)
    drawString(img, 0, targetRowCount - 1, statusText)
    return img
}

// Creates a black rectangle image with the given dimensions in pixels
func createImageBase(imagePixelWidth int, imagePixelHeight int) *image.RGBA {
    rect := image.Rect(0, 0, imagePixelWidth, imagePixelHeight)
    img := image.NewRGBA(rect)

    // Draw a black background onto the image
    draw.Draw(img, img.Bounds(), image.Black, image.Pt(0, 0), draw.Src)
    return img
}

// Returns a string containing the content of the status row to be displayed at the bottom
// of a 3270 image
func getStatusText(terminalImage TerminalImage, columns int, rows int) string {
    var buff strings.Builder

    buff.WriteString(terminalImage.Id)
    buff.WriteString(" - ")

    buff.WriteString(strconv.Itoa(columns))
    buff.WriteString("x")
    buff.WriteString(strconv.Itoa(rows))
    buff.WriteString(" - ")

    if terminalImage.Inbound {
        buff.WriteString("Inbound ")
    } else {
        buff.WriteString("Outbound - ")
        buff.WriteString(terminalImage.Aid)
    }
    return buff.String()
}

// Draws a string of text onto an image at the given column and row (x, y) coordinates
func drawString(img *image.RGBA, column int, row int, text string) {
    textColor := color.RGBA{0, 255, 0, 255}
    point := fixed.Point26_6{ X: fixed.I(column * FONT_WIDTH), Y: fixed.I((row + 1) * FONT_HEIGHT)}

    drawer := &font.Drawer{
        Dst:  img,
        Src:  image.NewUniform(textColor),
        Face: basicfont.Face7x13,
        Dot:  point,
    }
    drawer.DrawString(text)
}

// Writes a rendered 3270 terminal .png image to the filesystem
func WritePngImageToDisk(terminalImage TerminalImage, image *image.RGBA, fileSystem files.FileSystem, outputDirectory string) error {
    var err error = nil
    var fileWriter io.Writer

    pngFileName := fmt.Sprintf("%s.png", terminalImage.Id)

    // Create the .png file on the filesystem
    fileWriter, err = fileSystem.Create(filepath.Join(outputDirectory, pngFileName))
    if err == nil {
        err = png.Encode(fileWriter, image)
    }
    return err
}

// Converts a given JSON representation of a 3270 terminal into a Terminal struct
func convertJsonToTerminal(terminalJson string) (Terminal, error) {
	var terminal Terminal
	var err error = nil
	err = json.Unmarshal([]byte(terminalJson), &terminal)

	return terminal, err
}
