package client

import (
	"encoding/json"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/stream"
	"net"
)

var receive = make(chan *stream.Stream) //todo sync

type (
	Interface interface { //stickyBag() //by short connection way
		Connect(address string) bool
		SendJson(objectPtr any) (ok bool)                      //call send
		SendJsonWithHead(head string, objectPtr any) (ok bool) //call send
		SendWithHead(head, body *stream.Stream) (ok bool)      //call send
		Send(s *stream.Stream) (ok bool)
		Receive() *stream.Stream
		//MarshalIndent(objectPtr any) *stream.Stream
	}
	object struct {
		conn    net.Conn
		data    *stream.Stream
		err     error
		address string
		ok      bool
	}
)

func New() Interface {
	return &object{
		conn: nil,
		data: stream.New(),
		err:  nil,
	}
}

var Default = New()

func (o *object) MarshalIndent(objectPtr any) *stream.Stream {
	send, err := json.MarshalIndent(objectPtr, " ", " ")
	o.ok = mylog.Error(err)
	return stream.NewBytes(send)
}

func (o *object) SendJsonWithHead(head string, objectPtr any) (ok bool) {
	s := stream.NewString(head)
	marshalIndent := o.MarshalIndent(objectPtr)
	if !o.ok {
		return
	}
	return o.SendWithHead(s, marshalIndent)
}

func (o *object) SendJson(objectPtr any) (ok bool) {
	marshalIndent := o.MarshalIndent(objectPtr)
	if !o.ok {
		return
	}
	return o.Send(marshalIndent)
}

func (o *object) SendWithHead(head, body *stream.Stream) (ok bool) {
	buffer := stream.NewNil()
	buffer.Append(head, body)
	return o.Send(buffer)
}

func (o *object) stickyBag() bool {
	if !mylog.Error(o.conn.Close()) {
		return false
	}
	return o.Connect(o.address)
}

func (o *object) Connect(address string) bool {
	o.address = address
	o.conn, o.err = net.Dial("tcp", address)
	if !mylog.Error(o.err) {
		return false
	}
	mylog.Info("Client send to", address)
	return true
}

func (o *object) Receive() *stream.Stream {
	return <-receive
}

func (o *object) Send(s *stream.Stream) (ok bool) {
	if !mylog.Error2(o.conn.Write(s.Bytes())) {
		return false
	}
	mylog.Json("Client send stream", s.String())
	go func() {
		for {
			//o.data.Reset()
			n, err := o.conn.Read(o.data.Bytes())
			if err != nil || n == 0 {
				continue
			}
			data := make([]byte, n)
			copy(data, o.data.Bytes()[:n])
			s := stream.NewBytes(data)
			receive <- s //for sync every message
			if !o.stickyBag() {
				//todo handle sticky bag
			}
		}
	}()
	return true
}
