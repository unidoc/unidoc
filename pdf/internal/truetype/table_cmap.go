/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

// cmapTable represents a Character to Glyph Index Mapping Table (cmap).
// This table defines the mapping of character codes to the glyph index values used
// in the font.
// https://docs.microsoft.com/en-us/typography/opentype/spec/cmap
type cmapTable struct {
	version         uint16
	numTables       uint16
	encodingRecords []encodingRecord // len == numTables
}

type encodingRecord struct {
	platformID uint16
	encodingID uint16
	offset     offset32
}

/*
Regardless of the encoding scheme, character codes that do not correspond to any glyph in the font should be
mapped to glyph index 0. The glyph at this location must be a special glyph representing a missing character,
commonly known as .notdef.
*/

/*
There are 7 subtable formats.
*/
