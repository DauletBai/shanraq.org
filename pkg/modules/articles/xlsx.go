package articles

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Reading xlsx without a third-party library.
//
// The National Bank serves its statistics by exactly one machine-readable route —
// the "Export" button, which is an Excel workbook. Pulling a library into the
// project for it was not worth it: all we need are cell values, and a workbook is
// a zip holding a few XML files, every part of which the standard library
// parses.

// xlsxSheet is a sheet as a table of rows: row → column (A, B, …) → value.
type xlsxSheet []map[string]string

// readXLSX parses a workbook's first sheet.
func readXLSX(data []byte) (xlsxSheet, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("книга не открывается: %w", err)
	}

	shared, err := xlsxShared(zr)
	if err != nil {
		return nil, err
	}

	var sheetName string
	names := []string{}
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "xl/worksheets/sheet") && strings.HasSuffix(f.Name, ".xml") {
			names = append(names, f.Name)
		}
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("в книге нет листов")
	}
	sort.Strings(names)
	sheetName = names[0]

	raw, err := xlsxFile(zr, sheetName)
	if err != nil {
		return nil, err
	}

	var doc struct {
		Rows []struct {
			Cells []struct {
				Ref    string `xml:"r,attr"`
				Type   string `xml:"t,attr"`
				Value  string `xml:"v"`
				Inline string `xml:"is>t"`
			} `xml:"c"`
		} `xml:"sheetData>row"`
	}
	if err := xml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("разбор листа: %w", err)
	}

	out := make(xlsxSheet, 0, len(doc.Rows))
	for _, r := range doc.Rows {
		row := make(map[string]string, len(r.Cells))
		for _, c := range r.Cells {
			v := c.Value
			switch c.Type {
			case "s":
				// A reference into the shared string table: Excel does not store
				// the same text twice.
				if i, err := strconv.Atoi(v); err == nil && i >= 0 && i < len(shared) {
					v = shared[i]
				}
			case "inlineStr":
				v = c.Inline
			}
			row[xlsxCol(c.Ref)] = strings.TrimSpace(v)
		}
		out = append(out, row)
	}
	return out, nil
}

// xlsxShared reads a workbook's shared string table.
func xlsxShared(zr *zip.Reader) ([]string, error) {
	raw, err := xlsxFile(zr, "xl/sharedStrings.xml")
	if err != nil {
		return nil, nil // словаря может не быть, если весь текст встроенный
	}
	var doc struct {
		Items []struct {
			Runs []string `xml:"r>t"`
			Text string   `xml:"t"`
		} `xml:"si"`
	}
	if err := xml.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("разбор словаря строк: %w", err)
	}
	out := make([]string, 0, len(doc.Items))
	for _, it := range doc.Items {
		if len(it.Runs) > 0 {
			out = append(out, strings.Join(it.Runs, ""))
			continue
		}
		out = append(out, it.Text)
	}
	return out, nil
}

// xlsxFile takes one file out of the workbook.
func xlsxFile(zr *zip.Reader, name string) ([]byte, error) {
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return io.ReadAll(io.LimitReader(rc, 64<<20))
	}
	return nil, fmt.Errorf("в книге нет %s", name)
}

// xlsxColRe picks the letter part out of a cell address: "BC" from "BC12".
var xlsxColRe = regexp.MustCompile(`^[A-Z]+`)

// xlsxCol returns the column letter from a cell address.
func xlsxCol(ref string) string { return xlsxColRe.FindString(ref) }

// xlsxColNum turns a column letter into a number: A → 1, Z → 26, AA → 27. Needed
// for sorting: in the National Bank's workbook months run across the columns, and
// without this "AA" lands between "A" and "B".
func xlsxColNum(col string) int {
	n := 0
	for _, r := range col {
		if r < 'A' || r > 'Z' {
			return n
		}
		n = n*26 + int(r-'A') + 1
	}
	return n
}
