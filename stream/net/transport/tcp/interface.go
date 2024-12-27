package tcp

import (
	"bytes"
	"net"
)

var (
	SrcBufChan = make(chan []byte, 1)
	DstBufChan = make(chan []byte, 1)
)

type (
	Interface interface {
		TransportTCP(DstIP string, DstPort int)
	}
	object struct {
		ProtoCool string
		*bytes.Buffer
		BufSize int
		srcTransportCtx
		dstTransportCtx
	}
	srcTransportCtx struct {
		SrcConn        *net.TCPConn
		SrcAddr        *net.TCPAddr
		SrcTCPListener *net.TCPListener
	}
	dstTransportCtx struct {
		DstIP          string
		DstPort        int
		DstConn        *net.TCPConn
		DstAddr        *net.TCPAddr
		DstTCPListener *net.TCPListener
	}
)

func New() Interface {
	return &object{}
}
