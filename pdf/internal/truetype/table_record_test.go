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

// Test unmarshalling and marshalling table records.
func TestTableRecordsReadWrite(t *testing.T) {
	testcases := []struct {
		fontPath string
		expected []tableRecord
	}{
		{
			"../../creator/testdata/FreeSans.ttf",
			[]tableRecord{
				{
					tableTag: makeTag("FFTM"), // FontForge specific table.
					checksum: 1195616530,
					offset:   459736,
					length:   28,
				},
				{
					tableTag: makeTag("GDEF"),
					checksum: 31456477,
					offset:   433972,
					length:   1632,
				},
				{
					tableTag: makeTag("GPOS"),
					checksum: 4278766266,
					offset:   447632,
					length:   12102,
				},
				{
					tableTag: makeTag("GSUB"),
					checksum: 3391961157,
					offset:   435604,
					length:   12026,
				},
				{
					tableTag: makeTag("OS/2"),
					checksum: 3829110115,
					offset:   392,
					length:   86,
				},
				{
					tableTag: makeTag("cmap"),
					checksum: 4271469241,
					offset:   15376,
					length:   2526,
				},
				{
					tableTag: makeTag("cvt"),
					checksum: 2163321,
					offset:   17904,
					length:   4,
				},
				{
					tableTag: makeTag("gasp"),
					checksum: 4294901763,
					offset:   433964,
					length:   8,
				},
				{
					tableTag: makeTag("glyf"),
					checksum: 843000928,
					offset:   32816,
					length:   354716,
				},
				{
					tableTag: makeTag("head"),
					checksum: 3924650013,
					offset:   268,
					length:   54,
				},
				{
					tableTag: makeTag("hhea"),
					checksum: 124129540,
					offset:   324,
					length:   36,
				},

				{
					tableTag: makeTag("hmtx"),
					checksum: 2335681020,
					offset:   480,
					length:   14896,
				},
				{
					tableTag: makeTag("loca"),
					checksum: 537012616,
					offset:   17908,
					length:   14908,
				},
				{
					tableTag: makeTag("maxp"),
					checksum: 262341762,
					offset:   360,
					length:   32,
				},
				{
					tableTag: makeTag("name"),
					checksum: 2006447137,
					offset:   387532,
					length:   1521,
				},
				{
					tableTag: makeTag("post"),
					checksum: 964072869,
					offset:   389056,
					length:   44907,
				},
			},
		},
	}

	for _, tcase := range testcases {
		t.Logf("%s", tcase.fontPath)
		fnt, err := ParseFile(tcase.fontPath)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		assert.Equal(t, tcase.expected, fnt.trec.list)

		common.Log.Debug("Write table records")
		// Marshall to buffer.
		var buf bytes.Buffer
		bw := newByteWriter(&buf)
		err = fnt.writeTableRecords(bw)
		require.NoError(t, err)
		bw.flush()

		// Reload from buffer and check equality.
		br := newByteReader(bytes.NewReader(buf.Bytes()))
		trs, err := fnt.parseTableRecords(br)
		require.NoError(t, err)
		assert.Equal(t, fnt.trec.list, trs.list)
	}
}
