/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import "github.com/unidoc/unidoc/common"

// maxpTable represents the Maximum Profile (maxp) table.
// This table establishes the memory requirements for the font.
type maxpTable struct {
	// Version 0.5 and above:
	version   fixed
	numGlyphs uint16

	// Version 1.0 and above:
	maxPoints             uint16
	maxContours           uint16
	maxCompositePoints    uint16
	maxCompositeContours  uint16
	maxZones              uint16
	maxTwilightPoints     uint16
	maxStorage            uint16
	maxFunctionDefs       uint16
	maxInstructionDefs    uint16
	maxStackElements      uint16
	maxSizeOfInstructions uint16
	maxComponentElements  uint16
	maxComponentDepth     uint16
}

func (f *font) parseMaxp(r *byteReader) (*maxpTable, error) {
	_, has, err := f.seekToTable(r, "maxp")
	if err != nil {
		return nil, err
	}
	if !has {
		common.Log.Debug("maxp table not present")
		return nil, nil
	}

	t := &maxpTable{}

	err = r.read(&t.version, &t.numGlyphs)
	if err != nil {
		return nil, err
	}

	if t.version < 0x00010000 {
		common.Log.Debug("Range check error")
		return nil, errRangeCheck
	}

	err = r.read(&t.maxPoints, &t.maxContours, &t.maxCompositePoints, &t.maxCompositeContours)
	if err != nil {
		return nil, err
	}

	err = r.read(&t.maxZones, &t.maxTwilightPoints, &t.maxSizeOfInstructions, &t.maxFunctionDefs, &t.maxInstructionDefs)
	if err != nil {
		return nil, err
	}

	return t, r.read(&t.maxStackElements, &t.maxSizeOfInstructions, &t.maxComponentElements, &t.maxComponentDepth)
}

func (f *font) writeMaxp(w *byteWriter) error {
	if f.maxp == nil {
		return errRequiredField
	}
	t := f.maxp
	err := w.write(t.version, t.numGlyphs)
	if err != nil {
		return err
	}

	if t.version < 0x00010000 {
		common.Log.Debug("Range check error")
		return errRangeCheck
	}

	err = w.write(t.maxPoints, t.maxContours, t.maxCompositePoints, t.maxCompositeContours)
	if err != nil {
		return err
	}

	err = w.write(t.maxZones, t.maxTwilightPoints, t.maxStorage, t.maxFunctionDefs, t.maxInstructionDefs)
	if err != nil {
		return err
	}

	return w.write(t.maxStackElements, t.maxSizeOfInstructions, t.maxComponentElements, t.maxComponentDepth)
}
