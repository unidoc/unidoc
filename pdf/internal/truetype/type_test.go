/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"testing"
)

func TestFixedParts(t *testing.T) {
	tcases := []struct {
		val fixed
		a   uint16
		b   uint16
		f64 float64
	}{
		{
			fixed(0x00011000),
			0x0001,
			0x1000,
			1.0625, // FIXME/TODO(gunnsth): Should be 1.1 ?
		},
		{
			fixed(0x00005000),
			0x0000,
			0x5000,
			0.3125, // FIXME/TODO(gunnsth): Should be 0.5 ?
		},
		{
			fixed(0x00025000),
			0x0002,
			0x5000,
			2.3125, // FIXME/TODO(gunnsth): Should be 2.5?
		},
	}

	for _, tcase := range tcases {
		a, b := tcase.val.Parts()
		if a != tcase.a {
			t.Fatalf("%d != %d", a, tcase.a)
		}
		if b != tcase.b {
			t.Fatalf("%d != %d", b, tcase.b)
		}
		f64 := tcase.val.Float64()
		if f64 != tcase.f64 {
			t.Fatalf("%v != %v", f64, tcase.f64)
		}
	}
}
