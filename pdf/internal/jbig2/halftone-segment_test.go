package jbig2

import (
	"github.com/stretchr/testify/require"
	"github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/internal/jbig2/reader"
	"testing"
)

func TestHalftoneSegment(t *testing.T) {
	if testing.Verbose() {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	}

	t.Run("AnnexH", func(t *testing.T) {
		t.Run("S-7th", func(t *testing.T) {
			p := &Page{PageNumber: 1, Segments: make(map[int]*SegmentHeader)}
			d := &Document{
				Pages:          map[int]*Page{1: p},
				GlobalSegments: Globals(make(map[int]*SegmentHeader)),
			}
			p.Document = d

			patternData := []byte{
				// Header
				0x00, 0x00, 0x00, 0x05, 0x10, 0x01, 0x01, 0x00, 0x00, 0x00, 0x2D,

				// Data part
				0x01, 0x04, 0x04, 0x00, 0x00, 0x00, 0x0F, 0x20, 0xD1, 0x84,
				0x61, 0x18, 0x45, 0xF2, 0xF9, 0x7C, 0x8F, 0x11, 0xC3, 0x9E,
				0x45, 0xF2, 0xF9, 0x7D, 0x42, 0x85, 0x0A, 0xAA, 0x84, 0x62,
				0x2F, 0xEE, 0xEC, 0x44, 0x62, 0x22, 0x35, 0x2A, 0x0A, 0x83,
				0xB9, 0xDC, 0xEE, 0x77, 0x80,
			}

			halftoneData := []byte{
				// Header
				0x00, 0x00, 0x00, 0x06, 0x17, 0x20, 0x05, 0x01, 0x00, 0x00, 0x00, 0x57,

				// Data Part
				0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00,
				0x00, 0x0F, 0x00, 0x01, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0xAA, 0xAA, 0xAA, 0xAA,
				0x80, 0x08, 0x00, 0x80, 0x36, 0xD5, 0x55, 0x6B, 0x5A, 0xD4, 0x00, 0x40, 0x04, 0x2E,
				0xE9, 0x52, 0xD2, 0xD2, 0xD2, 0x8A, 0xA5, 0x4A, 0x00, 0x20, 0x02, 0x23, 0xE0, 0x95,
				0x24, 0xB4, 0x92, 0x8A, 0x4A, 0x92, 0x54, 0x92, 0xD2, 0x4A, 0x29, 0x2A, 0x49, 0x40,
				0x04, 0x00, 0x40,
			}
			common.Log.Debug("Pattern Data Length: %d", len(patternData))
			common.Log.Debug("Halftone Data Length: %d", len(halftoneData))
			var data []byte = append(patternData, halftoneData...)

			r := reader.New(data)
			// init by adding pattern dictionaryt
			getPattern := func(t *testing.T) *SegmentHeader {
				t.Helper()

				h, err := NewHeader(d, r, 0, OSequential)
				require.NoError(t, err)

				// h.SegmentDataStartOffset = 10
				// h.SegmentDataLength = 45

				common.Log.Debug("%#v", h)

				return h
			}

			ph := getPattern(t)

			offset := r.StreamPosition() + int64(ph.SegmentDataLength)

			h, err := NewHeader(d, r, offset, OSequential)
			require.NoError(t, err)

			common.Log.Debug("%#v", h)

			// h.SegmentDataStartOffset = 45 + 11 + 12

			h.RTSegments = append(h.RTSegments, ph)

			// s, err := h.subInputReader()
			// require.NoError(t, err)
			// hr := newHalftoneRegion(r)
			// require.NoError(t, hr.Init(h, s))
			hr, err := h.getSegmentData()
			require.NoError(t, err)

			bm, err := hr.(*HalftoneRegion).GetRegionBitmap()
			require.NoError(t, err)
			t.Logf("Pattern bitmap: %v", bm)
		})
	})

}
