package xlsx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream/txt"
)

type Sheet struct {
	Name  string
	Min   Ref
	Max   Ref
	Cells map[Ref]Cell
}

func Load(path string) ([]Sheet, error) {
	r := mylog.Check2(zip.OpenReader(path))

	defer r.Close()
	return load(&r.Reader)
}

func Read(in io.ReaderAt, size int64) ([]Sheet, error) {
	r := mylog.Check2(zip.NewReader(in, size))

	return load(r)
}

func load(r *zip.Reader) ([]Sheet, error) {
	var sheetNames []string
	var strs []string
	var files []*zip.File

	for _, f := range r.File {
		switch {
		case f.Name == "docProps/app.xml":
			sheetNames = mylog.Check2(loadSheetNames(f))
		case f.Name == "xl/sharedStrings.xml":
			strs = mylog.Check2(loadStrings(f))
		case strings.HasPrefix(f.Name, "xl/worksheets/sheet"):
			files = append(files, f)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return txt.NaturalLess(files[i].Name, files[j].Name, true)
	})
	sheets := make([]Sheet, 0, len(files))
	for i, f := range files {
		sheet := mylog.Check2(loadSheet(f, strs))
		if i < len(sheetNames) {
			sheet.Name = sheetNames[i]
		} else {
			sheet.Name = fmt.Sprintf("Sheet%d", i+1)
		}
		sheets = append(sheets, *sheet)
	}
	return sheets, nil
}

func loadSheetNames(f *zip.File) ([]string, error) {
	fr := mylog.Check2(f.Open())

	defer fr.Close()
	decoder := xml.NewDecoder(fr)
	var data struct {
		Names []string `xml:"TitlesOfParts>vector>lpstr"`
	}
	mylog.Check(decoder.Decode(&data))

	return data.Names, nil
}

func loadStrings(f *zip.File) ([]string, error) {
	fr := mylog.Check2(f.Open())

	defer fr.Close()
	decoder := xml.NewDecoder(fr)
	var data struct {
		SST []string `xml:"si>t"`
	}
	mylog.Check(decoder.Decode(&data))

	return data.SST, nil
}

func loadSheet(f *zip.File, strs []string) (*Sheet, error) {
	fr := mylog.Check2(f.Open())
	defer fr.Close()
	decoder := xml.NewDecoder(fr)
	var data struct {
		Cells []struct {
			Label string  `xml:"r,attr"`
			Type  string  `xml:"t,attr"`
			Value *string `xml:"v"`
		} `xml:"sheetData>row>c"`
	}
	mylog.Check(decoder.Decode(&data))

	sheet := &Sheet{
		Min:   Ref{Row: math.MaxInt32, Col: math.MaxInt32},
		Max:   Ref{},
		Cells: make(map[Ref]Cell, len(data.Cells)),
	}
	for _, one := range data.Cells {
		if one.Value == nil {
			continue
		}
		ref := ParseRef(one.Label)
		cell := Cell{Value: *one.Value}
		switch one.Type {
		case "s":
			v := mylog.Check2(strconv.Atoi(cell.Value))
			if v >= 0 && v < len(strs) {
				cell.Value = strs[v]
			} else {
				cell.Value = "#REF!"
			}
			cell.Type = String
		case "b":
			cell.Type = Boolean
		default:
			cell.Type = Number
		}
		if sheet.Min.Row > ref.Row {
			sheet.Min.Row = ref.Row
		}
		if sheet.Min.Col > ref.Col {
			sheet.Min.Col = ref.Col
		}
		if sheet.Max.Row < ref.Row {
			sheet.Max.Row = ref.Row
		}
		if sheet.Max.Col < ref.Col {
			sheet.Max.Col = ref.Col
		}
		sheet.Cells[ref] = cell
	}
	if sheet.Min.Row > sheet.Max.Row {
		sheet.Min.Row = sheet.Max.Row
	}
	if sheet.Min.Col > sheet.Max.Col {
		sheet.Min.Col = sheet.Max.Col
	}
	return sheet, nil
}
