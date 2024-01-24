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
	"log"
	"strconv"
	"strings"

	"github.com/galasa-dev/cli/pkg/embedded"
	galasaErrors "github.com/galasa-dev/cli/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
    PRIMARY_FONT_DIR  = "fonts/primary"
    FALLBACK_FONT_DIR = "fonts/fallback"
    FONT_WIDTH = 7
    FONT_HEIGHT = 13
)

var (
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
	drawer font.Drawer
}

func NewImageRenderer() ImageRenderer {
	renderer := new(ImageRendererImpl)

	fontFace := initRendererFonts()
	renderer.drawer = font.Drawer{
		Face: fontFace,
	}
	return renderer
}

func (renderer *ImageRendererImpl) RenderJsonBytesToImageFiles(jsonBinary []byte, writer ImageFileWriter) error {
	var err error
	var terminal Terminal

	terminal, err = convertJsonBytesToTerminal(jsonBinary)
	if err == nil {
		for _, terminalImage := range terminal.Images {
			image := renderer.renderTerminalImage(terminalImage)

			var pngImageBytes []byte
			pngImageBytes, err = encodeImageToPng(image)
			if err == nil {
				pngFileName := fmt.Sprintf("%s-%05d.png", terminal.Id, terminalImage.Sequence)
				err = writer.WriteImageFile(pngFileName, pngImageBytes)
			}
		}
	}
	return err
}

// Renders an RGBA image representation of a 3270 terminal and returns the rendered image
func (renderer *ImageRendererImpl) renderTerminalImage(terminalImage TerminalImage) *image.RGBA {
    targetColumnCount := terminalImage.ImageSize.Columns
    targetRowCount := terminalImage.ImageSize.Rows + 3

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
                textColor := getColor(field.ForegroundColor)
                renderer.drawString(img, column, row, string(char), textColor)
                column++
            }
        }
    }
    statusText := getStatusText(terminalImage, terminalImage.ImageSize.Columns, terminalImage.ImageSize.Rows)
    renderer.drawString(img, 0, targetRowCount - 2, statusText, DEFAULT_COLOR)
    return img
}

// Draws a string of text onto an image at the given column and row (x, y) coordinates
func (renderer *ImageRendererImpl) drawString(img *image.RGBA, column int, row int, text string, textColor color.RGBA) {
    startPoint := fixed.Point26_6{ X: fixed.I(column * FONT_WIDTH), Y: fixed.I((row + 1) * FONT_HEIGHT)}
	drawer := renderer.drawer

    drawer.Src = image.NewUniform(textColor)
	drawer.Dst = img
    drawer.Dot = startPoint

    drawer.DrawString(text)
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

// Loads all the fonts to be used in the renderer
func initRendererFonts() font.Face {
    var primaryFont font.Face

	fs := embedded.GetReadOnlyFileSystem()

    // Get the primary font to use in the renderer
    primaryFont = loadPrimaryFont(fs)

    fallbackFontFace := NewFallbackFontFace(primaryFont)
    loadFallbackFonts(fs, fallbackFontFace)

    return fallbackFontFace
}

// Loads the primary font to use in the renderer, defaulting to the built-in Face7x13 monospaced font
// if a primary font could not be loaded from the embedded filesystem
func loadPrimaryFont(fs embedded.ReadOnlyFileSystem) font.Face {
    var err error
    var loadedFonts []font.Face
    var primaryFont font.Face

    loadedFonts, err = loadFontsFromDirectory(fs, "fonts/primary")
	if err == nil && len(loadedFonts) > 0 {
        primaryFont = loadedFonts[0]
    } else {
		// Use a default monospaced font
		log.Println("Failed to load primary font, using a built-in font instead")
		primaryFont = basicfont.Face7x13
	}
    return primaryFont
}

// Loads any fallback fonts to use in the renderer when rendering glyphs that are not contained within the primary font
func loadFallbackFonts(fs embedded.ReadOnlyFileSystem, fallbackFontFace *FallbackFontFace) {
    var err error
    var loadedFonts []font.Face

    // Add any fallback fonts to use in the renderer
    loadedFonts, err = loadFontsFromDirectory(fs, "fonts/fallbacks")
	if err == nil {
        for _, font := range loadedFonts {
            fallbackFontFace.AddFallbackFont(font)
        }
	} else {
        // We can continue with the rendering, but certain characters may not appear correctly
        // if they are not included in the primary font
		log.Println("Failed to load substitute fonts, continuing without additional fonts")
	}
}