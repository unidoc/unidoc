/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package truetype

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/unidoc/unidoc/common"
)

// tableRecord represents table records, including name (tag) and file offset, size
// and checksum for integrity checking.
type tableRecord struct {
	tableTag tag
	checksum uint32
	offset   offset32
	length   uint32
}

func (tr *tableRecord) read(r *byteReader) error {
	return r.read(&tr.tableTag, &tr.checksum, &tr.offset, &tr.length)
}

func (tr tableRecord) write(w *byteWriter) error {
	return w.write(tr.tableTag, tr.checksum, tr.offset, tr.length)
}

// tableRecords represents a set of table records in a truetype font file.
// Includes a map by table name for quick lookup of records.
type tableRecords struct {
	list  []tableRecord
	trMap map[string]tableRecord
}

func (f *font) parseTableRecords(r *byteReader) (*tableRecords, error) {
	trs := &tableRecords{}

	numTables := int(f.ot.numTables)
	if numTables < 0 {
		common.Log.Debug("Invalid number of tables")
		return nil, errRangeCheck
	}

	if trs.trMap == nil {
		trs.trMap = map[string]tableRecord{}
	}

	for i := 0; i < numTables; i++ {
		var rec tableRecord
		err := rec.read(r)
		if err != nil {
			return nil, err
		}
		trs.list = append(trs.list, rec)
		trs.trMap[rec.tableTag.String()] = rec
	}

	return trs, nil
}

// seekToTable seeks to position font table `tableName` in `r` if it has the table.
// The table record is returned back when successful, otherwise is meaningless.
// The bool flag indicates that the table exists and should be at that position if there
// was no error.
func (f *font) seekToTable(r *byteReader, tableName string) (tr tableRecord, has bool, err error) {
	tr, has = f.trec.trMap[tableName]
	if !has {
		return tr, false, nil
	}

	err = r.Seek(int64(tr.offset))
	if err != nil {
		return tr, false, err
	}

	return tr, true, nil
}

func (f *font) writeTableRecords(w *byteWriter) error {
	if f.trec == nil {
		common.Log.Debug("Table records not set")
		return errRequiredField
	}

	for _, tr := range f.trec.list {
		err := tr.write(w)
		if err != nil {
			return err
		}
	}
	return nil
}

// HasTable returns true if there is a record of `tableName` in table records `trs`.
func (trs *tableRecords) HasTable(tableName string) bool {
	_, has := trs.trMap[strings.TrimSpace(tableName)]
	return has
}

func (trs *tableRecords) String() string {
	var buf bytes.Buffer
	for i, tr := range trs.list {
		buf.WriteString(fmt.Sprintf("Table record %d: %+v\n", i+1, tr))
		buf.WriteString(fmt.Sprintf("%s\n", tr.tableTag))
	}
	return buf.String()
}
