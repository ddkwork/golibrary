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
		err     error
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
		err:       nil,
	}
}

func (o *object) GetSrcAddrConn() bool { // 设置源地址和连接
	o.SrcConn, o.err = net.ListenUDP(o.ProtoCool, &net.UDPAddr{IP: net.IPv4zero, Port: o.DstPort})
	if !mylog.Error(o.err) {
		return false
	}
	o.BufSize, o.SrcAddr, o.err = o.SrcConn.ReadFromUDP(o.Bytes())
	return mylog.Error(o.err)
}

func (o *object) SetDstAddrConn() bool { // 设置目标地址和连接
	if !o.SetDstAddr() {
		return false
	}
	o.DstConn, o.err = net.DialUDP(o.ProtoCool, nil, o.DstAddr)
	return mylog.Error(o.err)
}

func (o *object) SetDstAddr() (ok bool) {
	if o.DstAddr == nil {
		return mylog.Error("DstAddr == nil")
	}
	if o.BufSize == 48 {
		o.DstAddr = &net.UDPAddr{IP: net.ParseIP(o.DstIP), Port: o.DstPort}
		return
	}
	o.DstAddr = &net.UDPAddr{IP: o.Bytes()[14:18], Port: o.DstPort}
	return true
}
