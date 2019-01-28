/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"io"
	"os"
)

// Font wraps font for outside access.
type Font struct {
	br *byteReader
	*font
}

// Parse parses the truetype font from `rs` and returns a new Font.
func Parse(rs io.ReadSeeker) (*Font, error) {
	r := newByteReader(rs)

	fnt, err := parseFont(r)
	if err != nil {
		return nil, err
	}

	return &Font{
		br:   r,
		font: fnt,
	}, nil
}

// ParseFile parses the truetype font from file given by path.
func ParseFile(filePath string) (*Font, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return Parse(f)
}

// ValidateFile validates the truetype font given by `filePath`.
func ValidateFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	br := newByteReader(f)
	fnt, err := parseFont(br)
	if err != nil {
		return err
	}

	return fnt.validate(br)
}

// Write writes the font to `w`.
func (f *Font) Write(w io.Writer) error {
	/*
		bw := newByteWriter(w)

		f.fnt.write(bw)

		err := f.offsetTable.Marshal(bw)
		if err != nil {
			return err
		}
	*/

	return nil
}
