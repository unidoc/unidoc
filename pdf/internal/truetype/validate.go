/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"bytes"
	"errors"
	"io"

	"github.com/unidoc/unidoc/common"
)

// validate font data model `f` in `r`. Checks if required tables are present and whether
// table checksums are correct.
func (f *font) validate(r *byteReader) error {
	if f.trec == nil {
		common.Log.Debug("Table records missing")
		return errRequiredField
	}
	if f.ot == nil {
		common.Log.Debug("Offsets table missing")
		return errRequiredField
	}
	if f.head == nil {
		common.Log.Debug("head table missing")
		return errRequiredField
	}

	// Validate the font.
	common.Log.Debug("Validating entire font")
	{
		err := r.Seek(0)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		_, err = io.Copy(&buf, r.reader)
		if err != nil {
			return err
		}

		data := buf.Bytes()

		headRec, ok := f.trec.trMap["head"]
		if !ok {
			common.Log.Debug("head not set")
			return errRequiredField
		}
		hoff := headRec.offset

		// set checksumAdjustment data to 0 in the head table.
		data[hoff+8] = 0
		data[hoff+9] = 0
		data[hoff+10] = 0
		data[hoff+11] = 0

		bw := newByteWriter(&bytes.Buffer{})
		bw.buffer.Write(data)

		checksum := bw.checksum()
		adjustment := 0xB1B0AFBA - checksum
		if f.head.checksumAdjustment != adjustment {
			return errors.New("file checksum mismatch")
		}
	}

	// Validate each table.
	common.Log.Debug("Validating font tables")
	for _, tr := range f.trec.list {
		common.Log.Debug("Validating %s", tr.tableTag.String())
		common.Log.Debug("%+v", tr)

		bw := newByteWriter(&bytes.Buffer{})

		if tr.offset < 0 || tr.length < 0 {
			common.Log.Debug("Range check error")
			return errRangeCheck
		}

		common.Log.Debug("Seeking to %d, to read %d bytes", tr.offset, tr.length)
		err := r.Seek(int64(tr.offset))
		if err != nil {
			return err
		}
		common.Log.Debug("Offset: %d", r.Offset())

		b := make([]byte, tr.length)
		_, err = io.ReadFull(r.reader, b)
		if err != nil {
			return err
		}
		common.Log.Debug("Read (%d)", len(b))
		// TODO(gunnsth): Validate head.
		if tr.tableTag.String() == "head" {
			// Set the checksumAdjustment to 0 so that head checksum is valid.
			if len(b) < 12 {
				return errors.New("head too short")
			}
			b[8], b[9], b[10], b[11] = 0, 0, 0, 0
		}

		_, err = bw.buffer.Write(b)
		if err != nil {
			return err
		}

		checksum := bw.checksum()
		if tr.checksum != checksum {
			common.Log.Debug("Invalid checksum (%d != %d)", checksum, tr.checksum)
			return errors.New("checksum incorrect")
		}

		if int(tr.length) != bw.bufferedLen() {
			common.Log.Debug("Length mismatch")
			return errRangeCheck
		}
	}

	return nil
}
