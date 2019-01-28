/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import "fmt"

// cvtTable represents the Control Value Table (cvt).
// This table contains a list of values that can be referenced by instructions.
type cvtTable struct {
	n      int // number of FWORD items that fits the size of the table.
	values []fword
}

func (t *cvtTable) Unmarshal(r *byteReader) error {
	if t.n == 0 {
		fmt.Printf("n == 0\n")
	}
	return r.readSlice(&t.values, t.n)
}

func (t cvtTable) Marshal(w *byteWriter) error {
	return w.writeSlice(t.values)
}
