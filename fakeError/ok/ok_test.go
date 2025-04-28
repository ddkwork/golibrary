package main

import (
	"github.com/ddkwork/golibrary/assert"
	"github.com/ddkwork/golibrary/fakeError"
	"github.com/ddkwork/golibrary/mylog"
	"image/color"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestName(t *testing.T) { //todo
	fakeError.Walk("")
}

var p = packet{
	Scheme:        "tcp",
	Method:        http.MethodGet,
	Host:          "www.baidu.com",
	Path:          "/cmsocket/",
	ContentType:   "json",
	ContentLength: 20,
	Status:        http.StatusText(http.StatusOK),
	Note:          "this is a note",
	Process:       "steam.exe",
	PadTime:       4,
}

func TestMarshalStruct(t *testing.T) {
	s := Struct[packet]{Data: p}
	for k, v := range s.MarshalRow(func(key string) (value string) {
		return ""
		switch key {
		case "ContentLength":
			return strconv.Itoa(p.ContentLength)
		case "PadTime":
			return p.PadTime.String()
		default:
			return ""
		}
	}) {
		mylog.Info(k, v) //fields *safemap.M[string, *Input]
	}
}

func TestUnmarshalStruct(t *testing.T) {
	s := Struct[*packet]{Data: &p}
	cells := s.MarshalRow(nil)
	row := s.UnmarshalRow(cells, func(key string) (value any) {
		return nil
	})
	assert.Equal(t, *row, p)

}

// mylog.Check2(strconv.Atoi(value))
type packet struct {
	Scheme        string        // 请求协议
	Method        string        // 请求方式
	Host          string        // 请求主机
	Path          string        // 请求路径
	ContentType   string        // 收发都有
	ContentLength int           // 收发都有
	Status        string        // 返回的状态码文本
	Note          string        // 注释
	Process       string        // 进程
	PadTime       time.Duration // 请求到返回耗时
}

type CellData struct {
	Value    string      // 单元格文本
	Tooltip  string      // 单元格提示信息
	Icon     []byte      // 单元格图标，格式支持：*giosvg.Icon, *widget.Icon, *widget.Image, image.Image
	FgColor  color.NRGBA // 单元格前景色,着色渲染单元格
	IsNasm   bool        // 是否是nasm汇编代码,为表头提供不同的着色渲染样式
	Disabled bool        // 是否显示表头或者body单元格，或者禁止编辑节点时候使用
	columID  int         // 列id,预计后期用于区域选中,每行的列数和表头是一样的，而表头单元格在init的时候填充了每列的id，那么body rows的列id只要在操作row节点的时候遍历一下表头的row node填充列id即可
	rowID    int         // 行id,预计后期用于区域选中,由layout遍历list渲染的index填充该节点的rowID
	//widget.Clickable             // todo 单元格点击事件,如果换成编辑框，那么编辑框需要支持单机事件和双击编辑回车事件,以及RichText高亮单元格
	//RichText                     // todo 单元格富文本
}
