package server

import (
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/src/stream"
	"net"
)

var receive = make(chan *stream.Stream)

type (
	Interface interface {
		ListenAndServer(address string) bool
		Receive() *stream.Stream
		Replay(data string) bool
	}
	object struct {
		data *stream.Stream
		l    net.Listener
		conn net.Conn
		err  error
	}
)

func (o *object) Replay(data string) bool {
	return mylog.Error2(o.conn.Write([]byte(data)))
}

func (o *object) Receive() *stream.Stream { return <-receive }

func (o *object) ListenAndServer(address string) bool {
	if !o.Listen(address) {
		return false
	}
	go o.Server()
	return true
}

func (o *object) Listen(address string) bool {
	o.l, o.err = net.Listen("tcp", address)
	if !mylog.Error(o.err) {
		return false
	}
	mylog.Info("Server Listen on", address)
	return true
}

func (o *object) Server() {
	//wg := sync.WaitGroup{}
	for {
		o.conn, o.err = o.l.Accept()
		if !mylog.Error(o.err) {
			continue
		}
		//o.data.Reset()
		//go func() {
		//    defer func() { mylog.Error(conn.Close()) }()
		n, err := o.conn.Read(o.data.Bytes())
		if err != nil || n == 0 {
			continue
		}
		data := make([]byte, n)
		copy(data, o.data.Bytes()[:n])
		s := stream.NewBytes(data)
		receive <- s
		//}()
	}
}

func New() Interface {
	return &object{
		data: stream.New(),
	}
}

var Default = New()
