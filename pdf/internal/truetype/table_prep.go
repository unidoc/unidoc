/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import "fmt"

// prepTable represents a Control Value Program table (prep).
// Consists of a set of TrueType instructions that will be executed whenever the font or point size
// or transformation matrix change and before each glyph is interpreted.
// Used for preparation (hence the name "prep").
type prepTable struct {
	n            int // number of instructions - the number of uint8 that fit the size of the table.
	instructions []uint8
}

func (t *prepTable) Unmarshal(r *byteReader) error {
	if t.n == 0 {
		fmt.Printf("n == 0\n")
	}

	return r.readSlice(&t.instructions, t.n)
}

func (t prepTable) Marshal(w *byteWriter) error {
	return w.writeSlice(t.instructions)
}
