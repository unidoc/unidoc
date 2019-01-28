/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"encoding/binary"
	"strings"
)

// GlyphName is a representation of a glyph name, e.g. from Adobe's glyph list.
type GlyphName string

// GlyphIndex or Glyph ID (GID) represent each glyph within a font.
type GlyphIndex uint16

/*
Types in truetype fonts:
https://docs.microsoft.com/en-us/typography/opentype/spec/otff

Data Type	Description
--------------------------------------------------------
uint8	  8-bit unsigned integer.
int8	  8-bit signed integer.
uint16	  16-bit unsigned integer.
int16	  16-bit signed integer.
uint24	  24-bit unsigned integer.
uint32	  32-bit unsigned integer.
int32	  32-bit signed integer.
Fixed	  32-bit signed fixed-point number (16.16)
FWORD	  int16 that describes a quantity in font design units.
UFWORD	  uint16 that describes a quantity in font design units.
F2DOT14	  16-bit signed fixed number with the low 14 bits of fraction (2.14).
LONGDATETIME
          Date represented in number of seconds since 12:00 midnight, January 1, 1904.
          The value is represented as a signed 64-bit integer.
Tag	      Array of four uint8s (length = 32 bits) used to identify a table,
          design-variation axis, script, language system, feature, or baseline
Offset16  Short offset to a table, same as uint16, NULL offset = 0x0000
Offset32  Long offset to a table, same as uint32, NULL offset = 0x00000000
*/

type fixed int32
type fword int16
type ufword uint16
type f2dot14 int16
type longdatetime int64
type tag [4]uint8
type offset16 uint16
type offset32 uint32

func (t tag) String() string {
	return strings.TrimSpace(string(t[:]))
}

// Parts returns the integral and decimal portions of `f`.
func (f fixed) Parts() (uint16, uint16) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(f))
	return binary.BigEndian.Uint16(b[0:2]), binary.BigEndian.Uint16(b[2:4])
}

// Float64 returns `f` as a float64.
func (f fixed) Float64() float64 {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(f))
	l, r := binary.BigEndian.Uint16(b[0:2]), binary.BigEndian.Uint16(b[2:4])
	integral := float64(int16(l))
	fraction := float64(r) / 65536.0
	return integral + fraction
}

func makeTag(s string) tag {
	bb := []byte(s[:])
	if len(bb) > 4 {
		// Trim to 4 bytes.
		bb = bb[:4]
	}
	if len(bb) < 4 {
		// Pad with spaces to fill 4 bytes.
		for i := 0; i < 4-len(bb); i++ {
			bb = append(bb, ' ')
		}
	}

	var t tag
	copy(t[:], bb)
	return t
}
