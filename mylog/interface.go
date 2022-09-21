package mylog

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	Interface interface {
		Error(err any) bool
		Error2(_ any, err error) bool
		HexDump(title string, msg any)     //hex buf todo support fn return []byte
		Hex(title string, msg any)         //hex value
		Info(title string, msg ...any)     //info
		Trace(title string, msg ...any)    //跟踪
		Warning(title string, msg ...any)  //警告
		MarshalJson(title string, msg any) //pb json todo rename
		Json(title string, msg ...any)     //pb json todo rename
		Success(title string, msg ...any)  //成功
		Struct(msg any)                    //结构体 todo indent ,add title
		Body() string
		Msg() string
		SetDebug(debug bool)
	}
	object struct {
		kind  kind
		title string
		msg   string
		body  string
		debug bool
	}
)

func (o *object) MarshalJson(title string, msg any) {
	indent, err := json.MarshalIndent(msg, "", " ")
	if !o.Error(err) {
		return
	}
	o.Info(title, string(indent))
}

func (o *object) SetDebug(debug bool) { o.debug = debug }
func (o *object) Msg() string         { return o.msg }
func (o *object) Body() string        { return o.body }

func New() Interface {
	return &object{
		kind:  -1,
		title: "",
		msg:   "",
		body:  "",
		debug: true,
	}
}
func init() {
	Trace("--------- title ---------", "------------------ info ------------------")
}

var Default = New()

func Assert(t *testing.T) *assert.Assertions { return assert.New(t) }
func Error(err any) bool                     { return Default.Error(err) }
func Error2(_ any, err error) bool           { return Default.Error2(nil, err) }
func HexDump(title string, msg any)          { Default.HexDump(title, msg) }
func Hex(title string, msg any)              { Default.Hex(title, msg) }
func Info(title string, msg ...any)          { Default.Info(title, msg) }
func Trace(title string, msg ...any)         { Default.Trace(title, msg) }
func Warning(title string, msg ...any)       { Default.Warning(title, msg) }
func MarshalJson(title string, msg any)      { Default.MarshalJson(title, msg) }
func Json(title string, msg ...any)          { Default.Json(title, msg) }
func Success(title string, msg ...any)       { Default.Success(title, msg) }
func Struct(msg any)                         { Default.Struct(msg) }
func Body() string                           { return Default.Body() }
func Msg() string                            { return Default.Msg() }
func SetDebug(debug bool)                    { Default.SetDebug(debug) }
