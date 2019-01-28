/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNameTable(t *testing.T) {
	testcases := []struct {
		fontPath   string
		numEntries int
		expected   map[int]string
	}{
		{
			"../../creator/testdata/FreeSans.ttf",
			24,
			map[int]string{
				0:  "Copyleft 2002, 2003, 2005 Free Software Foundation.",
				1:  "FreeSans",
				2:  "Medium",
				4:  "Free Sans",
				13: "The use of this font is granted subject to GNU General Public License.",
				19: "The quick brown fox jumps over the lazy dog.",
			},
		},
		{
			"../../creator/testdata/wts11.ttf",
			44,
			map[int]string{
				0:  "(C)Copyright Dr. Hann-Tzong Wang, 2002-2004.",
				1:  "HanWang KaiBold-Gb5",
				2:  "Regular",
				3:  "HanWang KaiBold-Gb5",
				4:  "HanWang KaiBold-Gb5",
				6:  "HanWang KaiBold-Gb5",
				7:  "HanWang KaiBold-Gb5 is a registered trademark of HtWang Graphics Laboratory",
				14: "http://www.gnu.org/licenses/gpl.txt",
			},
		},
		{
			"../../creator/testdata/roboto/Roboto-BoldItalic.ttf",
			26,
			map[int]string{
				0:  "Copyright 2011 Google Inc. All Rights Reserved.",
				1:  "Roboto",
				2:  "Bold Italic",
				3:  "Roboto Bold Italic",
				4:  "Roboto Bold Italic",
				5:  "Version 2.137; 2017",
				6:  "Roboto-BoldItalic",
				14: "http://www.apache.org/licenses/LICENSE-2.0",
			},
		},
	}

	for _, tcase := range testcases {
		t.Run(tcase.fontPath, func(t *testing.T) {
			f, err := os.Open(tcase.fontPath)
			assert.Equal(t, nil, err)
			defer f.Close()

			br := newByteReader(f)
			fnt, err := parseFont(br)
			assert.Equal(t, nil, err)
			require.NoError(t, err)

			require.NotNil(t, fnt)
			require.NotNil(t, fnt.name)
			require.NotNil(t, fnt.name.nameRecords)

			assert.Equal(t, tcase.numEntries, len(fnt.name.nameRecords))
			for nameID, expStr := range tcase.expected {
				assert.Equal(t, expStr, fnt.GetNameByID(nameID))
			}

			for _, nr := range fnt.name.nameRecords {
				t.Logf("%d/%d/%d - '%s'", nr.platformID, nr.encodingID, nr.nameID, nr.Decoded())
			}
		})
	}
}
