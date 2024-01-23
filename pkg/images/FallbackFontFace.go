/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package images

import (
	"image"
	"log"

	"github.com/galasa-dev/cli/pkg/files"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Implementation of the font.Face interface
type FallbackFontFace struct {
	primary   FontEntry
	fallbacks []FontEntry
}

type FontEntry struct {
	fontFace font.Face
	trueTypeFont *truetype.Font
}

// TODO: move these methods somewhere else?
func LoadTrueTypeFont(fontFilePath string) (font.Face, *truetype.Font, error) {
    var err error
    var fontBytes []byte
    var truetypeFont *truetype.Font
    var fontFace font.Face

    fs := files.NewOSFileSystem()
    fontBytes, err = fs.ReadBinaryFile(fontFilePath)
    if err == nil {
        truetypeFont, err = truetype.Parse(fontBytes)
        if err == nil {
            fontFace = truetype.NewFace(truetypeFont, &truetype.Options{})
        }
    }
    return fontFace, truetypeFont, err
}

func LoadOpenTypeFont(fontFilePath string) (font.Face, error) {
    var err error
    var fontBytes []byte
    var opentypeFont *opentype.Font
    var fontFace font.Face

    fs := files.NewOSFileSystem()
    fontBytes, err = fs.ReadBinaryFile(fontFilePath)
    if err == nil {
        opentypeFont, err = opentype.Parse(fontBytes)
        if err == nil {
            fontFace, err = opentype.NewFace(opentypeFont, nil)
        }
    }
    return fontFace, err
}

func initFontFaces() font.Face {
	var err error
	var fontFace font.Face
	var font *truetype.Font
	var primaryFont FontEntry

    fontFace, font, err = LoadTrueTypeFont("NotoSansMono.ttf")
	if err == nil {
		primaryFont = FontEntry{fontFace, font}
	} else {
		// Use a default mono font
		log.Println("Could not load primary font, using built-in font")
		primaryFont = FontEntry{basicfont.Face7x13, nil}
	}

	fallbackFonts := []FontEntry{}
    fontFace, err = LoadOpenTypeFont("NotoSansMonoCJKjp.otf")
	if err == nil {
		fallbackFonts = append(fallbackFonts, FontEntry{fontFace, nil})
	} else {
		log.Println("Failed to load substitute fonts, continuing without additional fonts")
	}

    return &FallbackFontFace{primary: primaryFont, fallbacks: fallbackFonts}
}

// Checks whether or not a given glyph is in a truetype font
func isInTrueTypeFont(font *truetype.Font, glyph rune) bool {
	return (font != nil) && (font.Index(glyph) != 0)
}

func (f *FallbackFontFace) Close() error {
	err := f.primary.fontFace.Close()
	for _, fallback := range f.fallbacks {
		fallbackErr := fallback.fontFace.Close()
		if fallbackErr != nil && err == nil {
			err = fallbackErr
			break
		}
	}
	return err
}

func (f *FallbackFontFace) Metrics() font.Metrics {
	return f.primary.fontFace.Metrics()
}

func (f *FallbackFontFace) Kern(r0 rune, r1 rune) fixed.Int26_6 {
	kern := f.primary.fontFace.Kern(r0, r1)
	if kern == 0 {
		for _, fallback := range f.fallbacks {
			fallbackKern := fallback.fontFace.Kern(r0, r1)
			if fallbackKern != 0 {
				kern = fallbackKern
				break
			}
		}
	}
	return kern
}

func (f *FallbackFontFace) Glyph(dot fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
	rect, img, point, advance, ok := f.primary.fontFace.Glyph(dot, r)
	if !isInTrueTypeFont(f.primary.trueTypeFont, r) || !ok {
		for _, fallback := range f.fallbacks {
			rect, img, point, advance, ok = fallback.fontFace.Glyph(dot, r)
			if isInTrueTypeFont(fallback.trueTypeFont, r) || ok {
				break
			}
		}
	}
	return rect, img, point, advance, ok
}

func (f *FallbackFontFace) GlyphAdvance(r rune) (fixed.Int26_6, bool) {
	advance, ok := f.primary.fontFace.GlyphAdvance(r)
	if !isInTrueTypeFont(f.primary.trueTypeFont, r) || !ok {
		for _, fallback := range f.fallbacks {
			fallbackAdv, fallbackOk := fallback.fontFace.GlyphAdvance(r)
			if isInTrueTypeFont(fallback.trueTypeFont, r) || fallbackOk {
				advance, ok = fallbackAdv, fallbackOk
				break
			}
		}
	}
	return advance, ok
}

func (f *FallbackFontFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	bounds, adv, ok := f.primary.fontFace.GlyphBounds(r)
	if !isInTrueTypeFont(f.primary.trueTypeFont, r) || !ok {
		for _, fallback := range f.fallbacks {
			fallbackBounds, fallbackAdvance, fallbackOk := fallback.fontFace.GlyphBounds(r)
			if isInTrueTypeFont(fallback.trueTypeFont, r) || ok {
				bounds, adv, ok = fallbackBounds, fallbackAdvance, fallbackOk
				break
			}
		}
	}
	return bounds, adv, ok
}