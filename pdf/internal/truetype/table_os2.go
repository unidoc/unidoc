/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import "github.com/unidoc/unidoc/common"

// os2Table represents the OS/2 metrics table. It consists of metrics and other data that are required.
type os2Table struct {
	// Version 0+
	version             uint16
	xAvgCharWidth       int16
	usWeightClass       uint16
	usWidthClass        uint16
	fsType              uint16
	ySubscriptXSize     int16
	ySubscriptYSize     int16
	ySubscriptXOffset   int16
	ySubscriptYOffset   int16
	ySuperscriptXSize   int16
	ySuperscriptYSize   int16
	ySuperscriptXOffset int16
	ySuperscriptYOffset int16
	yStrikeoutSize      int16
	yStrikeoutPosition  int16
	sFamilyClass        int16
	panose10            []uint8 // panose10 len = 10
	ulUnicodeRange1     uint32  // Bits 0-31.
	ulUnicodeRange2     uint32  // Bits 32-63.
	ulUnicodeRange3     uint32  // Bits 64-95.
	ulUnicodeRange4     uint32  // Bits 96-127.
	achVendId           tag
	fsSelection         uint16
	usFirstCharIndex    uint16
	usLastCharIndex     uint16
	sTypoAscender       int16
	sTypoDescender      int16
	sTypoLineGap        int16
	usWinAscent         uint16
	usWinDescent        uint16

	// Version 1-5.
	ulCodePageRange1 uint32 // Bits 0-31
	ulCodePageRange2 uint32 // Bits 32-63.

	// Version 2-5
	sxHeight      int16
	sCapHeight    int16
	usDefaultChar uint16
	usBreakChar   uint16
	usMaxContext  uint16

	// Version 5
	usLowerOpticalPointSize uint16
	usUpperOpticalPointSize uint16
}

func (f *font) parseOS2Table(r *byteReader) (*os2Table, error) {
	_, has, err := f.seekToTable(r, "OS/2")
	if err != nil {
		return nil, err
	}
	if !has {
		common.Log.Debug("OS/2 table not present")
		return nil, nil
	}

	t := &os2Table{}
	err = r.read(&t.version, &t.xAvgCharWidth, &t.usWeightClass, &t.usWidthClass, &t.fsType)
	if err != nil {
		return nil, err
	}

	if t.version > 10 {
		common.Log.Debug("OS/2 table version range error")
		return nil, errRangeCheck
	}

	err = r.read(&t.ySubscriptXSize, &t.ySubscriptYSize, &t.ySubscriptXOffset, &t.ySubscriptYOffset)
	if err != nil {
		return nil, err
	}

	err = r.read(&t.ySuperscriptXSize, &t.ySuperscriptYSize, &t.ySuperscriptXOffset, &t.ySuperscriptYOffset)
	if err != nil {
		return nil, err
	}

	err = r.read(&t.yStrikeoutSize, &t.yStrikeoutPosition, &t.sFamilyClass)
	if err != nil {
		return nil, err
	}

	err = r.readSlice(&t.panose10, 10)
	if err != nil {
		return nil, err
	}

	err = r.read(&t.ulUnicodeRange1, &t.ulUnicodeRange2, &t.ulUnicodeRange3, &t.ulUnicodeRange4)
	if err != nil {
		return nil, err
	}
	err = r.read(&t.achVendId, &t.fsSelection, &t.usFirstCharIndex, &t.usLastCharIndex, &t.sTypoAscender)
	if err != nil {
		return nil, err
	}
	err = r.read(&t.sTypoDescender, &t.sTypoLineGap, &t.usWinAscent, &t.usWinDescent)
	if err != nil {
		return nil, err
	}

	if t.version == 0 {
		return t, nil
	}

	// version >= 1.
	err = r.read(&t.ulCodePageRange1, &t.ulCodePageRange2)
	if err != nil {
		return nil, err
	}
	if t.version == 1 {
		return t, nil
	}

	// version 2-5.
	err = r.read(&t.sxHeight, &t.sCapHeight, &t.usDefaultChar, &t.usBreakChar, &t.usMaxContext)
	if err != nil {
		return nil, err
	}
	if t.version < 5 {
		return t, nil
	}

	// version >= 5.
	err = r.read(&t.usLowerOpticalPointSize, &t.usUpperOpticalPointSize)
	if err != nil {
		return nil, err
	}

	return t, nil
}
