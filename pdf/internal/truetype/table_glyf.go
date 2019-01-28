/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"errors"
	"strings"

	"github.com/unidoc/unidoc/common"
)

// glyfTable represents the Glyph Data table (glyf).
// Information that describes the glyphs in the font in the TrueType outline format.
//
// The 'glyf' table is comprised of a list of glyph data blocks, each of which provides
// the description for a single glyph. Glyphs are referenced by identifiers (glyph IDs),
// which are sequential integers beginning at zero. The total number of glyphs is specified
// by the numGlyphs field in the 'maxp' table. The 'glyf' table does not include any overall
// table header or records providing offsets to glyph data blocks. Rather, the 'loca' table
// provides an array of offsets, indexed by glyph IDs, which provide the location of each
// glyph data block within the 'glyf' table. Note that the 'glyf' table must always be used
// in conjunction with the 'loca' and 'maxp' tables.
// https://docs.microsoft.com/en-us/typography/opentype/spec/glyf
type glyfTable struct {
	descs []*glyphDescription
}

func (f *font) parseGlyf(r *byteReader) (*glyfTable, error) {
	if f.maxp == nil || f.loca == nil {
		common.Log.Debug("required field missing (glyf)")
		return nil, errRequiredField
	}

	tr, has, err := f.seekToTable(r, "glyf")
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil // table not found.
	}

	glyf := &glyfTable{}

	common.Log.Debug("parsing glyfs")
	common.Log.Debug("Number of glyphs: %d", f.maxp.numGlyphs)
	common.Log.Debug("Loca offset format: %d", f.head.indexToLocFormat)

	for i := 0; i < int(f.maxp.numGlyphs); i++ {
		gid := GlyphIndex(i)
		//common.Log.Debug("GID = %d", gid)
		gdOffset, gdLen, err := f.GetGlyphDataOffset(gid)
		if err != nil {
			return nil, err
		}

		if gdOffset > int64(tr.length) {
			common.Log.Debug("Range check error (glyf)")
			return nil, errRangeCheck
		}

		err = r.Seek(int64(tr.offset) + gdOffset)
		if err != nil {
			return nil, err
		}

		desc, err := f.parseGlyphDescription(r, gdLen)
		if err != nil {
			return nil, err
		}
		glyf.descs = append(glyf.descs, desc)
	}

	return glyf, nil
}

func (f *font) writeGlyf(w *byteWriter) error {
	if f.glyf == nil || f.maxp == nil || f.loca == nil {
		common.Log.Debug("glyf: required field missing (write)")
		return errRequiredField
	}

	if int(f.maxp.numGlyphs) != len(f.glyf.descs) {
		common.Log.Debug("Incorrect number of glyph descriptions")
		return errRangeCheck
	}

	for _, gd := range f.glyf.descs {
		err := gd.Write(w, f)
		if err != nil {
			return err
		}
	}

	return nil
}

type glyphDescription struct {
	header    glyfGlyphHeader
	simple    *simpleGlyphDescription
	composite *compositeGlyphDescription
}

func (d glyphDescription) IsSimple() bool {
	return d.simple != nil
}

func (f *font) parseGlyphDescription(r *byteReader, gdLen int64) (*glyphDescription, error) {
	var gh glyfGlyphHeader
	err := gh.read(r)
	if err != nil {
		return nil, err
	}
	common.Log.Trace("gh: %+v", gh)

	if gh.numberOfContours >= 0 {
		common.Log.Trace("simple glyph data, contours: %d", gh.numberOfContours)
		// Simple glyph.
		sgd, err := f.parseSimpleGlyphDescription(r, int(gh.numberOfContours))
		if err != nil {
			return nil, err
		}

		return &glyphDescription{
			header: gh,
			simple: sgd,
		}, nil
	}

	common.Log.Trace("composite glyph data")
	// Composite/compound glyph.
	cgd, err := f.parseCompositeGlyphDescription(r, gdLen)
	if err != nil {
		return nil, err
	}
	return &glyphDescription{
		header:    gh,
		composite: cgd,
	}, nil
}

func (d glyphDescription) Write(w *byteWriter, f *font) error {
	err := d.header.write(w)
	if err != nil {
		return err
	}

	if d.simple != nil {
		return d.simple.Write(w, f, int(d.header.numberOfContours))
	}

	// Composite.
	return d.composite.Write(w, f)
}

// glyfGlyphHeader represents the glyph header in the glyf table (one for each glyph).
type glyfGlyphHeader struct {
	// Header.
	numberOfContours int16
	xMin             int16
	yMin             int16
	xMax             int16
	yMax             int16
}

