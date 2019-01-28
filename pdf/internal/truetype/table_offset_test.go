/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/unidoc/unidoc/common"
)

// Test unmarshalling and marshalling offset table.
func TestOffsetTableReadWrite(t *testing.T) {
	testcases := []struct {
		fontPath string
		// Expected offset table parameters.
		expected offsetTable
	}{
		{
			"../../creator/testdata/FreeSans.ttf",
			offsetTable{
				sfntVersion:   0x10000, // opentype
				numTables:     16,
				searchRange:   256,
				entrySelector: 4,
				rangeShift:    0,
			},
		},
		{
			"../../creator/testdata/wts11.ttf",
			offsetTable{
				sfntVersion:   0x10000, // opentype
				numTables:     15,
				searchRange:   128,
				entrySelector: 3,
				rangeShift:    112,
			},
		},
		{
			"../../creator/testdata/roboto/Roboto-BoldItalic.ttf",
			offsetTable{
				sfntVersion:   0x10000, // opentype
				numTables:     18,
				searchRange:   256,
				entrySelector: 4,
				rangeShift:    32,
			},
		},
	}

	for _, tcase := range testcases {
		t.Logf("%s", tcase.fontPath)
		fnt, err := ParseFile(tcase.fontPath)
		require.NoError(t, err)
		assert.Equal(t, tcase.expected, *fnt.ot)

		common.Log.Debug("Write offset table")
		// Marshall to buffer.
		var buf bytes.Buffer
		bw := newByteWriter(&buf)
		err = fnt.writeOffsetTable(bw)
		require.NoError(t, err)
		bw.flush()

		// Reload from buffer.
		br := newByteReader(bytes.NewReader(buf.Bytes()))
		ot, err := fnt.parseOffsetTable(br)
		require.NoError(t, err)
		assert.Equal(t, fnt.ot, ot)
	}
}
