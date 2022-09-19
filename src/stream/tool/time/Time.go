package time

import (
	"strconv"
	"time"
)

type (
	Interface interface {
		GetTimeStamp13Bits() int64 //获取指定长度的时间戳
		GetTimeStamp() string      //获取时间戳
		GetTimeNowString() string
	}
	object struct{}
)

func New() Interface { return &object{} }

func (p *object) GetTimeNowString() string { return time.Now().Format("2006-01-02 15:04:05 ") }

func (p *object) GetTimeStamp13Bits() int64 { return time.Now().UnixNano() / 1000000 }

func (p *object) GetTimeStamp() string { return strconv.FormatInt(time.Now().UnixNano()/1000000, 10) }