func (h *glyfGlyphHeader) read(r *byteReader) error {
	return r.read(&h.numberOfContours, &h.xMin, &h.yMin, &h.xMax, &h.yMax)
}

func (h *glyfGlyphHeader) write(w *byteWriter) error {
	return w.write(h.numberOfContours, h.xMin, h.yMin, h.xMax, h.yMax)
}

// simpleGlyphFlag represents a flag data representation of a point in a simple glyph.
type simpleGlyphFlag uint8

const (
	onCurvePoint simpleGlyphFlag = (1 << iota)
	xShortVector
	yShortVector
	repeatFlag
	xIsSameOrPositiveVector
	yIsSameOrPositiveVector
	overlapSimple
	reserved
)

func (f simpleGlyphFlag) String() string {
	var flags []string
	if f&onCurvePoint != 0 {
		flags = append(flags, "onCurvePoint")
	}
	if f&xShortVector != 0 {
		flags = append(flags, "xShortVector")
	}
	if f&yShortVector != 0 {
		flags = append(flags, "yShortVector")
	}
	if f&repeatFlag != 0 {
		flags = append(flags, "repeatFlag")
	}
	if f&xIsSameOrPositiveVector != 0 {
		flags = append(flags, "xIsSameOrPositiveVector")
	}
	if f&yIsSameOrPositiveVector != 0 {
		flags = append(flags, "yIsSameOrPositiveVector")
	}
	if f&overlapSimple != 0 {
		flags = append(flags, "overlapSimple")
	}
	if f&reserved != 0 {
		flags = append(flags, "reserved")
	}
	return strings.Join(flags, "|")
}

// simpleGlyphDescription represents simple glyph descriptions (non composite glyphs).
// This is the table information needed when `numberOfContours >= 0`, i.e. not composite glyphs.
type simpleGlyphDescription struct {
	// list of point indices for the last point of each contour, in increasing numeric order.
	endPtsOfContours []uint16 // numberOfContours elements.

	instructionLength uint16
	instructions      []uint8 // instructionLength elements.
	// one flag byte element, one x-coordinate, and one y-coordinate for each point
	flags        []uint8  // variable length?
	xCoordinates []uint16 // Can be either 8 or 16 bits (depends on corresponding flag).
	yCoordinates []uint16 // Can be either 8 or 16 bits (depends on corresponding flag).
}

// parses description for a single simple glyph with `numContours` at current position in `r`.
func (f *font) parseSimpleGlyphDescription(r *byteReader, numContours int) (*simpleGlyphDescription, error) {
	if numContours == 0 {
		return nil, nil
	}

	if f.loca == nil {
		common.Log.Debug("loca not set")
		return nil, errRequiredField
	}

	var d simpleGlyphDescription

	err := r.readSlice(&d.endPtsOfContours, numContours)
	if err != nil {
		return nil, err
	}
	err = r.read(&d.instructionLength)
	if err != nil {
		return nil, err
	}

	err = r.readSlice(&d.instructions, int(d.instructionLength))
	if err != nil {
		return nil, err
	}

	// total number of points (all contours).
	numPoints := int(d.endPtsOfContours[numContours-1]) + 1

	common.Log.Trace("GID data - Number of points: %d", numPoints)

	// flags (one for each point).
	numFlags := 0
	for numFlags < numPoints {
		var flag uint8
		err := r.read(&flag)
		if err != nil {
			return nil, err
		}
		common.Log.Trace("flag: %d (%s)", flag, simpleGlyphFlag(flag).String())

		d.flags = append(d.flags, flag)
		numFlags++

		if simpleGlyphFlag(flag)&repeatFlag != 0 {
			// following byte specifies number of times this flag is to be repeated.
			var repeats uint8
			err := r.read(&repeats)
			if err != nil {
				return nil, err
			}
			common.Log.Trace("Repeats: %d", repeats)
			for i := 0; i < int(repeats); i++ {
				d.flags = append(d.flags, flag)
				numFlags++
			}
		}
	}
	if numFlags != numPoints {
		common.Log.Debug("Number of flags != number of points (%d != %d)", numFlags, numPoints)
		return nil, errors.New("numflags != numpoints")
	}
	common.Log.Trace("Number of flags: %d", numFlags)
	common.Log.Trace("Flags: % d\n", d.flags)
	common.Log.Trace("@Offset: %d", r.Offset())

	// x coordinates.
	common.Log.Trace("X Coordinates")
	var xLast uint16
	for _, flag := range d.flags {
		sflag := simpleGlyphFlag(flag)
		if sflag&xShortVector != 0 {
			var x uint8
			err := r.read(&x)
			if err != nil {
				return nil, err
			}
			common.Log.Trace("Short data - x=%d", x)
			d.xCoordinates = append(d.xCoordinates, uint16(x))
			xLast = uint16(x)
		} else {
			if sflag&xIsSameOrPositiveVector != 0 {
				common.Log.Trace("Long data - same as last")
				d.xCoordinates = append(d.xCoordinates, xLast)
			} else {
				var x uint16
				err := r.read(&x)
				if err != nil {
					return nil, err
				}
				common.Log.Trace("Long data - x=%d", x)
				d.xCoordinates = append(d.xCoordinates, x)
				xLast = x
			}
		}
	}

	// y coordinates.
	common.Log.Trace("Y Coordinates")
	var yLast uint16
	for _, flag := range d.flags {
		sflag := simpleGlyphFlag(flag)
		if sflag&yShortVector != 0 {
			var y uint8
			err := r.read(&y)
			if err != nil {
				return nil, err
			}
			common.Log.Trace("Short data - y=%d", y)
			d.yCoordinates = append(d.yCoordinates, uint16(y))
			yLast = uint16(y)
		} else {
			if sflag&yIsSameOrPositiveVector != 0 {
				common.Log.Trace("Long data - same as last")
				d.yCoordinates = append(d.yCoordinates, yLast)
			} else {
				var y uint16
				err := r.read(&y)
				if err != nil {
					return nil, err
				}
				common.Log.Trace("Long data - y=%d", y)
				d.yCoordinates = append(d.yCoordinates, y)
				yLast = y
			}
		}
	}

	return &d, nil
}

