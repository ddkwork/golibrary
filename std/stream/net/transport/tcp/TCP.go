package tcp

import (
	"bytes"
	"fmt"
	"net"

	"github.com/ddkwork/golibrary/std/mylog"
)

func (o *object) TransportTCP(DstIP string, DstPort int) {
	o.reset(DstIP, DstPort)
	o.getSrcAddrConn()
	defer func() {
		mylog.Check(o.SrcTCPListener == nil)
		mylog.Check(o.SrcTCPListener.Close())
		mylog.Check(o.DstTCPListener == nil)
		mylog.Check(o.DstTCPListener.Close())
	}()

	for {
		o.BufSize = mylog.Check2(o.SrcConn.Read(o.Bytes()))
		SrcBufChan <- o.Bytes()[:o.Len()]
		go o.do()
		mylog.Check2(o.SrcConn.Write(<-DstBufChan))
	}
}

func (o *object) do() {
	tcpConn := mylog.Check2(net.DialTCP(o.ProtoCool, o.DstAddr, nil))
	o.DstConn = tcpConn
	for {
		writeLen := mylog.Check2(o.DstConn.Write(<-SrcBufChan))
		o.BufSize = writeLen
		o.Reset()
		readLen := mylog.Check2(o.DstConn.Read(o.Bytes()))
		o.BufSize = readLen
		DstBufChan <- o.Bytes()[:o.Len()]
	}
}

func (o *object) reset(DstIP string, DstPort int) {
	*o = object{
		ProtoCool:       "tcp",
		Buffer:          bytes.NewBuffer(nil),
		BufSize:         0,
		srcTransportCtx: srcTransportCtx{},
		dstTransportCtx: dstTransportCtx{
			DstIP:          DstIP,
			DstPort:        DstPort,
			DstConn:        nil,
			DstAddr:        nil,
			DstTCPListener: nil,
		},
	}
}

func (o *object) getSrcAddrConn() {
	o.SrcAddr = mylog.Check2(net.ResolveTCPAddr(o.ProtoCool, "0.0.0.0"+":"+fmt.Sprint(o.DstPort)))
	o.SrcTCPListener = mylog.Check2(net.ListenTCP(o.ProtoCool, o.SrcAddr))
	o.SrcConn = mylog.Check2(o.SrcTCPListener.AcceptTCP())
}
