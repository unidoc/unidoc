/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

// TODO: Export only what unidoc needs:
// Encoding: rune <-> GID map.
// font flags:
//		IsFixedPitch, Serif, etc (Table 123 PDF32000_2008 - font flags)
//		FixedPitch() bool
//		Serif() bool
//		Symbolic() bool
//		Script() bool
//		Nonsymbolic() bool
//		Italic() bool
//		AllCap() bool
//		SmallCap() bool
//		ForceBold() bool
//      Need to be able to derive the font flags from the font to build a font descriptor
//
// Required table according to PDF32000_2008 (9.9 Embedded font programs - p. 299):
// “head”, “hhea”, “loca”, “maxp”, “cvt”, “prep”, “glyf”, “hmtx”, and “fpgm”. If used with a simple
// font dictionary, the font program shall additionally contain a cmap table defining one or more
// encodings, as discussed in 9.6.6.4, "Encodings for TrueType Fonts". If used with a CIDFont
// dictionary, the cmap table is not needed and shall not be present, since the mapping from
// character codes to glyph descriptions is provided separately.
//

// font is a data model for truetype fonts with basic access methods.
type font struct {
	ot   *offsetTable
	trec *tableRecords // table records (references other tables).
	head *headTable
	maxp *maxpTable
	hhea *hheaTable
	hmtx *hmtxTable
	loca *locaTable
	glyf *glyfTable
	name *nameTable
	os2  *os2Table
	post *postTable

	/*
	*fpgmTable
	*cmapTable
	 */
}

func (f font) numTables() int {
	return int(f.ot.numTables)
}

func parseFont(r *byteReader) (*font, error) {
	f := &font{}

	var err error

	f.ot, err = f.parseOffsetTable(r)
	if err != nil {
		return nil, err
	}

	f.trec, err = f.parseTableRecords(r)
	if err != nil {
		return nil, err
	}

	f.head, err = f.parseHead(r)
	if err != nil {
		return nil, err
	}

	f.maxp, err = f.parseMaxp(r)
	if err != nil {
		return nil, err
	}

	f.hhea, err = f.parseHhea(r)
	if err != nil {
		return nil, err
	}

	f.hmtx, err = f.parseHmtx(r)
	if err != nil {
		return nil, err
	}

	f.loca, err = f.parseLoca(r)
	if err != nil {
		return nil, err
	}

	f.glyf, err = f.parseGlyf(r)
	if err != nil {
		return nil, err
	}

	f.name, err = f.parseNameTable(r)
	if err != nil {
		return nil, err
	}

	f.os2, err = f.parseOS2Table(r)
	if err != nil {
		return nil, err
	}

	f.post, err = f.parsePost(r)
	if err != nil {
		return nil, err
	}
	/*
		if f.os2 != nil {
			fmt.Printf("OS2: %+v\n", *f.os2)
		}
	*/

	return f, nil
}

func (f *font) write(w *byteWriter) error {

	// TODO(gunnsth): Do in two steps:
	//    1. Write the content tables: head, hhea, etc in the expected order and keep track of the length, checksum for each.
	//    2. Generate the table records based on the information.
	//    3. Write out in final order: offset table, table records, head, ...
	//    4. Set checkAdjustment of head table based on checksumof entire file
	//    5. Write the final output

	err := f.writeOffsetTable(w)
	if err != nil {
		return err
	}

	err = f.writeTableRecords(w)
	if err != nil {
		return err
	}

	err = f.writeHead(w)
	if err != nil {
		return err
	}

	err = f.writeMaxp(w)
	if err != nil {
		return err
	}

	err = f.writeHhea(w)
	if err != nil {
		return err
	}

	err = f.writeLoca(w)
	if err != nil {
		return err
	}

	err = f.writeGlyf(w)
	if err != nil {
		return err
	}

	return nil
}
