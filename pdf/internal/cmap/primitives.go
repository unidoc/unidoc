/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cmap

type cmapObject interface {
}

type cmapName struct {
	Name string
}

type cmapOperand struct {
	Operand string
}

type cmapHexString struct {
	numBytes int // original number of bytes in the raw representation
	b        []byte
}

type cmapString struct {
	String string
}

type cmapArray struct {
	Array []cmapObject
}

type cmapDict struct {
	Dict map[string]cmapObject
}

type cmapFloat struct {
	val float64
}

type cmapInt struct {
	val int64
}

func makeDict() cmapDict {
	d := cmapDict{}
	d.Dict = map[string]cmapObject{}
	return d
}
