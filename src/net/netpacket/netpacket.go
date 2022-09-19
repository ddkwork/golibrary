package netpacket

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/net/httpClient"
	"github.com/hjson/hjson-go"
)

// ╔═════╤═══════════════════════════════════════╤════════════════════════════════════╤════════════════════════╤══════════════════════════════════╗
// ║index│api                                    │function                            │note                    │todo                              ║
// ╠═════╪═══════════════════════════════════════╪════════════════════════════════════╪════════════════════════╪══════════════════════════════════╣
// ║  1  │RegisterPacketHandles(handles []Handle)│Register Packet Handles as work pool│                        │                                  ║
// ║  2  │SetHandles(handles []Handle)           │RegisterPacketHandles means         │                        │                                  ║
// ║  3  │Handles() []Handle                     │work pool                           │work center             │                                  ║
// ║  4  │HandlePacket() (ok bool)               │rang work pool packet and post them │                        │                                  ║
// ║  5  │SetEvent(event any)            │set a event for any work            │                        │                                  ║
// ║  6  │SetEventsCap(cap int)                  │Set Events Cap                      │                        │                                  ║
// ║  7  │Events() <-chan any            │pop events                          │work life               │                                  ║
// ║  8  │HandleEvent()                          │rang event and handle them          │need handle by your self│                                  ║
// ║  9  │HttpClient() httpClient.Interface      │http client                         │                        │add udp wss etc                   ║
// ║  10 │                                       │                                    │                        │add Worn echo and saveDataBase api║
// ╚═════╧═══════════════════════════════════════╧════════════════════════════════════╧════════════════════════╧══════════════════════════════════╝
type (
	Interface interface {
		RegisterPacketHandles(handles []Handle) //注册或设置业务集
		//SetHandles(handles []Handle)
		Handles() []Handle              //业务集
		HandlePacket() (ok bool)        //遍历业务集并发出每个业务包请求
		SetEvent(event any)             //接受任意类型的消息事件
		SetEventsCap(cap int)           //设置消息事件数量
		Events() <-chan any             //信道先进先出队列存储事件信号
		HandleEvent()                   //todo 多态业务只要分态处理事件即可，即：这里处理如何退出主程序，或者永远不要退出,其余的接口签名通通 panic("implement me") 即可
		HttpClient() *httpClient.Object //http客户端接口（转发，cookie，表单，容错，各种请求方式等的封装），todo add udp，wss
		//todo add Worn echo and saveDataBase api
	}
	object struct {
		httpClient *httpClient.Object
		handles    []Handle
		events     chan any
	}
	Handle struct {
		PacketIndex int              //包序
		Fn          func() (ok bool) `json:"-"` //业务回调代码
		ReqUrl      string           //请求地址
		Info        string           //功能描述，包信息或预期获取内容
	}
)

func (o *object) RegisterPacketHandles(handles []Handle) { o.handles = handles }
func (o *object) SetHandles(handles []Handle)            { o.handles = handles } //means RegisterHandles
func (o *object) Handles() []Handle                      { return o.handles }
func (o *object) HandlePacket() (ok bool) {
	for i, handle := range o.Handles() {
		handle.PacketIndex = i + 1
		marshal, err := hjson.MarshalWithOptions(handle, hjson.EncoderOptions{
			Eol:            "",
			BracesSameLine: false,
			EmitRootBraces: false,
			QuoteAlways:    true,
			IndentBy:       "\t",
			AllowMinusZero: false,
			UnknownAsNull:  false,
		})
		if !mylog.Error(err) {
			return
		}
		mylog.Info("", string(marshal))
		o.HttpClient().Url(handle.ReqUrl)
		if !handle.Fn() {
			return mylog.Error("请检查发包数据结构是否正确")
		}
	}
	return true
}
func (o *object) SetEvent(event any)             { o.events <- event }
func (o *object) SetEventsCap(cap int)           { o.events = make(chan any, cap) }
func (o *object) Events() <-chan any             { return o.events }
func (o *object) HandleEvent()                   { panic("implement me") }
func (o *object) HttpClient() *httpClient.Object { return o.httpClient }

func New() Interface {
	return &object{
		httpClient: httpClient.New(),
		handles:    nil,
		events:     nil,
	}
}
