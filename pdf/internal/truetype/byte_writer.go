/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/unidoc/unidoc/common"
)

// byteWriter encapsulates io.Writer and provides methods to write binary data as fit for truetype fonts.
// Writes are buffered until flushed. Provides methods to calculate checksum of the current buffer.
type byteWriter struct {
	w   io.Writer
	len int64

	buffer bytes.Buffer
}

func newByteWriter(w io.Writer) *byteWriter {
	return &byteWriter{
		w: w,
	}
}

func (w *byteWriter) flush() error {
	b := w.buffer.Bytes()
	_, err := w.w.Write(b)
	if err != nil {
		return err
	}

	w.buffer.Reset()
	return nil
}

// bufferedLen returns the length of the current buffer.
func (w *byteWriter) bufferedLen() int {
	return w.buffer.Len()
}

// checksum returns the checksum of the current buffer.
func (w *byteWriter) checksum() uint32 {
	var sum uint32

	data := w.buffer.Bytes()

	if len(data) < 60 {
		common.Log.Debug("Data: % X", data)
	}
	common.Log.Debug("Data length: %d", len(data))
	sum = 0

	for i := 0; i < len(data); i += 4 {
		a := i
		b := i + 4
		if b > len(data) {
			b = len(data)
		}

		dup := make([]byte, 4)
		copy(dup, data[a:b])

		if b-a < 4 {
			for j := 0; j < b-a; j++ {
				dup = append(dup, 0) //
			}
		}

		val := binary.BigEndian.Uint32(dup)
		sum += val
	}

	return sum
}

func (w *byteWriter) writeSlice(slice interface{}) error {
	switch t := slice.(type) {
	case []uint8:
		for _, val := range t {
			err := w.writeUint8(val)
			if err != nil {
				return err
			}
		}
	case *[]uint8:
		for _, val := range *t {
			err := w.writeUint8(val)
			if err != nil {
				return err
			}
		}
	case []uint16:
		for _, val := range t {
			err := w.writeUint16(val)
			if err != nil {
				return err
			}
		}
	case *[]uint16:
		for _, val := range *t {
			err := w.writeUint16(val)
			if err != nil {
				return err
			}
		}
	default:
		fmt.Printf("Write type check error: %T (slice)\n", t)
		return errTypeCheck
	}
	return nil
}

// Write a series of values to `w`.
func (w *byteWriter) write(fields ...interface{}) error {
	for _, f := range fields {
		switch t := f.(type) {
		case uint8:
			err := w.writeUint8(t)
			if err != nil {
				return err
			}
		case uint16:
			err := w.writeUint16(t)
			if err != nil {
				return err
			}
		case int16:
			err := w.writeInt16(t)
			if err != nil {
				return err
			}
		case uint32:
			err := w.writeUint32(t)
			if err != nil {
				return err
			}
		case tag:
			err := w.writeTag(t)
			if err != nil {
				return err
			}
		case offset16:
			err := w.writeOffset16(t)
			if err != nil {
				return err
			}
		case offset32:
			err := w.writeOffset32(t)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("Write type check error: %T\n", t)
			return errTypeCheck
		}
	}

	return nil
}

func (w *byteWriter) writeUint8(vals ...uint8) error {
	err := binary.Write(&w.buffer, binary.BigEndian, vals)
	if err != nil {
		return err
	}
	w.len++
	return nil
}

func (w *byteWriter) writeUint16(vals ...uint16) error {
	err := binary.Write(&w.buffer, binary.BigEndian, vals)
	if err != nil {
		return err
	}
	w.len += 2
	return nil
}

func (w *byteWriter) writeInt16(vals ...int16) error {
	err := binary.Write(&w.buffer, binary.BigEndian, vals)
	if err != nil {
		return err
	}
	w.len += 2
	return nil
}

func (w *byteWriter) writeUint32(val uint32) error {
	err := binary.Write(&w.buffer, binary.BigEndian, val)
	if err != nil {
		return err
	}
	w.len += 4
	return nil
}

func (w *byteWriter) writeTag(val tag) error {
	err := binary.Write(&w.buffer, binary.BigEndian, val)
	if err != nil {
		return err
	}
	w.len += 4
	return nil
}

func (w *byteWriter) writeOffset16(val offset16) error {
	err := binary.Write(&w.buffer, binary.BigEndian, val)
	if err != nil {
		return err
	}
	w.len += 2
	return nil
}

func (w *byteWriter) writeOffset32(val offset32) error {
	err := binary.Write(&w.buffer, binary.BigEndian, val)
	if err != nil {
		return err
	}
	w.len += 4
	return nil
}
