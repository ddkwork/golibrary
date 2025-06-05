package netpacket

import (
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream"
	"github.com/ddkwork/golibrary/std/stream/net/httpClient"
)

type (
	Interface interface {
		RegisterPacketHandles(handles []Handle)

		Handles() []Handle
		HandlePacket() (ok bool)
		SetEvent(event any)
		SetEventsCap(cap int)
		Events() <-chan any
		HandleEvent()
		HttpClient() *httpClient.Client
	}
	object struct {
		httpClient *httpClient.Client
		handles    []Handle
		events     chan any
	}
	Handle struct {
		PacketIndex int
		Fn          func() (ok bool) `json:"-"`
		ReqUrl      string
		Info        string
	}
)

func (o *object) RegisterPacketHandles(handles []Handle) { o.handles = handles }
func (o *object) SetHandles(handles []Handle)            { o.handles = handles }
func (o *object) Handles() []Handle                      { return o.handles }

func (o *object) HandlePacket() (ok bool) {
	for i, handle := range o.Handles() {
		handle.PacketIndex = i + 1
		mylog.Info("", stream.MarshalJSON(handle))
		o.HttpClient().Url(handle.ReqUrl)
		if !handle.Fn() {
			mylog.Check("请检查发包数据结构是否正确")
		}
	}
	return true
}
func (o *object) SetEvent(event any)             { o.events <- event }
func (o *object) SetEventsCap(cap int)           { o.events = make(chan any, cap) }
func (o *object) Events() <-chan any             { return o.events }
func (o *object) HandleEvent()                   { panic("implement me") }
func (o *object) HttpClient() *httpClient.Client { return o.httpClient }

func New() Interface {
	return &object{
		httpClient: httpClient.New(),
		handles:    nil,
		events:     nil,
	}
}
