/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import "github.com/unidoc/unidoc/common"

// hheaTable represents the horizontal header table (hhea).
// This table contains information for horizontal layout.
// https://docs.microsoft.com/en-us/typography/opentype/spec/hhea
type hheaTable struct {
	majorVersion        uint16
	minorVersion        uint16
	ascender            fword
	descender           fword
	lineGap             fword
	advanceWidthMax     ufword
	minLeftSideBearing  fword
	minRightSideBearing fword
	xMaxExtent          fword
	caretSlopeRise      int16
	caretSlopeRun       int16
	caretOffset         int16
	metricDataFormat    int16
	numberOfHMetrics    uint16 // Number of hMetric entries in 'hmtx' table.
}

func (f *font) parseHhea(r *byteReader) (*hheaTable, error) {
	_, has, err := f.seekToTable(r, "hhea")
	if err != nil {
		return nil, err
	}
	if !has {
		common.Log.Debug("hhea table absent")
		return nil, nil
	}

	t := &hheaTable{}
	err = r.read(&t.majorVersion, &t.minorVersion)
	if err != nil {
		return nil, err
	}

	err = r.read(&t.ascender, &t.descender, &t.lineGap)
	if err != nil {
		return nil, err
	}

	err = r.read(&t.advanceWidthMax, &t.minLeftSideBearing, &t.minRightSideBearing, &t.xMaxExtent)
	if err != nil {
		return nil, err
	}

	err = r.read(&t.caretSlopeRise, &t.caretSlopeRun, &t.caretOffset)
	if err != nil {
		return nil, err
	}

	// Skip over reserved bytes.
	r.Skip(4 * 2)

	return t, r.read(&t.metricDataFormat, &t.numberOfHMetrics)
}

func (f *font) writeHhea(w *byteWriter) error {
	if f.hhea == nil {
		common.Log.Debug("hhea is nil - nothing to write")
		return nil
	}

	t := f.hhea
	err := w.write(t.majorVersion, t.minorVersion)
	if err != nil {
		return err
	}

	err = w.write(t.ascender, t.descender, t.lineGap)
	if err != nil {
		return err
	}

	err = w.write(t.advanceWidthMax, t.minLeftSideBearing, t.minRightSideBearing, t.xMaxExtent)
	if err != nil {
		return err
	}

	err = w.write(t.caretSlopeRise, t.caretSlopeRun, t.caretOffset)
	if err != nil {
		return err
	}

	reserved := int16(0)
	err = w.write(&reserved, &reserved, &reserved, &reserved)
	if err != nil {
		return err
	}

	return w.write(t.metricDataFormat, t.numberOfHMetrics)
}
