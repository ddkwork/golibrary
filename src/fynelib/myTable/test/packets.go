package main

import (
	"fmt"
	"github.com/ddkwork/golibrary/src/fynelib/myTable"
	"sort"
	"time"
)

type (
	a interface {
		myTable.Interface
		AddRow(info PacketInfo)
	}
	DecodedInfo struct {
		Head      string
		Pb2       string
		Pb3       string
		Tdf       string
		Taf       string
		Acc       string
		Text      string
		Json      string
		Websocket string
		Msgpack   string
		HexDump   string
	}
	PacketInfo struct {
		PacketIndex int
		Method      string
		Host        string
		Path        string
		ConnectType string
		Size        int
		PadTime     time.Duration
		StartTime   time.Time
		Status      string
		StatusCode  int
		Note        string
		Req         DecodedInfo
		Resp        DecodedInfo
	}
	packets struct {
		data    []PacketInfo
		colTemp [][]string
	}
)

func newPackets() *packets {
	p := &packets{
		data:    make([]PacketInfo, 0),
		colTemp: nil,
	}
	p.colTemp = make([][]string, p.ColumnLen())
	return p
}

func (p *packets) ColumnLen() int { return len(p.Header()) }
func (p *packets) ColumnWidths() []float32 {
	return nil
}
func (p *packets) Append(data any) {
	p.data = append(p.data, data.(PacketInfo))

}
func (p *packets) Len() int { return len(p.data) }

type (
	FieldName string
)

func (FieldName) PacketIndex() string { return "PacketIndex" }
func (FieldName) Method() string      { return "Method" }
func (FieldName) Host() string        { return "Host" }
func (FieldName) Path() string        { return "Path" }
func (FieldName) ConnectType() string { return "ConnectType" }
func (FieldName) Size() string        { return "Size" }
func (FieldName) PadTime() string     { return "PadTime" }
func (FieldName) StartTime() string   { return "StartTime" }
func (FieldName) Status() string      { return "Status" }
func (FieldName) StatusCode() string  { return "StatusCode" }
func (FieldName) Note() string        { return "Note" }

var fieldName FieldName

func (p packets) Header() []string {
	return []string{
		fieldName.PacketIndex(),
		fieldName.Method(),
		fieldName.Host(),
		fieldName.Path(),
		fieldName.ConnectType(),
		fieldName.Size(),
		fieldName.PadTime(),
		fieldName.StartTime(),
		fieldName.Status(),
		fieldName.StatusCode(),
		fieldName.Note(),
	}
}
func (p packets) Rows(id int) []string {
	if id < 0 || id >= p.Len() {
		return nil
	}
	return []string{
		fmt.Sprintf("%03d", p.data[id].PacketIndex),
		p.data[id].Method,
		p.data[id].Host,
		p.data[id].Path,
		p.data[id].ConnectType,
		fmt.Sprint(p.data[id].Size),
		p.data[id].PadTime.String(),
		p.data[id].StartTime.String(),
		p.data[id].Status,
		fmt.Sprint(p.data[id].StatusCode),
		p.data[id].Note,
	}
}
func (p *packets) Sort(field int, ascend bool) {
	switch p.Header()[field] {
	case fieldName.PacketIndex():
		sort.Slice(p, func(i, j int) bool {
			return p.data[i].PacketIndex > p.data[j].PacketIndex
		})
	case fieldName.Method():
		sort.Slice(p, func(i, j int) bool {
			return p.data[i].Method > p.data[j].Method
		})
	case fieldName.Host():
		sort.Slice(p, func(i, j int) bool {
			return p.data[i].Host > p.data[j].Host
		})
	case fieldName.Path():
		sort.Slice(p, func(i, j int) bool {
			return p.data[i].Path > p.data[j].Path
		})
	case fieldName.ConnectType():
		sort.Slice(p, func(i, j int) bool {
			return p.data[i].ConnectType > p.data[j].ConnectType
		})
	case fieldName.Size():
		sort.Slice(p, func(i, j int) bool {
			return p.data[i].Size > p.data[j].Size
		})
	case fieldName.PadTime():
		sort.Slice(p, func(i, j int) bool {
			return p.data[i].PadTime > p.data[j].PadTime
		})
	case fieldName.StartTime():
		sort.Slice(p, func(i, j int) bool {
			return p.data[i].StartTime.Unix() > p.data[j].StartTime.Unix()
		})
	case fieldName.Status():
		sort.Slice(p, func(i, j int) bool {
			return p.data[i].Status > p.data[j].Status
		})
	case fieldName.StatusCode():
		sort.Slice(p, func(i, j int) bool {
			return p.data[i].StatusCode > p.data[j].StatusCode
		})
	case fieldName.Note():
		sort.Slice(p, func(i, j int) bool {
			return p.data[i].Note > p.data[j].Note
		})
	}
}
func (p packets) Filter(kw string, i int) {
	//for lineIndex, s := range p.Rows(i) {
	//	println(s)
	//	if strings.Contains(kw, s) {
	//	}
	//}
	//tocdiag.list.ScrollTo(i)
	//tocdiag.list.SetSelection(i, true
	//save
	//delete
	//update
	//right meanu
	//copy clomn
}
