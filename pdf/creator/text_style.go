/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package creator

import (
	"github.com/unidoc/unidoc/pdf/model/fonts"
)

// TextStyle is a collection of properties that can be assigned to a chunk of text.
type TextStyle struct {
	// The color of the text.
	Color Color

	// The font the text will use.
	Font fonts.Font

	// The size of the font.
	FontSize float64
}

// NewTextStyle creates a new text style object which can be used with chunks
// of text. Uses default parameters: Helvetica, WinAnsiEncoding and wrap
// enabled with a wrap width of 100 points.
func NewTextStyle() TextStyle {
	return TextStyle{
		Color:    ColorRGBFrom8bit(0, 0, 0),
		Font:     defaultFont,
		FontSize: 10,
	}
}
