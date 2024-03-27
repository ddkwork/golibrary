package tcp

import (
	"bytes"
	"fmt"
	"net"

	"github.com/ddkwork/golibrary/mylog"
)

func (o *object) TransportTCP(DstIP string, DstPort int) {
	o.reset(DstIP, DstPort)
	if !o.getSrcAddrConn() {
		return
	}
	defer func() {
		if o.SrcTCPListener == nil {
			mylog.Error("SrcTCPListener == nil ")
			return
		}
		mylog.Error(o.SrcTCPListener.Close())

		if o.DstTCPListener == nil {
			mylog.Error("DstTCPListener == nil ")
			return
		}
		mylog.Error(o.DstTCPListener.Close())
	}()
	for {
		o.BufSize, o.err = o.SrcConn.Read(o.Bytes())
		if !mylog.Error(o.err) {
			return
		}
		SrcBufChan <- o.Bytes()[:o.Len()]
		go o.do()
		if !mylog.Error2(o.SrcConn.Write(<-DstBufChan)) {
			return
		}
	}
}

func (o *object) do() {
	o.DstConn, o.err = net.DialTCP(o.ProtoCool, o.DstAddr, nil)
	if !mylog.Error(o.err) {
		return
	}
	for {
		o.BufSize, o.err = o.DstConn.Write(<-SrcBufChan)
		if !mylog.Error(o.err) {
			return
		}
		o.Reset()
		o.BufSize, o.err = o.DstConn.Read(o.Bytes())
		if !mylog.Error(o.err) {
			return
		}
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

func (o *object) getSrcAddrConn() (ok bool) { // 设置源地址和连接
	o.SrcAddr, o.err = net.ResolveTCPAddr(o.ProtoCool, "0.0.0.0"+":"+fmt.Sprint(o.DstPort))
	if !mylog.Error(o.err) {
		return
	}
	o.SrcTCPListener, o.err = net.ListenTCP(o.ProtoCool, o.SrcAddr)
	if !mylog.Error(o.err) {
		return
	}
	o.SrcConn, o.err = o.SrcTCPListener.AcceptTCP()
	return mylog.Error(o.err)
}
