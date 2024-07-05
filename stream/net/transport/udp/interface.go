package udp

import (
	"bytes"
	"net"

	"github.com/ddkwork/golibrary/mylog"
)

var (
	SrcBufChan = make(chan []byte, 1)
	DstBufChan = make(chan []byte, 1)
)

type (
	Interface interface {
		TransportUDP(DstIP string, DstPort int)
		TransportUDPNoChan(DstIP string, DstPort int)
	}
	object struct {
		SrcConn   *net.UDPConn
		SrcAddr   *net.UDPAddr
		DstIP     string
		DstPort   int
		DstConn   *net.UDPConn
		DstAddr   *net.UDPAddr
		ProtoCool string
		*bytes.Buffer
		BufSize int
	}
)

func New() Interface {
	return &object{}
}

func (o *object) reset(DstIP string, DstPort int) {
	*o = object{
		SrcConn:   nil,
		SrcAddr:   nil,
		DstIP:     DstIP,
		DstPort:   DstPort,
		DstConn:   nil,
		DstAddr:   nil,
		ProtoCool: "udp",
		Buffer:    bytes.NewBuffer(nil),
		BufSize:   0,
	}
}

func (o *object) GetSrcAddrConn() {
	o.SrcConn = mylog.Check2(net.ListenUDP(o.ProtoCool, &net.UDPAddr{IP: net.IPv4zero, Port: o.DstPort}))
	n, addr := mylog.Check3(o.SrcConn.ReadFrom(o.Bytes()))
	o.BufSize, o.SrcAddr = n, addr.(*net.UDPAddr)
}

func (o *object) SetDstAddrConn() {
	o.SetDstAddr()
	o.DstConn = mylog.Check2(net.DialUDP(o.ProtoCool, nil, o.DstAddr))
}

func (o *object) SetDstAddr() {
	mylog.Check(o.DstAddr == nil)
	if o.BufSize == 48 {
		o.DstAddr = &net.UDPAddr{IP: net.ParseIP(o.DstIP), Port: o.DstPort}
		return
	}
	o.DstAddr = &net.UDPAddr{IP: o.Bytes()[14:18], Port: o.DstPort}
}
