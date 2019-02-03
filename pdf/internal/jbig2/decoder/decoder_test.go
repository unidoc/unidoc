package decoder

import (
	"github.com/stretchr/testify/assert"
	"github.com/unidoc/unidoc/common"
	"testing"
)

var data = []byte{
	0x97, 0x4A, 0x42, 0x32, 0x0D, 0x0A, 0x1A, 0x0A, 0x01, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x18, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
	0x00, 0x01, 0xE9, 0xCB, 0xF4, 0x00, 0x26, 0xAF, 0x04, 0xBF, 0xF0, 0x78, 0x2F, 0xE0, 0x00, 0x40,
	0x00, 0x00, 0x00, 0x01, 0x30, 0x00, 0x01, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x00, 0x40, 0x00,
	0x00, 0x00, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x02, 0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x1C, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00,
	0x00, 0x00, 0x02, 0xE5, 0xCD, 0xF8, 0x00, 0x79, 0xE0, 0x84, 0x10, 0x81, 0xF0, 0x82, 0x10, 0x86,
	0x10, 0x79, 0xF0, 0x00, 0x80, 0x00, 0x00, 0x00, 0x03, 0x07, 0x42, 0x00, 0x02, 0x01, 0x00, 0x00,
	0x00, 0x31, 0x00, 0x00, 0x00, 0x25, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00,
	0x00, 0x01, 0x00, 0x0C, 0x09, 0x00, 0x10, 0x00, 0x00, 0x00, 0x05, 0x01, 0x10, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0C, 0x40, 0x07, 0x08,
	0x70, 0x41, 0xD0, 0x00, 0x00, 0x00, 0x04, 0x27, 0x00, 0x01, 0x00, 0x00, 0x00, 0x2C, 0x00, 0x00,
	0x00, 0x36, 0x00, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0B, 0x00, 0x01,
	0x26, 0xA0, 0x71, 0xCE, 0xA7, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF8, 0xF0, 0x00, 0x00, 0x00, 0x05, 0x10, 0x01,
	0x01, 0x00, 0x00, 0x00, 0x2D, 0x01, 0x04, 0x04, 0x00, 0x00, 0x00, 0x0F, 0x20, 0xD1, 0x84, 0x61,
	0x18, 0x45, 0xF2, 0xF9, 0x7C, 0x8F, 0x11, 0xC3, 0x9E, 0x45, 0xF2, 0xF9, 0x7D, 0x42, 0x85, 0x0A,
	0xAA, 0x84, 0x62, 0x2F, 0xEE, 0xEC, 0x44, 0x62, 0x22, 0x35, 0x2A, 0x0A, 0x83, 0xB9, 0xDC, 0xEE,
	0x77, 0x80, 0x00, 0x00, 0x00, 0x06, 0x17, 0x20, 0x05, 0x01, 0x00, 0x00, 0x00, 0x57, 0x00, 0x00,
	0x00, 0x20, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x0F, 0x00, 0x01,
	0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x04, 0x00, 0x00, 0x00, 0xAA, 0xAA, 0xAA, 0xAA, 0x80, 0x08, 0x00, 0x80, 0x36, 0xD5, 0x55, 0x6B,
	0x5A, 0xD4, 0x00, 0x40, 0x04, 0x2E, 0xE9, 0x52, 0xD2, 0xD2, 0xD2, 0x8A, 0xA5, 0x4A, 0x00, 0x20,
	0x02, 0x23, 0xE0, 0x95, 0x24, 0xB4, 0x92, 0x8A, 0x4A, 0x92, 0x54, 0x92, 0xD2, 0x4A, 0x29, 0x2A,
	0x49, 0x40, 0x04, 0x00, 0x40, 0x00, 0x00, 0x00, 0x07, 0x31, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x08, 0x30, 0x00, 0x02, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x00, 0x40, 0x00,
	0x00, 0x00, 0x38, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x09, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x1B, 0x08, 0x00, 0x02, 0xFF, 0x00, 0x00, 0x00,
	0x02, 0x00, 0x00, 0x00, 0x02, 0x4F, 0xE7, 0x8C, 0x20, 0x0E, 0x1D, 0xC7, 0xCF, 0x01, 0x11, 0xC4,
	0xB2, 0x6F, 0xFF, 0xAC, 0x00, 0x00, 0x00, 0x0A, 0x07, 0x40, 0x00, 0x09, 0x02, 0x00, 0x00, 0x00,
	0x1F, 0x00, 0x00, 0x00, 0x25, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
	0x01, 0x00, 0x0C, 0x08, 0x00, 0x00, 0x00, 0x05, 0x8D, 0x6E, 0x5A, 0x12, 0x40, 0x85, 0xFF, 0xAC,
	0x00, 0x00, 0x00, 0x0B, 0x27, 0x00, 0x02, 0x00, 0x00, 0x00, 0x23, 0x00, 0x00, 0x00, 0x36, 0x00,
	0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x0B, 0x00, 0x08, 0x03, 0xFF, 0xFD,
	0xFF, 0x02, 0xFE, 0xFE, 0xFE, 0x04, 0xEE, 0xED, 0x87, 0xFB, 0xCB, 0x2B, 0xFF, 0xAC, 0x00, 0x00,
	0x00, 0x0C, 0x10, 0x01, 0x02, 0x00, 0x00, 0x00, 0x1C, 0x06, 0x04, 0x04, 0x00, 0x00, 0x00, 0x0F,
	0x90, 0x71, 0x6B, 0x6D, 0x99, 0xA7, 0xAA, 0x49, 0x7D, 0xF2, 0xE5, 0x48, 0x1F, 0xDC, 0x68, 0xBC,
	0x6E, 0x40, 0xBB, 0xFF, 0xAC, 0x00, 0x00, 0x00, 0x0D, 0x17, 0x20, 0x0C, 0x02, 0x00, 0x00, 0x00,
	0x3E, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x24, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00,
	0x0F, 0x00, 0x02, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x87, 0xCB, 0x82, 0x1E, 0x66, 0xA4, 0x14, 0xEB, 0x3C,
	0x4A, 0x15, 0xFA, 0xCC, 0xD6, 0xF3, 0xB1, 0x6F, 0x4C, 0xED, 0xBF, 0xA7, 0xBF, 0xFF, 0xAC, 0x00,
	0x00, 0x00, 0x0E, 0x31, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x30, 0x00,
	0x03, 0x00, 0x00, 0x00, 0x13, 0x00, 0x00, 0x00, 0x25, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x01, 0x00, 0x00,
	0x00, 0x00, 0x16, 0x08, 0x00, 0x02, 0xFF, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x4F,
	0xE7, 0x8D, 0x68, 0x1B, 0x14, 0x2F, 0x3F, 0xFF, 0xAC, 0x00, 0x00, 0x00, 0x11, 0x00, 0x21, 0x10,
	0x03, 0x00, 0x00, 0x00, 0x20, 0x08, 0x02, 0x02, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00,
	0x03, 0x00, 0x00, 0x00, 0x02, 0x4F, 0xE9, 0xD7, 0xD5, 0x90, 0xC3, 0xB5, 0x26, 0xA7, 0xFB, 0x6D,
	0x14, 0x98, 0x3F, 0xFF, 0xAC, 0x00, 0x00, 0x00, 0x12, 0x07, 0x20, 0x11, 0x03, 0x00, 0x00, 0x00,
	0x25, 0x00, 0x00, 0x00, 0x25, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x8C, 0x12, 0x00, 0x00, 0x00, 0x04, 0xA9, 0x5C, 0x8B, 0xF4, 0xC3, 0x7D, 0x96, 0x6A,
	0x28, 0xE5, 0x76, 0x8F, 0xFF, 0xAC, 0x00, 0x00, 0x00, 0x13, 0x31, 0x00, 0x03, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x14, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func TestDecoderStreamDecode(t *testing.T) {
	if testing.Verbose() {
		common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))
	}

	d := New()

	b, err := d.DecodeBytes(data)
	if assert.NoError(t, err) {
		t.Logf("%v", b)
	}
}

var i [][]int = [][]int{{0, 0, 0, 0}, {1, 1, 0, 0}, {2, 1, 0, 0}, {3, 0, 0, 0}, {4, 0, 0, 0}, {5, 0, 0, 0}, {6, 0, 0, 0}, {7, 0, 0, 0}, {8, 0, 0, 0}, {9, 0, 0, 0}, {10, 0, 0, 0}, {11, 0, 0, 0}, {12, 0, 0, 0}, {13, 0, 0, 0}, {14, 0, 0, 0}, {15, 0, 0, 0}, {16, 0, 0, 0}, {17, 0, 0, 0}, {18, 0, 0, 0}, {19, 0, 0, 0}, {20, 0, 0, 0}, {21, 0, 0, 0}, {22, 0, 0, 0}, {23, 0, 0, 0}, {24, 0, 0, 0}, {25, 0, 0, 0}, {26, 0, 0, 0}, {27, 0, 0, 0}, {28, 0, 0, 0}, {29, 0, 0, 0}, {30, 0, 0, 0}, {31, 0, 0, 0}, {259, 0, 2, 0}, {515, 0, 3, 0}, {523, 0, 7, 0}, {0, 0, 4294967295, 0}}
