/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"errors"

	"github.com/unidoc/unidoc/common"
)

// locaTable represents the Index to Location (loca) table.
// https://docs.microsoft.com/en-us/typography/opentype/spec/loca
type locaTable struct {
	// The extra entry at the end helps calculating the length of the last glyph data element.
	offsetsShort []offset16 // short format. (numGlyphs+1 entries).
	offsetsLong  []offset32 // long format. (numGlyphs+1 entries).
}

// GetGlyphDataOffset returns offset for glyph index `gid`. The offset is relative to
// the beginning of the glyf table.
func (f *font) GetGlyphDataOffset(gid GlyphIndex) (offset int64, len int64, err error) {
	if f.loca == nil || f.head == nil {
		common.Log.Debug("loca or head missing")
		return 0, 0, errRequiredField
	}
	if gid < 0 || int(gid) >= int(f.maxp.numGlyphs) {
		common.Log.Debug("invalid range")
		return 0, 0, errRangeCheck
	}

	short := f.head.indexToLocFormat == 0
	if short {
		offset1 := 2 * int64(f.loca.offsetsShort[gid])
		offset2 := 2 * int64(f.loca.offsetsShort[gid+1])
		return offset1, offset2 - offset1, nil
	}

	offset1 := int64(f.loca.offsetsLong[gid])
	offset2 := int64(f.loca.offsetsLong[gid+1])
	return offset1, offset2 - offset1, nil
}

func (f *font) parseLoca(r *byteReader) (*locaTable, error) {
	if f.head == nil || f.maxp == nil {
		common.Log.Debug("head or maxp not set - required missing")
		return nil, errRequiredField
	}

	_, has, err := f.seekToTable(r, "loca")
	if err != nil {
		return nil, err
	}
	if !has {
		common.Log.Debug("loca table not present")
		return nil, nil
	}

	if f.head.indexToLocFormat < 0 || f.head.indexToLocFormat > 1 {
		common.Log.Debug("Invalid index to loca value")
		return nil, errRangeCheck
	}

	loca := &locaTable{}

	numGlyphs := int(f.maxp.numGlyphs)
	isShort := f.head.indexToLocFormat == 0

	if isShort {
		err := r.readSlice(&loca.offsetsShort, numGlyphs+1)
		if err != nil {
			return nil, err
		}
		return loca, nil
	}

	err = r.readSlice(&loca.offsetsLong, numGlyphs+1)
	if err != nil {
		return nil, err
	}
	for i := 0; i < numGlyphs; i++ {
		offset := loca.offsetsLong[i]
		len := loca.offsetsLong[i+1] - loca.offsetsLong[i]
		if offset < 0 {
			common.Log.Debug("Invalid offset")
			return nil, errors.New("invalid indexToLoca offset")
		}
		if len < 0 {
			common.Log.Debug("Invalid length")
			return nil, errors.New("invalid indexToLoca len")
		}

	}

	return loca, nil
}

func (f *font) writeLoca(w *byteWriter) error {
	if f.loca == nil || f.head == nil || f.maxp == nil {
		return errRequiredField
	}
	numGlyphs := int(f.maxp.numGlyphs)
	isShort := f.head.indexToLocFormat == 0

	t := f.loca
	if isShort {
		if numGlyphs+1 != len(t.offsetsShort) {
			common.Log.Debug("Unexpected length")
		}
		return w.writeSlice(t.offsetsShort)
	}
	return w.writeSlice(t.offsetsLong)
}
