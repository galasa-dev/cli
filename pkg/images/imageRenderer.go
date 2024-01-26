/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"
	"strings"

	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	FONT_WIDTH  = 7
	FONT_HEIGHT = 13

	DEFAULT_COLOR = color.RGBA{0, 255, 0, 255}
	NEUTRAL       = color.RGBA{255, 255, 255, 255}
	RED           = color.RGBA{255, 0, 0, 255}
	GREEN         = color.RGBA{0, 255, 0, 255}
	BLUE          = color.RGBA{0, 0, 255, 255}
	PINK          = color.RGBA{255, 0, 204, 255}
	TURQUOISE     = color.RGBA{64, 224, 208, 255}
	YELLOW        = color.RGBA{255, 255, 0, 255}

	colors = map[string]color.RGBA{
		"d": DEFAULT_COLOR,
		"r": RED,
		"g": GREEN,
		"b": BLUE,
		"p": PINK,
		"t": TURQUOISE,
		"y": YELLOW,
		"n": NEUTRAL,
	}
)

type ImageRenderer interface {
	RenderJsonBytesToImageFiles(jsonBinary []byte, writer ImageFileWriter) error
}

type ImageRendererImpl struct {
}

func NewImageRenderer() ImageRenderer {
	renderer := new(ImageRendererImpl)
	return renderer
}

func (renderer *ImageRendererImpl) RenderJsonBytesToImageFiles(jsonBinary []byte, writer ImageFileWriter) error {
	var err error
	var terminal Terminal

	terminal, err = convertJsonBytesToTerminal(jsonBinary)
	if err == nil {
		for _, terminalImage := range terminal.Images {

			pngFileName := fmt.Sprintf("%s-%05d.png", terminal.Id, terminalImage.Sequence)

			var isWritable bool
			isWritable, err = writer.IsImageFileWritable(pngFileName)
			if err == nil {
				// Only render the image if we will be able to write it out later.
				// ie: Don't do all the work if the file already exists.
				// When the same root folder is scanned for a second time, we want to minimise the work done
				if isWritable {
					image := renderTerminalImage(terminalImage)

					var pngImageBytes []byte
					pngImageBytes, err = encodeImageToPng(image)
					if err == nil {
						err = writer.WriteImageFile(pngFileName, pngImageBytes)
					}
				}
			}
		}
	}
	return err
}

// Converts a JSON byte array representing one or more 3270 terminals into a Terminal object
func convertJsonBytesToTerminal(terminalJsonBytes []byte) (Terminal, error) {
	var terminal Terminal
	var err error = nil

	err = json.Unmarshal(terminalJsonBytes, &terminal)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_BAD_TERMINAL_JSON_FORMAT, err.Error())
	}
	return terminal, err
}

// Renders an RGBA image representation of a 3270 terminal and returns the rendered image
func renderTerminalImage(terminalImage TerminalImage) *image.RGBA {
	targetColumnCount := terminalImage.ImageSize.Columns
	targetRowCount := terminalImage.ImageSize.Rows + 3

	imagePixelWidth := targetColumnCount * FONT_WIDTH
	imagePixelHeight := targetRowCount * FONT_HEIGHT
	img := createImageBase(imagePixelWidth, imagePixelHeight)
	context := createImageDrawer(img)

	for _, field := range terminalImage.Fields {
		column := field.Column
		row := field.Row

		for _, contents := range field.Contents {
			// Field contents can sometimes span multiple rows, so draw each character individually,
			// adjusting the current row whenever the image column boundary is reached
			for _, char := range getCharacters(&contents) {
				if column >= targetColumnCount {
					column = 0
					row++
				}
				textColor := getColor(field.ForegroundColor)
				drawString(context, column, row, string(char), textColor)
				column++
			}
		}
	}
	statusText := getStatusText(terminalImage, terminalImage.ImageSize.Columns, terminalImage.ImageSize.Rows)
	drawString(context, 0, targetRowCount-2, statusText, DEFAULT_COLOR)
	return img
}

func getCharacters(fieldContents *FieldContents) []rune {
	var contents []rune
	if fieldContents.Characters != nil {
		// If the terminal JSON defines the contents of a field in characters,
		// then convert the character strings into runes
		for _, char := range fieldContents.Characters {
			contents = append(contents, []rune(char)...)
		}
	} else {
		contents = []rune(fieldContents.Text)
	}
	return contents
}

// Encodes a rendered 3270 terminal image in PNG format and returns the resulting byte array
func encodeImageToPng(image *image.RGBA) ([]byte, error) {
	var err error = nil

	buf := new(bytes.Buffer)
	err = png.Encode(buf, image)
	if err != nil {
		err = galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PNG_ENCODING_FAILED, err.Error())
	}

	return buf.Bytes(), err
}

// Creates a black rectangle image with the given dimensions in pixels
func createImageBase(imagePixelWidth int, imagePixelHeight int) *image.RGBA {
	rect := image.Rect(0, 0, imagePixelWidth, imagePixelHeight)
	img := image.NewRGBA(rect)

	// Draw a black background onto the image
	draw.Draw(img, img.Bounds(), image.Black, image.Pt(0, 0), draw.Src)
	return img
}

// Creates a drawer instance to draw text on a given image
func createImageDrawer(img *image.RGBA) *font.Drawer {
	drawer := &font.Drawer{
		Dst:  img,
		Face: initFontFaces(),
	}

	return drawer
}

// Draws a string of text onto an image at the given column and row (x, y) coordinates
func drawString(drawer *font.Drawer, column int, row int, text string, textColor color.RGBA) {
	startPoint := fixed.Point26_6{X: fixed.I(column * FONT_WIDTH), Y: fixed.I((row + 1) * FONT_HEIGHT)}

	drawer.Src = image.NewUniform(textColor)
	drawer.Dot = startPoint

	drawer.DrawString(text)
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

// Returns a color from the colors map that matches the given single-character identifier
func getColor(colorIdentifier string) color.RGBA {
	color := DEFAULT_COLOR
	if matchedColor, ok := colors[colorIdentifier]; ok {
		color = matchedColor
	}
	return color
}
