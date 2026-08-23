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

// Чтение xlsx без сторонней библиотеки.
//
// Нацбанк отдаёт статистику единственным машинным способом — кнопкой «Экспорт»,
// то есть книгой Excel. Ради неё тянуть в проект стороннюю библиотеку не
// стоило: нам нужны только значения ячеек, а книга — это zip с несколькими
// файлами XML, и всё нужное разбирается стандартной библиотекой.

// xlsxSheet — лист как таблица строк: строка → колонка (A, B, …) → значение.
type xlsxSheet []map[string]string

// readXLSX разбирает первый лист книги.
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
				// Ссылка в общий словарь строк: Excel не хранит один и тот же
				// текст дважды.
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

// xlsxShared читает общий словарь строк книги.
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

// xlsxFile достаёт один файл из книги.
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

// xlsxColRe выделяет буквенную часть адреса ячейки: из «BC12» — «BC».
var xlsxColRe = regexp.MustCompile(`^[A-Z]+`)

// xlsxCol возвращает букву колонки из адреса ячейки.
func xlsxCol(ref string) string { return xlsxColRe.FindString(ref) }

// xlsxColNum переводит букву колонки в номер: A → 1, Z → 26, AA → 27. Нужен
// для сортировки: в книге Нацбанка месяцы идут по колонкам, и «AA» без этого
// встаёт между «A» и «B».
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