func (d *simpleGlyphDescription) Write(w *byteWriter, f *font, numContours int) error {
	if f == nil || f.loca == nil {
		common.Log.Debug("sgd: required field missing (write)")
		return errRequiredField
	}

	err := w.writeSlice(d.endPtsOfContours)
	if err != nil {
		return err
	}

	err = w.write(d.instructionLength)
	if err != nil {
		return err
	}

	err = w.writeSlice(d.instructions)
	if err != nil {
		return err
	}

	if numContours > len(d.endPtsOfContours) {
		common.Log.Debug("range check error (numContours)")
		return errRangeCheck
	}

	numPoints := int(d.endPtsOfContours[numContours-1]) + 1
	if len(d.flags) != numPoints {
		common.Log.Debug("#flags != #points (%d/%d)", len(d.flags), numPoints)
		return errRangeCheck
	}

	// flags - packed.
	i := 0
	for i < len(d.flags) {
		flag := d.flags[i]
		var j int
		for j = i + 1; j < len(d.flags) && j-i <= 255; j++ {
			if d.flags[j] != flag {
				break
			}
		}

		repeats := uint8(j - i)
		if repeats > 1 {
			flag |= uint8(repeatFlag)
		}

		err = w.write(flag)
		if err != nil {
			return err
		}
		if repeats > 1 {
			err = w.write(repeats)
			if err != nil {
				return err
			}
		}

		i = j
	}

	if len(d.xCoordinates) != len(d.flags) {
		return errRangeCheck
	}

	// x coordinates.
	for i, x := range d.xCoordinates {
		flag := d.flags[i]

		if simpleGlyphFlag(flag)&xShortVector != 0 {
			err = w.write(uint8(x))
		} else {
			err = w.write(x)
		}
		if err != nil {
			return err
		}
	}

	// y coordinates.
	for i, y := range d.yCoordinates {
		flag := d.flags[i]

		if simpleGlyphFlag(flag)&yShortVector != 0 {
			err = w.write(uint8(y))
		} else {
			err = w.write(y)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

type compositeGlyphFlag uint16

const (
	arg1And2AreWords compositeGlyphFlag = (1 << iota) // If set, the args are 16-bit (uint16/int16), otherwise uint8/int8.
	argsAreXYValues                                   // If set, the args are signed xy values (otherwise unsigned).
	roundXYToGrid
	weHaveAScale
	_              // reserved
	moreComponents // Indicates at least one glyph following this one.
	weHaveAnXAndYScale
	weHaveATwoByTwo
	weHaveInstructions
	useMyMetrics
	overlapCompound
	scaledComponentOffset
	unscaledComponentOffset
)

func (f compositeGlyphFlag) IsSet(flag compositeGlyphFlag) bool {
	return f&flag != 0
}

func (f compositeGlyphFlag) String() string {
	var flags []string

	if f.IsSet(arg1And2AreWords) {
		flags = append(flags, "arg1And2AreWords")
	}
	if f.IsSet(argsAreXYValues) {
		flags = append(flags, "argsAreXYValues")
	}
	if f.IsSet(roundXYToGrid) {
		flags = append(flags, "roundXYToGrid")
	}
	if f.IsSet(weHaveAScale) {
		flags = append(flags, "weHaveAScale")
	}
	if f.IsSet(moreComponents) {
		flags = append(flags, "moreComponents")
	}
	if f.IsSet(weHaveAnXAndYScale) {
		flags = append(flags, "weHaveAnXAndYScale")
	}
	if f.IsSet(weHaveATwoByTwo) {
		flags = append(flags, "weHaveATwoByTwo")
	}
	if f.IsSet(weHaveInstructions) {
		flags = append(flags, "weHaveInstructions")
	}
	if f.IsSet(useMyMetrics) {
		flags = append(flags, "useMyMetrics")
	}
	if f.IsSet(overlapCompound) {
		flags = append(flags, "overlapCompound")
	}
	if f.IsSet(scaledComponentOffset) {
		flags = append(flags, "scaledComponentOffset")
	}
	if f.IsSet(unscaledComponentOffset) {
		flags = append(flags, "unscaledComponentOffset")
	}

	return strings.Join(flags, "|")
}

type compositeGlyphDescription struct {
	components   []compositeGlyphDescriptionComponent
	instructions []uint8
}

type compositeGlyphDescriptionComponent struct {
	flags      uint16
	glyphIndex uint16
	argument1  uint16 // uint8, int8, uint16 or int16.
	argument2  uint16 // uint8, int8, uint16 or int16.

	// Optional transformation flags.
	scale          *f2dot14 // same scale for x and y.
	scaleX, scaleY *f2dot14 // x and y scales
	a, b, c, d     *f2dot14 // 2x2
}

// gdLen is the length of the glyph data record according to the loca table.
// It is used to determine the length of the instructions following the record if present.
func (f *font) parseCompositeGlyphDescription(r *byteReader, gdLen int64) (*compositeGlyphDescription, error) {
	cgd := &compositeGlyphDescription{}

	start := r.Offset()

	instructionsFollow := false
	for {
		comp := compositeGlyphDescriptionComponent{}
		err := r.read(&comp.flags, &comp.glyphIndex)
		if err != nil {
			return nil, err
		}

		if compositeGlyphFlag(comp.flags).IsSet(arg1And2AreWords) {
			var arg1, arg2 uint16
			err := r.read(&arg1, &arg2)
			if err != nil {
				return nil, err
			}
			comp.argument1, comp.argument2 = arg1, arg2
		} else {
			var arg1, arg2 uint8
			err := r.read(&arg1, &arg2)
			if err != nil {
				return nil, err
			}
			comp.argument1, comp.argument2 = uint16(arg1), uint16(arg2)
		}

		if compositeGlyphFlag(comp.flags).IsSet(weHaveAScale) {
			var scale f2dot14
			err := r.read(&scale)
			if err != nil {
				return nil, err
			}
			comp.scale = &scale
		} else if compositeGlyphFlag(comp.flags).IsSet(weHaveAnXAndYScale) {
			var scaleX, scaleY f2dot14
			err := r.read(&scaleX, &scaleY)
			if err != nil {
				return nil, err
			}
			comp.scaleX, comp.scaleY = &scaleX, &scaleY
		} else if compositeGlyphFlag(comp.flags).IsSet(weHaveATwoByTwo) {
			var a, b, c, d f2dot14
			err := r.read(&a, &b, &c, &d)
			if err != nil {
				return nil, err
			}
			comp.a, comp.b, comp.c, comp.d = &a, &b, &c, &d
		}

		if compositeGlyphFlag(comp.flags).IsSet(weHaveInstructions) {
			instructionsFollow = true
		}

		if !compositeGlyphFlag(comp.flags).IsSet(moreComponents) {
			break
		}

		cgd.components = append(cgd.components, comp)
	}

	end := r.Offset()
	len := end - start

	if instructionsFollow {
		instructionLen := int64(gdLen) - len
		if instructionLen < 0 {
			common.Log.Debug("Read more than length in loca table showed")
			return nil, errors.New("read too far")
		}

		err := r.readSlice(&cgd.instructions, int(instructionLen))
		if err != nil {
			common.Log.Debug("Failed to read instructions")
			return nil, err
		}
	}

	return cgd, nil
}

func (cgd compositeGlyphDescription) Write(w *byteWriter, f *font) error {
	for _, comp := range cgd.components {
		err := w.write(comp.flags, comp.glyphIndex)
		if err != nil {
			return err
		}

		if comp.flags&uint16(arg1And2AreWords) != 0 {
			err := w.writeUint16(comp.argument1, comp.argument2)
			if err != nil {
				return err
			}
		} else {
			err := w.writeUint8(uint8(comp.argument1), uint8(comp.argument2))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
