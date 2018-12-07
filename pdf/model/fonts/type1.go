/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package fonts

import (
	"github.com/unidoc/unidoc/pdf/core"
	"github.com/unidoc/unidoc/pdf/internal/textencoding"
)

var _ Font = Type1Font{}

// Type1Font represents one of the built-in fonts and it is assumed that every reader has access to it.
type Type1Font struct {
	name    string
	metrics map[GlyphName]CharMetrics
	encoder textencoding.TextEncoder
}

// NewType1Font returns a new instance of the font with a default encoder set (WinAnsiEncoding).
func NewType1Font(name string, metrics map[GlyphName]CharMetrics) Type1Font {
	enc := textencoding.NewWinAnsiTextEncoder() // Default
	return NewType1FontWithEncoding(name, metrics, enc)
}

// NewType1FontWithEncoding returns a new instance of the font with a specified encoder.
func NewType1FontWithEncoding(name string, metrics map[GlyphName]CharMetrics, encoder textencoding.TextEncoder) Type1Font {
	return Type1Font{
		name:    name,
		metrics: metrics,
		encoder: encoder,
	}
}

// Name returns a PDF name of the font.
func (font Type1Font) Name() string {
	return font.name
}

// Encoder returns the font's text encoder.
func (font Type1Font) Encoder() textencoding.TextEncoder {
	return font.encoder
}

// GetGlyphCharMetrics returns character metrics for a given glyph.
func (font Type1Font) GetGlyphCharMetrics(glyph GlyphName) (CharMetrics, bool) {
	metrics, has := font.metrics[glyph]
	if !has {
		return metrics, false
	}

	return metrics, true
}

// ToPdfObject returns a primitive PDF object representation of the font.
func (font Type1Font) ToPdfObject() core.PdfObject {
	fontDict := core.MakeDict()
	fontDict.Set("Type", core.MakeName("Font"))
	fontDict.Set("Subtype", core.MakeName("Type1"))
	fontDict.Set("BaseFont", core.MakeName(font.name))
	fontDict.Set("Encoding", font.encoder.ToPdfObject())

	return core.MakeIndirectObject(fontDict)
}

// type1CommonGlyphs is list of common glyph names for some Type1. Used to unpack character metrics.
var type1CommonGlyphs = []textencoding.GlyphName{
	"A", "AE", "Aacute", "Abreve", "Acircumflex",
	"Adieresis", "Agrave", "Amacron", "Aogonek", "Aring",
	"Atilde", "B", "C", "Cacute", "Ccaron",
	"Ccedilla", "D", "Dcaron", "Dcroat", "Delta",
	"E", "Eacute", "Ecaron", "Ecircumflex", "Edieresis",
	"Edotaccent", "Egrave", "Emacron", "Eogonek", "Eth",
	"Euro", "F", "G", "Gbreve", "Gcommaaccent",
	"H", "I", "Iacute", "Icircumflex", "Idieresis",
	"Idotaccent", "Igrave", "Imacron", "Iogonek", "J",
	"K", "Kcommaaccent", "L", "Lacute", "Lcaron",
	"Lcommaaccent", "Lslash", "M", "N", "Nacute",
	"Ncaron", "Ncommaaccent", "Ntilde", "O", "OE",
	"Oacute", "Ocircumflex", "Odieresis", "Ograve", "Ohungarumlaut",
	"Omacron", "Oslash", "Otilde", "P", "Q",
	"R", "Racute", "Rcaron", "Rcommaaccent", "S",
	"Sacute", "Scaron", "Scedilla", "Scommaaccent", "T",
	"Tcaron", "Tcommaaccent", "Thorn", "U", "Uacute",
	"Ucircumflex", "Udieresis", "Ugrave", "Uhungarumlaut", "Umacron",
	"Uogonek", "Uring", "V", "W", "X",
	"Y", "Yacute", "Ydieresis", "Z", "Zacute",
	"Zcaron", "Zdotaccent", "a", "aacute", "abreve",
	"acircumflex", "acute", "adieresis", "ae", "agrave",
	"amacron", "ampersand", "aogonek", "aring", "asciicircum",
	"asciitilde", "asterisk", "at", "atilde", "b",
	"backslash", "bar", "braceleft", "braceright", "bracketleft",
	"bracketright", "breve", "brokenbar", "bullet", "c",
	"cacute", "caron", "ccaron", "ccedilla", "cedilla",
	"cent", "circumflex", "colon", "comma", "commaaccent",
	"copyright", "currency", "d", "dagger", "daggerdbl",
	"dcaron", "dcroat", "degree", "dieresis", "divide",
	"dollar", "dotaccent", "dotlessi", "e", "eacute",
	"ecaron", "ecircumflex", "edieresis", "edotaccent", "egrave",
	"eight", "ellipsis", "emacron", "emdash", "endash",
	"eogonek", "equal", "eth", "exclam", "exclamdown",
	"f", "fi", "five", "fl", "florin",
	"four", "fraction", "g", "gbreve", "gcommaaccent",
	"germandbls", "grave", "greater", "greaterequal", "guillemotleft",
	"guillemotright", "guilsinglleft", "guilsinglright", "h", "hungarumlaut",
	"hyphen", "i", "iacute", "icircumflex", "idieresis",
	"igrave", "imacron", "iogonek", "j", "k",
	"kcommaaccent", "l", "lacute", "lcaron", "lcommaaccent",
	"less", "lessequal", "logicalnot", "lozenge", "lslash",
	"m", "macron", "minus", "mu", "multiply",
	"n", "nacute", "ncaron", "ncommaaccent", "nine",
	"notequal", "ntilde", "numbersign", "o", "oacute",
	"ocircumflex", "odieresis", "oe", "ogonek", "ograve",
	"ohungarumlaut", "omacron", "one", "onehalf", "onequarter",
	"onesuperior", "ordfeminine", "ordmasculine", "oslash", "otilde",
	"p", "paragraph", "parenleft", "parenright", "partialdiff",
	"percent", "period", "periodcentered", "perthousand", "plus",
	"plusminus", "q", "question", "questiondown", "quotedbl",
	"quotedblbase", "quotedblleft", "quotedblright", "quoteleft", "quoteright",
	"quotesinglbase", "quotesingle", "r", "racute", "radical",
	"rcaron", "rcommaaccent", "registered", "ring", "s",
	"sacute", "scaron", "scedilla", "scommaaccent", "section",
	"semicolon", "seven", "six", "slash", "space",
	"sterling", "summation", "t", "tcaron", "tcommaaccent",
	"thorn", "three", "threequarters", "threesuperior", "tilde",
	"trademark", "two", "twosuperior", "u", "uacute",
	"ucircumflex", "udieresis", "ugrave", "uhungarumlaut", "umacron",
	"underscore", "uogonek", "uring", "v", "w",
	"x", "y", "yacute", "ydieresis", "yen",
	"z", "zacute", "zcaron", "zdotaccent", "zero",
}
