/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import "github.com/unidoc/unidoc/common"

type hmtxTable struct {
	hMetrics         []longHorMetric // length is numberOfHMetrics from hhea table.
	leftSideBearings []int16         // length is numGlyphs - numberOfHmetrics from maxp and hhea tables.
}

type longHorMetric struct {
	advanceWidth uint16
	lsb          int16
}

func (f *font) parseHmtx(r *byteReader) (*hmtxTable, error) {
	if f.maxp == nil || f.hhea == nil {
		common.Log.Debug("maxp or hhea table missing")
		return nil, errRequiredField
	}

	_, has, err := f.seekToTable(r, "htmx")
	if err != nil {
		return nil, err
	}
	if !has {
		common.Log.Debug("hmtx table absent")
		return nil, nil
	}

	t := &hmtxTable{}

	numberOfHMetrics := int(f.hhea.numberOfHMetrics)
	for i := 0; i < numberOfHMetrics; i++ {
		var lhm longHorMetric
		err := r.read(&lhm.advanceWidth, &lhm.lsb)
		if err != nil {
			return nil, err
		}

		t.hMetrics = append(t.hMetrics, lhm)
	}

	lsbLen := int(f.maxp.numGlyphs) - numberOfHMetrics
	if lsbLen < 0 {
		common.Log.Debug("ERROR: Negative length")
		return nil, errRangeCheck
	}

	err = r.readSlice(&t.leftSideBearings, lsbLen)
	if err != nil {
		return nil, err
	}

	return t, nil
}
